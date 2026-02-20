# Queries in dependency mod

query "dep_query" {
  title       = "Dependency Query"
  description = "Query from the dependency mod"
  sql         = "SELECT 'dependency' as source, 100 as value"
}

query "dep_shared" {
  title       = "Dependency Shared Query"
  description = "Query shared across controls in dependency mod"
  sql         = "SELECT 'pass' as status, 'dep_resource' as resource, 'Dependency check' as reason"
}

query "dep_parameterized" {
  title       = "Dependency Parameterized Query"
  description = "Query with parameters in dependency mod"
  sql         = "SELECT * FROM dep_data WHERE category = $1"

  param "category" {
    description = "Category filter"
    default     = "default_category"
  }
}

# Same-named query as dep_mod_2 for collision testing
query "shared_name" {
  title       = "Shared Name in Dep Mod"
  description = "Query with same short name as dep_mod_2.query.shared_name"
  sql         = "SELECT 'dep_mod' as source, 'shared' as type"
}
