# HCL Resource Tags

## Introduction
Hcl resources use several different tags, which can be quite confusing. This document will attempt to explain the purpose of each tag.

## `hcl` tags
These are used for implicit HCL parsing of the resource. We mostly try to use these the parse resources, but sometimes for complex resources, manual parsing is required.

## `cty` tags
These use used to control the cty representation of the resource in the evaluation context. Any property which may be referenced within the HCL _must_ have a cty tag

## `column` tags
These are used to control which properties are shown as columns in the introspection tables. 

## `JSON` tags
In Steampipe, these are only set on resources which implement `DashboardLeafNode`. They are used to control how the resource is serialised in a SteampipeSnapshot
where they appear under the `properties` property.

In Powerpipe, the json tags are required for all resources which may be  

NOTE: `HclResourceImpl` has empty json tags for all properties, to avoid them being serialised in the snapshot. 

