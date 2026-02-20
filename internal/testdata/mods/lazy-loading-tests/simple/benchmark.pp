# Simple benchmark for baseline lazy loading tests

benchmark "simple" {
  title       = "Simple Benchmark"
  description = "A simple benchmark containing two controls"
  tags        = local.common_tags

  children = [
    control.inline_sql,
    control.uses_query
  ]
}
