benchmark "control_reused" {
  title         = "Top Benchmark"
  description   = "Benchmark with 3 children benchmarks which run the same control"
  children      = [
    benchmark.test1,
    benchmark.test2,
    benchmark.test3
  ]
}

benchmark "test1" {
  title         = "Benchmark 1"
  description   = "Benchmark 1"
  children      = [
    control.delay
  ]
}

benchmark "test2" {
  title         = "Benchmark 2"
  description   = "Benchmark 2"
  children      = [
    control.delay
  ]
}

benchmark "test3" {
  title         = "Benchmark 3"
  description   = "Benchmark 3"
  children      = [
    control.delay
  ]
}

control "delay" {
  title         = "Control to sleep"
  description   = "Control to sleep"
  query         = query.sleep
  severity      = "high"
}

query "sleep" {
  title = "Sleep for 10"
  sql = "select pg_sleep(5), 'ok' as status, 'pp' as resource, 'test' as reason"
}