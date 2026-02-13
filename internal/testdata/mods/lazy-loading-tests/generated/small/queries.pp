# Generated queries for lazy loading testing

query "query_0" {
  title       = "Query 0"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 0 as id, 'query_0' as name"
  tags        = local.common_tags
}

query "query_1" {
  title       = "Query 1 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      1 as id,
      'query_1' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_2" {
  title       = "Parameterized Query 2"
  description = "Query with parameters for lazy loading test"
  sql         = "SELECT * FROM data WHERE region = $1 AND count > $2"

  param "region" {
    description = "Region filter"
    default     = var.test_region
  }

  param "threshold" {
    description = "Threshold value"
    default     = var.test_threshold
  }
}

query "query_3" {
  title       = "Control Query 3"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_3' as resource, 'Query 3 check passed' as reason"
}

query "query_4" {
  title       = "Query 4"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 4 as id, 'query_4' as name"
  tags        = local.common_tags
}

query "query_5" {
  title       = "Query 5 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      5 as id,
      'query_5' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_6" {
  title       = "Parameterized Query 6"
  description = "Query with parameters for lazy loading test"
  sql         = "SELECT * FROM data WHERE region = $1 AND count > $2"

  param "region" {
    description = "Region filter"
    default     = var.test_region
  }

  param "threshold" {
    description = "Threshold value"
    default     = var.test_threshold
  }
}

query "query_7" {
  title       = "Control Query 7"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_7' as resource, 'Query 7 check passed' as reason"
}

query "query_8" {
  title       = "Query 8"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 8 as id, 'query_8' as name"
  tags        = local.common_tags
}

query "query_9" {
  title       = "Query 9 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      9 as id,
      'query_9' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_10" {
  title       = "Parameterized Query 10"
  description = "Query with parameters for lazy loading test"
  sql         = "SELECT * FROM data WHERE region = $1 AND count > $2"

  param "region" {
    description = "Region filter"
    default     = var.test_region
  }

  param "threshold" {
    description = "Threshold value"
    default     = var.test_threshold
  }
}

query "query_11" {
  title       = "Control Query 11"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_11' as resource, 'Query 11 check passed' as reason"
}

query "query_12" {
  title       = "Query 12"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 12 as id, 'query_12' as name"
  tags        = local.common_tags
}

query "query_13" {
  title       = "Query 13 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      13 as id,
      'query_13' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_14" {
  title       = "Parameterized Query 14"
  description = "Query with parameters for lazy loading test"
  sql         = "SELECT * FROM data WHERE region = $1 AND count > $2"

  param "region" {
    description = "Region filter"
    default     = var.test_region
  }

  param "threshold" {
    description = "Threshold value"
    default     = var.test_threshold
  }
}

query "query_15" {
  title       = "Control Query 15"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_15' as resource, 'Query 15 check passed' as reason"
}

query "query_16" {
  title       = "Query 16"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 16 as id, 'query_16' as name"
  tags        = local.common_tags
}

query "query_17" {
  title       = "Query 17 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      17 as id,
      'query_17' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_18" {
  title       = "Parameterized Query 18"
  description = "Query with parameters for lazy loading test"
  sql         = "SELECT * FROM data WHERE region = $1 AND count > $2"

  param "region" {
    description = "Region filter"
    default     = var.test_region
  }

  param "threshold" {
    description = "Threshold value"
    default     = var.test_threshold
  }
}

query "query_19" {
  title       = "Control Query 19"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_19' as resource, 'Query 19 check passed' as reason"
}

