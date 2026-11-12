package api

import (
	"fmt"
	"log/slog"
	"net"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/turbot/pipe-fittings/v2/perr"
	localconstants "github.com/turbot/powerpipe/internal/constants"
	"github.com/turbot/powerpipe/internal/service/api/common"
)

var allowedHostPattern = regexp.MustCompile(`^\.?[a-z0-9_-]+(\.[a-z0-9_-]+)*$`)

// ParseAllowedHosts parses the --allowed-hosts value: host names separated by commas. A name that starts with a
// dot, such as .example.com, also matches its subdomains.
func ParseAllowedHosts(value string) ([]string, error) {
	var hosts []string
	for _, host := range strings.Split(value, ",") {
		host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
		if host == "" {
			continue
		}
		if !allowedHostPattern.MatchString(host) {
			return nil, fmt.Errorf("invalid --%s value '%s': expected a host name such as example.com, or .example.com to also match its subdomains", localconstants.ArgAllowedHosts, host)
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

// allowedHostMiddleware rejects a request for a host other than this machine or the allowed hosts. It is for a
// server that listens on localhost only. A web page can point its own domain at 127.0.0.1 (DNS rebinding) and then
// call a local server as a page on the same origin, which the browser permits. Its requests still carry its own
// domain in the Host header.
func allowedHostMiddleware(allowedHosts []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := hostname(c.Request.Host)
		// a request with no Host is not from a browser
		if host == "" || isLoopbackHost(host) || isAllowedHost(host, allowedHosts) {
			c.Next()
			return
		}
		slog.Warn("rejected a request for a host that is not allowed", "host", c.Request.Host)
		common.AbortWithError(c, perr.ForbiddenWithMessage(fmt.Sprintf("host '%s' is not allowed: with --listen local the server accepts only localhost. To accept another host name, such as the domain of a proxy, set --%s.", host, localconstants.ArgAllowedHosts)))
	}
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func isAllowedHost(host string, allowedHosts []string) bool {
	for _, allowed := range allowedHosts {
		if host == strings.TrimPrefix(allowed, ".") || (strings.HasPrefix(allowed, ".") && strings.HasSuffix(host, allowed)) {
			return true
		}
	}
	return false
}
