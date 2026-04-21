package version

var (
	Version   = "dev"
	BuildDate = "unknown"
)

// Audit log version info
type AuditLogVersionInfo struct {
	Version   string
	BuildDate string
}

func GetAuditLogVersionInfo() AuditLogVersionInfo {
	return AuditLogVersionInfo{Version: Version, BuildDate: BuildDate}
}

const DefaultReleasesURL = "https://github.com/dubin555/azlake/releases"

type AuditResponse struct {
	UpgradeURL      *string
	UpgradeRecommended bool
}

type LatestVersionResponse struct {
	CheckPassed     bool
	LatestVersion   string
	CurrentIsLatest bool
	Alerts          []interface{}
}

func CheckLatestVersion(currentVersion string, releasesURL string) (*LatestVersionResponse, error) {
	return &LatestVersionResponse{
		CheckPassed:     true,
		LatestVersion:   currentVersion,
		CurrentIsLatest: true,
	}, nil
}

func AuditCheck(releasesURL string) (*AuditResponse, error) {
	return &AuditResponse{}, nil
}
