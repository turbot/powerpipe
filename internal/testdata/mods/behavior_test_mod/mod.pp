mod "behavior_test" {
  title       = "Behavior Test Mod"
  description = "Comprehensive test mod covering all resource types"
}

# Variables
variable "region" {
  type        = string
  default     = "us-east-1"
  description = "Region variable for testing"
}

variable "count" {
  type    = number
  default = 10
}

variable "enabled" {
  type    = bool
  default = true
}

# Locals
locals {
  common_tags = {
    test    = "true"
    env     = "test"
    version = "1.0"
  }

  sql_prefix = "SELECT"
}

# Queries
query "simple" {
  title       = "Simple Query"
  description = "A simple SELECT query"
  sql         = "SELECT 1 as value"
  tags        = local.common_tags
}

query "parameterized" {
  title       = "Parameterized Query"
  description = "Query with parameters"
  sql         = "SELECT * FROM data WHERE region = $1 AND count > $2"

  param "region" {
    description = "Region filter"
    default     = var.region
  }

  param "min_count" {
    description = "Minimum count"
    default     = var.count
  }
}

query "for_control" {
  sql = "SELECT 'pass' as status, 'resource1' as resource, 'ok' as reason"
}

query "for_table" {
  sql = "SELECT 'row1' as col1, 'val1' as col2 UNION ALL SELECT 'row2', 'val2'"
}
