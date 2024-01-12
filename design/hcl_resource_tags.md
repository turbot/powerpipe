`# HCL Resource Tags

## Introduction
Hcl resources use several different tags, which can be quite confusing. This document will attempt to explain the purpose of each tag.

## `hcl` tags
These are used for implicit HCL parsing of the resource. We mostly try to use these the parse resources, but sometimes for complex resources, manual parsing is required.

## `cty` tags
These use used to control the cty representation of the resource in the evaluation context. Any property which may be referenced within the HCL _must_ have a cty tag

## `column` tags
These are used to control which properties are shown as columns in the introspection tables.

## `JSON` tags
In steampipe, these are only set on resources which implement `DashboardLeafNode` (and `Control`. They are used to add _resource specific_ properties to a SteampipeSnapshot
where they appear under the `properties` property. NOTE: teh HclResourceImple class currently has empty JSON tags for all properties to avoid including them in the snapshot

```
    Resource modconfig.DashboardLeafNode     `json:"properties,omitempty"
    Control *modconfig.Control               `json:"properties,omitempty"`
```

DashboardLeafNode is implemented by all resources which are leaf nodes in the dashboard tree.
Benchmark
Control
Dashboard
DashboardCard
DashboardChart
DashboardContainer
DashboardEdge
DashboardFlow
DashboardGraph
DashboardHierarchy
DashboardImage
DashboardInput
DashboardNode
DashboardTable
DashboardText
DashboardWith