package constants

import (
	"fmt"

	"github.com/turbot/powerpipe/internal/version"
)

// DashboardListenAddresses is an arrays is listen addresses which Steampipe accepts
var DashboardListenAddresses = []string{"localhost", "127.0.0.1"}

const (
	DashboardServerDefaultPort    = 9194
	DashboardAssetsImageRefFormat = "us-docker.pkg.dev/steampipe/steampipe/assets:%s"
)

// TODO explicitly define dashboard asset version
var (
	DashboardAssetsImageRef = fmt.Sprintf(DashboardAssetsImageRefFormat, version.VersionString)
)
