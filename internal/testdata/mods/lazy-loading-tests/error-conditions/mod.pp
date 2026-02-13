mod "lazy_error_conditions" {
  title       = "Error Conditions Lazy Loading Test"
  description = "Test mod with files designed to trigger specific error paths"
}

# Valid resources for reference
query "valid_query" {
  title = "Valid Query"
  sql   = "SELECT 'pass' as status, 'resource' as resource, 'OK' as reason"
}

control "valid_control" {
  title = "Valid Control"
  sql   = "SELECT 'pass' as status"
}
