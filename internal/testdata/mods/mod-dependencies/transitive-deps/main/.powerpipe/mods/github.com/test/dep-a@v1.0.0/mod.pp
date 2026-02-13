mod "dep_a" {
  title       = "Dependency A"
  description = "Dependency that itself has a dependency (dep_leaf)"
  version     = "1.0.0"

  require {
    mod "github.com/test/dep-leaf" {
      version = "v1.0.0"
    }
  }
}

query "dep_a_query" {
  title = "Dep A Query"
  sql   = "SELECT 'dep_a' as source"
}

# Control that uses transitive dependency
control "dep_a_control" {
  title = "Dep A Control"
  query = dep_leaf.query.leaf_query
}

benchmark "dep_a_benchmark" {
  title    = "Dep A Benchmark"
  children = [control.dep_a_control]
}
