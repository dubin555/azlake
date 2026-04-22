package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/swag"
	"github.com/dubin555/azlake/pkg/api/apigen"
	"github.com/dubin555/azlake/pkg/azcat"
	"github.com/dubin555/azlake/pkg/version"
)

// Controller implements apigen.ServerInterface
type Controller struct {
	Catalog *azcat.Catalog
}

// NewController creates a new Controller
func NewController(catalog *azcat.Catalog) *Controller {
	return &Controller{Catalog: catalog}
}

// Verify Controller implements ServerInterface at compile time
var _ apigen.ServerInterface = &Controller{}

func (c *Controller) GetAuthCapabilities(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ExternalPrincipalLogin(w http.ResponseWriter, r *http.Request, body apigen.ExternalPrincipalLoginJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetExternalPrincipal(w http.ResponseWriter, r *http.Request, params apigen.GetExternalPrincipalParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetTokenFromMailbox(w http.ResponseWriter, r *http.Request, mailbox string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ReleaseTokenToMailbox(w http.ResponseWriter, r *http.Request, loginRequestToken string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetTokenRedirect(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListGroups(w http.ResponseWriter, r *http.Request, params apigen.ListGroupsParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreateGroup(w http.ResponseWriter, r *http.Request, body apigen.CreateGroupJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetGroupACL(w http.ResponseWriter, r *http.Request, groupId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) SetGroupACL(w http.ResponseWriter, r *http.Request, body apigen.SetGroupACLJSONRequestBody, groupId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListGroupMembers(w http.ResponseWriter, r *http.Request, groupId string, params apigen.ListGroupMembersParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteGroupMembership(w http.ResponseWriter, r *http.Request, groupId string, userId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) AddGroupMembership(w http.ResponseWriter, r *http.Request, groupId string, userId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListGroupPolicies(w http.ResponseWriter, r *http.Request, groupId string, params apigen.ListGroupPoliciesParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DetachPolicyFromGroup(w http.ResponseWriter, r *http.Request, groupId string, policyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) AttachPolicyToGroup(w http.ResponseWriter, r *http.Request, groupId string, policyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request, body apigen.LoginJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListPolicies(w http.ResponseWriter, r *http.Request, params apigen.ListPoliciesParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreatePolicy(w http.ResponseWriter, r *http.Request, body apigen.CreatePolicyJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeletePolicy(w http.ResponseWriter, r *http.Request, policyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetPolicy(w http.ResponseWriter, r *http.Request, policyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) UpdatePolicy(w http.ResponseWriter, r *http.Request, body apigen.UpdatePolicyJSONRequestBody, policyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListUsers(w http.ResponseWriter, r *http.Request, params apigen.ListUsersParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request, body apigen.CreateUserJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request, userId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request, userId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListUserCredentials(w http.ResponseWriter, r *http.Request, userId string, params apigen.ListUserCredentialsParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreateCredentials(w http.ResponseWriter, r *http.Request, userId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteCredentials(w http.ResponseWriter, r *http.Request, userId string, accessKeyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetCredentials(w http.ResponseWriter, r *http.Request, userId string, accessKeyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteUserExternalPrincipal(w http.ResponseWriter, r *http.Request, userId string, params apigen.DeleteUserExternalPrincipalParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreateUserExternalPrincipal(w http.ResponseWriter, r *http.Request, body apigen.CreateUserExternalPrincipalJSONRequestBody, userId string, params apigen.CreateUserExternalPrincipalParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListUserExternalPrincipals(w http.ResponseWriter, r *http.Request, userId string, params apigen.ListUserExternalPrincipalsParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListUserGroups(w http.ResponseWriter, r *http.Request, userId string, params apigen.ListUserGroupsParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListUserPolicies(w http.ResponseWriter, r *http.Request, userId string, params apigen.ListUserPoliciesParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DetachPolicyFromUser(w http.ResponseWriter, r *http.Request, userId string, policyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) AttachPolicyToUser(w http.ResponseWriter, r *http.Request, userId string, policyId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetConfig(w http.ResponseWriter, r *http.Request) {
	v := version.Version
	writeJSON(w, http.StatusOK, apigen.Config{
		VersionConfig: &apigen.VersionConfig{
			Version: &v,
		},
		StorageConfig: &apigen.StorageConfig{
			BlockstoreType:                   "azure",
			PreSignSupport:                   false,
			PreSignMultipartUpload:           ptrBool(false),
			ImportSupport:                    false,
			ImportValidityRegex:              "az://.*",
			DefaultNamespacePrefix:           ptrString("az://"),
			PreSignSupportUi:                 false,
			BlockstoreNamespaceExample:       "az://my-container",
			BlockstoreNamespaceValidityRegex: "^az://",
		},
	})
}

func (c *Controller) GetGarbageCollectionConfig(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetStorageConfig(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetLakeFSVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, apigen.VersionConfig{
		Version: ptrString(version.Version),
	})
}

func (c *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) PullIcebergTable(w http.ResponseWriter, r *http.Request, body apigen.PullIcebergTableJSONRequestBody, catalog string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) PushIcebergTable(w http.ResponseWriter, r *http.Request, body apigen.PushIcebergTableJSONRequestBody, catalog string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetLicense(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) OauthCallback(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListRepositories(w http.ResponseWriter, r *http.Request, params apigen.ListRepositoriesParams) {
	after := ""
	if params.After != nil {
		after = string(*params.After)
	}
	amount := 100
	if params.Amount != nil {
		amount = int(*params.Amount)
	}
	repos, hasMore, err := c.Catalog.ListRepositories(after, amount)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	results := make([]apigen.Repository, 0, len(repos))
	for _, repo := range repos {
		results = append(results, apigen.Repository{
			CreationDate:     repo.CreationDate.Unix(),
			DefaultBranch:    repo.DefaultBranch,
			Id:               repo.Name,
			StorageNamespace: repo.StorageNamespace,
		})
	}
	nextOffset := ""
	if hasMore && len(results) > 0 {
		nextOffset = results[len(results)-1].Id
	}
	writeJSON(w, http.StatusOK, apigen.RepositoryList{
		Pagination: apigen.Pagination{
			HasMore:    hasMore,
			MaxPerPage: 100,
			Results:    len(results),
			NextOffset: nextOffset,
		},
		Results: results,
	})
}

func (c *Controller) CreateRepository(w http.ResponseWriter, r *http.Request, body apigen.CreateRepositoryJSONRequestBody, params apigen.CreateRepositoryParams) {
	defaultBranch := swag.StringValue(body.DefaultBranch)
	if defaultBranch == "" {
		defaultBranch = "main"
	}
	repo, err := c.Catalog.CreateRepository(body.Name, body.StorageNamespace, "", defaultBranch)
	if err != nil {
		if errors.Is(err, azcat.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, err.Error())
		} else {
			writeError(w, r, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, apigen.Repository{
		CreationDate:     repo.CreationDate.Unix(),
		DefaultBranch:    repo.DefaultBranch,
		Id:               repo.Name,
		StorageNamespace: repo.StorageNamespace,
	})
}

func (c *Controller) DeleteRepository(w http.ResponseWriter, r *http.Request, repository string, params apigen.DeleteRepositoryParams) {
	if err := c.Catalog.DeleteRepository(repository); err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) GetRepository(w http.ResponseWriter, r *http.Request, repository string) {
	repo, err := c.Catalog.GetRepository(repository)
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, apigen.Repository{
		CreationDate:     repo.CreationDate.Unix(),
		DefaultBranch:    repo.DefaultBranch,
		Id:               repo.Name,
		StorageNamespace: repo.StorageNamespace,
	})
}

func (c *Controller) ListRepositoryRuns(w http.ResponseWriter, r *http.Request, repository string, params apigen.ListRepositoryRunsParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetRun(w http.ResponseWriter, r *http.Request, repository string, runId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListRunHooks(w http.ResponseWriter, r *http.Request, repository string, runId string, params apigen.ListRunHooksParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetRunHookOutput(w http.ResponseWriter, r *http.Request, repository string, runId string, hookRunId string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) InternalDeleteBranchProtectionRule(w http.ResponseWriter, r *http.Request, body apigen.InternalDeleteBranchProtectionRuleJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) InternalGetBranchProtectionRules(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) InternalCreateBranchProtectionRule(w http.ResponseWriter, r *http.Request, body apigen.InternalCreateBranchProtectionRuleJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreateBranchProtectionRulePreflight(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListBranches(w http.ResponseWriter, r *http.Request, repository string, params apigen.ListBranchesParams) {
	after := ""
	if params.After != nil {
		after = string(*params.After)
	}
	amount := 100
	if params.Amount != nil {
		amount = int(*params.Amount)
	}
	branches, hasMore, err := c.Catalog.ListBranches(repository, after, amount)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	results := make([]apigen.Ref, 0, len(branches))
	for _, b := range branches {
		results = append(results, apigen.Ref{Id: b.Name, CommitId: b.CommitID})
	}
	nextOffset := ""
	if hasMore && len(results) > 0 {
		nextOffset = results[len(results)-1].Id
	}
	writeJSON(w, http.StatusOK, apigen.RefList{
		Pagination: apigen.Pagination{HasMore: hasMore, MaxPerPage: 100, Results: len(results), NextOffset: nextOffset},
		Results:    results,
	})
}

func (c *Controller) CreateBranch(w http.ResponseWriter, r *http.Request, body apigen.CreateBranchJSONRequestBody, repository string) {
	branch, err := c.Catalog.CreateBranch(repository, body.Name, body.Source)
	if err != nil {
		if errors.Is(err, azcat.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, err.Error())
		} else if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, branch.CommitID)
}

func (c *Controller) DeleteBranch(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.DeleteBranchParams) {
	if err := c.Catalog.DeleteBranch(repository, branch); err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusBadRequest, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) GetBranch(w http.ResponseWriter, r *http.Request, repository string, branch string) {
	b, err := c.Catalog.GetBranch(repository, branch)
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, apigen.Ref{Id: b.Name, CommitId: b.CommitID})
}

func (c *Controller) ResetBranch(w http.ResponseWriter, r *http.Request, body apigen.ResetBranchJSONRequestBody, repository string, branch string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CherryPick(w http.ResponseWriter, r *http.Request, body apigen.CherryPickJSONRequestBody, repository string, branch string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) Commit(w http.ResponseWriter, r *http.Request, body apigen.CommitJSONRequestBody, repository string, branch string, params apigen.CommitParams) {
	metadata := make(map[string]string)
	if body.Metadata != nil {
		metadata = body.Metadata.AdditionalProperties
	}
	commit, err := c.Catalog.Commit(repository, branch, body.Message, "azlake", metadata)
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, commitToAPI(commit))
}

func (c *Controller) CommitAsync(w http.ResponseWriter, r *http.Request, body apigen.CommitAsyncJSONRequestBody, repository string, branch string, params apigen.CommitAsyncParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CommitAsyncStatus(w http.ResponseWriter, r *http.Request, repository string, branch string, id string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DiffBranch(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.DiffBranchParams) {
	diffs, err := c.Catalog.DiffBranch(repository, branch)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	results := make([]apigen.Diff, 0, len(diffs))
	for _, d := range diffs {
		results = append(results, apigen.Diff{Path: d.Path, Type: d.Type, SizeBytes: &d.SizeBytes})
	}
	writeJSON(w, http.StatusOK, apigen.DiffList{
		Pagination: apigen.Pagination{HasMore: false, MaxPerPage: 100, Results: len(results)},
		Results:    results,
	})
}

func (c *Controller) HardResetBranch(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.HardResetBranchParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ImportCancel(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.ImportCancelParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ImportStatus(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.ImportStatusParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ImportStart(w http.ResponseWriter, r *http.Request, body apigen.ImportStartJSONRequestBody, repository string, branch string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteObject(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.DeleteObjectParams) {
	if err := c.Catalog.DeleteObject(repository, branch, string(params.Path)); err != nil {
		writeError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) UploadObject(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.UploadObjectParams) {
	obj, err := c.Catalog.UploadObject(repository, branch, string(params.Path), r.Body)
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, apigen.ObjectStats{
		Path:      obj.Path,
		PathType:  "object",
		SizeBytes: &obj.SizeBytes,
		Checksum:  obj.Checksum,
		Mtime:     obj.Mtime.Unix(),
	})
}

func (c *Controller) StageObject(w http.ResponseWriter, r *http.Request, body apigen.StageObjectJSONRequestBody, repository string, branch string, params apigen.StageObjectParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CopyObject(w http.ResponseWriter, r *http.Request, body apigen.CopyObjectJSONRequestBody, repository string, branch string, params apigen.CopyObjectParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteObjects(w http.ResponseWriter, r *http.Request, body apigen.DeleteObjectsJSONRequestBody, repository string, branch string, params apigen.DeleteObjectsParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) UploadObjectPreflight(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.UploadObjectPreflightParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) UpdateObjectUserMetadata(w http.ResponseWriter, r *http.Request, body apigen.UpdateObjectUserMetadataJSONRequestBody, repository string, branch string, params apigen.UpdateObjectUserMetadataParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) RevertBranch(w http.ResponseWriter, r *http.Request, body apigen.RevertBranchJSONRequestBody, repository string, branch string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetPhysicalAddress(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.GetPhysicalAddressParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) LinkPhysicalAddress(w http.ResponseWriter, r *http.Request, body apigen.LinkPhysicalAddressJSONRequestBody, repository string, branch string, params apigen.LinkPhysicalAddressParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreatePresignMultipartUpload(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.CreatePresignMultipartUploadParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) AbortPresignMultipartUpload(w http.ResponseWriter, r *http.Request, body apigen.AbortPresignMultipartUploadJSONRequestBody, repository string, branch string, uploadId string, params apigen.AbortPresignMultipartUploadParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CompletePresignMultipartUpload(w http.ResponseWriter, r *http.Request, body apigen.CompletePresignMultipartUploadJSONRequestBody, repository string, branch string, uploadId string, params apigen.CompletePresignMultipartUploadParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) UploadPart(w http.ResponseWriter, r *http.Request, body apigen.UploadPartJSONRequestBody, repository string, branch string, uploadId string, partNumber int, params apigen.UploadPartParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) UploadPartCopy(w http.ResponseWriter, r *http.Request, body apigen.UploadPartCopyJSONRequestBody, repository string, branch string, uploadId string, partNumber int, params apigen.UploadPartCopyParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreateCommitRecord(w http.ResponseWriter, r *http.Request, body apigen.CreateCommitRecordJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetCommit(w http.ResponseWriter, r *http.Request, repository string, commitId string) {
	commit, err := c.Catalog.GetCommit(repository, commitId)
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, commitToAPI(commit))
}

func (c *Controller) DumpStatus(w http.ResponseWriter, r *http.Request, repository string, params apigen.DumpStatusParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DumpSubmit(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) PrepareGarbageCollectionCommits(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) PrepareGarbageCollectionCommitsAsync(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) PrepareGarbageCollectionCommitsStatus(w http.ResponseWriter, r *http.Request, repository string, params apigen.PrepareGarbageCollectionCommitsStatusParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) PrepareGarbageCollectionUncommitted(w http.ResponseWriter, r *http.Request, body apigen.PrepareGarbageCollectionUncommittedJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) InternalDeleteGarbageCollectionRules(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) InternalGetGarbageCollectionRules(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) InternalSetGarbageCollectionRules(w http.ResponseWriter, r *http.Request, body apigen.InternalSetGarbageCollectionRulesJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) SetGarbageCollectionRulesPreflight(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteRepositoryMetadata(w http.ResponseWriter, r *http.Request, body apigen.DeleteRepositoryMetadataJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetRepositoryMetadata(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) SetRepositoryMetadata(w http.ResponseWriter, r *http.Request, body apigen.SetRepositoryMetadataJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetMetaRange(w http.ResponseWriter, r *http.Request, repository string, metaRange string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetMetadataObject(w http.ResponseWriter, r *http.Request, repository string, pType string, objectId string, params apigen.GetMetadataObjectParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetRange(w http.ResponseWriter, r *http.Request, repository string, pRange string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListPullRequests(w http.ResponseWriter, r *http.Request, repository string, params apigen.ListPullRequestsParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreatePullRequest(w http.ResponseWriter, r *http.Request, body apigen.CreatePullRequestJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetPullRequest(w http.ResponseWriter, r *http.Request, repository string, pullRequest string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) UpdatePullRequest(w http.ResponseWriter, r *http.Request, body apigen.UpdatePullRequestJSONRequestBody, repository string, pullRequest string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) MergePullRequest(w http.ResponseWriter, r *http.Request, repository string, pullRequest string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DumpRefs(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) RestoreRefs(w http.ResponseWriter, r *http.Request, body apigen.RestoreRefsJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) CreateSymlinkFile(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.CreateSymlinkFileParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DiffRefs(w http.ResponseWriter, r *http.Request, repository string, leftRef string, rightRef string, params apigen.DiffRefsParams) {
	leftObjects, _, _ := c.Catalog.ListObjects(repository, leftRef, "", "", 0)
	rightObjects, _, _ := c.Catalog.ListObjects(repository, rightRef, "", "", 0)

	leftMap := make(map[string]*azcat.ObjectEntry)
	for _, obj := range leftObjects {
		leftMap[obj.Path] = obj
	}
	rightMap := make(map[string]*azcat.ObjectEntry)
	for _, obj := range rightObjects {
		rightMap[obj.Path] = obj
	}

	var results []apigen.Diff
	// Added or changed in right
	for path, rObj := range rightMap {
		lObj, exists := leftMap[path]
		if !exists {
			results = append(results, apigen.Diff{Path: path, Type: "added", SizeBytes: &rObj.SizeBytes})
		} else if lObj.Checksum != rObj.Checksum {
			results = append(results, apigen.Diff{Path: path, Type: "changed", SizeBytes: &rObj.SizeBytes})
		}
	}
	// Removed from right
	for path, lObj := range leftMap {
		if _, exists := rightMap[path]; !exists {
			results = append(results, apigen.Diff{Path: path, Type: "removed", SizeBytes: &lObj.SizeBytes})
		}
	}

	writeJSON(w, http.StatusOK, apigen.DiffList{
		Pagination: apigen.Pagination{HasMore: false, MaxPerPage: 100, Results: len(results)},
		Results:    results,
	})
}

func (c *Controller) LogCommits(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.LogCommitsParams) {
	after := ""
	if params.After != nil {
		after = string(*params.After)
	}
	amount := 100
	if params.Amount != nil {
		amount = int(*params.Amount)
	}
	commits, hasMore, err := c.Catalog.LogCommits(repository, ref, after, amount)
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	results := make([]apigen.Commit, 0, len(commits))
	for _, cm := range commits {
		results = append(results, commitToAPI(cm))
	}
	nextOffset := ""
	if hasMore && len(results) > 0 {
		nextOffset = results[len(results)-1].Id
	}
	writeJSON(w, http.StatusOK, apigen.CommitList{
		Pagination: apigen.Pagination{HasMore: hasMore, MaxPerPage: 100, Results: len(results), NextOffset: nextOffset},
		Results:    results,
	})
}

func (c *Controller) GetObject(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.GetObjectParams) {
	rc, err := c.Catalog.GetObjectContent(repository, ref, string(params.Path))
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	defer rc.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, rc)
}

func (c *Controller) HeadObject(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.HeadObjectParams) {
	obj, err := c.Catalog.GetObject(repository, ref, string(params.Path))
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", obj.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", obj.SizeBytes))
	w.Header().Set("ETag", fmt.Sprintf(`"%s"`, obj.Checksum))
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) ListObjects(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.ListObjectsParams) {
	prefix := ""
	if params.Prefix != nil {
		prefix = string(*params.Prefix)
	}
	after := ""
	if params.After != nil {
		after = string(*params.After)
	}
	amount := 100
	if params.Amount != nil {
		amount = int(*params.Amount)
	}
	// Handle "directory" listing — collect unique prefixes
	delimiter := ""
	if params.Delimiter != nil {
		delimiter = string(*params.Delimiter)
	}
	objects, hasMore, err := c.Catalog.ListObjects(repository, ref, prefix, after, amount)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	results := make([]apigen.ObjectStats, 0)
	seenPrefixes := make(map[string]bool)
	for _, obj := range objects {
		if delimiter != "" && delimiter == "/" {
			// Check if object is in a subdirectory relative to prefix
			rest := strings.TrimPrefix(obj.Path, prefix)
			if idx := strings.Index(rest, "/"); idx >= 0 {
				dirPrefix := prefix + rest[:idx+1]
				if !seenPrefixes[dirPrefix] {
					seenPrefixes[dirPrefix] = true
					results = append(results, apigen.ObjectStats{
						Path:     dirPrefix,
						PathType: "common_prefix",
					})
				}
				continue
			}
		}
		results = append(results, apigen.ObjectStats{
			Path:            obj.Path,
			PathType:        "object",
			SizeBytes:       &obj.SizeBytes,
			Checksum:        obj.Checksum,
			ContentType:     &obj.ContentType,
			Mtime:           obj.Mtime.Unix(),
			PhysicalAddress: obj.PhysicalAddr,
		})
	}
	nextOffset := ""
	if hasMore && len(results) > 0 {
		nextOffset = results[len(results)-1].Path
	}
	writeJSON(w, http.StatusOK, apigen.ObjectStatsList{
		Pagination: apigen.Pagination{HasMore: hasMore, MaxPerPage: 100, Results: len(results), NextOffset: nextOffset},
		Results:    results,
	})
}

func (c *Controller) StatObject(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.StatObjectParams) {
	obj, err := c.Catalog.GetObject(repository, ref, string(params.Path))
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, apigen.ObjectStats{
		Path:            obj.Path,
		PathType:        "object",
		SizeBytes:       &obj.SizeBytes,
		Checksum:        obj.Checksum,
		ContentType:     &obj.ContentType,
		Mtime:           obj.Mtime.Unix(),
		PhysicalAddress: obj.PhysicalAddr,
	})
}

func (c *Controller) GetUnderlyingProperties(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.GetUnderlyingPropertiesParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) FindMergeBase(w http.ResponseWriter, r *http.Request, repository string, sourceRef string, destinationBranch string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) MergeIntoBranch(w http.ResponseWriter, r *http.Request, body apigen.MergeIntoBranchJSONRequestBody, repository string, sourceRef string, destinationBranch string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) MergeIntoBranchAsync(w http.ResponseWriter, r *http.Request, body apigen.MergeIntoBranchAsyncJSONRequestBody, repository string, sourceRef string, destinationBranch string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) MergeIntoBranchAsyncStatus(w http.ResponseWriter, r *http.Request, repository string, sourceRef string, destinationBranch string, id string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) RestoreStatus(w http.ResponseWriter, r *http.Request, repository string, params apigen.RestoreStatusParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) RestoreSubmit(w http.ResponseWriter, r *http.Request, body apigen.RestoreSubmitJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetBranchProtectionRules(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) SetBranchProtectionRules(w http.ResponseWriter, r *http.Request, body apigen.SetBranchProtectionRulesJSONRequestBody, repository string, params apigen.SetBranchProtectionRulesParams) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) DeleteGCRules(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetGCRules(w http.ResponseWriter, r *http.Request, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) SetGCRules(w http.ResponseWriter, r *http.Request, body apigen.SetGCRulesJSONRequestBody, repository string) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) ListTags(w http.ResponseWriter, r *http.Request, repository string, params apigen.ListTagsParams) {
	after := ""
	if params.After != nil {
		after = string(*params.After)
	}
	amount := 100
	if params.Amount != nil {
		amount = int(*params.Amount)
	}
	tags, hasMore, err := c.Catalog.ListTags(repository, after, amount)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	results := make([]apigen.Ref, 0, len(tags))
	for name, commitID := range tags {
		results = append(results, apigen.Ref{Id: name, CommitId: commitID})
	}
	writeJSON(w, http.StatusOK, apigen.RefList{
		Pagination: apigen.Pagination{HasMore: hasMore, MaxPerPage: 100, Results: len(results)},
		Results:    results,
	})
}

func (c *Controller) CreateTag(w http.ResponseWriter, r *http.Request, body apigen.CreateTagJSONRequestBody, repository string) {
	commitID, err := c.Catalog.CreateTag(repository, body.Id, body.Ref)
	if err != nil {
		if errors.Is(err, azcat.ErrAlreadyExists) {
			writeError(w, r, http.StatusConflict, err.Error())
		} else {
			writeError(w, r, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, apigen.Ref{Id: body.Id, CommitId: commitID})
}

func (c *Controller) DeleteTag(w http.ResponseWriter, r *http.Request, repository string, tag string, params apigen.DeleteTagParams) {
	if err := c.Catalog.DeleteTag(repository, tag); err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) GetTag(w http.ResponseWriter, r *http.Request, repository string, tag string) {
	commitID, err := c.Catalog.GetTag(repository, tag)
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, apigen.Ref{Id: tag, CommitId: commitID})
}

func (c *Controller) SetupCommPrefs(w http.ResponseWriter, r *http.Request, body apigen.SetupCommPrefsJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetSetupState(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, apigen.SetupState{
		State:            ptrString("initialized"),
		CommPrefsMissing: ptrBool(false),
	})
}

func (c *Controller) Setup(w http.ResponseWriter, r *http.Request, body apigen.SetupJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) PostStatsEvents(w http.ResponseWriter, r *http.Request, body apigen.PostStatsEventsJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) StsLogin(w http.ResponseWriter, r *http.Request, body apigen.StsLoginJSONRequestBody) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetUsageReportSummary(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusNotImplemented, "not implemented in azlake")
}

func (c *Controller) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, apigen.CurrentUser{
		User: apigen.User{
			Id: "admin",
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func ptrString(s string) *string { return &s }
func ptrBool(b bool) *bool       { return &b }

// GetObjectSASURL generates a temporary SAS URL for direct browser/DuckDB access to Azure Blob.
// This is a custom endpoint not in the OpenAPI spec.
func (c *Controller) GetObjectSASURL(w http.ResponseWriter, r *http.Request) {
	repository := chi.URLParam(r, "repository")
	ref := chi.URLParam(r, "ref")
	path := r.URL.Query().Get("path")

	if path == "" {
		writeError(w, r, http.StatusBadRequest, "missing path parameter")
		return
	}

	// Resolve the object to find its storage key
	obj, err := c.Catalog.GetObject(repository, ref, path)
	if err != nil {
		if errors.Is(err, azcat.ErrNotFound) {
			writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Check if the storage backend supports SAS URLs
	storage := c.Catalog.Storage()
	sasBackend, ok := storage.(azcat.SASCapable)
	if !ok {
		// Fallback: return the REST API URL for httpfs mode
		writeJSON(w, http.StatusOK, map[string]string{
			"sas_url": "",
			"mode":    "httpfs",
		})
		return
	}

	// Find the storage key from physical address
	storageKey := azcat.FindStorageKey(obj.PhysicalAddr, repository)
	sasURL, err := sasBackend.GetSASURL(repository, storageKey, 1*time.Hour)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{
			"sas_url": "",
			"mode":    "httpfs",
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"sas_url": sasURL,
		"mode":    "sas",
	})
}

func commitToAPI(cm *azcat.Commit) apigen.Commit {
	meta := apigen.Commit_Metadata{}
	if cm.Metadata != nil {
		meta.AdditionalProperties = cm.Metadata
	}
	return apigen.Commit{
		Id:           cm.ID,
		Message:      cm.Message,
		Committer:    cm.Committer,
		CreationDate: cm.CreationDate.Unix(),
		Metadata:     &meta,
		Parents:      cm.Parents,
	}
}

