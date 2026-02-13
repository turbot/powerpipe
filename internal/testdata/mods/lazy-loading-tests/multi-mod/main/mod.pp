mod "lazy_multi_mod_main" {
  title       = "Multi-Mod Main"
  description = "Main mod with dependencies for lazy loading cross-mod tests"

  require {
    mod "github.com/test/dep-mod" {
      version = "v1.0.0"
    }
    mod "github.com/test/dep-mod-2" {
      version = "v2.0.0"
    }
  }
}

# Local resources
query "main_query" {
  title = "Main Mod Query"
  sql   = "SELECT 'main' as source, 42 as value"
}

control "main_control" {
  title       = "Main Control"
  description = "Control defined in main mod"
  sql         = "SELECT 'pass' as status, 'main_resource' as resource, 'Main mod check' as reason"
}

# Control referencing dependency mod's query
control "uses_dep_query" {
  title       = "Uses Dependency Query"
  description = "Control that uses a query from the dependency mod"
  query       = dep_mod.query.dep_query
}

# Benchmark mixing local and dependency controls
benchmark "mixed_sources" {
  title       = "Mixed Sources Benchmark"
  description = "Benchmark with controls from main and dependency mods"
  children = [
    control.main_control,
    control.uses_dep_query,
    dep_mod.control.dep_control
  ]
}

dashboard "main_dashboard" {
  title       = "Main Dashboard"
  description = "Dashboard using resources from main and dependency mods"

  card {
    title = "Main Value"
    width = 6
    sql   = query.main_query.sql
  }

  card {
    title = "Dep Value"
    width = 6
    sql   = dep_mod.query.dep_query.sql
  }
}

# Same-named query as in dep_mod and dep_mod_2 for collision testing
query "shared_name" {
  title       = "Shared Name in Main"
  description = "Query with same short name as dep_mod.query.shared_name"
  sql         = "SELECT 'main' as source, 'shared' as type"
}

# Benchmark including nested benchmark from dep_mod
benchmark "includes_dep_benchmark" {
  title       = "Includes Dependency Benchmark"
  description = "Benchmark that includes a benchmark from dependency mod"
  children = [
    dep_mod.benchmark.dep_benchmark
  ]
}

# Benchmark spanning all mods
benchmark "cross_mod_all" {
  title       = "Cross-Mod All Benchmark"
  description = "Benchmark with children from main, dep_mod, and dep_mod_2"
  children = [
    control.main_control,
    dep_mod.control.dep_control,
    dep_mod_2.control.utility_control
  ]
}

# Control referencing dep_mod_2 query
control "uses_dep2_query" {
  title       = "Uses Dep Mod 2 Query"
  description = "Control that uses a query from dep_mod_2"
  query       = dep_mod_2.query.utility_query
}

# Top-level benchmark to test hierarchy traversal
benchmark "top_level_spanning" {
  title       = "Top Level Spanning Benchmark"
  description = "Top-level benchmark with nested cross-mod children"
  children = [
    benchmark.mixed_sources,
    benchmark.includes_dep_benchmark,
    dep_mod.benchmark.dep_nested
  ]
}
