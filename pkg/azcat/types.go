package azcat

import "time"

// Repository represents a versioned data repository
type Repository struct {
	Name             string    `json:"name"`
	StorageNamespace string    `json:"storage_namespace"`
	StorageID        string    `json:"storage_id"`
	DefaultBranch    string    `json:"default_branch"`
	CreationDate     time.Time `json:"creation_date"`
	ReadOnly         bool      `json:"read_only"`
}

// Branch represents a mutable reference to a commit
type Branch struct {
	Name     string `json:"name"`
	CommitID string `json:"commit_id"` // head commit
}

// Commit represents an immutable snapshot
type Commit struct {
	ID           string            `json:"id"`
	Message      string            `json:"message"`
	Committer    string            `json:"committer"`
	CreationDate time.Time         `json:"creation_date"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Parents      []string          `json:"parents,omitempty"`
	// MetaRangeID is not used in azlake's simple model
}

// ObjectEntry represents a file/object in the repository
type ObjectEntry struct {
	Path         string            `json:"path"`
	PhysicalAddr string            `json:"physical_address,omitempty"`
	Checksum     string            `json:"checksum"`
	SizeBytes    int64             `json:"size_bytes"`
	ContentType  string            `json:"content_type,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Mtime        time.Time         `json:"mtime"`
}

// DiffEntry represents a change between two refs
type DiffEntry struct {
	Path      string `json:"path"`
	Type      string `json:"type"` // "added", "removed", "changed"
	SizeBytes int64  `json:"size_bytes,omitempty"`
}

// Pagination helper
type Pagination struct {
	HasMore    bool   `json:"has_more"`
	NextOffset string `json:"next_offset"`
	Results    int    `json:"results"`
	MaxPerPage int    `json:"max_per_page"`
}
