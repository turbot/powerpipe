mod "benchmark_only" {
  title = "Benchmark Only Mod"
}

control "ctrl_1" {
  sql = "SELECT 'pass'"
}

control "ctrl_2" {
  sql = "SELECT 'pass'"
}

benchmark "parent" {
  title = "Parent Benchmark"
  children = [
    benchmark.child_1,
    benchmark.child_2
  ]
}

benchmark "child_1" {
  title = "Child 1"
  children = [control.ctrl_1]
}

benchmark "child_2" {
  title = "Child 2"
  children = [control.ctrl_2]
}
