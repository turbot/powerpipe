benchmark "benchmark_with_tags" {
  title       = "Benchmark With Tags"
  description = "A benchmark with various tags for comparison testing"

  tags = {
    service  = "test_service"
    category = "comparison"
    type     = "acceptance_test"
  }

  children = [
    control.control_with_tags,
    control.control_without_tags,
    benchmark.nested_benchmark
  ]
}

benchmark "nested_benchmark" {
  title       = "Nested Benchmark"
  description = "A nested benchmark for hierarchy testing"

  tags = {
    level = "nested"
  }

  children = [
    control.nested_control
  ]
}

benchmark "benchmark_without_tags" {
  title       = "Benchmark Without Tags"
  description = "A benchmark without tags"

  children = [
    control.control_without_tags
  ]
}

control "control_with_tags" {
  title       = "Control With Tags"
  description = "A control with tags"
  sql         = "SELECT 'ok' as status, 'test' as resource, 'Test result' as reason"

  tags = {
    control_type = "test"
    severity     = "low"
  }
}

control "control_without_tags" {
  title       = "Control Without Tags"
  description = "A control without tags"
  sql         = "SELECT 'ok' as status, 'test2' as resource, 'Another test' as reason"
}

control "nested_control" {
  title       = "Nested Control"
  description = "A control in a nested benchmark"
  sql         = "SELECT 'ok' as status, 'nested' as resource, 'Nested test' as reason"

  tags = {
    level = "nested"
  }
}
