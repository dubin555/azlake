package azcat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidValue  = errors.New("invalid value")
)

// Catalog provides versioned data operations backed by a KV store
type Catalog struct {
	kv      KV
	storage ObjectStorage
}

// NewCatalog creates a new Catalog with the given KV backend and object storage
func NewCatalog(kv KV, storage ObjectStorage) *Catalog {
	return &Catalog{kv: kv, storage: storage}
}

func (c *Catalog) Close() error {
	return c.kv.Close()
}

// Storage returns the underlying ObjectStorage backend
func (c *Catalog) Storage() ObjectStorage {
	return c.storage
}

// ── Repository ──

func repoKey(name string) string { return "repos/" + name }

func (c *Catalog) CreateRepository(name, storageNamespace, storageID, defaultBranch string) (*Repository, error) {
	if name == "" {
		return nil, fmt.Errorf("repository name: %w", ErrInvalidValue)
	}
	if defaultBranch == "" {
		defaultBranch = "main"
	}
	key := repoKey(name)
	if _, err := c.kv.Get(key); err == nil {
		return nil, fmt.Errorf("repository %q: %w", name, ErrAlreadyExists)
	}
	repo := &Repository{
		Name:             name,
		StorageNamespace: storageNamespace,
		StorageID:        storageID,
		DefaultBranch:    defaultBranch,
		CreationDate:     time.Now(),
	}
	if err := c.kv.SetJSON(key, repo); err != nil {
		return nil, err
	}
	// Create the initial empty commit
	initCommit := &Commit{
		ID:           generateID("commit"),
		Message:      "Repository created",
		Committer:    "system",
		CreationDate: time.Now(),
	}
	if err := c.kv.SetJSON(commitKey(name, initCommit.ID), initCommit); err != nil {
		return nil, err
	}
	// Create the default branch pointing to the initial commit
	branch := &Branch{Name: defaultBranch, CommitID: initCommit.ID}
	if err := c.kv.SetJSON(branchKey(name, defaultBranch), branch); err != nil {
		return nil, err
	}
	return repo, nil
}

func (c *Catalog) GetRepository(name string) (*Repository, error) {
	var repo Repository
	if err := c.kv.GetJSON(repoKey(name), &repo); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("repository %q: %w", name, ErrNotFound)
		}
		return nil, err
	}
	return &repo, nil
}

func (c *Catalog) ListRepositories(after string, amount int) ([]*Repository, bool, error) {
	if amount <= 0 {
		amount = 100
	}
	var repos []*Repository
	prefix := "repos/"
	scanAfter := ""
	if after != "" {
		scanAfter = repoKey(after)
	}
	err := c.kv.Scan(prefix, scanAfter, amount+1, func(key string, value []byte) error {
		// Only match top-level repo keys (repos/{name}, not repos/{name}/branches/...)
		rest := strings.TrimPrefix(key, prefix)
		if strings.Contains(rest, "/") {
			return nil
		}
		var repo Repository
		if err := unmarshalJSON(value, &repo); err != nil {
			return nil
		}
		repos = append(repos, &repo)
		return nil
	})
	if err != nil {
		return nil, false, err
	}
	hasMore := len(repos) > amount
	if hasMore {
		repos = repos[:amount]
	}
	return repos, hasMore, nil
}

func (c *Catalog) DeleteRepository(name string) error {
	if _, err := c.GetRepository(name); err != nil {
		return err
	}
	// Delete all keys under repos/{name}
	if err := c.kv.DeletePrefix("repos/" + name + "/"); err != nil {
		return err
	}
	return c.kv.Delete(repoKey(name))
}

// ── Branches ──

func branchKey(repo, branch string) string {
	return fmt.Sprintf("repos/%s/branches/%s", repo, branch)
}

