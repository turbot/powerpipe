mod "missing_dep_mod" {
  title       = "Missing Dep Mod"
  description = "Mod that requires a non-existent dependency"

  require {
    mod "github.com/test/nonexistent-mod" {
      version = "v1.0.0"
    }
  }
}

query "local_query" {
  title = "Local Query"
  sql   = "SELECT 'missing_dep_mod' as source"
}

# Control that references missing dependency
control "uses_missing" {
  title       = "Uses Missing Dependency"
  description = "This should fail when trying to resolve nonexistent_mod.query.some_query"
  query       = nonexistent_mod.query.some_query
}
