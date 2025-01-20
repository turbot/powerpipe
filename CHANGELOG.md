## v1.1.0 [2025-01-20]
_Whats new_
- Add support for installing mods from GitLab repositories. ([#656](https://github.com/turbot/powerpipe/issues/656))

_Dependencies_
- Upgrade `crypto`, `net` and `go-git` packages to remediate critical and high vulnerabilities.

## v1.0.1 [2024-11-21]
_Whats new_
- Update uploaded Pipes snapshot URL to include `/powerpipe`. ([#577](https://github.com/turbot/powerpipe/issues/577))

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