func (c *Catalog) CreateBranch(repo, name, sourceRef string) (*Branch, error) {
	if _, err := c.GetRepository(repo); err != nil {
		return nil, err
	}
	key := branchKey(repo, name)
	if _, err := c.kv.Get(key); err == nil {
		return nil, fmt.Errorf("branch %q: %w", name, ErrAlreadyExists)
	}
	// Resolve sourceRef — could be a branch name or commit ID
	commitID, err := c.resolveRef(repo, sourceRef)
	if err != nil {
		return nil, fmt.Errorf("source ref %q: %w", sourceRef, err)
	}
	branch := &Branch{Name: name, CommitID: commitID}
	if err := c.kv.SetJSON(key, branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (c *Catalog) GetBranch(repo, name string) (*Branch, error) {
	var branch Branch
	if err := c.kv.GetJSON(branchKey(repo, name), &branch); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("branch %q: %w", name, ErrNotFound)
		}
		return nil, err
	}
	return &branch, nil
}

func (c *Catalog) ListBranches(repo, after string, amount int) ([]*Branch, bool, error) {
	if amount <= 0 {
		amount = 100
	}
	var branches []*Branch
	prefix := fmt.Sprintf("repos/%s/branches/", repo)
	scanAfter := ""
	if after != "" {
		scanAfter = branchKey(repo, after)
	}
	err := c.kv.Scan(prefix, scanAfter, amount+1, func(key string, value []byte) error {
		var b Branch
		if err := unmarshalJSON(value, &b); err != nil {
			return nil
		}
		branches = append(branches, &b)
		return nil
	})
	if err != nil {
		return nil, false, err
	}
	hasMore := len(branches) > amount
	if hasMore {
		branches = branches[:amount]
	}
	return branches, hasMore, nil
}

func (c *Catalog) DeleteBranch(repo, name string) error {
	r, err := c.GetRepository(repo)
	if err != nil {
		return err
	}
	if r.DefaultBranch == name {
		return fmt.Errorf("cannot delete default branch %q", name)
	}
	if _, err := c.GetBranch(repo, name); err != nil {
		return err
	}
	// Also clean up staging area for this branch
	_ = c.kv.DeletePrefix(stagingKey(repo, name, ""))
	return c.kv.Delete(branchKey(repo, name))
}

// ── Commits ──

func commitKey(repo, id string) string {
	return fmt.Sprintf("repos/%s/commits/%s", repo, id)
}

func (c *Catalog) GetCommit(repo, commitID string) (*Commit, error) {
	var commit Commit
	if err := c.kv.GetJSON(commitKey(repo, commitID), &commit); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("commit %q: %w", commitID, ErrNotFound)
		}
		return nil, err
	}
	return &commit, nil
}

// Commit creates a new commit from staged changes on a branch
func (c *Catalog) Commit(repo, branch, message, committer string, metadata map[string]string) (*Commit, error) {
	br, err := c.GetBranch(repo, branch)
	if err != nil {
		return nil, err
	}

	// Collect staged entries
	staged, err := c.listStaging(repo, branch)
	if err != nil {
		return nil, err
	}
	if len(staged) == 0 {
		return nil, fmt.Errorf("no changes to commit")
	}

	commitID := generateID("commit")
	commit := &Commit{
		ID:           commitID,
		Message:      message,
		Committer:    committer,
		CreationDate: time.Now(),
		Metadata:     metadata,
		Parents:      []string{br.CommitID},
	}

	// Move staged entries to committed objects under this commit
	for _, entry := range staged {
		objKey := objectKey(repo, commitID, entry.Path)
		if err := c.kv.SetJSON(objKey, entry); err != nil {
			return nil, err
		}
	}

	// Also copy forward all objects from parent commit that are NOT overwritten
	parentObjects, _ := c.listObjects(repo, br.CommitID, "", "", 0)
	stagedPaths := make(map[string]bool)
	for _, s := range staged {
		stagedPaths[s.Path] = true
	}
	for _, obj := range parentObjects {
		if !stagedPaths[obj.Path] {
			if err := c.kv.SetJSON(objectKey(repo, commitID, obj.Path), obj); err != nil {
				return nil, err
			}
		}
	}

	// Save commit
	if err := c.kv.SetJSON(commitKey(repo, commitID), commit); err != nil {
		return nil, err
	}

	// Update branch head
	br.CommitID = commitID
	if err := c.kv.SetJSON(branchKey(repo, branch), br); err != nil {
		return nil, err
	}

	// Clear staging area
	_ = c.kv.DeletePrefix(stagingKey(repo, branch, ""))

	return commit, nil
}

