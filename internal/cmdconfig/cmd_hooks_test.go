package cmdconfig

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"
)

// The plugin SDK supplies the diagnostics env var names and validation, so these
// tests pin what a user sees whichever SDK version Powerpipe is built with.
func TestValidateConfigDiagnosticsEnv(t *testing.T) {
	cases := map[string]struct {
		env         map[string]string
		wantWarning bool
		wantError   bool
	}{
		"nothing set":            {env: nil},
		"legacy var warns":       {env: map[string]string{"STEAMPIPE_DIAGNOSTICS_LEVEL": "ALL"}, wantWarning: true},
		"current var is quiet":   {env: map[string]string{"STEAMPIPE_DIAGNOSTIC_LEVEL": "NONE"}},
		"bad value is an error":  {env: map[string]string{"STEAMPIPE_DIAGNOSTIC_LEVEL": "loud"}, wantError: true},
		"bad legacy value warns": {env: map[string]string{"STEAMPIPE_DIAGNOSTICS_LEVEL": "loud"}, wantWarning: true, wantError: true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// t.Setenv restores the original value on cleanup; the unset makes "not set" real.
			for _, k := range []string{"STEAMPIPE_DIAGNOSTICS_LEVEL", "STEAMPIPE_DIAGNOSTIC_LEVEL"} {
				t.Setenv(k, "")
				os.Unsetenv(k)
			}
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			viper.Set(constants.ArgTelemetry, constants.TelemetryNone)
			t.Cleanup(func() { viper.Set(constants.ArgTelemetry, nil) })

			res := validateConfig(nil)
			gotWarning := false
			for _, w := range res.Warnings {
				if strings.Contains(w, "STEAMPIPE_DIAGNOSTICS_LEVEL is deprecated") {
					gotWarning = true
				}
			}
			if gotWarning != tc.wantWarning {
				t.Errorf("deprecation warning: got %v, want %v (warnings %v)", gotWarning, tc.wantWarning, res.Warnings)
			}
			if (res.Error != nil) != tc.wantError {
				t.Errorf("error: got %v, want error %v", res.Error, tc.wantError)
			}
		})
	}
}
