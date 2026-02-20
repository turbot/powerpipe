# Utilities in dependency mod 2

query "utility_query" {
  title       = "Utility Query"
  description = "Utility query from dep_mod_2"
  sql         = "SELECT 'dep_mod_2' as source, 'utility' as type"
}

# Same-named query as dep_mod for collision testing
query "shared_name" {
  title       = "Shared Name in Dep Mod 2"
  description = "Query with same short name as dep_mod.query.shared_name"
  sql         = "SELECT 'dep_mod_2' as source, 'shared' as type"
}

control "utility_control" {
  title       = "Utility Control"
  description = "Control from dep_mod_2"
  sql         = "SELECT 'pass' as status, 'util_resource' as resource, 'Utility check' as reason"
}

benchmark "utility_benchmark" {
  title       = "Utility Benchmark"
  description = "Benchmark from dep_mod_2"
  children = [
    control.utility_control
  ]
}
