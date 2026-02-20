// Circular dependency: A -> B -> C -> A

benchmark "benchmark_a" {
  title = "Benchmark A (depends on C)"
  children = [
    benchmark.benchmark_b,
  ]
}

benchmark "benchmark_b" {
  title = "Benchmark B"
  children = [
    benchmark.benchmark_c,
  ]
}

benchmark "benchmark_c" {
  title = "Benchmark C (creates cycle back to A)"
  children = [
    benchmark.benchmark_a,
  ]
}
