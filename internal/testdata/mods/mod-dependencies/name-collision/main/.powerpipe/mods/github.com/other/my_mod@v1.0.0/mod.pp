mod "my_mod" {
  title   = "My Mod (other org)"
  version = "1.0.0"
}

query "other_query" {
  title = "Other Query from other/my_mod"
  sql   = "SELECT 'other_org' as org"
}
