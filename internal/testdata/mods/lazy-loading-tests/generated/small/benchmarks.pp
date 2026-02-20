# Generated benchmarks for lazy loading testing

benchmark "benchmark_0" {
  title       = "Flat Benchmark 0"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_0,
    control.control_1,
    control.control_2,
    control.control_3,
    control.control_4,
    control.control_5
  ]
  tags = local.common_tags
}

benchmark "benchmark_1" {
  title       = "Flat Benchmark 1"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_6,
    control.control_7,
    control.control_8,
    control.control_9,
    control.control_10,
    control.control_11
  ]
  tags = local.common_tags
}

benchmark "nested_root" {
  title       = "Nested Hierarchy Root"
  description = "Root of 3-level deep hierarchy"
  children = [
    benchmark.nested_level_1
  ]
  tags = local.common_tags
}

benchmark "nested_level_1" {
  title       = "Nested Level 1"
  description = "Intermediate level in hierarchy"
  children = [
    benchmark.nested_level_2
  ]
}

benchmark "nested_level_2" {
  title       = "Nested Level 2"
  description = "Intermediate level in hierarchy"
  children = [
    benchmark.nested_level_3
  ]
}

benchmark "nested_level_3" {
  title       = "Nested Level 3 (Leaf)"
  description = "Deepest level with controls"
  children = [
    control.control_27,
    control.control_28,
    control.control_29
  ]
}

benchmark "wide_root" {
  title       = "Wide Benchmark Root"
  description = "Wide benchmark with many child benchmarks"
  children = [
    benchmark.benchmark_0,
    benchmark.benchmark_1
  ]
  tags = local.common_tags
}

