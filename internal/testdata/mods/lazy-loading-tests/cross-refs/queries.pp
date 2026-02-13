# Queries with various reference patterns

# Base query (no references)
query "base" {
  title = "Base Query"
  sql   = "SELECT 1 as id, 'base' as name"
}

# Query referencing local
query "uses_local" {
  title       = "Query Using Local"
  description = "Query that uses local.sql_template"
  sql         = local.sql_template
  tags        = local.common_tags
}

# Query with parameters referencing variables
query "parameterized" {
  title       = "Parameterized Query"
  description = "Query with param defaults from variables"
  sql         = "SELECT * FROM resources WHERE region = $1 AND count > $2"

  param "region" {
    description = "Region filter"
    default     = var.default_region
  }

  param "min_count" {
    description = "Minimum count threshold"
    default     = var.threshold
  }
}

# Query for control results
query "control_result" {
  title = "Control Result Query"
  sql   = "SELECT 'pass' as status, 'cross_ref_resource' as resource, 'Cross reference check passed' as reason"
}

# Query that will be referenced by multiple resources
query "shared" {
  title       = "Shared Query"
  description = "Query referenced by multiple controls and dashboard panels"
  sql         = "SELECT 'shared_data' as type, 42 as value"
}
