mod "my_mod" {
  title   = "My Mod (test org)"
  version = "1.0.0"
}

query "test_query" {
  title = "Test Query from test/my-mod"
  sql   = "SELECT 'test_org' as org"
}
