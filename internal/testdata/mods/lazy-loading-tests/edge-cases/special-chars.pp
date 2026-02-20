# Resources with special characters in names (dots, dashes, underscores)

query "query_with_underscore_name" {
  title       = "Query With Underscore"
  description = "Name contains underscores"
  sql         = "SELECT 'underscore_test' as name"
}

query "complex_name_v2_final" {
  title       = "Complex Name v2 Final"
  description = "Name with version-like pattern"
  sql         = "SELECT 'v2' as version"
}

control "control_aws_ec2_check" {
  title       = "AWS EC2 Check"
  description = "Name mimics real-world naming patterns"
  sql         = "SELECT 'pass' as status, 'i-1234567890abcdef0' as resource, 'EC2 instance OK' as reason"
}

control "control_123_numeric_prefix" {
  title       = "Numeric in Name"
  description = "Name starts with letters but has numbers"
  sql         = "SELECT 'pass' as status, 'numeric_123' as resource, 'Numeric name OK' as reason"
}

benchmark "benchmark_v1_0_0" {
  title       = "Benchmark Version 1.0.0"
  description = "Benchmark with version-style name"
  children = [
    control.control_aws_ec2_check,
    control.control_123_numeric_prefix
  ]
}

dashboard "dashboard_2024_q4_report" {
  title       = "2024 Q4 Report Dashboard"
  description = "Name with year and quarter"

  card {
    title = "Q4 Summary"
    sql   = query.query_with_underscore_name.sql
  }
}
