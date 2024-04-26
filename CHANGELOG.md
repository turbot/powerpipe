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