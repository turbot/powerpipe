package dashboardserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	typeHelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/steampipeconfig"
	localcmdconfig "github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/dashboardassets"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/dashboardworkspace"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

func buildServerMetadataPayload(workspaceResources *modconfig.ResourceMaps, cloudMetadata *steampipeconfig.CloudMetadata) ([]byte, error) {
	installedMods := make(map[string]*ModMetadata)
	for _, mod := range workspaceResources.Mods {
		// Ignore current mod
		if mod.FullName == workspaceResources.Mod.FullName {
			continue
		}
		installedMods[mod.FullName] = &ModMetadata{
			Title:     typeHelpers.SafeString(mod.Title),
			FullName:  mod.FullName,
			ShortName: mod.ShortName,
		}
	}

	cliVersion := app_specific.AppVersion.String()

	// when in local mode, we need to hack the response to include the version of the assets and not the version of the cli
	// this is because the UI depends on the version of the assets to be equal to the version it gets from this response
	// since during development, the cli version is always timestamped, we need to hack the response
	if localcmdconfig.IsLocal() {
		versionFile, err := dashboardassets.LoadDashboardAssetVersion()
		if err != nil {
			return nil, sperr.WrapWithMessage(err, "could not load dashboard assets version file")
		}

		cliVersion = versionFile.Version
	}

	payload := ServerMetadataPayload{
		Action: "server_metadata",
		Metadata: ServerMetadata{
			CLI: DashboardCLIMetadata{
				Version: cliVersion,
			},
			InstalledMods: installedMods,
			Telemetry:     viper.GetString(constants.ArgTelemetry),
		},
	}

	defaultDatabase, defaultSearchPathConfig := db_client.GetDefaultDatabaseConfig()
	searchPath, err := getSearchPathMetadata(context.Background(), defaultDatabase, defaultSearchPathConfig)
	if err != nil {
		return nil, err
	}
	payload.Metadata.SearchPath = searchPath

	if mod := workspaceResources.Mod; mod != nil {
		payload.Metadata.Mod = &ModMetadata{
			Title:     typeHelpers.SafeString(mod.Title),
			FullName:  mod.FullName,
			ShortName: mod.ShortName,
		}
	}
	// if telemetry is enabled, send cloud metadata
	if payload.Metadata.Telemetry != constants.TelemetryNone {
		payload.Metadata.Cloud = cloudMetadata
	}

	return json.Marshal(payload)
}

func buildDashboardMetadataPayload(ctx context.Context, dashboard modconfig.ModTreeItem, w *dashboardworkspace.WorkspaceEvents) ([]byte, error) {
	defaultDatabase, defaultSearchPathConfig := db_client.GetDefaultDatabaseConfig()
	database, searchPathConfig, err := db_client.GetDatabaseConfigForResource(dashboard, w.Mod, defaultDatabase, defaultSearchPathConfig)
	if err != nil {
		return nil, err
	}

	payload := DashboardMetadataPayload{
		Action: "dashboard_metadata",
		Metadata: DashboardMetadata{
			Database: database,
		},
	}

	searchPath, err := getSearchPathMetadata(ctx, database, searchPathConfig)
	if err != nil {
		return nil, err
	}
	payload.Metadata.SearchPath = searchPath

	return json.Marshal(payload)
}

func getSearchPathMetadata(ctx context.Context, database string, searchPathConfig backend.SearchPathConfig) (*SearchPathMetadata, error) {
	// if backend supports search path, get it
	client, err := db_client.NewClientMap().GetOrCreate(ctx, database, searchPathConfig)
	if err != nil {
		return nil, err

	}
	//  close the client after we are done
	defer client.Close(ctx)

	if sp, ok := client.Backend.(backend.SearchPathProvider); ok {
		return &SearchPathMetadata{
			OriginalSearchPath:   sp.OriginalSearchPath(),
			ResolvedSearchPath:   sp.ResolvedSearchPath(),
			ConfiguredSearchPath: searchPathConfig.SearchPath,
			SearchPathPrefix:     searchPathConfig.SearchPathPrefix,
		}, nil
	}

	return nil, nil
}

func addBenchmarkChildren(benchmark *modconfig.Benchmark, recordTrunk bool, trunk []string, trunks map[string][][]string) []ModAvailableBenchmark {
	var children []ModAvailableBenchmark
	for _, child := range benchmark.GetChildren() {
		switch t := child.(type) {
		case *modconfig.Benchmark:
			childTrunk := make([]string, len(trunk)+1)
			copy(childTrunk, trunk)
			childTrunk[len(childTrunk)-1] = t.FullName
			if recordTrunk {
				trunks[t.FullName] = append(trunks[t.FullName], childTrunk)
			}
			availableBenchmark := ModAvailableBenchmark{
				Title:     t.GetTitle(),
				FullName:  t.FullName,
				ShortName: t.ShortName,
				Tags:      t.Tags,
				Children:  addBenchmarkChildren(t, recordTrunk, childTrunk, trunks),
			}
			children = append(children, availableBenchmark)
		}
	}
	return children
}

