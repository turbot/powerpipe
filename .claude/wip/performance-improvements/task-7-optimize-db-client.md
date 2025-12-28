# Task 7: Optimize Database Client Creation

## Objective

Move database client creation to run concurrently with other initialization tasks, reducing blocking time during server startup.

## Context

- Current `db_client.NewDbClient()` blocks server initialization
- Database connection can take 200-500ms depending on the backend
- This time could overlap with other initialization work
- Change is entirely within Powerpipe codebase

## Dependencies

### Prerequisites
- Task 4 (Baseline Measurement) - Need baseline data for comparison

### Files to Modify
- `internal/initialisation/init_data.go`
- `internal/dashboardexecute/executor.go` (if lazy initialization needed)

## Implementation Details

### 1. Current Implementation

```go
// internal/initialisation/init_data.go

func (i *InitData) Init(ctx context.Context, args ...string) {
    // ... early setup ...

    // BLOCKING: Database client creation
    connectionString, err := csp.GetConnectionString()
    if err != nil {
        i.Result.Error = err
        return
    }
    client, err := db_client.NewDbClient(ctx, connectionString, opts...)
    if err != nil {
        i.Result.Error = err
        return
    }
    i.DefaultClient = client

    // More work that could have happened in parallel...
    validationWarnings := validateModRequirementsRecursively(i.Workspace.Mod, client)
    // ...
}
```

### 2. Optimized Implementation

```go
// internal/initialisation/init_data.go

package initialisation

import (
    "context"
    "sync"
    // ... other imports
)

// clientResult holds the async result of database client creation
type clientResult struct {
    client *db_client.DbClient
    err    error
}

func (i *InitData) Init(ctx context.Context, args ...string) {
    defer func() {
        if r := recover(); r != nil {
            i.Result.Error = helpers.ToError(r)
        }
        if i.Result.Error == nil {
            i.Result.Error = ctx.Err()
        }
    }()

    slog.Info("Initializing...")

    if i.Workspace == nil {
        i.Result.Error = sperr.WrapWithRootMessage(
            error_helpers.InvalidStateError,
            "InitData.Init called before setting up WorkspaceEvents",
        )
        return
    }

    statushooks.SetStatus(ctx, "Initializing")

    // Start database client creation in background
    clientChan := make(chan clientResult, 1)
    var clientWg sync.WaitGroup
    clientWg.Add(1)

    go func() {
        defer clientWg.Done()
        defer close(clientChan)

        csp, searchPathConfig, err := db_client.GetDefaultDatabaseConfig(i.Workspace.Mod)
        if err != nil {
            clientChan <- clientResult{err: err}
            return
        }

        // Store these for later use
        i.DefaultDatabase = csp
        i.DefaultSearchPathConfig = searchPathConfig

        var opts []backend.BackendOption
        if !searchPathConfig.Empty() {
            opts = append(opts, backend.WithSearchPathConfig(searchPathConfig))
        }

        connectionString, err := csp.GetConnectionString()
        if err != nil {
            clientChan <- clientResult{err: err}
            return
        }

        client, err := db_client.NewDbClient(ctx, connectionString, opts...)
        clientChan <- clientResult{client: client, err: err}
    }()

    // Meanwhile, do other initialization work that doesn't need DB

    // Initialize telemetry (doesn't need DB)
    shutdownTelemetry, err := telemetry.Init(app_specific.AppName)
    if err != nil {
        i.Result.AddWarnings(err.Error())
    } else {
        i.ShutdownTelemetry = shutdownTelemetry
    }

    // Install mod dependencies if needed (doesn't need DB)
    if viper.GetBool(constants.ArgModInstall) {
        statushooks.SetStatus(ctx, "Installing workspace dependencies")
        slog.Info("Installing workspace dependencies")
        opts := modinstaller.NewInstallOpts(i.Workspace.Mod)
        opts.UpdateStrategy = viper.GetString(constants.ArgPull)
        opts.Force = true
        _, err := modinstaller.InstallWorkspaceDependencies(ctx, opts)
        if err != nil {
            i.Result.Error = err
            return
        }
    }

    // Now wait for database client
    statushooks.SetStatus(ctx, "Connecting to database")
    result := <-clientChan

    if result.err != nil {
        i.Result.Error = result.err
        return
    }
    i.DefaultClient = result.client

    // Validation needs the client, so must be after
    validationWarnings := validateModRequirementsRecursively(i.Workspace.Mod, i.DefaultClient)
    i.Result.AddWarnings(validationWarnings...)

    // Create dashboard executor
    clientMap := db_client.NewClientMap().Add(i.DefaultClient, i.DefaultSearchPathConfig)
    dashboardexecute.Executor = dashboardexecute.NewDashboardExecutor(
        clientMap,
        i.DefaultDatabase,
        i.DefaultSearchPathConfig,
    )
}
```

