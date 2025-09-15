v1.4.0
_Whats new_
- Add support for ducklake backend. Update pipe-fittings to v2.7.0. ([#760](https://github.com/turbot/pipe-fittings/issues/760))
- DbClient.Close calls Backend.Close if it supports it
- Update server command to call initData.Cleanup
- Optimise getSearchPathMetadata to not connect to backend unless needed


## v1.3.1 [2025-08-25]
_Bug Fixes_
- Fix issue where the `database` argument from a query resource was not respected. ([#829](https://github.com/turbot/powerpipe/issues/829))
- Fix issue where the default config path was not resolved correctly. The default is the `mod-location`, followed by the `$POWERPIPE_INSTALL_DIR/config`. Also the `POWERPIPE_CONFIG_PATH` environment variable was not respected. ([#898](https://github.com/turbot/powerpipe/issues/898))
- Fix issue where `pie/donut` charts were not rendering correctly on boolean values. ([#433](https://github.com/turbot/powerpipe/issues/433))

_Dependencies_
- Upgrade `hashicorp/go-getter`, `sha.js` and `cipher-base` to remediate critical and high vulnerabilities.

## v1.3.0 [2025-08-05]
_Whats new_
- Compiled with Go 1.24.

_Bug Fixes_
- Fix issue where the `--where` arg was not correctly filtering the benchmarks/controls when JSON path expressions were passed. ([#740](https://github.com/turbot/powerpipe/issues/740))
- Fix error handling for JSON output/export incase of database failures. ([#665](https://github.com/turbot/powerpipe/issues/665))

_Dependencies_
- Upgrade `form-data` and `go-viper/mapstructure/v2` packages to remediate critical and high vulnerabilities.

## v1.2.8 [2025-07-10]
_Bug Fixes_
- Fix regression where powerpipe was failing to run detections from dependant mods if steampipe service was not running. ([#824](https://github.com/turbot/powerpipe/issues/824))
- Fix issue where powerpipe server was failing to resolve the default pipes token. ([#818](https://github.com/turbot/powerpipe/issues/818))
- Fix issue where query run snapshot output was returning a 0 exit code even incase of query failures. ([#816](https://github.com/turbot/powerpipe/issues/816))

_Dependencies_
- Upgrade `pbkdf2` package to remediate critical vulnerabilities.

## v1.2.7 [2025-05-21]
_Bug Fixes_
* Fix issue where powerpipe was exposing the server port to the internet even when listen was local. ([#761](https://github.com/turbot/powerpipe/issues/761))
* Fix issue where benchmark/control json output was not printing the complete JSON to stdout. ([#791](https://github.com/turbot/powerpipe/issues/791))
* Fix issue where the `from` property in Tailpipe connections was getting parsed incorrectly. ([#790](https://github.com/turbot/powerpipe/issues/790))

## v1.2.6 [2025-05-08]
_Bug Fixes_
* Fix issue where powerpipe was failing to run detections from dependant mods if steampipe service was not running. ([#788](https://github.com/turbot/powerpipe/issues/788))
* Fix issue where detection benchmarks were not showing date range selector in powerpipe server. ([#789](https://github.com/turbot/powerpipe/issues/789))
* Fix issue where powerpipe was failing to export detection benchmarks. ([#796](https://github.com/turbot/powerpipe/issues/796))

## v1.2.5 [2025-04-14]
_Bug Fixes_
* Fix the table data cells displaying a vertical scroll in some browsers. ([#765](https://github.com/turbot/powerpipe/issues/765))
* Fix issue where setting custom relative duration not shown in presets does not update the relative timestamp. ([#769](https://github.com/turbot/powerpipe/issues/769))

## v1.2.4 [2025-04-04]
_Bug Fixes_
* Table headers for numeric columns are now right-aligned like the row-level cell content. ([#755](https://github.com/turbot/powerpipe/issues/755)) 

## v1.2.3 [2025-04-02]
_Bug Fixes_
- Fix the supported export formats in `powerpipe query run` command. ([#539](https://github.com/turbot/powerpipe/issues/539))
- Ensure dashboard UI Table component fetches external link from registry when used rather than during module import. ([#720](https://github.com/turbot/powerpipe/issues/720))

_Dependencies_
- Upgrade `containerd` and `golang.org/x/net` packages to remediate moderate vulnerabilities.

## v1.2.2 [2025-02-05]
_Bug Fixes_
- When DashboardServer executes a dashboard, ensure search path prefix is respected. ([#717](https://github.com/turbot/powerpipe/issues/717))

## v1.2.1 [2025-02-04]
_Bug Fixes_
- Fix backend support if the database is specified by a connection string. ([#713](https://github.com/turbot/powerpipe/issues/713))
- Improve search path button config popover to handle narrower screens. ([#711](https://github.com/turbot/powerpipe/issues/711))
- Dashboard UI sending `changed_input` field at wrong level in `input_changed` event. ([#708](https://github.com/turbot/powerpipe/issues/708))

## v1.2.0 [2025-01-30]
_Whats new_
- Add support for `tailpipe` detections and detection benchmarks.
- Add `tailpipe` connection type.
- Add `detection` command. 
- Add default column support for tables. ([#567](https://github.com/turbot/powerpipe/issues/567))
- Allow multiple benchmarks/controls to be filtered/grouped on a single dashboard. ([#588](https://github.com/turbot/powerpipe/issues/588))
- Add support for equal / not_equal / in / not_in filters. ([#594](https://github.com/turbot/powerpipe/issues/594))
- Show documentation in Benchmark UI. ([#591](https://github.com/turbot/powerpipe/issues/591))
- Combine snap / open button into single split button. ([#606](https://github.com/turbot/powerpipe/issues/606))
- Add table column selection to table panel controls. ([#631](https://github.com/turbot/powerpipe/issues/631))
- Add table row view. ([#636](https://github.com/turbot/powerpipe/issues/636))
- Add support for draggable split pane for dashboard UI right content section. ([#648](https://github.com/turbot/powerpipe/issues/648))
- Update default date_range input to be 7 days. ([#641](https://github.com/turbot/powerpipe/issues/641))
- Add cell filtering to tables. ([#662](https://github.com/turbot/powerpipe/issues/662))
- Add href support to detection column data. ([#642](https://github.com/turbot/powerpipe/issues/642))
- Add table filter side panel for regular tables. ([#670](https://github.com/turbot/powerpipe/issues/670))
- Add 500 MB limit for opening snapshots. ([#671](https://github.com/turbot/powerpipe/issues/671))
- Add control row detail side panel view. ([#672](https://github.com/turbot/powerpipe/issues/672))
- Add heatmap chart. ([#673](https://github.com/turbot/powerpipe/issues/673))

## v1.1.0 [2025-01-20]
_Whats new_
- Add support for installing mods from GitLab repositories. ([#656](https://github.com/turbot/powerpipe/issues/656))

_Dependencies_
- Upgrade `crypto`, `net` and `go-git` packages to remediate critical and high vulnerabilities.

## v1.0.1 [2024-11-21]
_Bug fixes_
- Fix minor spelling issue in query help output. ([#542](https://github.com/turbot/powerpipe/issues/542))
- Update error message to inform users to run `make dashboard_assets` when dashboard assets are not present. ([#524](https://github.com/turbot/powerpipe/issues/524))

## v1.0.0 [2024-10-22]
_Whats new_
- `connection` resource to manage credentials. [Documentation](https://powerpipe.io/docs/reference/config-files/connection).
- `database` property has been added to [mod](https://powerpipe.io/docs/powerpipe-hcl/mod). A database can be a connection reference, connection string, or Pipes workspace to query.

_Deprecations_
- Deprecated `database` CLI arg. See [Setting the Database](https://powerpipe.io/docs/run#selecting-a-database) for the new syntax to set the database.
- Deprecated `POWERPIPE_DATABASE` env var. See [Setting the Database](https://powerpipe.io/docs/run#selecting-a-database) for the new syntax to set the database.
- Deprecated `database` workspace profile arg. See [Setting the Database](https://powerpipe.io/docs/run#selecting-a-database) for the new syntax to set the database.

## v0.4.4 [2024-10-01]
_Bug fixes_
- Fixed the issue where the search path setting was not being retained while navigating to a different dashboard. ([#325](https://github.com/turbot/powerpipe/issues/325))

## v0.4.3 [2024-08-22]
_Whats new_
- Add `json` extension support for duckdb backends. ([#467](https://github.com/turbot/powerpipe/issues/467))

_Bug fixes_
- Fix incorrect mod name in powerpipe help command output. ([#471](https://github.com/turbot/powerpipe/issues/471))

## v0.4.2 [2024-08-14]
_Whats new_
- Compiled with Go 1.22. ([#448](https://github.com/turbot/powerpipe/issues/448))

_Bug fixes_
- Fix issue where CLI notifications were interfering with the Powerpipe JSON outputs resulting in invalid JSON outputs. ([#452](https://github.com/turbot/powerpipe/issues/452))
- Fix issue where powerpipe crashed when running a benchmark with `--dry-run` flag set. ([#455](https://github.com/turbot/powerpipe/issues/455))

## v0.4.1 [2024-07-26]
_Bug fixes_
- Fix issue where the `--arg` flag was not working for control and query runs. ([#439](https://github.com/turbot/powerpipe/issues/439))
- Fix data inconsistency issue in snapshot output when the same control was included in multiple benchmarks. ([#436](https://github.com/turbot/powerpipe/issues/436))

## v0.4.0 [2024-06-10]
_Whats new_
- Update JSON output to handle duplicate column names - append a unique suffix to duplicate column names. ([#375](https://github.com/turbot/powerpipe/issues/375))
- Update snapshot schema and dashboard event schema versions to `20240607`.  ([#406](https://github.com/turbot/powerpipe/issues/406))

_Bug fixes_
- When generating a snapshot  from a benchmark run, the row data is empty any of the rows are in error. ([#366](https://github.com/turbot/powerpipe/issues/366))
- Update mod install to only install or update mods which are command targets (and their dependencies). Set default pull mode for install is latest if there is a target, and minimal if no target is given. ([#381](https://github.com/turbot/powerpipe/issues/381))
- Fix incorrect help message for output in powerpipe benchmark/control run. ([#367](https://github.com/turbot/powerpipe/issues/367))
- Fix issue where `POWERPIPE_PORT` env var was not being honoured. ([#362](https://github.com/turbot/powerpipe/issues/362))
- Update timing metadata output to rename `duration` field to `duration_ms` for consistency with steampipe. ([#368](https://github.com/turbot/powerpipe/issues/368))
- Dashboard graph should not crash if an invalid edge category color is provided. ([#364](https://github.com/turbot/powerpipe/issues/364))
- Dashboard flow/hierarchy components should show panel controls. ([#363](https://github.com/turbot/powerpipe/issues/363))

## v0.3.1 [2024-05-23]
_Bug fixes_
- Respect the app version defined `powerpipe` block of the mod `require` block. ([#405](https://github.com/turbot/pipe-fittings/issues/405))
- Dashboard UI should handle graph categories containing `resource_name` rather than `name`. ([#360](https://github.com/turbot/powerpipe/issues/360))

## v0.3.0 [2024-05-14]
_Whats new_
- Add support for installing mods from a branch or from the local file system. ([#285](https://github.com/turbot/powerpipe/issues/285))

To install from a branch:
```
powerpipe mod install github.com/turbot/steampipe-mod-aws-well-architected#main
```
To reference a mod in the local file system:
```
powerpipe mod install ../mods/local_mod_folder
```

- Add `--pull` flag to `mod`, `dashboard` and `benchmark` commands to control the mod update strategy. ([#352](https://github.com/turbot/powerpipe/issues/352)). Possible update strategies are:

  - `full` - check branch and tags for both latest and accuracy
  - `latest` - update everything to latest, but only branches - not tags - are commit checked (which is the same as latest)
  - `development` - update branches and broken constraints to latest, leave satisfied constraints unchanged
  - `minimal` - only update broken constraints, do not check branches for new commits

## v0.2.0 [2024-04-26]
_Whats new_
- It is now possible to set a timeout for benchmark and dashboard execution. These can be set:
    - in the workspace using config `dashboard_timeout` and `benchmark_timeout`
    - using the `--dashboard-timeout` flag for the the `dashboard run` and `server` commands 
    - using the `--benchmark-timeout` flag for the `benchmark run` commands. 
    - using the environment variables `POWERPIPE_DASHBOARD_TIMEOUT` and `POWERPIPE_BENCHMARK_TIMEOUT` respectively.
  ([#336](https://github.com/turbot/powerpipe/issues/336))
- Support installing private mods using a github app token. ([#381](https://github.com/turbot/pipe-fittings/issues/381)).
- Improve layout of filter and grouping components for control tags and dimensions. ([#263](https://github.com/turbot/powerpipe/issues/263))
- Remove `dashboard input list` and `dashboard input show` commands.
- Add thousands separator to numeric values in dashboard tables. ([#315](https://github.com/turbot/powerpipe/issues/315))
- Only show benchmark cards for statuses that are contained in the current filter and add status to filter on card click. ([#322](https://github.com/turbot/powerpipe/issues/322))

_Bug fixes_
- When calling mod update, respect the argument (if any) and only update specified mods. ([#331](https://github.com/turbot/powerpipe/issues/331))
- Fix `mod update` display of updates to transitive dependencies. ([#288](https://github.com/turbot/powerpipe/issues/288))

## v0.1.3 [2024-03-18]
_Bug fixes_
- When exporting or displaying a check run as a snapshot, ensure the top level panel has a valid summary. ([#274](https://github.com/turbot/powerpipe/issues/274))
- Update `mod list` output to include `resource_name` and `mod` fields. 

## v0.1.2 [2024-03-15]
_Whats new_
- Optimize workspace load time for large workspaces with multiple dependent mods. ([#365](https://github.com/turbot/pipe-fittings/issues/365))

## v0.1.1 [2024-03-07]
_Bug fixes_
* Fix CLI available version check. ([#250](https://github.com/turbot/powerpipe/issues/250))
* Notify when `mod install` creates a default mod. ([#246](https://github.com/turbot/powerpipe/issues/246))
* Remove newline from end of `mod install` output.  ([#247](https://github.com/turbot/powerpipe/issues/247))
* Fix issue where `asff` output was always missing the first row.


## v0.1.0 [2024-03-06]
Introducing Powerpipe - Dashboards for DevOps.

*Benchmarks* - 5,000+ open-source controls from CIS, NIST, PCI, HIPAA, FedRamp and more. Run instantly on your machine or as part of your deployment pipeline.

*Relationship Diagrams* - The only dashboarding tool designed from the ground up to visualize DevOps data. Explore your cloud,understand relationships and drill down to the details.

*Dashboards & Reports* - High level dashboards provide a quick management view. Reports highlight misconfigurations and attention areas. Filter, pivot and snapshot results.

*Code, not clicks* - Our dashboards are code. Version controlled, composable, shareable, easy to edit - designed for the way you work. Join our open-source community!

Learn more at:
- Website - https://powerpipe.io
- Docs - https://powerpipe.io/docs
- Hub - https://hub.powerpipe.io
- Introduction - https://powerpipe.io/blog/introducing-powerpipe