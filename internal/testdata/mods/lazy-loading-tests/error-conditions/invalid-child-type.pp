# Benchmark with invalid child types
# This should trigger type validation errors

# Query (cannot be a benchmark child directly in many contexts)
query "not_a_control" {
  title = "Query That Is Not A Control"
  sql   = "SELECT 1"
}

# Dashboard (definitely not a valid benchmark child)
dashboard "not_a_benchmark_child" {
  title = "Dashboard Is Not Valid Benchmark Child"

  card {
    sql = "SELECT 1"
  }
}

# Benchmark trying to use query as child (invalid)
# Note: Some systems may allow this, adjust if needed
benchmark "invalid_query_child" {
  title       = "Invalid Query Child"
  description = "Benchmark with query as direct child (may be invalid)"
  children = [
    query.not_a_control
  ]
}

# Benchmark trying to use dashboard as child (invalid)
benchmark "invalid_dashboard_child" {
  title       = "Invalid Dashboard Child"
  description = "Benchmark trying to include a dashboard as child"
  children = [
    dashboard.not_a_benchmark_child
  ]
}
