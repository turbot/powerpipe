mod "mod_with_db"{
  title = "Mod with database defined"
  description = "This is a simple mod used for testing the database precedence. This mod has a database(sqlite) specified in the mod.pp."
  database = connection.sqlite.albums
}