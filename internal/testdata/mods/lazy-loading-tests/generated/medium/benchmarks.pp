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

benchmark "benchmark_2" {
  title       = "Flat Benchmark 2"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_12,
    control.control_13,
    control.control_14,
    control.control_15,
    control.control_16,
    control.control_17
  ]
  tags = local.common_tags
}

benchmark "benchmark_3" {
  title       = "Flat Benchmark 3"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_18,
    control.control_19,
    control.control_20,
    control.control_21,
    control.control_22,
    control.control_23
  ]
  tags = local.common_tags
}

benchmark "benchmark_4" {
  title       = "Flat Benchmark 4"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_24,
    control.control_25,
    control.control_26,
    control.control_27,
    control.control_28,
    control.control_29
  ]
  tags = local.common_tags
}

benchmark "benchmark_5" {
  title       = "Flat Benchmark 5"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_30,
    control.control_31,
    control.control_32,
    control.control_33,
    control.control_34,
    control.control_35
  ]
  tags = local.common_tags
}

benchmark "benchmark_6" {
  title       = "Flat Benchmark 6"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_36,
    control.control_37,
    control.control_38,
    control.control_39,
    control.control_40,
    control.control_41
  ]
  tags = local.common_tags
}

benchmark "benchmark_7" {
  title       = "Flat Benchmark 7"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_42,
    control.control_43,
    control.control_44,
    control.control_45,
    control.control_46,
    control.control_47
  ]
  tags = local.common_tags
}

benchmark "benchmark_8" {
  title       = "Flat Benchmark 8"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_48,
    control.control_49,
    control.control_50,
    control.control_51,
    control.control_52,
    control.control_53
  ]
  tags = local.common_tags
}

benchmark "benchmark_9" {
  title       = "Flat Benchmark 9"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_54,
    control.control_55,
    control.control_56,
    control.control_57,
    control.control_58,
    control.control_59
  ]
  tags = local.common_tags
}

benchmark "benchmark_10" {
  title       = "Flat Benchmark 10"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_60,
    control.control_61,
    control.control_62,
    control.control_63,
    control.control_64,
    control.control_65
  ]
  tags = local.common_tags
}

benchmark "benchmark_11" {
  title       = "Flat Benchmark 11"
  description = "Flat benchmark with 6 controls"
  children = [
    control.control_66,
    control.control_67,
    control.control_68,
    control.control_69,
    control.control_70,
    control.control_71
  ]
  tags = local.common_tags
}

benchmark "nested_root" {
  title       = "Nested Hierarchy Root"
  description = "Root of 5-level deep hierarchy"
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
  title       = "Nested Level 3"
  description = "Intermediate level in hierarchy"
  children = [
    benchmark.nested_level_4
  ]
}

benchmark "nested_level_4" {
  title       = "Nested Level 4"
  description = "Intermediate level in hierarchy"
  children = [
    benchmark.nested_level_5
  ]
}

benchmark "nested_level_5" {
  title       = "Nested Level 5 (Leaf)"
  description = "Deepest level with controls"
  children = [
    control.control_147,
    control.control_148,
    control.control_149
  ]
}

benchmark "wide_root" {
  title       = "Wide Benchmark Root"
  description = "Wide benchmark with many child benchmarks"
  children = [
    benchmark.benchmark_0,
    benchmark.benchmark_1,
    benchmark.benchmark_2,
    benchmark.benchmark_3,
    benchmark.benchmark_4,
    benchmark.benchmark_5,
    benchmark.benchmark_6,
    benchmark.benchmark_7,
    benchmark.benchmark_8,
    benchmark.benchmark_9,
    benchmark.benchmark_10,
    benchmark.benchmark_11
  ]
  tags = local.common_tags
}

