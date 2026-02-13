mod "circular_b" {
  title       = "Circular B"
  description = "Mod B in circular dependency test"

  require {
    mod "github.com/test/circular-a" {
      version = "v1.0.0"
    }
  }
}

query "query_b" {
  title = "Query B"
  sql   = "SELECT 'circular_b' as source"
}

# Control that references circular_a (creating the cycle)
control "control_b" {
  title = "Control B"
  query = circular_a.query.query_a
}
