package db_common

import (
	"github.com/turbot/pipe-fittings/app_specific"
	"strings"
)

// TODO think about app specific stuff here
func IsClientAppName(appName string) bool {
	return strings.HasPrefix(appName, app_specific.ClientConnectionAppNamePrefix) && !strings.HasPrefix(appName, app_specific.ClientSystemConnectionAppNamePrefix)
}

func IsClientSystemAppName(appName string) bool {
	return strings.HasPrefix(appName, app_specific.ClientSystemConnectionAppNamePrefix)
}

func IsServiceAppName(appName string) bool {
	return strings.HasPrefix(appName, app_specific.ServiceConnectionAppNamePrefix)
}
