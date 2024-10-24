variable "database" {
  type    = connection
  default = connection.steampipe.with_search_path_prefix
}

mod "mod_with_db_steampipe_search_path"{
  title = "Mod with database defined which is a steampipe connection ref"
  description = "This is a simple mod used for testing the database precedence. This mod has a database(steampipe connection) specified in the mod.pp. The steampipe connection has a search_path_prefix defined."
  database = var.database
}

