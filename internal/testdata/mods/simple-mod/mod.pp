mod "simple_test" {
  title = "Simple Test Mod"
}

query "simple_query" {
  sql = "SELECT 1"
}

dashboard "simple_dashboard" {
  title = "Simple Dashboard"

  card {
    sql = "SELECT 1 as value"
  }
}
