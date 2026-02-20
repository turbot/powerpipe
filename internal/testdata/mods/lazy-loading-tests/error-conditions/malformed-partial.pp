# Syntactically valid but semantically wrong configurations
# These should parse but fail during resolution or execution

# Control with SQL that will fail at runtime (bad column names for control)
control "bad_columns" {
  title       = "Bad Column Names"
  description = "SQL returns wrong column names for a control"
  sql         = "SELECT 'active' as wrong_status, 'item' as wrong_resource"
}

# Query with param referenced but not defined
query "undefined_param_usage" {
  title       = "Undefined Param Usage"
  description = "Uses $1 but has no param defined"
  sql         = "SELECT * FROM data WHERE id = $1"
}

# Control with args for query but query has no params
control "args_for_paramless_query" {
  title       = "Args For Paramless Query"
  description = "Provides args but referenced query has no params"
  query       = query.valid_query

  args = {
    unused_param = "this won't work"
  }
}

# Benchmark with duplicate children
benchmark "duplicate_children" {
  title       = "Duplicate Children Benchmark"
  description = "Same control listed twice"
  children = [
    control.valid_control,
    control.valid_control
  ]
}

# Query referencing undefined local
query "undefined_local" {
  title       = "Undefined Local Reference"
  description = "References a local that doesn't exist"
  sql         = local.undefined_sql_template
}

# Control with severity not matching enum values
# Note: HCL parser may catch this at parse time
control "invalid_severity" {
  title       = "Invalid Severity Value"
  description = "Severity value not in allowed set"
  sql         = "SELECT 'pass' as status"
  severity    = "super_critical_extreme"
}
