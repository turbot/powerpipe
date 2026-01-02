// Resources that reference non-existent resources

control "control_missing_query" {
  title = "Control with missing query"
  query = query.nonexistent_query
}

benchmark "benchmark_missing_children" {
  title = "Benchmark with missing children"
  children = [
    control.nonexistent_control_1,
    control.nonexistent_control_2,
  ]
}

dashboard "dashboard_missing_child" {
  title = "Dashboard with missing chart"

  chart "chart_missing_query" {
    title = "Chart with missing query"
    query = query.another_nonexistent
  }
}
