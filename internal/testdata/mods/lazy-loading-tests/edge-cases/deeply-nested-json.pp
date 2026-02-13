# Resources with deeply nested JSON in params and locals

locals {
  deep_config = {
    level_1 = {
      level_2 = {
        level_3 = {
          level_4 = {
            level_5 = {
              value   = "deeply_nested"
              enabled = true
              count   = 42
            }
          }
        }
      }
    }
    metadata = {
      version = "1.0.0"
      author  = "test"
      tags    = ["deep", "nested", "config"]
    }
  }

  array_config = {
    items = [
      {
        id   = 1
        name = "first"
        children = [
          { child_id = 11, child_name = "first_child_1" },
          { child_id = 12, child_name = "first_child_2" }
        ]
      },
      {
        id   = 2
        name = "second"
        children = [
          { child_id = 21, child_name = "second_child_1" }
        ]
      }
    ]
  }
}

query "deep_json_query" {
  title       = "Deep JSON Query"
  description = "Query with deeply nested JSON param defaults"
  sql         = "SELECT $1::jsonb as config"

  param "config" {
    description = "Deeply nested configuration"
    default = {
      database = {
        connection = {
          host     = "localhost"
          port     = 5432
          ssl_mode = "prefer"
          pool = {
            min_size = 5
            max_size = 20
            timeout  = 30
          }
        }
        retry = {
          attempts = 3
          backoff = {
            initial_ms = 100
            max_ms     = 5000
            multiplier = 2
          }
        }
      }
    }
  }
}

dashboard "json_dashboard" {
  title = "JSON Config Dashboard"

  card {
    title = "Nested Value"
    sql   = "SELECT 'deeply_nested' as value"
  }
}
