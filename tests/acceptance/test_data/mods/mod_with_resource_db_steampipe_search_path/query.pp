query "search_path" {
  database = connection.steampipe.with_search_path_prefix_bar
  sql = <<-EOQ
    SHOW search_path;
  EOQ
}