// LogCommits returns commits reachable from a ref, walking parent links
func (c *Catalog) LogCommits(repo, ref string, after string, amount int) ([]*Commit, bool, error) {
	if amount <= 0 {
		amount = 100
	}
	commitID, err := c.resolveRef(repo, ref)
	if err != nil {
		return nil, false, err
	}

	var commits []*Commit
	current := commitID
	pastAfter := after == ""

	for current != "" && len(commits) < amount+1 {
		commit, err := c.GetCommit(repo, current)
		if err != nil {
			break
		}
		if pastAfter {
			commits = append(commits, commit)
		}
		if !pastAfter && current == after {
			pastAfter = true
		}
		if len(commit.Parents) > 0 {
			current = commit.Parents[0]
		} else {
			break
		}
	}

	hasMore := len(commits) > amount
	if hasMore {
		commits = commits[:amount]
	}
	return commits, hasMore, nil
}

// ── Objects ──

func stagingKey(repo, branch, path string) string {
	return fmt.Sprintf("repos/%s/staging/%s/%s", repo, branch, path)
}

func objectKey(repo, commitID, path string) string {
	return fmt.Sprintf("repos/%s/objects/%s/%s", repo, commitID, path)
}

// UploadObject stages an object on a branch
func (c *Catalog) UploadObject(repo, branch, path string, content io.Reader) (*ObjectEntry, error) {
	if _, err := c.GetBranch(repo, branch); err != nil {
		return nil, err
	}

	// Store content via storage backend
	storageKey := branch + "/" + path
	physAddr, size, checksum, err := c.storage.Put(context.Background(), repo, storageKey, content)
	if err != nil {
		return nil, fmt.Errorf("storing object: %w", err)
	}

	entry := &ObjectEntry{
		Path:         path,
		PhysicalAddr: physAddr,
		Checksum:     checksum,
		SizeBytes:    size,
		Mtime:        time.Now(),
	}
	if err := c.kv.SetJSON(stagingKey(repo, branch, path), entry); err != nil {
		return nil, err
	}
	return entry, nil
}

// GetObject retrieves an object entry from a ref (branch or commit)
func (c *Catalog) GetObject(repo, ref, path string) (*ObjectEntry, error) {
	// Try as branch first — check staging then head commit
	branch, err := c.GetBranch(repo, ref)
	if err == nil {
		// Check staging
		var entry ObjectEntry
		if err := c.kv.GetJSON(stagingKey(repo, ref, path), &entry); err == nil {
			return &entry, nil
		}
		// Check head commit
		if err := c.kv.GetJSON(objectKey(repo, branch.CommitID, path), &entry); err == nil {
			return &entry, nil
		}
		return nil, fmt.Errorf("object %q: %w", path, ErrNotFound)
	}
	// Try as commit ID
	var entry ObjectEntry
	if err := c.kv.GetJSON(objectKey(repo, ref, path), &entry); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("object %q: %w", path, ErrNotFound)
		}
		return nil, err
	}
	return &entry, nil
}