func buildAvailableDashboardsPayload(workspaceResources *modconfig.ResourceMaps) ([]byte, error) {

	payload := AvailableDashboardsPayload{
		Action:     "available_dashboards",
		Dashboards: make(map[string]ModAvailableDashboard),
		Benchmarks: make(map[string]ModAvailableBenchmark),
		Snapshots:  workspaceResources.Snapshots,
	}

	// if workspace resources has a mod, populate dashboards and benchmarks
	if workspaceResources.Mod != nil {
		// build a map of the dashboards provided by each mod

		// iterate over the dashboards for the top level mod - this will include the dashboards from dependency mods
		for _, dashboard := range workspaceResources.Mod.ResourceMaps.Dashboards {
			mod := dashboard.Mod
			// add this dashboard
			payload.Dashboards[dashboard.FullName] = ModAvailableDashboard{
				Title:       typeHelpers.SafeString(dashboard.Title),
				FullName:    dashboard.FullName,
				ShortName:   dashboard.ShortName,
				Tags:        dashboard.Tags,
				ModFullName: mod.FullName,
			}
		}

		benchmarkTrunks := make(map[string][][]string)
		for _, benchmark := range workspaceResources.Mod.ResourceMaps.Benchmarks {
			if benchmark.IsAnonymous() {
				continue
			}

			// Find any benchmarks who have a parent that is a mod - we consider these top-level
			isTopLevel := false
			for _, parent := range benchmark.GetParents() {
				switch parent.(type) {
				case *modconfig.Mod:
					isTopLevel = true
				}
			}

			mod := benchmark.Mod
			trunk := []string{benchmark.FullName}

			if isTopLevel {
				benchmarkTrunks[benchmark.FullName] = [][]string{trunk}
			}

			availableBenchmark := ModAvailableBenchmark{
				Title:       benchmark.GetTitle(),
				FullName:    benchmark.FullName,
				ShortName:   benchmark.ShortName,
				Tags:        benchmark.Tags,
				IsTopLevel:  isTopLevel,
				Children:    addBenchmarkChildren(benchmark, isTopLevel, trunk, benchmarkTrunks),
				ModFullName: mod.FullName,
			}

			payload.Benchmarks[benchmark.FullName] = availableBenchmark
		}
		for benchmarkName, trunks := range benchmarkTrunks {
			if foundBenchmark, ok := payload.Benchmarks[benchmarkName]; ok {
				foundBenchmark.Trunks = trunks
				payload.Benchmarks[benchmarkName] = foundBenchmark
			}
		}
	}

	return json.Marshal(payload)
}

func buildWorkspaceErrorPayload(e *dashboardevents.WorkspaceError) ([]byte, error) {
	payload := ErrorPayload{
		Action: "workspace_error",
		Error:  e.Error.Error(),
	}
	return json.Marshal(payload)
}

func buildControlCompletePayload(event *dashboardevents.ControlComplete) ([]byte, error) {
	payload := ControlEventPayload{
		Action:      "control_complete",
		Control:     event.Control,
		Name:        event.Name,
		Progress:    event.Progress,
		ExecutionId: event.ExecutionId,
		Timestamp:   event.Timestamp,
	}
	return json.Marshal(payload)
}

func buildControlErrorPayload(event *dashboardevents.ControlError) ([]byte, error) {
	payload := ControlEventPayload{
		Action:      "control_error",
		Control:     event.Control,
		Name:        event.Name,
		Progress:    event.Progress,
		ExecutionId: event.ExecutionId,
		Timestamp:   event.Timestamp,
	}
	return json.Marshal(payload)
}

func buildLeafNodeUpdatedPayload(event *dashboardevents.LeafNodeUpdated) ([]byte, error) {
	payload := LeafNodeUpdatedPayload{
		SchemaVersion: fmt.Sprintf("%d", LeafNodeUpdatedSchemaVersion),
		Action:        "leaf_node_updated",
		DashboardNode: event.LeafNode,
		ExecutionId:   event.ExecutionId,
		Timestamp:     event.Timestamp,
	}
	return json.Marshal(payload)
}

func buildExecutionStartedPayload(event *dashboardevents.ExecutionStarted) ([]byte, error) {
	payload := ExecutionStartedPayload{
		SchemaVersion: fmt.Sprintf("%d", ExecutionStartedSchemaVersion),
		Action:        "execution_started",
		ExecutionId:   event.ExecutionId,
		Panels:        event.Panels,
		Layout:        event.Root.AsTreeNode(),
		Inputs:        event.Inputs,
		Variables:     event.Variables,
		StartTime:     event.StartTime,
	}
	return json.Marshal(payload)
}

func buildExecutionErrorPayload(event *dashboardevents.ExecutionError) ([]byte, error) {
	payload := ExecutionErrorPayload{
		Action:    "execution_error",
		Error:     event.Error.Error(),
		Timestamp: event.Timestamp,
	}
	return json.Marshal(payload)
}

func buildExecutionCompletePayload(event *dashboardevents.ExecutionComplete) ([]byte, error) {
	snap := dashboardexecute.ExecutionCompleteToSnapshot(event)
	payload := &ExecutionCompletePayload{
		Action:        "execution_complete",
		SchemaVersion: fmt.Sprintf("%d", ExecutionCompletePayloadSchemaVersion),
		ExecutionId:   event.ExecutionId,
		Snapshot:      snap,
	}
	return json.Marshal(payload)
}

func buildDisplaySnapshotPayload(snap map[string]any) ([]byte, error) {
	payload := &DisplaySnapshotPayload{
		Action:        "execution_complete",
		SchemaVersion: fmt.Sprintf("%d", ExecutionCompletePayloadSchemaVersion),
		Snapshot:      snap,
	}
	return json.Marshal(payload)
}

func buildInputValuesClearedPayload(event *dashboardevents.InputValuesCleared) ([]byte, error) {
	payload := InputValuesClearedPayload{
		Action:        "input_values_cleared",
		ClearedInputs: event.ClearedInputs,
		ExecutionId:   event.ExecutionId,
	}
	return json.Marshal(payload)
}
