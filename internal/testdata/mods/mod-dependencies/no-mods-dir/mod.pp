mod "no_mods_dir" {
  title       = "No Mods Dir"
  description = "Mod with no .powerpipe/mods directory"
}

query "local_query" {
  title = "Local Query"
  sql   = "SELECT 'no_mods_dir' as source"
}

control "local_control" {
  title = "Local Control"
  sql   = "SELECT 'pass' as status"
}
