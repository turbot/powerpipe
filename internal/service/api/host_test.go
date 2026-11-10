package api

import (
	"slices"
	"strings"
	"testing"
)

func TestParseAllowedHosts(t *testing.T) {
	tests := map[string]struct {
		value   string
		want    []string
		wantErr bool
	}{
		"empty":                  {value: ""},
		"names and domains":      {value: " .Pipes.Turbot.com , example.com. ,,", want: []string{".pipes.turbot.com", "example.com"}},
		"single label":           {value: "my-proxy", want: []string{"my-proxy"}},
		"URL":                    {value: "https://example.com", wantErr: true},
		"port":                   {value: "example.com:8080", wantErr: true},
		"wildcard":               {value: "*.example.com", wantErr: true},
		"space in a name":        {value: "exa mple.com", wantErr: true},
		"empty label":            {value: "example..com", wantErr: true},
		"one bad name in a list": {value: "example.com,example.com/path", wantErr: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseAllowedHosts(test.value)
			if test.wantErr {
				if err == nil || !strings.Contains(err.Error(), "--allowed-hosts") {
					t.Fatalf("expected an error that names --allowed-hosts, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, test.want) {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestIsAllowedHost(t *testing.T) {
	allowed := []string{".pipes.turbot.com", "proxy.example.com"}
	tests := map[string]bool{
		"pipes.turbot.com":                           true,
		"acme-aws.usea1.dashboard.pipes.turbot.com":  true,
		"proxy.example.com":                          true,
		"evilpipes.turbot.com":                       false,
		"pipes.turbot.com.evil.example":              false,
		"a.proxy.example.com":                        false,
		"example.com":                                false,
		"turbot.com":                                 false,
		"acme-aws.usea1.dashboard.pipes.turbot.co":   false,
		"acme-aws.usea1.dashboard.pipes.turbot.com2": false,
	}
	for host, want := range tests {
		if got := isAllowedHost(host, allowed); got != want {
			t.Errorf("isAllowedHost(%q) = %v, want %v", host, got, want)
		}
	}
}

func TestIsLoopbackHost(t *testing.T) {
	tests := map[string]bool{
		"localhost":         true,
		"127.0.0.1":         true,
		"127.1.2.3":         true,
		"::1":               true,
		"0.0.0.0":           false,
		"10.0.0.1":          false,
		"localhost.example": false,
		"foo.localhost":     false,
		"evil.example":      false,
	}
	for host, want := range tests {
		if got := isLoopbackHost(host); got != want {
			t.Errorf("isLoopbackHost(%q) = %v, want %v", host, got, want)
		}
	}
}
