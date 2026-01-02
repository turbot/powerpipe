mod "lazy_multi_mod_errors" {
  title       = "Multi-Mod Error Cases"
  description = "Mod with references to missing/invalid dependencies"
}

# Control referencing non-existent mod
control "missing_mod_ref" {
  title       = "Missing Mod Reference"
  description = "Control referencing a query from non-existent mod"
  query       = nonexistent_mod.query.some_query
}

# Benchmark referencing non-existent control in non-existent mod
benchmark "missing_mod_benchmark" {
  title       = "Missing Mod Benchmark"
  description = "Benchmark with children from non-existent mod"
  children = [
    nonexistent_mod.control.some_control
  ]
}

# Local query for self-reference test
query "local_query" {
  title = "Local Query"
  sql   = "SELECT 'local' as source"
}

# Control with self-reference (local query)
control "self_ref" {
  title = "Self Reference"
  query = query.local_query
}
