package api

import (
	"net/http"

	"github.com/dubin555/azlake/pkg/api/apigen"
)

// Controller implements apigen.ServerInterface with stub methods
type Controller struct{}

// NewController creates a new Controller
func NewController() *Controller {
	return &Controller{}
}

// Verify Controller implements ServerInterface at compile time
var _ apigen.ServerInterface = &Controller{}

func (c *Controller) GetAuthCapabilities(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ExternalPrincipalLogin(w http.ResponseWriter, r *http.Request, body apigen.ExternalPrincipalLoginJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetExternalPrincipal(w http.ResponseWriter, r *http.Request, params apigen.GetExternalPrincipalParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetTokenFromMailbox(w http.ResponseWriter, r *http.Request, mailbox string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ReleaseTokenToMailbox(w http.ResponseWriter, r *http.Request, loginRequestToken string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetTokenRedirect(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListGroups(w http.ResponseWriter, r *http.Request, params apigen.ListGroupsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateGroup(w http.ResponseWriter, r *http.Request, body apigen.CreateGroupJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetGroupACL(w http.ResponseWriter, r *http.Request, groupId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) SetGroupACL(w http.ResponseWriter, r *http.Request, body apigen.SetGroupACLJSONRequestBody, groupId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListGroupMembers(w http.ResponseWriter, r *http.Request, groupId string, params apigen.ListGroupMembersParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteGroupMembership(w http.ResponseWriter, r *http.Request, groupId string, userId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) AddGroupMembership(w http.ResponseWriter, r *http.Request, groupId string, userId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListGroupPolicies(w http.ResponseWriter, r *http.Request, groupId string, params apigen.ListGroupPoliciesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DetachPolicyFromGroup(w http.ResponseWriter, r *http.Request, groupId string, policyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) AttachPolicyToGroup(w http.ResponseWriter, r *http.Request, groupId string, policyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request, body apigen.LoginJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListPolicies(w http.ResponseWriter, r *http.Request, params apigen.ListPoliciesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreatePolicy(w http.ResponseWriter, r *http.Request, body apigen.CreatePolicyJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeletePolicy(w http.ResponseWriter, r *http.Request, policyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetPolicy(w http.ResponseWriter, r *http.Request, policyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) UpdatePolicy(w http.ResponseWriter, r *http.Request, body apigen.UpdatePolicyJSONRequestBody, policyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListUsers(w http.ResponseWriter, r *http.Request, params apigen.ListUsersParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request, body apigen.CreateUserJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request, userId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request, userId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListUserCredentials(w http.ResponseWriter, r *http.Request, userId string, params apigen.ListUserCredentialsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateCredentials(w http.ResponseWriter, r *http.Request, userId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteCredentials(w http.ResponseWriter, r *http.Request, userId string, accessKeyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetCredentials(w http.ResponseWriter, r *http.Request, userId string, accessKeyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteUserExternalPrincipal(w http.ResponseWriter, r *http.Request, userId string, params apigen.DeleteUserExternalPrincipalParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateUserExternalPrincipal(w http.ResponseWriter, r *http.Request, body apigen.CreateUserExternalPrincipalJSONRequestBody, userId string, params apigen.CreateUserExternalPrincipalParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListUserExternalPrincipals(w http.ResponseWriter, r *http.Request, userId string, params apigen.ListUserExternalPrincipalsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListUserGroups(w http.ResponseWriter, r *http.Request, userId string, params apigen.ListUserGroupsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListUserPolicies(w http.ResponseWriter, r *http.Request, userId string, params apigen.ListUserPoliciesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DetachPolicyFromUser(w http.ResponseWriter, r *http.Request, userId string, policyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) AttachPolicyToUser(w http.ResponseWriter, r *http.Request, userId string, policyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetGarbageCollectionConfig(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetStorageConfig(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetLakeFSVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) PullIcebergTable(w http.ResponseWriter, r *http.Request, body apigen.PullIcebergTableJSONRequestBody, catalog string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) PushIcebergTable(w http.ResponseWriter, r *http.Request, body apigen.PushIcebergTableJSONRequestBody, catalog string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetLicense(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) OauthCallback(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListRepositories(w http.ResponseWriter, r *http.Request, params apigen.ListRepositoriesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateRepository(w http.ResponseWriter, r *http.Request, body apigen.CreateRepositoryJSONRequestBody, params apigen.CreateRepositoryParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteRepository(w http.ResponseWriter, r *http.Request, repository string, params apigen.DeleteRepositoryParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetRepository(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListRepositoryRuns(w http.ResponseWriter, r *http.Request, repository string, params apigen.ListRepositoryRunsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetRun(w http.ResponseWriter, r *http.Request, repository string, runId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListRunHooks(w http.ResponseWriter, r *http.Request, repository string, runId string, params apigen.ListRunHooksParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetRunHookOutput(w http.ResponseWriter, r *http.Request, repository string, runId string, hookRunId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) InternalDeleteBranchProtectionRule(w http.ResponseWriter, r *http.Request, body apigen.InternalDeleteBranchProtectionRuleJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) InternalGetBranchProtectionRules(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) InternalCreateBranchProtectionRule(w http.ResponseWriter, r *http.Request, body apigen.InternalCreateBranchProtectionRuleJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateBranchProtectionRulePreflight(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListBranches(w http.ResponseWriter, r *http.Request, repository string, params apigen.ListBranchesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateBranch(w http.ResponseWriter, r *http.Request, body apigen.CreateBranchJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteBranch(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.DeleteBranchParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetBranch(w http.ResponseWriter, r *http.Request, repository string, branch string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ResetBranch(w http.ResponseWriter, r *http.Request, body apigen.ResetBranchJSONRequestBody, repository string, branch string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CherryPick(w http.ResponseWriter, r *http.Request, body apigen.CherryPickJSONRequestBody, repository string, branch string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) Commit(w http.ResponseWriter, r *http.Request, body apigen.CommitJSONRequestBody, repository string, branch string, params apigen.CommitParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CommitAsync(w http.ResponseWriter, r *http.Request, body apigen.CommitAsyncJSONRequestBody, repository string, branch string, params apigen.CommitAsyncParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CommitAsyncStatus(w http.ResponseWriter, r *http.Request, repository string, branch string, id string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DiffBranch(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.DiffBranchParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) HardResetBranch(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.HardResetBranchParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ImportCancel(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.ImportCancelParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ImportStatus(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.ImportStatusParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ImportStart(w http.ResponseWriter, r *http.Request, body apigen.ImportStartJSONRequestBody, repository string, branch string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteObject(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.DeleteObjectParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) UploadObject(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.UploadObjectParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) StageObject(w http.ResponseWriter, r *http.Request, body apigen.StageObjectJSONRequestBody, repository string, branch string, params apigen.StageObjectParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CopyObject(w http.ResponseWriter, r *http.Request, body apigen.CopyObjectJSONRequestBody, repository string, branch string, params apigen.CopyObjectParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteObjects(w http.ResponseWriter, r *http.Request, body apigen.DeleteObjectsJSONRequestBody, repository string, branch string, params apigen.DeleteObjectsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) UploadObjectPreflight(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.UploadObjectPreflightParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) UpdateObjectUserMetadata(w http.ResponseWriter, r *http.Request, body apigen.UpdateObjectUserMetadataJSONRequestBody, repository string, branch string, params apigen.UpdateObjectUserMetadataParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) RevertBranch(w http.ResponseWriter, r *http.Request, body apigen.RevertBranchJSONRequestBody, repository string, branch string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetPhysicalAddress(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.GetPhysicalAddressParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) LinkPhysicalAddress(w http.ResponseWriter, r *http.Request, body apigen.LinkPhysicalAddressJSONRequestBody, repository string, branch string, params apigen.LinkPhysicalAddressParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreatePresignMultipartUpload(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.CreatePresignMultipartUploadParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) AbortPresignMultipartUpload(w http.ResponseWriter, r *http.Request, body apigen.AbortPresignMultipartUploadJSONRequestBody, repository string, branch string, uploadId string, params apigen.AbortPresignMultipartUploadParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CompletePresignMultipartUpload(w http.ResponseWriter, r *http.Request, body apigen.CompletePresignMultipartUploadJSONRequestBody, repository string, branch string, uploadId string, params apigen.CompletePresignMultipartUploadParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) UploadPart(w http.ResponseWriter, r *http.Request, body apigen.UploadPartJSONRequestBody, repository string, branch string, uploadId string, partNumber int, params apigen.UploadPartParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) UploadPartCopy(w http.ResponseWriter, r *http.Request, body apigen.UploadPartCopyJSONRequestBody, repository string, branch string, uploadId string, partNumber int, params apigen.UploadPartCopyParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateCommitRecord(w http.ResponseWriter, r *http.Request, body apigen.CreateCommitRecordJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetCommit(w http.ResponseWriter, r *http.Request, repository string, commitId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DumpStatus(w http.ResponseWriter, r *http.Request, repository string, params apigen.DumpStatusParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DumpSubmit(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) PrepareGarbageCollectionCommits(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) PrepareGarbageCollectionCommitsAsync(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) PrepareGarbageCollectionCommitsStatus(w http.ResponseWriter, r *http.Request, repository string, params apigen.PrepareGarbageCollectionCommitsStatusParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) PrepareGarbageCollectionUncommitted(w http.ResponseWriter, r *http.Request, body apigen.PrepareGarbageCollectionUncommittedJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) InternalDeleteGarbageCollectionRules(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) InternalGetGarbageCollectionRules(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) InternalSetGarbageCollectionRules(w http.ResponseWriter, r *http.Request, body apigen.InternalSetGarbageCollectionRulesJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) SetGarbageCollectionRulesPreflight(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteRepositoryMetadata(w http.ResponseWriter, r *http.Request, body apigen.DeleteRepositoryMetadataJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetRepositoryMetadata(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) SetRepositoryMetadata(w http.ResponseWriter, r *http.Request, body apigen.SetRepositoryMetadataJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetMetaRange(w http.ResponseWriter, r *http.Request, repository string, metaRange string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetMetadataObject(w http.ResponseWriter, r *http.Request, repository string, pType string, objectId string, params apigen.GetMetadataObjectParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetRange(w http.ResponseWriter, r *http.Request, repository string, pRange string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListPullRequests(w http.ResponseWriter, r *http.Request, repository string, params apigen.ListPullRequestsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreatePullRequest(w http.ResponseWriter, r *http.Request, body apigen.CreatePullRequestJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetPullRequest(w http.ResponseWriter, r *http.Request, repository string, pullRequest string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) UpdatePullRequest(w http.ResponseWriter, r *http.Request, body apigen.UpdatePullRequestJSONRequestBody, repository string, pullRequest string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) MergePullRequest(w http.ResponseWriter, r *http.Request, repository string, pullRequest string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DumpRefs(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) RestoreRefs(w http.ResponseWriter, r *http.Request, body apigen.RestoreRefsJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateSymlinkFile(w http.ResponseWriter, r *http.Request, repository string, branch string, params apigen.CreateSymlinkFileParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DiffRefs(w http.ResponseWriter, r *http.Request, repository string, leftRef string, rightRef string, params apigen.DiffRefsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) LogCommits(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.LogCommitsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetObject(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.GetObjectParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) HeadObject(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.HeadObjectParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListObjects(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.ListObjectsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) StatObject(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.StatObjectParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetUnderlyingProperties(w http.ResponseWriter, r *http.Request, repository string, ref string, params apigen.GetUnderlyingPropertiesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) FindMergeBase(w http.ResponseWriter, r *http.Request, repository string, sourceRef string, destinationBranch string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) MergeIntoBranch(w http.ResponseWriter, r *http.Request, body apigen.MergeIntoBranchJSONRequestBody, repository string, sourceRef string, destinationBranch string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) MergeIntoBranchAsync(w http.ResponseWriter, r *http.Request, body apigen.MergeIntoBranchAsyncJSONRequestBody, repository string, sourceRef string, destinationBranch string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) MergeIntoBranchAsyncStatus(w http.ResponseWriter, r *http.Request, repository string, sourceRef string, destinationBranch string, id string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) RestoreStatus(w http.ResponseWriter, r *http.Request, repository string, params apigen.RestoreStatusParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) RestoreSubmit(w http.ResponseWriter, r *http.Request, body apigen.RestoreSubmitJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetBranchProtectionRules(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) SetBranchProtectionRules(w http.ResponseWriter, r *http.Request, body apigen.SetBranchProtectionRulesJSONRequestBody, repository string, params apigen.SetBranchProtectionRulesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteGCRules(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetGCRules(w http.ResponseWriter, r *http.Request, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) SetGCRules(w http.ResponseWriter, r *http.Request, body apigen.SetGCRulesJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) ListTags(w http.ResponseWriter, r *http.Request, repository string, params apigen.ListTagsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) CreateTag(w http.ResponseWriter, r *http.Request, body apigen.CreateTagJSONRequestBody, repository string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) DeleteTag(w http.ResponseWriter, r *http.Request, repository string, tag string, params apigen.DeleteTagParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetTag(w http.ResponseWriter, r *http.Request, repository string, tag string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) SetupCommPrefs(w http.ResponseWriter, r *http.Request, body apigen.SetupCommPrefsJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetSetupState(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) Setup(w http.ResponseWriter, r *http.Request, body apigen.SetupJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) PostStatsEvents(w http.ResponseWriter, r *http.Request, body apigen.PostStatsEventsJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) StsLogin(w http.ResponseWriter, r *http.Request, body apigen.StsLoginJSONRequestBody) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetUsageReportSummary(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *Controller) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

