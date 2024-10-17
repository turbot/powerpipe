mod "mod_with_db_in_resource"{
  title = "Mod with database defined in mod and resource"
  description = "This is a simple mod used for testing the database precedence. This mod has a database(steampipe) specified in mod and a database(sqlite) specified in the query resource."
  database = connection.steampipe.default
}