### 3. Alternative: Lazy Client Creation

If database is only needed on first dashboard execution:

```go
// internal/initialisation/init_data.go

type InitData struct {
    // ... existing fields ...

    // Lazy database client
    clientOnce sync.Once
    clientErr  error
}

func (i *InitData) GetOrCreateClient(ctx context.Context) (*db_client.DbClient, error) {
    i.clientOnce.Do(func() {
        csp, searchPathConfig, err := db_client.GetDefaultDatabaseConfig(i.Workspace.Mod)
        if err != nil {
            i.clientErr = err
            return
        }

        i.DefaultDatabase = csp
        i.DefaultSearchPathConfig = searchPathConfig

        var opts []backend.BackendOption
        if !searchPathConfig.Empty() {
            opts = append(opts, backend.WithSearchPathConfig(searchPathConfig))
        }

        connectionString, err := csp.GetConnectionString()
        if err != nil {
            i.clientErr = err
            return
        }

        i.DefaultClient, i.clientErr = db_client.NewDbClient(ctx, connectionString, opts...)
    })

    return i.DefaultClient, i.clientErr
}
```

### 4. Update Server Command

```go
// internal/cmd/server.go

func runServerCmd(cmd *cobra.Command, _ []string) {
    // ... validation ...

    // Initialize workspace (includes async DB client)
    modInitData := initialisation.NewInitData[*resources.Dashboard](ctx, cmd)
    error_helpers.FailOnError(modInitData.Result.Error)
    defer modInitData.Cleanup(ctx)

    // These can proceed without waiting for DB client to be ready
    err := dashboardassets.Ensure(ctx)
    error_helpers.FailOnError(err)

    webSocket := melody.New()
    dashboardServer, err := dashboardserver.NewServer(ctx, modInitData, webSocket)
    error_helpers.FailOnError(err)

    // API server can start - DB will be ready by first dashboard request
    powerpipeService, err := api.NewAPIService(ctx,
        api.WithWebSocket(webSocket),
        api.WithWorkspace(modInitData.Workspace),
        api.WithHTTPPortAndListenConfig(serverPort, serverListen),
    )
    // ...
}
```

### 5. Add Tests

```go
// internal/initialisation/init_data_test.go

func TestInitDataAsyncClientCreation(t *testing.T) {
    // Mock workspace setup
    ctx := context.Background()

    // Measure time to "ready" state
    start := time.Now()

    // Create init data with async client
    initData := &InitData{
        Workspace: testWorkspace,
        Result:    &InitResult{},
    }
    initData.Init(ctx)

    elapsed := time.Since(start)

    // Should not be blocked by slow DB connection
    // (In reality, test with mock DB that has artificial delay)
    assert.NoError(t, initData.Result.Error)
    assert.NotNil(t, initData.DefaultClient)

    t.Logf("Init completed in %v", elapsed)
}

func TestInitDataClientAvailableBeforeUse(t *testing.T) {
    ctx := context.Background()

    initData := setupTestInitData(t)
    initData.Init(ctx)

    // Client should be available
    assert.NotNil(t, initData.DefaultClient)

    // Should be able to use client
    _, err := initData.DefaultClient.Backend.Ping(ctx)
    assert.NoError(t, err)
}
```

### 6. Measure Performance Improvement

```bash
# Time server startup
POWERPIPE_TIMING=detailed powerpipe server --port 19033 &
PID=$!

# Wait for ready, measure time
# ...

kill $PID

# Compare with baseline
```

## Acceptance Criteria

- [ ] Database client creation runs concurrently with telemetry init
- [ ] Database client creation runs concurrently with mod installation (if enabled)
- [ ] Server startup doesn't block on database connection
- [ ] Database client is available before first dashboard request
- [ ] Error handling works correctly for async client creation
- [ ] Unit tests verify async behavior
- [ ] No race conditions (verify with `go test -race`)
- [ ] Performance improvement measured and documented

## Expected Performance Improvement

| Scenario | Baseline | After | Improvement |
|----------|----------|-------|-------------|
| Local DB (fast) | 200ms | 50ms | ~75% |
| Remote DB (slow) | 500ms | 100ms | ~80% |
| With mod install | 800ms | 400ms | ~50% |

Note: Improvement depends on how much other work can overlap with DB connection.

## Notes

- Must ensure client is ready before any dashboard execution
- Consider connection pooling for better resource usage
- May need to handle connection errors gracefully at request time
- Test with slow database backends to validate improvement
