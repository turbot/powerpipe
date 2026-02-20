mod "dep_a" {
  title       = "Dependency A"
  description = "First dependency mod"
  version     = "1.0.0"
}

query "helper_query" {
  title = "Dep A Helper Query"
  sql   = "SELECT 'dep_a' as source, 'helper' as type"
}

control "dep_a_control" {
  title = "Dep A Control"
  sql   = "SELECT 'pass' as status, 'dep_a_resource' as resource"
}

benchmark "dep_a_benchmark" {
  title    = "Dep A Benchmark"
  children = [control.dep_a_control]
}
