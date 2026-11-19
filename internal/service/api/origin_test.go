package api

import (
	"net/http/httptest"
	"testing"
)

func TestHostname(t *testing.T) {
	tests := map[string]string{
		"LocalHost:9033":   "localhost",
		"localhost":        "localhost",
		"[::1]:9033":       "::1",
		"127.0.0.1:9033":   "127.0.0.1",
		"example.com.":     "example.com",
		"example.com.:443": "example.com",
		"":                 "",
	}
	for hostPort, want := range tests {
		if got := hostname(hostPort); got != want {
			t.Errorf("hostname(%q) = %q, want %q", hostPort, got, want)
		}
	}
}

func TestWebSocketOriginCheck(t *testing.T) {
	tests := []struct {
		name, host, origin string
		allowedHosts       []string
		want               bool
	}{
		{name: "not a browser", host: "localhost:9033", want: true},
		{name: "same origin", host: "localhost:9033", origin: "http://localhost:9033", want: true},
		{name: "same origin, loopback IP", host: "127.0.0.1:9033", origin: "http://127.0.0.1:9033", want: true},
		{name: "same origin, IPv6 loopback", host: "[::1]:9033", origin: "http://[::1]:9033", want: true},
		{name: "same origin, host in another case", host: "LocalHost:9033", origin: "http://localhost:9033", want: true},
		{name: "https page behind a TLS proxy", host: "powerpipe.example.com", origin: "https://powerpipe.example.com", want: true},
		{name: "page on an allowed host, behind a proxy that sends 127.0.0.1 as the Host", host: "127.0.0.1:9033", origin: "https://proxy.example.com", allowedHosts: []string{"proxy.example.com"}, want: true},
		{name: "dashboard UI dev server, localhost allowed", host: "localhost:9033", origin: "http://localhost:3000", allowedHosts: []string{"localhost"}, want: true},
		{name: "Pipes UI, Pipes domain allowed", host: "acme-aws.usea1.dashboard.pipes.turbot.com", origin: "https://pipes.turbot.com", allowedHosts: []string{".pipes.turbot.com"}, want: true},
		{name: "page on another localhost port", host: "localhost:9033", origin: "http://localhost:8080", want: false},
		{name: "dashboard UI dev server, localhost not allowed", host: "localhost:9033", origin: "http://localhost:3000", want: false},
		{name: "page on the same site", host: "acme-aws.usea1.dashboard.pipes.turbot.com", origin: "https://pipes.turbot.com", want: false},
		{name: "loopback IP page, localhost server", host: "localhost:9033", origin: "http://127.0.0.1:9033", want: false},
		{name: "page on another site", host: "localhost:9033", origin: "https://evil.example", want: false},
		{name: "page on a site with a similar name", host: "localhost:9033", origin: "http://localhost.evil.example:9033", want: false},
		{name: "page on a subdomain of an allowed host", host: "127.0.0.1:9033", origin: "https://a.proxy.example.com", allowedHosts: []string{"proxy.example.com"}, want: false},
		{name: "sandboxed page", host: "localhost:9033", origin: "null", want: false},
		{name: "sandboxed page, no Host", origin: "null", want: false},
		{name: "invalid origin", host: "localhost:9033", origin: "http://[::1", want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/ws", nil)
			r.Host = test.host
			if test.origin != "" {
				r.Header.Set("Origin", test.origin)
			}
			if got := webSocketOriginCheck(test.allowedHosts)(r); got != test.want {
				t.Errorf("check() = %v, want %v", got, test.want)
			}
		})
	}
}
