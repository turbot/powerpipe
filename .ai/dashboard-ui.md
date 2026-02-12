# Dashboard UI (React/TypeScript)

## Tech Stack

- **React 18** with TypeScript 4.5, built via **Craco** (CRA override)
- **State**: React Context + `useReducer` (no Redux)
- **Routing**: `react-router-dom` v6
- **WebSocket**: `react-use-websocket`
- **Charts**: ECharts (`echarts-for-react`)
- **Tables**: `@tanstack/react-table` + `@tanstack/react-virtual`
- **Graphs**: React Flow (`reactflow`) + dagre layout
- **Styling**: Tailwind CSS 3 with custom CSS variable themes
- **Storybook**: Available for component development (`yarn storybook`)
- **Node**: >= 20 required

## Source Tree

```
ui/dashboard/src/
├── index.tsx              Entry point (provider tree)
├── App.tsx                Routes + DashboardProvider
├── hooks/                 State management & side effects
│   ├── useDashboard.tsx          Root provider composition (7 nested providers)
│   ├── useDashboardState.tsx     Main reducer (IDashboardContext)
│   ├── useDashboardExecution.tsx WebSocket lifecycle + snapshot loading
│   ├── useDashboardWebSocket.ts  WebSocket connection (react-use-websocket)
│   ├── useDashboardWebSocketEventHandler.ts  Event buffering (500ms flush)
│   ├── useDashboardInputs.tsx    Input/filter state
│   ├── useDashboardSearchPath.tsx Search path state
│   ├── useDashboardDatetimeRange.tsx DateTime range state
│   ├── useDashboardPanelDetail.tsx Side panel state
│   ├── useDashboardSearch.tsx    Global search state
│   ├── useTheme.tsx              Light/dark theme
│   ├── useBreakpoint.tsx         Responsive breakpoints
│   └── useAnalytics.tsx          Analytics tracking
├── components/
│   ├── dashboards/        Dashboard-specific components
│   │   ├── charts/        AreaChart, BarChart, ColumnChart, DonutChart, etc.
│   │   ├── flows/         Sankey, Flow
│   │   ├── graphs/        ForceDirectedGraph, Graph
│   │   ├── hierarchies/   Tree, Hierarchy
│   │   ├── inputs/        DateInput, SelectInput, ComboInput, etc.
│   │   ├── layout/        Dashboard, Container, Panel, Grid
│   │   ├── grouping/      Benchmark, CheckPanel, DetectionBenchmark
│   │   ├── Card/, Table/, Text/, Image/, Error/
│   │   └── index.ts       Component registry (getComponent/registerComponent)
│   └── DashboardHeader, DashboardList, DashboardSearch, etc.
├── utils/
│   ├── registerComponents.ts  Registers all panel types in registry
│   ├── dashboardEventHandlers.ts  WebSocket event processing + schema migration
│   ├── state.ts               State builders
│   └── data.ts, color.ts, url.ts, snapshot.ts, ...
├── types/                 TypeScript type definitions
├── constants/             Schema versions, icon mappings
└── styles/                Tailwind config, CSS themes
```

## Provider Hierarchy (outermost to innermost)

```
BrowserRouter
  └─ ThemeProvider (light/dark)
      └─ ErrorBoundary
          └─ BreakpointProvider (responsive)
              └─ AnalyticsProvider
                  └─ DashboardProvider (composes 7 inner providers):
                      ├─ DashboardThemeProvider
                      ├─ DashboardSearchProvider
                      ├─ DashboardStateProvider  ← main reducer
                      ├─ DashboardInputsProvider
                      ├─ DashboardSearchPathProvider
                      ├─ DashboardDatetimeRangeProvider
                      ├─ DashboardPanelDetailProvider
                      └─ DashboardExecutionProvider  ← WebSocket
```

## State Management

`useDashboardState.tsx` defines the main reducer with `IDashboardContext` state type.

Key state fields:
- `dataMode`: `"live"` | `"cli_snapshot"` | `"cloud_snapshot"`
- `state`: `"running"` | `"complete"` | `"error"`
- `panelsMap`: All panel data keyed by name
- `panelsLog`: Execution logs per panel
- `dashboard`: Current dashboard definition
- `selectedDashboard`: Currently selected dashboard
- `dashboards` / `dashboardsMap`: Available dashboards
- `snapshot`: Execution snapshot for completed runs
- `execution_id`: Current execution UUID
- `progress`: Execution progress (0-100)
- `metadata`: Server metadata
- `error`: Current error state

Key dispatch actions (in `DashboardActions` enum):
- `SERVER_METADATA` - Server metadata received
- `DASHBOARD_METADATA` - Dashboard-specific metadata
- `AVAILABLE_DASHBOARDS` - List of runnable dashboards
- `EXECUTION_STARTED` - Execution begins, panel structure received
- `EXECUTION_COMPLETE` - All panels done, snapshot ready
- `CONTROLS_UPDATED` - Batch of control updates
- `LEAF_NODES_COMPLETE` - Batch of completed leaf nodes
- `LEAF_NODES_UPDATED` - Batch of updated leaf nodes

## WebSocket Event Handling

### Event Buffering (`useDashboardWebSocketEventHandler.ts`)

Rapid events from server are buffered to prevent UI thrashing:
- `CONTROL_COMPLETE` / `CONTROL_ERROR` events → buffered in array
- `LEAF_NODE_COMPLETE` events → buffered with timestamp
- `LEAF_NODE_UPDATED` events → buffered
- Buffer flushed every **500ms** via `setInterval`
- Other events (`execution_started`, `execution_complete`) → dispatched immediately

