variable "sqlite_database" {
  type    = connection
  default = connection.sqlite.albums
}

mod "local" {
  title = "mod_with_db_in_require_var"
  require {
    mod "github.com/pskrbasu/powerpipe-mod-db-var" {
      version = "*"
      args = {
        database_connection = var.sqlite_database
      }
    }
  }
}