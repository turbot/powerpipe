mod "circular_a" {
  title       = "Circular A"
  description = "Mod A in circular dependency test"

  require {
    mod "github.com/test/circular-b" {
      version = "v1.0.0"
    }
  }
}

query "query_a" {
  title = "Query A"
  sql   = "SELECT 'circular_a' as source"
}

# Control that references circular_b
control "control_a" {
  title = "Control A"
  query = circular_b.query.query_b
}
