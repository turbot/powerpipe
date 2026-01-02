# Generated controls for lazy loading testing

control "control_0" {
  title       = "Control 0 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_0' as resource, 'Control 0 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_1" {
  title       = "Control 1 (Query Ref)"
  description = "Control referencing query.query_3"
  query       = query.query_3
  severity    = "medium"
  tags        = local.common_tags
}

control "control_2" {
  title       = "Control 2 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_2' as resource, 'Control 2 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_3" {
  title       = "Control 3 (Query Ref)"
  description = "Control referencing query.query_9"
  query       = query.query_9
  severity    = "critical"
  tags        = local.common_tags
}

control "control_4" {
  title       = "Control 4 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_4' as resource, 'Control 4 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_5" {
  title       = "Control 5 (Query Ref)"
  description = "Control referencing query.query_15"
  query       = query.query_15
  severity    = "medium"
  tags        = local.common_tags
}

control "control_6" {
  title       = "Control 6 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_6' as resource, 'Control 6 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_7" {
  title       = "Control 7 (Query Ref)"
  description = "Control referencing query.query_21"
  query       = query.query_21
  severity    = "critical"
  tags        = local.common_tags
}

control "control_8" {
  title       = "Control 8 (Query Ref)"
  description = "Control referencing query.query_24"
  query       = query.query_24
  severity    = "low"
  tags        = local.common_tags
}

control "control_9" {
  title       = "Control 9 (Query Ref)"
  description = "Control referencing query.query_27"
  query       = query.query_27
  severity    = "medium"
  tags        = local.common_tags
}

control "control_10" {
  title       = "Control 10 (Query Ref)"
  description = "Control referencing query.query_30"
  query       = query.query_30
  severity    = "high"
  tags        = local.common_tags
}

control "control_11" {
  title       = "Control 11 (Query Ref)"
  description = "Control referencing query.query_33"
  query       = query.query_33
  severity    = "critical"
  tags        = local.common_tags
}

control "control_12" {
  title       = "Control 12 (Query Ref)"
  description = "Control referencing query.query_36"
  query       = query.query_36
  severity    = "low"
  tags        = local.common_tags
}

control "control_13" {
  title       = "Control 13 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_13' as resource, 'Control 13 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_14" {
  title       = "Control 14 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_14' as resource, 'Control 14 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_15" {
  title       = "Control 15 (Query Ref)"
  description = "Control referencing query.query_45"
  query       = query.query_45
  severity    = "critical"
  tags        = local.common_tags
}

control "control_16" {
  title       = "Control 16 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_16' as resource, 'Control 16 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_17" {
  title       = "Control 17 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_17' as resource, 'Control 17 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_18" {
  title       = "Control 18 (Query Ref)"
  description = "Control referencing query.query_54"
  query       = query.query_54
  severity    = "high"
  tags        = local.common_tags
}

control "control_19" {
  title       = "Control 19 (Query Ref)"
  description = "Control referencing query.query_57"
  query       = query.query_57
  severity    = "critical"
  tags        = local.common_tags
}

control "control_20" {
  title       = "Control 20 (Query Ref)"
  description = "Control referencing query.query_60"
  query       = query.query_60
  severity    = "low"
  tags        = local.common_tags
}

control "control_21" {
  title       = "Control 21 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_21' as resource, 'Control 21 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_22" {
  title       = "Control 22 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_22' as resource, 'Control 22 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_23" {
  title       = "Control 23 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_23' as resource, 'Control 23 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_24" {
  title       = "Control 24 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_24' as resource, 'Control 24 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_25" {
  title       = "Control 25 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_25' as resource, 'Control 25 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_26" {
  title       = "Control 26 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_26' as resource, 'Control 26 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_27" {
  title       = "Control 27 (Query Ref)"
  description = "Control referencing query.query_81"
  query       = query.query_81
  severity    = "critical"
  tags        = local.common_tags
}

control "control_28" {
  title       = "Control 28 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_28' as resource, 'Control 28 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_29" {
  title       = "Control 29 (Query Ref)"
  description = "Control referencing query.query_87"
  query       = query.query_87
  severity    = "medium"
  tags        = local.common_tags
}

control "control_30" {
  title       = "Control 30 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_30' as resource, 'Control 30 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_31" {
  title       = "Control 31 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_31' as resource, 'Control 31 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_32" {
  title       = "Control 32 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_32' as resource, 'Control 32 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_33" {
  title       = "Control 33 (Query Ref)"
  description = "Control referencing query.query_99"
  query       = query.query_99
  severity    = "medium"
  tags        = local.common_tags
}

control "control_34" {
  title       = "Control 34 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_34' as resource, 'Control 34 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_35" {
  title       = "Control 35 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_35' as resource, 'Control 35 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_36" {
  title       = "Control 36 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_36' as resource, 'Control 36 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_37" {
  title       = "Control 37 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_37' as resource, 'Control 37 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_38" {
  title       = "Control 38 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_38' as resource, 'Control 38 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_39" {
  title       = "Control 39 (Query Ref)"
  description = "Control referencing query.query_17"
  query       = query.query_17
  severity    = "critical"
  tags        = local.common_tags
}

control "control_40" {
  title       = "Control 40 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_40' as resource, 'Control 40 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_41" {
  title       = "Control 41 (Query Ref)"
  description = "Control referencing query.query_23"
  query       = query.query_23
  severity    = "medium"
  tags        = local.common_tags
}

control "control_42" {
  title       = "Control 42 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_42' as resource, 'Control 42 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_43" {
  title       = "Control 43 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_43' as resource, 'Control 43 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_44" {
  title       = "Control 44 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_44' as resource, 'Control 44 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_45" {
  title       = "Control 45 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_45' as resource, 'Control 45 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_46" {
  title       = "Control 46 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_46' as resource, 'Control 46 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_47" {
  title       = "Control 47 (Query Ref)"
  description = "Control referencing query.query_41"
  query       = query.query_41
  severity    = "critical"
  tags        = local.common_tags
}

control "control_48" {
  title       = "Control 48 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_48' as resource, 'Control 48 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_49" {
  title       = "Control 49 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_49' as resource, 'Control 49 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_50" {
  title       = "Control 50 (Query Ref)"
  description = "Control referencing query.query_50"
  query       = query.query_50
  severity    = "high"
  tags        = local.common_tags
}

control "control_51" {
  title       = "Control 51 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_51' as resource, 'Control 51 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_52" {
  title       = "Control 52 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_52' as resource, 'Control 52 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_53" {
  title       = "Control 53 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_53' as resource, 'Control 53 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_54" {
  title       = "Control 54 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_54' as resource, 'Control 54 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_55" {
  title       = "Control 55 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_55' as resource, 'Control 55 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_56" {
  title       = "Control 56 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_56' as resource, 'Control 56 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_57" {
  title       = "Control 57 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_57' as resource, 'Control 57 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_58" {
  title       = "Control 58 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_58' as resource, 'Control 58 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_59" {
  title       = "Control 59 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_59' as resource, 'Control 59 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_60" {
  title       = "Control 60 (Query Ref)"
  description = "Control referencing query.query_80"
  query       = query.query_80
  severity    = "low"
  tags        = local.common_tags
}

control "control_61" {
  title       = "Control 61 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_61' as resource, 'Control 61 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_62" {
  title       = "Control 62 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_62' as resource, 'Control 62 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_63" {
  title       = "Control 63 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_63' as resource, 'Control 63 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_64" {
  title       = "Control 64 (Query Ref)"
  description = "Control referencing query.query_92"
  query       = query.query_92
  severity    = "low"
  tags        = local.common_tags
}

control "control_65" {
  title       = "Control 65 (Query Ref)"
  description = "Control referencing query.query_95"
  query       = query.query_95
  severity    = "medium"
  tags        = local.common_tags
}

control "control_66" {
  title       = "Control 66 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_66' as resource, 'Control 66 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_67" {
  title       = "Control 67 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_67' as resource, 'Control 67 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_68" {
  title       = "Control 68 (Query Ref)"
  description = "Control referencing query.query_4"
  query       = query.query_4
  severity    = "low"
  tags        = local.common_tags
}

control "control_69" {
  title       = "Control 69 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_69' as resource, 'Control 69 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_70" {
  title       = "Control 70 (Query Ref)"
  description = "Control referencing query.query_10"
  query       = query.query_10
  severity    = "high"
  tags        = local.common_tags
}

control "control_71" {
  title       = "Control 71 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_71' as resource, 'Control 71 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_72" {
  title       = "Control 72 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_72' as resource, 'Control 72 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_73" {
  title       = "Control 73 (Query Ref)"
  description = "Control referencing query.query_19"
  query       = query.query_19
  severity    = "medium"
  tags        = local.common_tags
}

control "control_74" {
  title       = "Control 74 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_74' as resource, 'Control 74 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_75" {
  title       = "Control 75 (Query Ref)"
  description = "Control referencing query.query_25"
  query       = query.query_25
  severity    = "critical"
  tags        = local.common_tags
}

control "control_76" {
  title       = "Control 76 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_76' as resource, 'Control 76 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_77" {
  title       = "Control 77 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_77' as resource, 'Control 77 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_78" {
  title       = "Control 78 (Query Ref)"
  description = "Control referencing query.query_34"
  query       = query.query_34
  severity    = "high"
  tags        = local.common_tags
}

control "control_79" {
  title       = "Control 79 (Query Ref)"
  description = "Control referencing query.query_37"
  query       = query.query_37
  severity    = "critical"
  tags        = local.common_tags
}

control "control_80" {
  title       = "Control 80 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_80' as resource, 'Control 80 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_81" {
  title       = "Control 81 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_81' as resource, 'Control 81 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_82" {
  title       = "Control 82 (Query Ref)"
  description = "Control referencing query.query_46"
  query       = query.query_46
  severity    = "high"
  tags        = local.common_tags
}

control "control_83" {
  title       = "Control 83 (Query Ref)"
  description = "Control referencing query.query_49"
  query       = query.query_49
  severity    = "critical"
  tags        = local.common_tags
}

control "control_84" {
  title       = "Control 84 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_84' as resource, 'Control 84 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_85" {
  title       = "Control 85 (Query Ref)"
  description = "Control referencing query.query_55"
  query       = query.query_55
  severity    = "medium"
  tags        = local.common_tags
}

control "control_86" {
  title       = "Control 86 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_86' as resource, 'Control 86 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_87" {
  title       = "Control 87 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_87' as resource, 'Control 87 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_88" {
  title       = "Control 88 (Query Ref)"
  description = "Control referencing query.query_64"
  query       = query.query_64
  severity    = "low"
  tags        = local.common_tags
}

control "control_89" {
  title       = "Control 89 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_89' as resource, 'Control 89 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_90" {
  title       = "Control 90 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_90' as resource, 'Control 90 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_91" {
  title       = "Control 91 (Query Ref)"
  description = "Control referencing query.query_73"
  query       = query.query_73
  severity    = "critical"
  tags        = local.common_tags
}

control "control_92" {
  title       = "Control 92 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_92' as resource, 'Control 92 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_93" {
  title       = "Control 93 (Query Ref)"
  description = "Control referencing query.query_79"
  query       = query.query_79
  severity    = "medium"
  tags        = local.common_tags
}

control "control_94" {
  title       = "Control 94 (Query Ref)"
  description = "Control referencing query.query_82"
  query       = query.query_82
  severity    = "high"
  tags        = local.common_tags
}

control "control_95" {
  title       = "Control 95 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_95' as resource, 'Control 95 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_96" {
  title       = "Control 96 (Query Ref)"
  description = "Control referencing query.query_88"
  query       = query.query_88
  severity    = "low"
  tags        = local.common_tags
}

control "control_97" {
  title       = "Control 97 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_97' as resource, 'Control 97 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_98" {
  title       = "Control 98 (Query Ref)"
  description = "Control referencing query.query_94"
  query       = query.query_94
  severity    = "high"
  tags        = local.common_tags
}

control "control_99" {
  title       = "Control 99 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_99' as resource, 'Control 99 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_100" {
  title       = "Control 100 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_100' as resource, 'Control 100 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_101" {
  title       = "Control 101 (Query Ref)"
  description = "Control referencing query.query_3"
  query       = query.query_3
  severity    = "medium"
  tags        = local.common_tags
}

control "control_102" {
  title       = "Control 102 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_102' as resource, 'Control 102 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_103" {
  title       = "Control 103 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_103' as resource, 'Control 103 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_104" {
  title       = "Control 104 (Query Ref)"
  description = "Control referencing query.query_12"
  query       = query.query_12
  severity    = "low"
  tags        = local.common_tags
}

control "control_105" {
  title       = "Control 105 (Query Ref)"
  description = "Control referencing query.query_15"
  query       = query.query_15
  severity    = "medium"
  tags        = local.common_tags
}

control "control_106" {
  title       = "Control 106 (Query Ref)"
  description = "Control referencing query.query_18"
  query       = query.query_18
  severity    = "high"
  tags        = local.common_tags
}

control "control_107" {
  title       = "Control 107 (Query Ref)"
  description = "Control referencing query.query_21"
  query       = query.query_21
  severity    = "critical"
  tags        = local.common_tags
}

control "control_108" {
  title       = "Control 108 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_108' as resource, 'Control 108 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_109" {
  title       = "Control 109 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_109' as resource, 'Control 109 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_110" {
  title       = "Control 110 (Query Ref)"
  description = "Control referencing query.query_30"
  query       = query.query_30
  severity    = "high"
  tags        = local.common_tags
}

control "control_111" {
  title       = "Control 111 (Query Ref)"
  description = "Control referencing query.query_33"
  query       = query.query_33
  severity    = "critical"
  tags        = local.common_tags
}

control "control_112" {
  title       = "Control 112 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_112' as resource, 'Control 112 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_113" {
  title       = "Control 113 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_113' as resource, 'Control 113 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_114" {
  title       = "Control 114 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_114' as resource, 'Control 114 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_115" {
  title       = "Control 115 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_115' as resource, 'Control 115 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_116" {
  title       = "Control 116 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_116' as resource, 'Control 116 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_117" {
  title       = "Control 117 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_117' as resource, 'Control 117 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_118" {
  title       = "Control 118 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_118' as resource, 'Control 118 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_119" {
  title       = "Control 119 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_119' as resource, 'Control 119 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_120" {
  title       = "Control 120 (Query Ref)"
  description = "Control referencing query.query_60"
  query       = query.query_60
  severity    = "low"
  tags        = local.common_tags
}

control "control_121" {
  title       = "Control 121 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_121' as resource, 'Control 121 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_122" {
  title       = "Control 122 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_122' as resource, 'Control 122 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_123" {
  title       = "Control 123 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_123' as resource, 'Control 123 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_124" {
  title       = "Control 124 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_124' as resource, 'Control 124 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_125" {
  title       = "Control 125 (Query Ref)"
  description = "Control referencing query.query_75"
  query       = query.query_75
  severity    = "medium"
  tags        = local.common_tags
}

control "control_126" {
  title       = "Control 126 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_126' as resource, 'Control 126 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_127" {
  title       = "Control 127 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_127' as resource, 'Control 127 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_128" {
  title       = "Control 128 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_128' as resource, 'Control 128 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_129" {
  title       = "Control 129 (Query Ref)"
  description = "Control referencing query.query_87"
  query       = query.query_87
  severity    = "medium"
  tags        = local.common_tags
}

control "control_130" {
  title       = "Control 130 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_130' as resource, 'Control 130 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_131" {
  title       = "Control 131 (Query Ref)"
  description = "Control referencing query.query_93"
  query       = query.query_93
  severity    = "critical"
  tags        = local.common_tags
}

control "control_132" {
  title       = "Control 132 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_132' as resource, 'Control 132 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_133" {
  title       = "Control 133 (Query Ref)"
  description = "Control referencing query.query_99"
  query       = query.query_99
  severity    = "medium"
  tags        = local.common_tags
}

control "control_134" {
  title       = "Control 134 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_134' as resource, 'Control 134 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_135" {
  title       = "Control 135 (Query Ref)"
  description = "Control referencing query.query_5"
  query       = query.query_5
  severity    = "critical"
  tags        = local.common_tags
}

control "control_136" {
  title       = "Control 136 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_136' as resource, 'Control 136 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_137" {
  title       = "Control 137 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_137' as resource, 'Control 137 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_138" {
  title       = "Control 138 (Query Ref)"
  description = "Control referencing query.query_14"
  query       = query.query_14
  severity    = "high"
  tags        = local.common_tags
}

control "control_139" {
  title       = "Control 139 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_139' as resource, 'Control 139 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_140" {
  title       = "Control 140 (Query Ref)"
  description = "Control referencing query.query_20"
  query       = query.query_20
  severity    = "low"
  tags        = local.common_tags
}

control "control_141" {
  title       = "Control 141 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_141' as resource, 'Control 141 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_142" {
  title       = "Control 142 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_142' as resource, 'Control 142 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_143" {
  title       = "Control 143 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_143' as resource, 'Control 143 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_144" {
  title       = "Control 144 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_144' as resource, 'Control 144 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_145" {
  title       = "Control 145 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_145' as resource, 'Control 145 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_146" {
  title       = "Control 146 (Query Ref)"
  description = "Control referencing query.query_38"
  query       = query.query_38
  severity    = "high"
  tags        = local.common_tags
}

control "control_147" {
  title       = "Control 147 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_147' as resource, 'Control 147 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_148" {
  title       = "Control 148 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_148' as resource, 'Control 148 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_149" {
  title       = "Control 149 (Query Ref)"
  description = "Control referencing query.query_47"
  query       = query.query_47
  severity    = "medium"
  tags        = local.common_tags
}

