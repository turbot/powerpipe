mod "name_collision_main" {
  title       = "Name Collision Main"
  description = "Tests for mod name collision scenarios"

  require {
    mod "github.com/test/my-mod" {
      version = "v1.0.0"
    }
    mod "github.com/other/my_mod" {
      version = "v1.0.0"
    }
  }
}

query "main_query" {
  title = "Main Query"
  sql   = "SELECT 'name_collision_main' as source"
}
