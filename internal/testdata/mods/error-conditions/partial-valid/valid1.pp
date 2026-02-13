// Valid resources that should load correctly

query "valid_query_1" {
  sql = "select 1 as value"
}

query "valid_query_2" {
  sql = "select 2 as value"
}

control "valid_control" {
  title = "Valid Control"
  sql = "select 'ok' as status, 'test' as resource"
}

benchmark "valid_benchmark" {
  title = "Valid Benchmark"
  children = [
    control.valid_control,
  ]
}

dashboard "valid_dashboard" {
  title = "Valid Dashboard"

  card {
    sql = "select 42 as value"
  }
}
