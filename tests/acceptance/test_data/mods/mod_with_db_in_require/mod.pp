

mod "local" {
  title = "mod_with_db_in_require_var"
  require {
    mod "github.com/pskrbasu/powerpipe-mod-db-var" {
      version = "*"
      args = {
        database = connection.sqlite.albums
      }
    }
  }
}