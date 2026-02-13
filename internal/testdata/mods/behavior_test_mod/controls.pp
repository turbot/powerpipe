# Controls and Benchmarks

control "basic" {
  title       = "Basic Control"
  description = "A simple control for testing"
  sql         = "SELECT 'pass' as status, 'resource1' as resource, 'All good' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "uses_query" {
  title       = "Control Using Query"
  description = "Control that references a query"
  query       = query.for_control
  severity    = "medium"
}

control "with_params" {
  title       = "Control With Parameters"
  description = "Control with query parameters"
  sql         = "SELECT 'pass' as status, $1 as resource, 'Checked' as reason"

  param "resource_name" {
    description = "Resource name to check"
    default     = "test-resource"
  }
}

control "nested_1" {
  title = "Nested Control 1"
  sql   = "SELECT 'pass' as status"
}

control "nested_2" {
  title = "Nested Control 2"
  sql   = "SELECT 'pass' as status"
}

control "nested_3" {
  title = "Nested Control 3"
  sql   = "SELECT 'pass' as status"
}

# Top-level benchmark with nested children
benchmark "top" {
  title       = "Top Level Benchmark"
  description = "Root benchmark for testing hierarchy"
  tags        = local.common_tags

  children = [
    benchmark.child_a,
    benchmark.child_b,
    control.basic
  ]
}

benchmark "child_a" {
  title    = "Child Benchmark A"
  children = [
    control.nested_1,
    control.nested_2
  ]
}

benchmark "child_b" {
  title    = "Child Benchmark B"
  children = [
    control.nested_3
  ]
}

# Flat benchmark (no nesting)
benchmark "flat" {
  title = "Flat Benchmark"
  children = [
    control.basic,
    control.uses_query,
    control.with_params
  ]
}
