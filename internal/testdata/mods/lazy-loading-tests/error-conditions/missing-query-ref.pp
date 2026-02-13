# Control referencing non-existent query
# This should trigger a "resource not found" error during resolution

control "refs_missing_query" {
  title       = "Control Refs Missing Query"
  description = "This control references a query that does not exist"
  query       = query.nonexistent_query
  severity    = "high"
}

control "refs_missing_in_benchmark" {
  title       = "Control for Benchmark with Missing Ref"
  description = "Valid control, but benchmark refs missing resource"
  sql         = "SELECT 'pass' as status"
}

# Benchmark referencing non-existent control
benchmark "refs_missing_control" {
  title = "Benchmark Refs Missing Control"
  children = [
    control.refs_missing_in_benchmark,
    control.control_that_does_not_exist
  ]
}
