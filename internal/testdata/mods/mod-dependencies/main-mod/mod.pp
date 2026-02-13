mod "main_mod" {
  title       = "Main Mod"
  description = "Main mod for dependency tests"

  require {
    mod "github.com/test/dep-a" {
      version = "v1.0.0"
    }
    mod "github.com/test/dep-b" {
      version = "v1.0.0"
    }
  }
}

# Local query
query "main_query" {
  title = "Main Query"
  sql   = "SELECT 'main' as source"
}

# Control using local query
control "local_control" {
  title = "Local Control"
  query = query.main_query
}

# Control using dependency mod's query
control "uses_dep_a_query" {
  title = "Uses Dep A Query"
  query = dep_a.query.helper_query
}

# Control using another dependency mod's query
control "uses_dep_b_query" {
  title = "Uses Dep B Query"
  query = dep_b.query.helper_query
}

# Benchmark mixing local and dependency controls
benchmark "mixed_benchmark" {
  title = "Mixed Sources Benchmark"
  children = [
    control.local_control,
    control.uses_dep_a_query,
    dep_a.control.dep_a_control,
    dep_b.control.dep_b_control
  ]
}

# Dashboard using resources from main and dependency mods
dashboard "main_dashboard" {
  title = "Main Dashboard"

  card {
    title = "Main Value"
    sql   = query.main_query.sql
  }

  card {
    title = "Dep A Value"
    sql   = dep_a.query.helper_query.sql
  }
}
