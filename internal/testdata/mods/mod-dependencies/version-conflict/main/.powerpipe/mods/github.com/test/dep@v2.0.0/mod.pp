mod "dep" {
  title   = "Dependency v2.0.0"
  version = "2.0.0"
}

query "dep_query" {
  title = "Dep Query v2"
  sql   = "SELECT 'dep_v2' as source, '2.0.0' as version"
}

# New resource only in v2
query "new_in_v2" {
  title = "New in v2"
  sql   = "SELECT 'new_v2' as source"
}
