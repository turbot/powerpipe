query "sqlite_db_query" {
  sql = <<-EOQ
    SELECT COUNT(*) AS "Total Albums" FROM albums;
  EOQ
}