mod "complex_test" {
  title = "Complex Test Mod"
}

variable "region" {
  type    = string
  default = "us-east-1"
}

locals {
  common_tags = {
    test = "true"
  }
}

query "parameterized_query" {
  sql = "SELECT * FROM table WHERE region = $1"
  param "region" {
    default = var.region
  }
}

dashboard "complex_dashboard" {
  title = "Complex Dashboard"

  input "selection" {
    type = "select"
    sql  = "SELECT DISTINCT name FROM options"
  }

  container {
    card {
      sql = query.parameterized_query.sql
      args = [self.input.selection.value]
    }

    chart {
      type = "bar"
      sql  = "SELECT * FROM metrics"
    }
  }
}

control "test_control" {
  title = "Test Control"
  sql   = "SELECT 'ok' as status"
}

benchmark "test_benchmark" {
  title    = "Test Benchmark"
  children = [control.test_control]
}
