# Controls and benchmarks in dependency mod

control "dep_control" {
  title       = "Dependency Control"
  description = "Control from the dependency mod"
  query       = query.dep_shared
  severity    = "medium"
}

control "dep_control_inline" {
  title       = "Dependency Inline Control"
  description = "Control with inline SQL in dependency mod"
  sql         = "SELECT 'pass' as status, 'inline_dep' as resource, 'Inline dependency' as reason"
  severity    = "low"
}

control "dep_control_param" {
  title       = "Dependency Parameterized Control"
  description = "Control using parameterized query from dependency mod"
  query       = query.dep_parameterized
}

benchmark "dep_benchmark" {
  title       = "Dependency Benchmark"
  description = "Benchmark defined in dependency mod"
  children = [
    control.dep_control,
    control.dep_control_inline
  ]
}

benchmark "dep_nested" {
  title       = "Dependency Nested Benchmark"
  description = "Nested benchmark in dependency mod"
  children = [
    benchmark.dep_benchmark,
    control.dep_control_param
  ]
}
