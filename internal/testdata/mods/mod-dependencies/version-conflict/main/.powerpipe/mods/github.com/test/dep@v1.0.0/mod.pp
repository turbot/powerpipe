mod "dep" {
  title   = "Dependency v1.0.0"
  version = "1.0.0"
}

query "dep_query" {
  title = "Dep Query v1"
  sql   = "SELECT 'dep_v1' as source, '1.0.0' as version"
}
