variable "database" {
  type    = connection
  default = connection.sqlite.albums
}

mod "mod_with_db_var"{
  title = "Mod with database defined through a variable"
  description = "This is a simple mod used for testing the database precedence. This mod has a database(sqlite) specified in the mod.pp. through a variable."
  database = var.database
}