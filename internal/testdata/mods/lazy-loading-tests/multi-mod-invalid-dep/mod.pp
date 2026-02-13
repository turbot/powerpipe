mod "lazy_multi_mod_invalid" {
  title       = "Multi-Mod Invalid Dependency"
  description = "Mod with invalid dependency structure"

  require {
    mod "github.com/test/invalid-dep" {
      version = "v1.0.0"
    }
  }
}

query "main_query" {
  title = "Main Query"
  sql   = "SELECT 'main' as source"
}
