# Circular benchmark references: A → B → C → A
# This should trigger a cycle detection error

benchmark "circular_a" {
  title       = "Circular Benchmark A"
  description = "First benchmark in circular chain"
  children = [
    benchmark.circular_b
  ]
}

benchmark "circular_b" {
  title       = "Circular Benchmark B"
  description = "Second benchmark in circular chain"
  children = [
    benchmark.circular_c
  ]
}

benchmark "circular_c" {
  title       = "Circular Benchmark C"
  description = "Third benchmark in circular chain - references back to A"
  children = [
    benchmark.circular_a
  ]
}

# Self-referencing benchmark
benchmark "self_reference" {
  title       = "Self Referencing Benchmark"
  description = "Benchmark that references itself"
  children = [
    benchmark.self_reference
  ]
}
