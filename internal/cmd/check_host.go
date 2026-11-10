package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/error_helpers"
	"github.com/turbot/pipe-fittings/v2/export"
	"github.com/turbot/pipe-fittings/v2/statushooks"
	"github.com/turbot/powerpipe/internal/benchmarkrun"
	localconstants "github.com/turbot/powerpipe/internal/constants"
)

// hostPollInterval is how often `benchmark run --host` reads the state of its run
const hostPollInterval = time.Second

// hostUnsupportedFlags are the benchmark run flags that act on this machine: with --host the server loads
// the mod, connects to the database and runs the controls, so these flags would have no effect.
var hostUnsupportedFlags = []string{
	constants.ArgModLocation,
	constants.ArgDatabase,
	constants.ArgSearchPath,
	constants.ArgSearchPathPrefix,
	constants.ArgVariable,
	constants.ArgVarFile,
	constants.ArgModInstall,
	constants.ArgPull,
	constants.ArgDatabaseQueryTimeout,
	constants.ArgMaxParallel,
	constants.ArgWhere,
	constants.ArgTag,
	constants.ArgDryRun,
	constants.ArgShare,
	constants.ArgSnapshot,
}

// runBenchmarkOnHost runs the benchmark on the powerpipe server at --host. It writes the snapshot of the run
// as `powerpipe benchmark run` does for --output=snapshot and --export=<file>.pps, and sets the same exit code.
func runBenchmarkOnHost(ctx context.Context, cmd *cobra.Command, args []string) {
	client, exportManager, err := validateHostArgs(cmd, args)
	if err != nil {
		exitCode = constants.ExitCodeInsufficientOrWrongInputs
		error_helpers.ShowError(ctx, err)
		return
	}

	statushooks.Show(ctx)
	defer statushooks.Done(ctx)
	statushooks.SetStatus(ctx, fmt.Sprintf("Running %s on %s…", args[0], viper.GetString(constants.ArgHost)))

	timeout := time.Duration(viper.GetInt(constants.ArgBenchmarkTimeout)) * time.Second
	run, runURL, err := client.Start(ctx, args[0], timeout)
	if err != nil {
		exitCode = constants.ExitCodeInitializationFailed
		error_helpers.ShowError(ctx, err)
		return
	}
	// remove the run from the server when the command ends, also when it is cancelled
	defer func() {
		deleteCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var apiErr *benchmarkrun.APIError
		if err := client.Delete(deleteCtx, runURL); err != nil && !(errors.As(err, &apiErr) && apiErr.StatusCode == 404) {
			error_helpers.ShowWarning(fmt.Sprintf("failed to remove the run from the server: %s", err))
		}
	}()

	run, err = waitForHostRun(ctx, client, runURL)
	if err != nil {
		exitCode = constants.ExitCodeSnapshotCreationFailed
		if error_helpers.IsContextCanceled(ctx) {
			exitCode = constants.ExitCodeOperationCancelled
		}
		error_helpers.ShowError(ctx, err)
		return
	}
	if run.Status == benchmarkrun.StatusCancelled {
		exitCode = constants.ExitCodeOperationCancelled
		error_helpers.ShowError(ctx, fmt.Errorf("the run was cancelled on the server"))
		return
	}
	statushooks.Done(ctx)

	// as a local run counts them
	alarms, errorCount := 0, 0
	if run.Error != "" {
		errorCount++
		error_helpers.ShowError(ctx, errors.New(run.Error))
	}
	if run.Summary != nil {
		alarms, errorCount = run.Summary.Alarm, run.Summary.Error
	}

	if output := viper.GetString(constants.ArgOutput); output != constants.OutputFormatNone {
		if err := writeHostSnapshot(ctx, client, runURL, os.Stdout); err != nil {
			error_helpers.ShowError(ctx, err)
			errorCount++
		}
	}

	exportMsg, err := exportManager.DoExport(ctx, run.Benchmark, hostSnapshot{client: client, runURL: runURL}, viper.GetStringSlice(constants.ArgExport))
	if err != nil {
		error_helpers.ShowError(ctx, err)
		errorCount++
	}
	if len(exportMsg) > 0 && viper.GetBool(constants.ArgProgress) {
		fmt.Printf("\n%s\n", strings.Join(exportMsg, "\n")) //nolint:forbidigo // we want to print
	}

	exitCode = getExitCode(alarms, errorCount)
}

