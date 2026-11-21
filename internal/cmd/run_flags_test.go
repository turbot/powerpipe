package cmd

import (
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/powerpipe/internal/resources"
)

// Runs never install or update mods (#1211). The run commands still accept --mod-install and --pull, so a
// script that passes them keeps working, but both are deprecated: hidden from help, with a warning when used.
func TestRunCommandsDeprecateModInstallAndPull(t *testing.T) {
	runCmds := map[string]*cobra.Command{
		"benchmark run":           checkCmd[*resources.Benchmark](),
		"control run":             checkCmd[*resources.Control](),
		"dashboard run":           dashboardRunCmd(),
		"detection run":           detectionRunCmd[*resources.Detection](),
		"detection benchmark run": detectionRunCmd[*resources.DetectionBenchmark](),
	}
	for name, cmd := range runCmds {
		for _, flagName := range []string{constants.ArgModInstall, constants.ArgPull} {
			flag := cmd.Flags().Lookup(flagName)
			if flag == nil {
				t.Errorf("%s: --%s is no longer accepted", name, flagName)
				continue
			}
			if flag.Deprecated == "" {
				t.Errorf("%s: --%s is not deprecated", name, flagName)
			}
		}
		if usage := cmd.Flags().FlagUsages(); strings.Contains(usage, "--mod-install") || strings.Contains(usage, "--pull") {
			t.Errorf("%s: help still lists --mod-install or --pull:\n%s", name, usage)
		}
		cmd.Flags().SetOutput(io.Discard) // the deprecation warnings
		if err := cmd.ParseFlags([]string{"--mod-install=false", "--pull", "latest"}); err != nil {
			t.Errorf("%s: --mod-install=false --pull latest no longer parse: %v", name, err)
		}
	}
}

// mod install and mod update still install and update mods, so their --pull stays a working flag.
func TestModCommandsKeepPull(t *testing.T) {
	for name, cmd := range map[string]*cobra.Command{"mod install": modInstallCmd(), "mod update": modUpdateCmd()} {
		flag := cmd.Flags().Lookup(constants.ArgPull)
		if flag == nil {
			t.Errorf("%s: --pull is missing", name)
			continue
		}
		if flag.Deprecated != "" {
			t.Errorf("%s: --pull is deprecated", name)
		}
	}
}
