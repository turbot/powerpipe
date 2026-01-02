# Duplicate resource names
# This should trigger duplicate name errors

# First definition of duplicate_query
query "duplicate_query" {
  title = "Duplicate Query - First Definition"
  sql   = "SELECT 'first' as version"
}

# Second definition of duplicate_query (duplicate!)
query "duplicate_query" {
  title = "Duplicate Query - Second Definition"
  sql   = "SELECT 'second' as version"
}

# First definition of duplicate_control
control "duplicate_control" {
  title = "Duplicate Control - First"
  sql   = "SELECT 'pass' as status"
}

# Second definition of duplicate_control (duplicate!)
control "duplicate_control" {
  title = "Duplicate Control - Second"
  sql   = "SELECT 'fail' as status"
}

# Duplicate benchmark
benchmark "duplicate_benchmark" {
  title    = "Duplicate Benchmark - First"
  children = []
}

benchmark "duplicate_benchmark" {
  title    = "Duplicate Benchmark - Second"
  children = []
}
