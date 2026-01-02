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