// ListObjects lists objects visible from a ref (committed + staged for branches)
func (c *Catalog) ListObjects(repo, ref, prefix, after string, amount int) ([]*ObjectEntry, bool, error) {
	if amount <= 0 {
		amount = 100
	}

	// Merge committed objects + staging for branches
	objectMap := make(map[string]*ObjectEntry)

	// Resolve to get committed objects
	branch, branchErr := c.GetBranch(repo, ref)
	commitID := ref
	if branchErr == nil {
		commitID = branch.CommitID
	}

	// List committed objects
	committed, _ := c.listObjects(repo, commitID, prefix, "", 0)
	for _, obj := range committed {
		objectMap[obj.Path] = obj
	}

	// If this is a branch, overlay staging
	if branchErr == nil {
		staged, _ := c.listStaging(repo, ref)
		for _, obj := range staged {
			if prefix == "" || strings.HasPrefix(obj.Path, prefix) {
				objectMap[obj.Path] = obj
			}
		}
	}

	// Convert to sorted slice
	var objects []*ObjectEntry
	for _, obj := range objectMap {
		if after != "" && obj.Path <= after {
			continue
		}
		objects = append(objects, obj)
	}
	sort.Slice(objects, func(i, j int) bool { return objects[i].Path < objects[j].Path })

	hasMore := len(objects) > amount
	if hasMore {
		objects = objects[:amount]
	}
	return objects, hasMore, nil
}

// DeleteObject removes an object from staging
func (c *Catalog) DeleteObject(repo, branch, path string) error {
	if _, err := c.GetBranch(repo, branch); err != nil {
		return err
	}
	// Also try to delete from storage
	storageKey := branch + "/" + path
	_ = c.storage.Delete(context.Background(), repo, storageKey)
	return c.kv.Delete(stagingKey(repo, branch, path))
}

// GetObjectContent retrieves the actual content of an object as a stream
func (c *Catalog) GetObjectContent(repo, ref, path string) (io.ReadCloser, error) {
	// Resolve the storage key by checking staging then committed
	branch, err := c.GetBranch(repo, ref)
	if err == nil {
		// Check staging first
		var entry ObjectEntry
		if err := c.kv.GetJSON(stagingKey(repo, ref, path), &entry); err == nil {
			storageKey := ref + "/" + path
			return c.storage.Get(context.Background(), repo, storageKey)
		}
		// Check head commit
		if err := c.kv.GetJSON(objectKey(repo, branch.CommitID, path), &entry); err == nil {
			// For committed objects, the storage key was branch/path at upload time
			// but PhysicalAddr has the actual location
			return c.storage.Get(context.Background(), repo, FindStorageKey(entry.PhysicalAddr, repo))
		}
		return nil, fmt.Errorf("object %q: %w", path, ErrNotFound)
	}
	// Try as commit ID
	var entry ObjectEntry
	if err := c.kv.GetJSON(objectKey(repo, ref, path), &entry); err != nil {
		return nil, fmt.Errorf("object %q: %w", path, ErrNotFound)
	}
	return c.storage.Get(context.Background(), repo, FindStorageKey(entry.PhysicalAddr, repo))
}

// FindStorageKey extracts the storage key from a physical address.
// For local: "/home/user/.azlake/objects/repo/branch/path" → "branch/path"
// For azure: "az://container/repo/branch/path" → "branch/path"
func FindStorageKey(physAddr, repo string) string {
	// Try to find repo name in the path and take everything after it
	if idx := strings.Index(physAddr, repo + "/"); idx >= 0 {
		return physAddr[idx+len(repo)+1:]
	}
	return physAddr
}

// ── Tags ──

func tagKey(repo, tag string) string {
	return fmt.Sprintf("repos/%s/tags/%s", repo, tag)
}

func (c *Catalog) CreateTag(repo, name, ref string) (string, error) {
	if _, err := c.GetRepository(repo); err != nil {
		return "", err
	}
	key := tagKey(repo, name)
	if _, err := c.kv.Get(key); err == nil {
		return "", fmt.Errorf("tag %q: %w", name, ErrAlreadyExists)
	}
	commitID, err := c.resolveRef(repo, ref)
	if err != nil {
		return "", err
	}
	if err := c.kv.Set(key, []byte(commitID)); err != nil {
		return "", err
	}
	return commitID, nil
}

