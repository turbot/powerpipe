# Controls with various reference patterns

# Control referencing a query
control "refs_query" {
  title       = "Control Referencing Query"
  description = "Control that references query.control_result"
  query       = query.control_result
  severity    = "medium"
  tags        = local.common_tags
}

# Control referencing shared query
control "refs_shared_query" {
  title       = "Control Referencing Shared Query"
  description = "Control that references query.shared (also used elsewhere)"
  query       = query.shared
  severity    = "low"
}

# Control with inline SQL referencing local
control "refs_local" {
  title       = "Control Referencing Local"
  description = "Control using local in tags"
  sql         = "SELECT 'pass' as status, 'local_ref' as resource, 'Local reference works' as reason"
  tags        = local.common_tags
}

# Control used by multiple benchmarks
control "shared_control" {
  title       = "Shared Control"
  description = "Control referenced by multiple benchmarks"
  sql         = "SELECT 'pass' as status, 'shared_control' as resource, 'Shared control passed' as reason"
  severity    = "high"
}

# Control with parameterized query
control "with_param_query" {
  title       = "Control With Parameterized Query"
  description = "Control that uses the parameterized query"
  query       = query.parameterized
  args = {
    region    = var.default_region
    min_count = var.threshold
  }
}

# Benchmarks that share controls
benchmark "group_a" {
  title       = "Group A Benchmark"
  description = "First benchmark using shared control"
  children = [
    control.refs_query,
    control.shared_control
  ]
}

benchmark "group_b" {
  title       = "Group B Benchmark"
  description = "Second benchmark using shared control"
  children = [
    control.refs_shared_query,
    control.shared_control
  ]
}

benchmark "root" {
  title       = "Cross Reference Root"
  description = "Root benchmark showing cross-reference patterns"
  children = [
    benchmark.group_a,
    benchmark.group_b,
    control.refs_local,
    control.with_param_query
  ]
}
