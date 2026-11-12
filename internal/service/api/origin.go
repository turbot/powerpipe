package api

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

// webSocketOriginCheck returns the check for the Origin of a WebSocket connection. It accepts a client that is not a
// browser (it sends no Origin), a page on the same origin as the server (the same host and port), and a page on an
// allowed host, such as the domain of a proxy, or localhost for dashboard UI development. Browsers let any page open
// a WebSocket to any server, so without this check a page that a user visits could connect to their server, run
// dashboards and read the results. As in the default check of gorilla/websocket, the scheme is not compared: behind a
// proxy that terminates TLS, the page is https and the server sees http.
func webSocketOriginCheck(allowedHosts []string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		if u, err := url.Parse(origin); err == nil && u.Host != "" && (strings.EqualFold(u.Host, r.Host) || isAllowedHost(hostname(u.Host), allowedHosts)) {
			return true
		}
		slog.Warn("rejected a WebSocket connection from a page on another origin", "origin", origin, "host", r.Host)
		return false
	}
}

// hostname returns the lower case host name of a host with an optional port, such as localhost:9033 or [::1]:9033
func hostname(hostPort string) string {
	return strings.TrimSuffix(strings.ToLower((&url.URL{Host: hostPort}).Hostname()), ".")
}
