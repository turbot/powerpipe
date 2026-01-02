# Simple queries for baseline lazy loading tests

query "simple_count" {
  title       = "Simple Count Query"
  description = "Returns a simple count value"
  sql         = "SELECT 42 as count"
  tags        = local.common_tags
}

query "simple_status" {
  title       = "Simple Status Query"
  description = "Returns a status string"
  sql         = "SELECT 'active' as status"
  tags        = local.common_tags
}

query "control_query" {
  title       = "Control Query"
  description = "Query designed to be used by controls"
  sql         = "SELECT 'pass' as status, 'resource_1' as resource, 'All checks passed' as reason"
}
