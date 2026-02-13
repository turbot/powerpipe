# Simple controls for baseline lazy loading tests

control "inline_sql" {
  title       = "Control with Inline SQL"
  description = "Control using inline SQL rather than query reference"
  sql         = "SELECT 'pass' as status, 'test_resource' as resource, 'Inline SQL works' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "uses_query" {
  title       = "Control Using Query Reference"
  description = "Control that references a query resource"
  query       = query.control_query
  severity    = "medium"
  tags        = local.common_tags
}
