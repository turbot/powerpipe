mod "version_conflict_main" {
  title       = "Version Conflict Main"
  description = "Tests version conflict scenarios"

  require {
    mod "github.com/test/dep" {
      version = "v1.0.0"
    }
  }
}

query "main_query" {
  title = "Main Query"
  sql   = "SELECT 'version_conflict_main' as source"
}

control "uses_dep" {
  title = "Uses Dep"
  query = dep.query.dep_query
}
