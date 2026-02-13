mod "dep_leaf" {
  title       = "Dependency Leaf"
  description = "Leaf dependency (no further dependencies)"
  version     = "1.0.0"
}

query "leaf_query" {
  title = "Leaf Query"
  sql   = "SELECT 'dep_leaf' as source, 'leaf_level' as level"
}

control "leaf_control" {
  title = "Leaf Control"
  sql   = "SELECT 'pass' as status, 'leaf_resource' as resource"
}
