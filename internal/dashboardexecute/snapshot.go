package dashboardexecute

import (
	"context"
	"fmt"
	"log"

	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/db_client"
	"github.com/turbot/powerpipe/internal/initialisation"
)

func GenerateSnapshot(ctx context.Context, target string, initData *initialisation.InitData, inputs map[string]any) (snapshot *steampipeconfig.SteampipeSnapshot, err error) {
	w := initData.WorkspaceEvents

	parsedName, err := modconfig.ParseResourceName(target)
	if err != nil {
		return nil, err
	}
	// no session for manual execution
	sessionId := ""
	errorChannel := make(chan error)
	resultChannel := make(chan *steampipeconfig.SteampipeSnapshot)
	dashboardEventHandler := func(ctx context.Context, event dashboardevents.DashboardEvent) {
		handleDashboardEvent(ctx, event, resultChannel, errorChannel)
	}
	w.RegisterDashboardEventHandler(ctx, dashboardEventHandler)
	// clear event handlers again in case another snapshot will be generated in this run
	defer w.UnregisterDashboardEventHandlers()

	// all runtime dependencies must be resolved before execution (i.e. inputs must be passed in)
	Executor.interactive = false
	clientMap := map[string]*db_client.DbClient{initData.Client.GetConnectionString(): initData.Client}
	if err := Executor.ExecuteDashboard(ctx, sessionId, target, inputs, w, clientMap); err != nil {
		return nil, err
	}

	select {
	case err = <-errorChannel:
		return nil, err
	case snapshot = <-resultChannel:
		// set the filename root of the snapshot
		fileRootName := parsedName.ToFullNameWithMod(w.Mod.ShortName)

		snapshot.FileNameRoot = fileRootName
		//  return the context error (if any) to ensure we respect cancellation
		return snapshot, ctx.Err()
	}
}

func handleDashboardEvent(_ context.Context, event dashboardevents.DashboardEvent, resultChannel chan *steampipeconfig.SteampipeSnapshot, errorChannel chan error) {
	switch e := event.(type) {
	case *dashboardevents.ExecutionError:
		errorChannel <- e.Error
	case *dashboardevents.ExecutionComplete:
		log.Println("[TRACE] execution complete event", *e)
		snap := ExecutionCompleteToSnapshot(e)
		resultChannel <- snap
	}
}

// ExecutionCompleteToSnapshot transforms the ExecutionComplete event into a SteampipeSnapshot
func ExecutionCompleteToSnapshot(event *dashboardevents.ExecutionComplete) *steampipeconfig.SteampipeSnapshot {
	return &steampipeconfig.SteampipeSnapshot{
		SchemaVersion: fmt.Sprintf("%d", steampipeconfig.SteampipeSnapshotSchemaVersion),
		Panels:        event.Panels,
		Layout:        event.Root.AsTreeNode(),
		Inputs:        event.Inputs,
		Variables:     event.Variables,
		SearchPath:    event.SearchPath,
		StartTime:     event.StartTime,
		EndTime:       event.EndTime,
		Title:         event.Root.GetTitle(),
	}
}
