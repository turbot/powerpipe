# Deep hierarchy benchmark chain: benchmark_0 → benchmark_1 → ... → benchmark_10 → control

control "leaf_control" {
  title       = "Leaf Control"
  description = "Control at the deepest level of the hierarchy"
  sql         = "SELECT 'pass' as status, 'deep_resource' as resource, 'Reached leaf at depth 10' as reason"
  severity    = "high"
}

# Level 10 (deepest benchmark, contains the leaf control)
benchmark "level_10" {
  title       = "Level 10 Benchmark"
  description = "Deepest benchmark level"
  children    = [control.leaf_control]
}

# Level 9
benchmark "level_9" {
  title       = "Level 9 Benchmark"
  description = "Depth level 9"
  children    = [benchmark.level_10]
}

# Level 8
benchmark "level_8" {
  title       = "Level 8 Benchmark"
  description = "Depth level 8"
  children    = [benchmark.level_9]
}

# Level 7
benchmark "level_7" {
  title       = "Level 7 Benchmark"
  description = "Depth level 7"
  children    = [benchmark.level_8]
}

# Level 6
benchmark "level_6" {
  title       = "Level 6 Benchmark"
  description = "Depth level 6"
  children    = [benchmark.level_7]
}

# Level 5
benchmark "level_5" {
  title       = "Level 5 Benchmark"
  description = "Depth level 5"
  children    = [benchmark.level_6]
}

# Level 4
benchmark "level_4" {
  title       = "Level 4 Benchmark"
  description = "Depth level 4"
  children    = [benchmark.level_5]
}

# Level 3
benchmark "level_3" {
  title       = "Level 3 Benchmark"
  description = "Depth level 3"
  children    = [benchmark.level_4]
}

# Level 2
benchmark "level_2" {
  title       = "Level 2 Benchmark"
  description = "Depth level 2"
  children    = [benchmark.level_3]
}

# Level 1
benchmark "level_1" {
  title       = "Level 1 Benchmark"
  description = "Depth level 1"
  children    = [benchmark.level_2]
}

# Level 0 (root benchmark)
benchmark "root" {
  title       = "Root Benchmark"
  description = "Root of the deep hierarchy (11 levels total)"
  children    = [benchmark.level_1]
}

# Additional deep hierarchy with branching at each level
control "branch_a_leaf" {
  title = "Branch A Leaf"
  sql   = "SELECT 'pass' as status, 'branch_a' as resource, 'Branch A leaf' as reason"
}

control "branch_b_leaf" {
  title = "Branch B Leaf"
  sql   = "SELECT 'pass' as status, 'branch_b' as resource, 'Branch B leaf' as reason"
}

benchmark "branching_level_3" {
  title    = "Branching Level 3"
  children = [control.branch_a_leaf, control.branch_b_leaf]
}

benchmark "branching_level_2" {
  title    = "Branching Level 2"
  children = [benchmark.branching_level_3]
}

benchmark "branching_level_1" {
  title    = "Branching Level 1"
  children = [benchmark.branching_level_2]
}

benchmark "branching_root" {
  title       = "Branching Root"
  description = "Hierarchy with branches at the leaf level"
  children    = [benchmark.branching_level_1]
}
