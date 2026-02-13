# Invalid dependency - has resources but no mod.pp
query "orphan_query" {
  title = "Orphan Query"
  sql   = "SELECT 'orphan' as source"
}
