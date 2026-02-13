mod "dep_b" {
  title       = "Dependency B"
  description = "Independent dependency (no transitive deps)"
  version     = "1.0.0"
}

query "dep_b_query" {
  title = "Dep B Query"
  sql   = "SELECT 'dep_b' as source"
}

control "dep_b_control" {
  title = "Dep B Control"
  sql   = "SELECT 'pass' as status, 'dep_b_resource' as resource"
}
