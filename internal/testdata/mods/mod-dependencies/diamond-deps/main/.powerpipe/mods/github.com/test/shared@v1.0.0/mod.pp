mod "shared" {
  title   = "Shared Dependency"
  version = "1.0.0"
}

query "shared_query" {
  title = "Shared Query"
  sql   = "SELECT 'shared' as source, 'diamond_bottom' as position"
}

control "shared_control" {
  title = "Shared Control"
  sql   = "SELECT 'pass' as status, 'shared_resource' as resource"
}
