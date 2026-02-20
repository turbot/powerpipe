mod "empty_mods_dir" {
  title       = "Empty Mods Dir"
  description = "Mod with empty .powerpipe/mods directory"
}

query "local_query" {
  title = "Local Query"
  sql   = "SELECT 'empty_mods_dir' as source"
}

control "local_control" {
  title = "Local Control"
  sql   = "SELECT 'pass' as status"
}
