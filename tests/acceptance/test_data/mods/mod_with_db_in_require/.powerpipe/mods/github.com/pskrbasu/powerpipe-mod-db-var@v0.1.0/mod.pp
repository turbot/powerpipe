var "database" {
  type = connection
  default = connection.steampipe.default
}

mod "mod_with_db_var"{
  title = "Mod with database defined"
  description = "This is a simple mod used for testing the database precedence. This mod has a database(steampipe) specified in the mod.pp."
  database = connection.steampipe.default
}