### WebSocket Connection (`useDashboardWebSocket.ts`)

Uses `react-use-websocket` library.

URL resolution:
- **Development**: `ws://localhost:9033/ws`
- **Production**: Derives from `window.location` (http→ws, https→wss)

Reconnection: max 10 attempts, 3000ms interval.

### Socket Actions (client → server)

```typescript
SocketActions = {
  CLEAR_DASHBOARD: "clear_dashboard",
  GET_AVAILABLE_DASHBOARDS: "get_available_dashboards",
  GET_SERVER_METADATA: "get_server_metadata",
  SELECT_DASHBOARD: "select_dashboard",
  INPUT_CHANGED: "input_changed",
}
```

### Server Events (server → client)

- `available_dashboards` - List of runnable dashboards
- `execution_started` - Dashboard execution began; includes panel metadata and layout
- `leaf_node_updated` - Panel result ready; includes data/status
- `leaf_node_complete` - Panel execution finished
- `execution_complete` - All panels complete; includes full snapshot
- `execution_error` - Runtime error
- `control_complete` / `control_error` - Individual control results
- `workspace_error` - Mod parse/load error

## Routes

```
/                          Live dashboard view (WebSocket streaming)
/:dashboard_name           Select specific dashboard
/snapshot/:dashboard_name  View saved snapshot (read-only, no re-execution)
```

Query params on snapshot route encode: inputs, datetime range, search_path_prefix.

## Theming

Two CSS variable themes defined in `styles/index.css`:
- `.theme-steampipe-default` (light)
- `.theme-steampipe-dark` (dark)

Custom color variables:
- Control colors: `--color-alert`, `--color-ok`, `--color-info`, `--color-skip`, `--color-severity`
- Layout: `--color-dashboard`, `--color-dashboard-panel`
- Text: `--color-foreground`, `--color-foreground-light`, `-lighter`, `-lightest`
- Tables: `--color-table-border`, `--color-table-divide`, `--color-table-head`
- Scale: `--color-black-scale-1` through `--color-black-scale-8`

Tailwind plugins: `@tailwindcss/forms`, `@tailwindcss/typography`.

## Component Registry

Components registered dynamically in `utils/registerComponents.ts` into a map in `components/dashboards/index.ts`:

```typescript
const componentsMap = {};
const getComponent = (key: string) => componentsMap[key];
const registerComponent = (key: string, component) => { componentsMap[key] = component; };
```

The layout renderer looks up components by panel type string (e.g., `"card"`, `"chart"`, `"table"`). This allows extensibility without hardcoded imports.

Registered types: Panel, Container, Dashboard, all chart types (Area, Bar, Column, Donut, Heatmap, Line, Pie), Flow, Sankey, Graph, ForceDirectedGraph, Tree, Hierarchy, all input types, Benchmark, DetectionBenchmark, Table, Text, Image, Card, Error.

## Build & Dev Workflow

```bash
cd ui/dashboard
yarn install
yarn start           # Dev server on http://localhost:3000 (proxies to Go backend on :9033)
yarn test            # Jest + React Testing Library
yarn storybook       # Component playground on http://localhost:6006
yarn build           # Production build → build/
```

### Craco Configuration (`craco.config.js`)

- **WebAssembly**: `experiments.asyncWebAssembly` enabled
- **Path alias**: `@powerpipe` → `src/`
- **Node polyfills**: buffer, crypto, path, stream, vm (via `ProvidePlugin`)
- **Circular dependency detection**: `CircularDependencyPlugin` fails build on circular imports in `/src`

### Schema Versions

Event schema versions tracked in `utils/dashboardEventHandlers.ts`:
```
EXECUTION_SCHEMA_VERSION_20220614
EXECUTION_SCHEMA_VERSION_20220929
EXECUTION_SCHEMA_VERSION_20221222
EXECUTION_SCHEMA_VERSION_20240130
EXECUTION_SCHEMA_VERSION_20240607
EXECUTION_SCHEMA_VERSION_20241125
```

Schema migration functions handle backwards compatibility between versions.

## Key Type Definitions (`types/index.ts`)

- `DashboardDataMode`: `"live"` | `"cli_snapshot"` | `"cloud_snapshot"`
- `DashboardRunState`: `"running"` | `"complete"` | `"error"`
- `DashboardPanelType`: `"dashboard"` | `"card"` | `"table"` | `"chart"` | `"input"` | `"graph"` | `"hierarchy"` | `"flow"` | `"benchmark"` | `"control"` | `"detection"` | `"image"` | `"text"` | `"error"` | `"with"` | `"edge"`
- `PanelLog`: `{ error, executionTime, isDependency, prefix, status, timestamp, title }`
- `ReceivedSocketMessagePayload`: `{ action: string, [key: string]: any }`

## Data Flow

```
User Action (click dashboard, change input)
    ↓
React Component / Event Handler
    ↓
WebSocket send (react-use-websocket)
    ↓
Go backend receives, executes query/dashboard
    ↓
Go backend sends event via WebSocket
    ↓
useDashboardWebSocket receives message
    ↓
useDashboardWebSocketEventHandler buffers (500ms)
    ↓
Dispatch action to useDashboardState reducer
    ↓
State updated, Context re-renders subscribed components
    ↓
Component renders with new data (ECharts, ReactFlow, etc.)
```
