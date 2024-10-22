query "steampipe_db_query" {
  sql = "select cache_enabled from steampipe_server_settings"
}

query "sqlite_db_query" {
  sql = <<-EOQ
    SELECT COUNT(*) AS "Total Albums" FROM albums;
  EOQ
}