// validateHostArgs checks the args and flags for a run on a server, and returns a client for the server and
// the export manager for the --export flags
func validateHostArgs(cmd *cobra.Command, args []string) (*benchmarkrun.Client, *export.Manager, error) {
	if len(args) != 1 {
		return nil, nil, fmt.Errorf("--%s runs one benchmark at a time", constants.ArgHost)
	}
	var unsupported []string
	for _, name := range hostUnsupportedFlags {
		if cmd.Flags().Changed(name) {
			unsupported = append(unsupported, "--"+name)
		}
	}
	if len(unsupported) > 0 {
		return nil, nil, fmt.Errorf("%s cannot be used with --%s: the server loads the mod and runs the benchmark", strings.Join(unsupported, ", "), constants.ArgHost)
	}
	switch viper.GetString(constants.ArgOutput) {
	case constants.OutputFormatNone, constants.OutputFormatSnapshot, localconstants.OutputFormatPpSnapshotShort:
	default:
		return nil, nil, fmt.Errorf("--%s supports --output=snapshot or --output=none", constants.ArgHost)
	}
	exportManager := export.NewManager()
	if err := exportManager.Register(&hostSnapshotExporter{}); err != nil {
		return nil, nil, err
	}
	if err := exportManager.ValidateExportFormat(viper.GetStringSlice(constants.ArgExport)); err != nil {
		return nil, nil, fmt.Errorf("--%s supports only snapshot (pps) export: %w", constants.ArgHost, err)
	}
	client, err := benchmarkrun.NewClient(viper.GetString(constants.ArgHost))
	if err != nil {
		return nil, nil, err
	}
	return client, exportManager, nil
}

// waitForHostRun reads the run until it is no longer running
func waitForHostRun(ctx context.Context, client *benchmarkrun.Client, runURL *url.URL) (benchmarkrun.Run, error) {
	ticker := time.NewTicker(hostPollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return benchmarkrun.Run{}, ctx.Err()
		case <-ticker.C:
		}
		run, err := client.Get(ctx, runURL)
		if err != nil {
			return benchmarkrun.Run{}, fmt.Errorf("failed to read the run from the server: %w", err)
		}
		if run.Status != benchmarkrun.StatusRunning {
			return run, nil
		}
	}
}

func writeHostSnapshot(ctx context.Context, client *benchmarkrun.Client, runURL *url.URL, w io.Writer) error {
	body, err := client.Snapshot(ctx, runURL)
	if err != nil {
		return err
	}
	defer body.Close()
	_, err = io.Copy(w, body)
	return err
}

// hostSnapshot is the snapshot of a run on a server, as the source of an export
type hostSnapshot struct {
	client *benchmarkrun.Client
	runURL *url.URL
}

func (hostSnapshot) IsExportSourceData() {}

// hostSnapshotExporter exports the snapshot of a run on a server, as the snapshot export of a local run does
type hostSnapshotExporter struct{}

func (*hostSnapshotExporter) Export(ctx context.Context, input export.ExportSourceData, destPath string) error {
	snapshot, ok := input.(hostSnapshot)
	if !ok {
		return fmt.Errorf("the snapshot exporter for --%s cannot export %T", constants.ArgHost, input)
	}
	body, err := snapshot.client.Snapshot(ctx, snapshot.runURL)
	if err != nil {
		return err
	}
	defer body.Close()
	return export.Write(destPath, body)
}

func (*hostSnapshotExporter) FileExtension() string {
	return localconstants.SnapshotExtension
}

func (*hostSnapshotExporter) Name() string {
	return constants.OutputFormatSnapshot
}

func (*hostSnapshotExporter) Alias() string {
	return localconstants.OutputFormatPpSnapshotShort
}
