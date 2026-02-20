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

query "query_20" {
  title       = "Query 20"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 20 as id, 'query_20' as name"
  tags        = local.common_tags
}

query "query_21" {
  title       = "Query 21 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      21 as id,
      'query_21' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_22" {
  title       = "Parameterized Query 22"
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

query "query_23" {
  title       = "Control Query 23"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_23' as resource, 'Query 23 check passed' as reason"
}

query "query_24" {
  title       = "Query 24"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 24 as id, 'query_24' as name"
  tags        = local.common_tags
}

query "query_25" {
  title       = "Query 25 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      25 as id,
      'query_25' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_26" {
  title       = "Parameterized Query 26"
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

query "query_27" {
  title       = "Control Query 27"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_27' as resource, 'Query 27 check passed' as reason"
}

query "query_28" {
  title       = "Query 28"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 28 as id, 'query_28' as name"
  tags        = local.common_tags
}

query "query_29" {
  title       = "Query 29 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      29 as id,
      'query_29' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_30" {
  title       = "Parameterized Query 30"
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

query "query_31" {
  title       = "Control Query 31"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_31' as resource, 'Query 31 check passed' as reason"
}

query "query_32" {
  title       = "Query 32"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 32 as id, 'query_32' as name"
  tags        = local.common_tags
}

query "query_33" {
  title       = "Query 33 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      33 as id,
      'query_33' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_34" {
  title       = "Parameterized Query 34"
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

query "query_35" {
  title       = "Control Query 35"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_35' as resource, 'Query 35 check passed' as reason"
}

query "query_36" {
  title       = "Query 36"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 36 as id, 'query_36' as name"
  tags        = local.common_tags
}

query "query_37" {
  title       = "Query 37 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      37 as id,
      'query_37' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_38" {
  title       = "Parameterized Query 38"
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

query "query_39" {
  title       = "Control Query 39"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_39' as resource, 'Query 39 check passed' as reason"
}

query "query_40" {
  title       = "Query 40"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 40 as id, 'query_40' as name"
  tags        = local.common_tags
}

query "query_41" {
  title       = "Query 41 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      41 as id,
      'query_41' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_42" {
  title       = "Parameterized Query 42"
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

query "query_43" {
  title       = "Control Query 43"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_43' as resource, 'Query 43 check passed' as reason"
}

query "query_44" {
  title       = "Query 44"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 44 as id, 'query_44' as name"
  tags        = local.common_tags
}

query "query_45" {
  title       = "Query 45 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      45 as id,
      'query_45' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_46" {
  title       = "Parameterized Query 46"
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

query "query_47" {
  title       = "Control Query 47"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_47' as resource, 'Query 47 check passed' as reason"
}

query "query_48" {
  title       = "Query 48"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 48 as id, 'query_48' as name"
  tags        = local.common_tags
}

query "query_49" {
  title       = "Query 49 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      49 as id,
      'query_49' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_50" {
  title       = "Parameterized Query 50"
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

query "query_51" {
  title       = "Control Query 51"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_51' as resource, 'Query 51 check passed' as reason"
}

query "query_52" {
  title       = "Query 52"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 52 as id, 'query_52' as name"
  tags        = local.common_tags
}

query "query_53" {
  title       = "Query 53 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      53 as id,
      'query_53' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_54" {
  title       = "Parameterized Query 54"
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

query "query_55" {
  title       = "Control Query 55"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_55' as resource, 'Query 55 check passed' as reason"
}

query "query_56" {
  title       = "Query 56"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 56 as id, 'query_56' as name"
  tags        = local.common_tags
}

query "query_57" {
  title       = "Query 57 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      57 as id,
      'query_57' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_58" {
  title       = "Parameterized Query 58"
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

query "query_59" {
  title       = "Control Query 59"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_59' as resource, 'Query 59 check passed' as reason"
}

query "query_60" {
  title       = "Query 60"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 60 as id, 'query_60' as name"
  tags        = local.common_tags
}

query "query_61" {
  title       = "Query 61 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      61 as id,
      'query_61' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_62" {
  title       = "Parameterized Query 62"
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

query "query_63" {
  title       = "Control Query 63"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_63' as resource, 'Query 63 check passed' as reason"
}

query "query_64" {
  title       = "Query 64"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 64 as id, 'query_64' as name"
  tags        = local.common_tags
}

query "query_65" {
  title       = "Query 65 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      65 as id,
      'query_65' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_66" {
  title       = "Parameterized Query 66"
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

query "query_67" {
  title       = "Control Query 67"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_67' as resource, 'Query 67 check passed' as reason"
}

query "query_68" {
  title       = "Query 68"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 68 as id, 'query_68' as name"
  tags        = local.common_tags
}

query "query_69" {
  title       = "Query 69 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      69 as id,
      'query_69' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_70" {
  title       = "Parameterized Query 70"
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

query "query_71" {
  title       = "Control Query 71"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_71' as resource, 'Query 71 check passed' as reason"
}

query "query_72" {
  title       = "Query 72"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 72 as id, 'query_72' as name"
  tags        = local.common_tags
}

query "query_73" {
  title       = "Query 73 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      73 as id,
      'query_73' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_74" {
  title       = "Parameterized Query 74"
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

query "query_75" {
  title       = "Control Query 75"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_75' as resource, 'Query 75 check passed' as reason"
}

query "query_76" {
  title       = "Query 76"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 76 as id, 'query_76' as name"
  tags        = local.common_tags
}

query "query_77" {
  title       = "Query 77 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      77 as id,
      'query_77' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_78" {
  title       = "Parameterized Query 78"
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

query "query_79" {
  title       = "Control Query 79"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_79' as resource, 'Query 79 check passed' as reason"
}

query "query_80" {
  title       = "Query 80"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 80 as id, 'query_80' as name"
  tags        = local.common_tags
}

query "query_81" {
  title       = "Query 81 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      81 as id,
      'query_81' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_82" {
  title       = "Parameterized Query 82"
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

query "query_83" {
  title       = "Control Query 83"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_83' as resource, 'Query 83 check passed' as reason"
}

query "query_84" {
  title       = "Query 84"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 84 as id, 'query_84' as name"
  tags        = local.common_tags
}

query "query_85" {
  title       = "Query 85 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      85 as id,
      'query_85' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_86" {
  title       = "Parameterized Query 86"
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

query "query_87" {
  title       = "Control Query 87"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_87' as resource, 'Query 87 check passed' as reason"
}

query "query_88" {
  title       = "Query 88"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 88 as id, 'query_88' as name"
  tags        = local.common_tags
}

query "query_89" {
  title       = "Query 89 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      89 as id,
      'query_89' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_90" {
  title       = "Parameterized Query 90"
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

query "query_91" {
  title       = "Control Query 91"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_91' as resource, 'Query 91 check passed' as reason"
}

query "query_92" {
  title       = "Query 92"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 92 as id, 'query_92' as name"
  tags        = local.common_tags
}

query "query_93" {
  title       = "Query 93 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      93 as id,
      'query_93' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_94" {
  title       = "Parameterized Query 94"
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

query "query_95" {
  title       = "Control Query 95"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_95' as resource, 'Query 95 check passed' as reason"
}

query "query_96" {
  title       = "Query 96"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 96 as id, 'query_96' as name"
  tags        = local.common_tags
}

query "query_97" {
  title       = "Query 97 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      97 as id,
      'query_97' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_98" {
  title       = "Parameterized Query 98"
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

query "query_99" {
  title       = "Control Query 99"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_99' as resource, 'Query 99 check passed' as reason"
}

query "query_100" {
  title       = "Query 100"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 100 as id, 'query_100' as name"
  tags        = local.common_tags
}

query "query_101" {
  title       = "Query 101 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      101 as id,
      'query_101' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_102" {
  title       = "Parameterized Query 102"
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

query "query_103" {
  title       = "Control Query 103"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_103' as resource, 'Query 103 check passed' as reason"
}

query "query_104" {
  title       = "Query 104"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 104 as id, 'query_104' as name"
  tags        = local.common_tags
}

query "query_105" {
  title       = "Query 105 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      105 as id,
      'query_105' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_106" {
  title       = "Parameterized Query 106"
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

query "query_107" {
  title       = "Control Query 107"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_107' as resource, 'Query 107 check passed' as reason"
}

query "query_108" {
  title       = "Query 108"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 108 as id, 'query_108' as name"
  tags        = local.common_tags
}

query "query_109" {
  title       = "Query 109 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      109 as id,
      'query_109' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_110" {
  title       = "Parameterized Query 110"
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

query "query_111" {
  title       = "Control Query 111"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_111' as resource, 'Query 111 check passed' as reason"
}

query "query_112" {
  title       = "Query 112"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 112 as id, 'query_112' as name"
  tags        = local.common_tags
}

query "query_113" {
  title       = "Query 113 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      113 as id,
      'query_113' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_114" {
  title       = "Parameterized Query 114"
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

query "query_115" {
  title       = "Control Query 115"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_115' as resource, 'Query 115 check passed' as reason"
}

query "query_116" {
  title       = "Query 116"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 116 as id, 'query_116' as name"
  tags        = local.common_tags
}

query "query_117" {
  title       = "Query 117 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      117 as id,
      'query_117' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_118" {
  title       = "Parameterized Query 118"
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

query "query_119" {
  title       = "Control Query 119"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_119' as resource, 'Query 119 check passed' as reason"
}

query "query_120" {
  title       = "Query 120"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 120 as id, 'query_120' as name"
  tags        = local.common_tags
}

query "query_121" {
  title       = "Query 121 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      121 as id,
      'query_121' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_122" {
  title       = "Parameterized Query 122"
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

query "query_123" {
  title       = "Control Query 123"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_123' as resource, 'Query 123 check passed' as reason"
}

query "query_124" {
  title       = "Query 124"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 124 as id, 'query_124' as name"
  tags        = local.common_tags
}

query "query_125" {
  title       = "Query 125 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      125 as id,
      'query_125' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_126" {
  title       = "Parameterized Query 126"
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

query "query_127" {
  title       = "Control Query 127"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_127' as resource, 'Query 127 check passed' as reason"
}

query "query_128" {
  title       = "Query 128"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 128 as id, 'query_128' as name"
  tags        = local.common_tags
}

query "query_129" {
  title       = "Query 129 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      129 as id,
      'query_129' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_130" {
  title       = "Parameterized Query 130"
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

query "query_131" {
  title       = "Control Query 131"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_131' as resource, 'Query 131 check passed' as reason"
}

query "query_132" {
  title       = "Query 132"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 132 as id, 'query_132' as name"
  tags        = local.common_tags
}

query "query_133" {
  title       = "Query 133 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      133 as id,
      'query_133' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_134" {
  title       = "Parameterized Query 134"
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

query "query_135" {
  title       = "Control Query 135"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_135' as resource, 'Query 135 check passed' as reason"
}

query "query_136" {
  title       = "Query 136"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 136 as id, 'query_136' as name"
  tags        = local.common_tags
}

query "query_137" {
  title       = "Query 137 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      137 as id,
      'query_137' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_138" {
  title       = "Parameterized Query 138"
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

query "query_139" {
  title       = "Control Query 139"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_139' as resource, 'Query 139 check passed' as reason"
}

query "query_140" {
  title       = "Query 140"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 140 as id, 'query_140' as name"
  tags        = local.common_tags
}

query "query_141" {
  title       = "Query 141 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      141 as id,
      'query_141' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_142" {
  title       = "Parameterized Query 142"
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

query "query_143" {
  title       = "Control Query 143"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_143' as resource, 'Query 143 check passed' as reason"
}

query "query_144" {
  title       = "Query 144"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 144 as id, 'query_144' as name"
  tags        = local.common_tags
}

query "query_145" {
  title       = "Query 145 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      145 as id,
      'query_145' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_146" {
  title       = "Parameterized Query 146"
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

query "query_147" {
  title       = "Control Query 147"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_147' as resource, 'Query 147 check passed' as reason"
}

query "query_148" {
  title       = "Query 148"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 148 as id, 'query_148' as name"
  tags        = local.common_tags
}

query "query_149" {
  title       = "Query 149 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      149 as id,
      'query_149' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_150" {
  title       = "Parameterized Query 150"
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

query "query_151" {
  title       = "Control Query 151"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_151' as resource, 'Query 151 check passed' as reason"
}

query "query_152" {
  title       = "Query 152"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 152 as id, 'query_152' as name"
  tags        = local.common_tags
}

query "query_153" {
  title       = "Query 153 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      153 as id,
      'query_153' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_154" {
  title       = "Parameterized Query 154"
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

query "query_155" {
  title       = "Control Query 155"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_155' as resource, 'Query 155 check passed' as reason"
}

query "query_156" {
  title       = "Query 156"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 156 as id, 'query_156' as name"
  tags        = local.common_tags
}

query "query_157" {
  title       = "Query 157 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      157 as id,
      'query_157' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_158" {
  title       = "Parameterized Query 158"
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

query "query_159" {
  title       = "Control Query 159"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_159' as resource, 'Query 159 check passed' as reason"
}

query "query_160" {
  title       = "Query 160"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 160 as id, 'query_160' as name"
  tags        = local.common_tags
}

query "query_161" {
  title       = "Query 161 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      161 as id,
      'query_161' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_162" {
  title       = "Parameterized Query 162"
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

query "query_163" {
  title       = "Control Query 163"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_163' as resource, 'Query 163 check passed' as reason"
}

query "query_164" {
  title       = "Query 164"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 164 as id, 'query_164' as name"
  tags        = local.common_tags
}

query "query_165" {
  title       = "Query 165 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      165 as id,
      'query_165' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_166" {
  title       = "Parameterized Query 166"
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

query "query_167" {
  title       = "Control Query 167"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_167' as resource, 'Query 167 check passed' as reason"
}

query "query_168" {
  title       = "Query 168"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 168 as id, 'query_168' as name"
  tags        = local.common_tags
}

query "query_169" {
  title       = "Query 169 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      169 as id,
      'query_169' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_170" {
  title       = "Parameterized Query 170"
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

query "query_171" {
  title       = "Control Query 171"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_171' as resource, 'Query 171 check passed' as reason"
}

query "query_172" {
  title       = "Query 172"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 172 as id, 'query_172' as name"
  tags        = local.common_tags
}

query "query_173" {
  title       = "Query 173 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      173 as id,
      'query_173' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_174" {
  title       = "Parameterized Query 174"
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

query "query_175" {
  title       = "Control Query 175"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_175' as resource, 'Query 175 check passed' as reason"
}

query "query_176" {
  title       = "Query 176"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 176 as id, 'query_176' as name"
  tags        = local.common_tags
}

query "query_177" {
  title       = "Query 177 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      177 as id,
      'query_177' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_178" {
  title       = "Parameterized Query 178"
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

query "query_179" {
  title       = "Control Query 179"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_179' as resource, 'Query 179 check passed' as reason"
}

query "query_180" {
  title       = "Query 180"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 180 as id, 'query_180' as name"
  tags        = local.common_tags
}

query "query_181" {
  title       = "Query 181 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      181 as id,
      'query_181' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_182" {
  title       = "Parameterized Query 182"
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

query "query_183" {
  title       = "Control Query 183"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_183' as resource, 'Query 183 check passed' as reason"
}

query "query_184" {
  title       = "Query 184"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 184 as id, 'query_184' as name"
  tags        = local.common_tags
}

query "query_185" {
  title       = "Query 185 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      185 as id,
      'query_185' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_186" {
  title       = "Parameterized Query 186"
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

query "query_187" {
  title       = "Control Query 187"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_187' as resource, 'Query 187 check passed' as reason"
}

query "query_188" {
  title       = "Query 188"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 188 as id, 'query_188' as name"
  tags        = local.common_tags
}

query "query_189" {
  title       = "Query 189 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      189 as id,
      'query_189' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_190" {
  title       = "Parameterized Query 190"
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

query "query_191" {
  title       = "Control Query 191"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_191' as resource, 'Query 191 check passed' as reason"
}

query "query_192" {
  title       = "Query 192"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 192 as id, 'query_192' as name"
  tags        = local.common_tags
}

query "query_193" {
  title       = "Query 193 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      193 as id,
      'query_193' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_194" {
  title       = "Parameterized Query 194"
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

query "query_195" {
  title       = "Control Query 195"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_195' as resource, 'Query 195 check passed' as reason"
}

query "query_196" {
  title       = "Query 196"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 196 as id, 'query_196' as name"
  tags        = local.common_tags
}

query "query_197" {
  title       = "Query 197 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      197 as id,
      'query_197' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_198" {
  title       = "Parameterized Query 198"
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

query "query_199" {
  title       = "Control Query 199"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_199' as resource, 'Query 199 check passed' as reason"
}

query "query_200" {
  title       = "Query 200"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 200 as id, 'query_200' as name"
  tags        = local.common_tags
}

query "query_201" {
  title       = "Query 201 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      201 as id,
      'query_201' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_202" {
  title       = "Parameterized Query 202"
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

query "query_203" {
  title       = "Control Query 203"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_203' as resource, 'Query 203 check passed' as reason"
}

query "query_204" {
  title       = "Query 204"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 204 as id, 'query_204' as name"
  tags        = local.common_tags
}

query "query_205" {
  title       = "Query 205 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      205 as id,
      'query_205' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_206" {
  title       = "Parameterized Query 206"
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

query "query_207" {
  title       = "Control Query 207"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_207' as resource, 'Query 207 check passed' as reason"
}

query "query_208" {
  title       = "Query 208"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 208 as id, 'query_208' as name"
  tags        = local.common_tags
}

query "query_209" {
  title       = "Query 209 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      209 as id,
      'query_209' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_210" {
  title       = "Parameterized Query 210"
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

query "query_211" {
  title       = "Control Query 211"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_211' as resource, 'Query 211 check passed' as reason"
}

query "query_212" {
  title       = "Query 212"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 212 as id, 'query_212' as name"
  tags        = local.common_tags
}

query "query_213" {
  title       = "Query 213 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      213 as id,
      'query_213' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_214" {
  title       = "Parameterized Query 214"
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

query "query_215" {
  title       = "Control Query 215"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_215' as resource, 'Query 215 check passed' as reason"
}

query "query_216" {
  title       = "Query 216"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 216 as id, 'query_216' as name"
  tags        = local.common_tags
}

query "query_217" {
  title       = "Query 217 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      217 as id,
      'query_217' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_218" {
  title       = "Parameterized Query 218"
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

query "query_219" {
  title       = "Control Query 219"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_219' as resource, 'Query 219 check passed' as reason"
}

query "query_220" {
  title       = "Query 220"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 220 as id, 'query_220' as name"
  tags        = local.common_tags
}

query "query_221" {
  title       = "Query 221 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      221 as id,
      'query_221' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_222" {
  title       = "Parameterized Query 222"
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

query "query_223" {
  title       = "Control Query 223"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_223' as resource, 'Query 223 check passed' as reason"
}

query "query_224" {
  title       = "Query 224"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 224 as id, 'query_224' as name"
  tags        = local.common_tags
}

query "query_225" {
  title       = "Query 225 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      225 as id,
      'query_225' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_226" {
  title       = "Parameterized Query 226"
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

query "query_227" {
  title       = "Control Query 227"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_227' as resource, 'Query 227 check passed' as reason"
}

query "query_228" {
  title       = "Query 228"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 228 as id, 'query_228' as name"
  tags        = local.common_tags
}

query "query_229" {
  title       = "Query 229 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      229 as id,
      'query_229' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_230" {
  title       = "Parameterized Query 230"
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

query "query_231" {
  title       = "Control Query 231"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_231' as resource, 'Query 231 check passed' as reason"
}

query "query_232" {
  title       = "Query 232"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 232 as id, 'query_232' as name"
  tags        = local.common_tags
}

query "query_233" {
  title       = "Query 233 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      233 as id,
      'query_233' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_234" {
  title       = "Parameterized Query 234"
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

query "query_235" {
  title       = "Control Query 235"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_235' as resource, 'Query 235 check passed' as reason"
}

query "query_236" {
  title       = "Query 236"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 236 as id, 'query_236' as name"
  tags        = local.common_tags
}

query "query_237" {
  title       = "Query 237 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      237 as id,
      'query_237' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_238" {
  title       = "Parameterized Query 238"
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

query "query_239" {
  title       = "Control Query 239"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_239' as resource, 'Query 239 check passed' as reason"
}

query "query_240" {
  title       = "Query 240"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 240 as id, 'query_240' as name"
  tags        = local.common_tags
}

query "query_241" {
  title       = "Query 241 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      241 as id,
      'query_241' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_242" {
  title       = "Parameterized Query 242"
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

query "query_243" {
  title       = "Control Query 243"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_243' as resource, 'Query 243 check passed' as reason"
}

query "query_244" {
  title       = "Query 244"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 244 as id, 'query_244' as name"
  tags        = local.common_tags
}

query "query_245" {
  title       = "Query 245 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      245 as id,
      'query_245' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_246" {
  title       = "Parameterized Query 246"
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

query "query_247" {
  title       = "Control Query 247"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_247' as resource, 'Query 247 check passed' as reason"
}

query "query_248" {
  title       = "Query 248"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 248 as id, 'query_248' as name"
  tags        = local.common_tags
}

query "query_249" {
  title       = "Query 249 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      249 as id,
      'query_249' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_250" {
  title       = "Parameterized Query 250"
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

query "query_251" {
  title       = "Control Query 251"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_251' as resource, 'Query 251 check passed' as reason"
}

query "query_252" {
  title       = "Query 252"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 252 as id, 'query_252' as name"
  tags        = local.common_tags
}

query "query_253" {
  title       = "Query 253 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      253 as id,
      'query_253' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_254" {
  title       = "Parameterized Query 254"
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

query "query_255" {
  title       = "Control Query 255"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_255' as resource, 'Query 255 check passed' as reason"
}

query "query_256" {
  title       = "Query 256"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 256 as id, 'query_256' as name"
  tags        = local.common_tags
}

query "query_257" {
  title       = "Query 257 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      257 as id,
      'query_257' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_258" {
  title       = "Parameterized Query 258"
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

query "query_259" {
  title       = "Control Query 259"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_259' as resource, 'Query 259 check passed' as reason"
}

query "query_260" {
  title       = "Query 260"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 260 as id, 'query_260' as name"
  tags        = local.common_tags
}

query "query_261" {
  title       = "Query 261 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      261 as id,
      'query_261' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_262" {
  title       = "Parameterized Query 262"
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

query "query_263" {
  title       = "Control Query 263"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_263' as resource, 'Query 263 check passed' as reason"
}

query "query_264" {
  title       = "Query 264"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 264 as id, 'query_264' as name"
  tags        = local.common_tags
}

query "query_265" {
  title       = "Query 265 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      265 as id,
      'query_265' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_266" {
  title       = "Parameterized Query 266"
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

query "query_267" {
  title       = "Control Query 267"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_267' as resource, 'Query 267 check passed' as reason"
}

query "query_268" {
  title       = "Query 268"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 268 as id, 'query_268' as name"
  tags        = local.common_tags
}

query "query_269" {
  title       = "Query 269 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      269 as id,
      'query_269' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_270" {
  title       = "Parameterized Query 270"
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

query "query_271" {
  title       = "Control Query 271"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_271' as resource, 'Query 271 check passed' as reason"
}

query "query_272" {
  title       = "Query 272"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 272 as id, 'query_272' as name"
  tags        = local.common_tags
}

query "query_273" {
  title       = "Query 273 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      273 as id,
      'query_273' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_274" {
  title       = "Parameterized Query 274"
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

query "query_275" {
  title       = "Control Query 275"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_275' as resource, 'Query 275 check passed' as reason"
}

query "query_276" {
  title       = "Query 276"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 276 as id, 'query_276' as name"
  tags        = local.common_tags
}

query "query_277" {
  title       = "Query 277 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      277 as id,
      'query_277' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_278" {
  title       = "Parameterized Query 278"
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

query "query_279" {
  title       = "Control Query 279"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_279' as resource, 'Query 279 check passed' as reason"
}

query "query_280" {
  title       = "Query 280"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 280 as id, 'query_280' as name"
  tags        = local.common_tags
}

query "query_281" {
  title       = "Query 281 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      281 as id,
      'query_281' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_282" {
  title       = "Parameterized Query 282"
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

query "query_283" {
  title       = "Control Query 283"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_283' as resource, 'Query 283 check passed' as reason"
}

query "query_284" {
  title       = "Query 284"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 284 as id, 'query_284' as name"
  tags        = local.common_tags
}

query "query_285" {
  title       = "Query 285 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      285 as id,
      'query_285' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_286" {
  title       = "Parameterized Query 286"
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

query "query_287" {
  title       = "Control Query 287"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_287' as resource, 'Query 287 check passed' as reason"
}

query "query_288" {
  title       = "Query 288"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 288 as id, 'query_288' as name"
  tags        = local.common_tags
}

query "query_289" {
  title       = "Query 289 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      289 as id,
      'query_289' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_290" {
  title       = "Parameterized Query 290"
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

query "query_291" {
  title       = "Control Query 291"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_291' as resource, 'Query 291 check passed' as reason"
}

query "query_292" {
  title       = "Query 292"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 292 as id, 'query_292' as name"
  tags        = local.common_tags
}

query "query_293" {
  title       = "Query 293 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      293 as id,
      'query_293' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_294" {
  title       = "Parameterized Query 294"
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

query "query_295" {
  title       = "Control Query 295"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_295' as resource, 'Query 295 check passed' as reason"
}

query "query_296" {
  title       = "Query 296"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 296 as id, 'query_296' as name"
  tags        = local.common_tags
}

query "query_297" {
  title       = "Query 297 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      297 as id,
      'query_297' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_298" {
  title       = "Parameterized Query 298"
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

query "query_299" {
  title       = "Control Query 299"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_299' as resource, 'Query 299 check passed' as reason"
}

query "query_300" {
  title       = "Query 300"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 300 as id, 'query_300' as name"
  tags        = local.common_tags
}

query "query_301" {
  title       = "Query 301 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      301 as id,
      'query_301' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_302" {
  title       = "Parameterized Query 302"
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

query "query_303" {
  title       = "Control Query 303"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_303' as resource, 'Query 303 check passed' as reason"
}

query "query_304" {
  title       = "Query 304"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 304 as id, 'query_304' as name"
  tags        = local.common_tags
}

query "query_305" {
  title       = "Query 305 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      305 as id,
      'query_305' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_306" {
  title       = "Parameterized Query 306"
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

query "query_307" {
  title       = "Control Query 307"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_307' as resource, 'Query 307 check passed' as reason"
}

query "query_308" {
  title       = "Query 308"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 308 as id, 'query_308' as name"
  tags        = local.common_tags
}

query "query_309" {
  title       = "Query 309 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      309 as id,
      'query_309' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_310" {
  title       = "Parameterized Query 310"
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

query "query_311" {
  title       = "Control Query 311"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_311' as resource, 'Query 311 check passed' as reason"
}

query "query_312" {
  title       = "Query 312"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 312 as id, 'query_312' as name"
  tags        = local.common_tags
}

query "query_313" {
  title       = "Query 313 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      313 as id,
      'query_313' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_314" {
  title       = "Parameterized Query 314"
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

query "query_315" {
  title       = "Control Query 315"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_315' as resource, 'Query 315 check passed' as reason"
}

query "query_316" {
  title       = "Query 316"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 316 as id, 'query_316' as name"
  tags        = local.common_tags
}

query "query_317" {
  title       = "Query 317 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      317 as id,
      'query_317' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_318" {
  title       = "Parameterized Query 318"
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

query "query_319" {
  title       = "Control Query 319"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_319' as resource, 'Query 319 check passed' as reason"
}

query "query_320" {
  title       = "Query 320"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 320 as id, 'query_320' as name"
  tags        = local.common_tags
}

query "query_321" {
  title       = "Query 321 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      321 as id,
      'query_321' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_322" {
  title       = "Parameterized Query 322"
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

query "query_323" {
  title       = "Control Query 323"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_323' as resource, 'Query 323 check passed' as reason"
}

query "query_324" {
  title       = "Query 324"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 324 as id, 'query_324' as name"
  tags        = local.common_tags
}

query "query_325" {
  title       = "Query 325 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      325 as id,
      'query_325' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_326" {
  title       = "Parameterized Query 326"
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

query "query_327" {
  title       = "Control Query 327"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_327' as resource, 'Query 327 check passed' as reason"
}

query "query_328" {
  title       = "Query 328"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 328 as id, 'query_328' as name"
  tags        = local.common_tags
}

query "query_329" {
  title       = "Query 329 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      329 as id,
      'query_329' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_330" {
  title       = "Parameterized Query 330"
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

query "query_331" {
  title       = "Control Query 331"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_331' as resource, 'Query 331 check passed' as reason"
}

query "query_332" {
  title       = "Query 332"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 332 as id, 'query_332' as name"
  tags        = local.common_tags
}

query "query_333" {
  title       = "Query 333 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      333 as id,
      'query_333' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_334" {
  title       = "Parameterized Query 334"
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

query "query_335" {
  title       = "Control Query 335"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_335' as resource, 'Query 335 check passed' as reason"
}

query "query_336" {
  title       = "Query 336"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 336 as id, 'query_336' as name"
  tags        = local.common_tags
}

query "query_337" {
  title       = "Query 337 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      337 as id,
      'query_337' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_338" {
  title       = "Parameterized Query 338"
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

query "query_339" {
  title       = "Control Query 339"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_339' as resource, 'Query 339 check passed' as reason"
}

query "query_340" {
  title       = "Query 340"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 340 as id, 'query_340' as name"
  tags        = local.common_tags
}

query "query_341" {
  title       = "Query 341 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      341 as id,
      'query_341' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_342" {
  title       = "Parameterized Query 342"
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

query "query_343" {
  title       = "Control Query 343"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_343' as resource, 'Query 343 check passed' as reason"
}

query "query_344" {
  title       = "Query 344"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 344 as id, 'query_344' as name"
  tags        = local.common_tags
}

query "query_345" {
  title       = "Query 345 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      345 as id,
      'query_345' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_346" {
  title       = "Parameterized Query 346"
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

query "query_347" {
  title       = "Control Query 347"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_347' as resource, 'Query 347 check passed' as reason"
}

query "query_348" {
  title       = "Query 348"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 348 as id, 'query_348' as name"
  tags        = local.common_tags
}

query "query_349" {
  title       = "Query 349 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      349 as id,
      'query_349' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_350" {
  title       = "Parameterized Query 350"
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

query "query_351" {
  title       = "Control Query 351"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_351' as resource, 'Query 351 check passed' as reason"
}

query "query_352" {
  title       = "Query 352"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 352 as id, 'query_352' as name"
  tags        = local.common_tags
}

query "query_353" {
  title       = "Query 353 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      353 as id,
      'query_353' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_354" {
  title       = "Parameterized Query 354"
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

query "query_355" {
  title       = "Control Query 355"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_355' as resource, 'Query 355 check passed' as reason"
}

query "query_356" {
  title       = "Query 356"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 356 as id, 'query_356' as name"
  tags        = local.common_tags
}

query "query_357" {
  title       = "Query 357 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      357 as id,
      'query_357' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_358" {
  title       = "Parameterized Query 358"
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

query "query_359" {
  title       = "Control Query 359"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_359' as resource, 'Query 359 check passed' as reason"
}

query "query_360" {
  title       = "Query 360"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 360 as id, 'query_360' as name"
  tags        = local.common_tags
}

query "query_361" {
  title       = "Query 361 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      361 as id,
      'query_361' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_362" {
  title       = "Parameterized Query 362"
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

query "query_363" {
  title       = "Control Query 363"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_363' as resource, 'Query 363 check passed' as reason"
}

query "query_364" {
  title       = "Query 364"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 364 as id, 'query_364' as name"
  tags        = local.common_tags
}

query "query_365" {
  title       = "Query 365 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      365 as id,
      'query_365' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_366" {
  title       = "Parameterized Query 366"
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

query "query_367" {
  title       = "Control Query 367"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_367' as resource, 'Query 367 check passed' as reason"
}

query "query_368" {
  title       = "Query 368"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 368 as id, 'query_368' as name"
  tags        = local.common_tags
}

query "query_369" {
  title       = "Query 369 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      369 as id,
      'query_369' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_370" {
  title       = "Parameterized Query 370"
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

query "query_371" {
  title       = "Control Query 371"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_371' as resource, 'Query 371 check passed' as reason"
}

query "query_372" {
  title       = "Query 372"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 372 as id, 'query_372' as name"
  tags        = local.common_tags
}

query "query_373" {
  title       = "Query 373 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      373 as id,
      'query_373' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_374" {
  title       = "Parameterized Query 374"
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

query "query_375" {
  title       = "Control Query 375"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_375' as resource, 'Query 375 check passed' as reason"
}

query "query_376" {
  title       = "Query 376"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 376 as id, 'query_376' as name"
  tags        = local.common_tags
}

query "query_377" {
  title       = "Query 377 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      377 as id,
      'query_377' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_378" {
  title       = "Parameterized Query 378"
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

query "query_379" {
  title       = "Control Query 379"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_379' as resource, 'Query 379 check passed' as reason"
}

query "query_380" {
  title       = "Query 380"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 380 as id, 'query_380' as name"
  tags        = local.common_tags
}

query "query_381" {
  title       = "Query 381 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      381 as id,
      'query_381' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_382" {
  title       = "Parameterized Query 382"
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

query "query_383" {
  title       = "Control Query 383"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_383' as resource, 'Query 383 check passed' as reason"
}

query "query_384" {
  title       = "Query 384"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 384 as id, 'query_384' as name"
  tags        = local.common_tags
}

query "query_385" {
  title       = "Query 385 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      385 as id,
      'query_385' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_386" {
  title       = "Parameterized Query 386"
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

query "query_387" {
  title       = "Control Query 387"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_387' as resource, 'Query 387 check passed' as reason"
}

query "query_388" {
  title       = "Query 388"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 388 as id, 'query_388' as name"
  tags        = local.common_tags
}

query "query_389" {
  title       = "Query 389 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      389 as id,
      'query_389' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_390" {
  title       = "Parameterized Query 390"
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

query "query_391" {
  title       = "Control Query 391"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_391' as resource, 'Query 391 check passed' as reason"
}

query "query_392" {
  title       = "Query 392"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 392 as id, 'query_392' as name"
  tags        = local.common_tags
}

query "query_393" {
  title       = "Query 393 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      393 as id,
      'query_393' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_394" {
  title       = "Parameterized Query 394"
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

query "query_395" {
  title       = "Control Query 395"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_395' as resource, 'Query 395 check passed' as reason"
}

query "query_396" {
  title       = "Query 396"
  description = "Simple query for lazy loading test"
  sql         = "SELECT 396 as id, 'query_396' as name"
  tags        = local.common_tags
}

query "query_397" {
  title       = "Query 397 with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      397 as id,
      'query_397' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

query "query_398" {
  title       = "Parameterized Query 398"
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

query "query_399" {
  title       = "Control Query 399"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_399' as resource, 'Query 399 check passed' as reason"
}