func (c *Catalog) GetTag(repo, name string) (string, error) {
	data, err := c.kv.Get(tagKey(repo, name))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", fmt.Errorf("tag %q: %w", name, ErrNotFound)
		}
		return "", err
	}
	return string(data), nil
}

func (c *Catalog) ListTags(repo, after string, amount int) (map[string]string, bool, error) {
	if amount <= 0 {
		amount = 100
	}
	tags := make(map[string]string)
	prefix := fmt.Sprintf("repos/%s/tags/", repo)
	scanAfter := ""
	if after != "" {
		scanAfter = tagKey(repo, after)
	}
	count := 0
	err := c.kv.Scan(prefix, scanAfter, amount+1, func(key string, value []byte) error {
		tagName := strings.TrimPrefix(key, prefix)
		tags[tagName] = string(value)
		count++
		return nil
	})
	if err != nil {
		return nil, false, err
	}
	hasMore := count > amount
	if hasMore {
		// Remove the extra one
		// (tags is a map, so we just check count)
	}
	return tags, hasMore, nil
}

func (c *Catalog) DeleteTag(repo, name string) error {
	if _, err := c.GetTag(repo, name); err != nil {
		return err
	}
	return c.kv.Delete(tagKey(repo, name))
}

// ── Diff ──

// DiffBranch returns staged changes on a branch (uncommitted diff)
func (c *Catalog) DiffBranch(repo, branch string) ([]*DiffEntry, error) {
	staged, err := c.listStaging(repo, branch)
	if err != nil {
		return nil, err
	}
	var diffs []*DiffEntry
	for _, entry := range staged {
		diffs = append(diffs, &DiffEntry{
			Path:      entry.Path,
			Type:      "added",
			SizeBytes: entry.SizeBytes,
		})
	}
	return diffs, nil
}

// ── Helpers ──

// resolveRef resolves a ref (branch name, tag name, or commit ID) to a commit ID
func (c *Catalog) resolveRef(repo, ref string) (string, error) {
	// Try branch
	branch, err := c.GetBranch(repo, ref)
	if err == nil {
		return branch.CommitID, nil
	}
	// Try tag
	commitID, err := c.GetTag(repo, ref)
	if err == nil {
		return commitID, nil
	}
	// Try as commit ID directly
	if _, err := c.GetCommit(repo, ref); err == nil {
		return ref, nil
	}
	return "", fmt.Errorf("ref %q: %w", ref, ErrNotFound)
}

func (c *Catalog) listStaging(repo, branch string) ([]*ObjectEntry, error) {
	prefix := stagingKey(repo, branch, "")
	var entries []*ObjectEntry
	err := c.kv.Scan(prefix, "", 0, func(key string, value []byte) error {
		var entry ObjectEntry
		if err := unmarshalJSON(value, &entry); err != nil {
			return nil
		}
		entries = append(entries, &entry)
		return nil
	})
	return entries, err
}

func (c *Catalog) listObjects(repo, commitID, prefix, after string, limit int) ([]*ObjectEntry, error) {
	scanPrefix := fmt.Sprintf("repos/%s/objects/%s/%s", repo, commitID, prefix)
	var entries []*ObjectEntry
	err := c.kv.Scan(scanPrefix, "", limit, func(key string, value []byte) error {
		var entry ObjectEntry
		if err := unmarshalJSON(value, &entry); err != nil {
			return nil
		}
		if after != "" && entry.Path <= after {
			return nil
		}
		entries = append(entries, &entry)
		return nil
	})
	return entries, err
}

func generateID(prefix string) string {
	return fmt.Sprintf("%s_%x", prefix, time.Now().UnixNano())
}

func unmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
