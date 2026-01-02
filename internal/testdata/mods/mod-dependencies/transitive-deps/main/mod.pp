mod "transitive_main" {
  title       = "Transitive Main"
  description = "Main mod for transitive dependency tests"

  require {
    mod "github.com/test/dep-a" {
      version = "v1.0.0"
    }
    mod "github.com/test/dep-b" {
      version = "v1.0.0"
    }
  }
}

query "main_query" {
  title = "Main Query"
  sql   = "SELECT 'transitive_main' as source"
}

# Control using transitive dependency (dep_a -> dep_leaf)
control "uses_transitive" {
  title       = "Uses Transitive Dependency"
  description = "This control uses a query from dep_leaf through dep_a"
  query       = dep_leaf.query.leaf_query
}

# Benchmark including controls from all levels
benchmark "all_levels" {
  title = "All Dependency Levels"
  children = [
    control.uses_transitive,
    dep_a.control.dep_a_control,
    dep_b.control.dep_b_control
  ]
}
