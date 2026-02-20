# Resources with empty or minimal configurations

# Query with only required field
query "minimal_query" {
  sql = "SELECT 1"
}

# Control with only required field
control "minimal_control" {
  sql = "SELECT 'pass' as status"
}

# Benchmark with empty children (valid but has no controls)
benchmark "empty_children" {
  title    = "Empty Benchmark"
  children = []
}

# Benchmark with single child
benchmark "single_child" {
  title    = "Single Child Benchmark"
  children = [control.minimal_control]
}

# Dashboard with minimal content
dashboard "minimal_dashboard" {
  title = "Minimal"

  card {
    sql = "SELECT 1"
  }
}

# Dashboard with empty container
dashboard "empty_containers" {
  title = "Empty Containers"

  container {
    title = "Empty Container"
  }

  container {
    title = "Container with single card"

    card {
      sql = "SELECT 1"
    }
  }
}

# Query with empty tags
query "empty_tags" {
  title = "Empty Tags Query"
  sql   = "SELECT 1"
  tags  = {}
}

# Control without severity (uses default)
control "no_severity" {
  title       = "No Severity Control"
  description = "Control without explicit severity"
  sql         = "SELECT 'pass' as status, 'resource' as resource, 'OK' as reason"
}
