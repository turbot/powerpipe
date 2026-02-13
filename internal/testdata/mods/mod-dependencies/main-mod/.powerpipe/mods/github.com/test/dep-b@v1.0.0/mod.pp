mod "dep_b" {
  title       = "Dependency B"
  description = "Second dependency mod"
  version     = "1.0.0"
}

query "helper_query" {
  title = "Dep B Helper Query"
  sql   = "SELECT 'dep_b' as source, 'helper' as type"
}

control "dep_b_control" {
  title = "Dep B Control"
  sql   = "SELECT 'pass' as status, 'dep_b_resource' as resource"
}

benchmark "dep_b_benchmark" {
  title    = "Dep B Benchmark"
  children = [control.dep_b_control]
}
