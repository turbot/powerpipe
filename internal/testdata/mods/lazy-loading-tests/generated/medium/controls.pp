# Generated controls for lazy loading testing

control "control_0" {
  title       = "Control 0 (Query Ref)"
  description = "Control referencing query.query_0"
  query       = query.query_0
  severity    = "low"
  tags        = local.common_tags
}

control "control_1" {
  title       = "Control 1 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_1' as resource, 'Control 1 passed' as reason"
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
  title       = "Control 5 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_5' as resource, 'Control 5 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_6" {
  title       = "Control 6 (Query Ref)"
  description = "Control referencing query.query_18"
  query       = query.query_18
  severity    = "high"
  tags        = local.common_tags
}

control "control_7" {
  title       = "Control 7 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_7' as resource, 'Control 7 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_8" {
  title       = "Control 8 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_8' as resource, 'Control 8 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_9" {
  title       = "Control 9 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_9' as resource, 'Control 9 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_10" {
  title       = "Control 10 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_10' as resource, 'Control 10 passed' as reason"
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
  title       = "Control 12 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_12' as resource, 'Control 12 passed' as reason"
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
  title       = "Control 14 (Query Ref)"
  description = "Control referencing query.query_42"
  query       = query.query_42
  severity    = "high"
  tags        = local.common_tags
}

control "control_15" {
  title       = "Control 15 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_15' as resource, 'Control 15 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_16" {
  title       = "Control 16 (Query Ref)"
  description = "Control referencing query.query_48"
  query       = query.query_48
  severity    = "low"
  tags        = local.common_tags
}

control "control_17" {
  title       = "Control 17 (Query Ref)"
  description = "Control referencing query.query_51"
  query       = query.query_51
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
  title       = "Control 19 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_19' as resource, 'Control 19 passed' as reason"
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
  title       = "Control 23 (Query Ref)"
  description = "Control referencing query.query_69"
  query       = query.query_69
  severity    = "critical"
  tags        = local.common_tags
}

control "control_24" {
  title       = "Control 24 (Query Ref)"
  description = "Control referencing query.query_72"
  query       = query.query_72
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
  title       = "Control 29 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_29' as resource, 'Control 29 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_30" {
  title       = "Control 30 (Query Ref)"
  description = "Control referencing query.query_90"
  query       = query.query_90
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
  title       = "Control 32 (Query Ref)"
  description = "Control referencing query.query_96"
  query       = query.query_96
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
  title       = "Control 35 (Query Ref)"
  description = "Control referencing query.query_5"
  query       = query.query_5
  severity    = "critical"
  tags        = local.common_tags
}

control "control_36" {
  title       = "Control 36 (Query Ref)"
  description = "Control referencing query.query_8"
  query       = query.query_8
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
  title       = "Control 39 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_39' as resource, 'Control 39 passed' as reason"
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
  title       = "Control 41 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_41' as resource, 'Control 41 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_42" {
  title       = "Control 42 (Query Ref)"
  description = "Control referencing query.query_26"
  query       = query.query_26
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
  title       = "Control 44 (Query Ref)"
  description = "Control referencing query.query_32"
  query       = query.query_32
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
  title       = "Control 47 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_47' as resource, 'Control 47 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_48" {
  title       = "Control 48 (Query Ref)"
  description = "Control referencing query.query_44"
  query       = query.query_44
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
  title       = "Control 51 (Query Ref)"
  description = "Control referencing query.query_53"
  query       = query.query_53
  severity    = "critical"
  tags        = local.common_tags
}

control "control_52" {
  title       = "Control 52 (Query Ref)"
  description = "Control referencing query.query_56"
  query       = query.query_56
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
  title       = "Control 56 (Query Ref)"
  description = "Control referencing query.query_68"
  query       = query.query_68
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
  title       = "Control 62 (Query Ref)"
  description = "Control referencing query.query_86"
  query       = query.query_86
  severity    = "high"
  tags        = local.common_tags
}

control "control_63" {
  title       = "Control 63 (Query Ref)"
  description = "Control referencing query.query_89"
  query       = query.query_89
  severity    = "critical"
  tags        = local.common_tags
}

control "control_64" {
  title       = "Control 64 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_64' as resource, 'Control 64 passed' as reason"
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
  title       = "Control 70 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_70' as resource, 'Control 70 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_71" {
  title       = "Control 71 (Query Ref)"
  description = "Control referencing query.query_13"
  query       = query.query_13
  severity    = "critical"
  tags        = local.common_tags
}

control "control_72" {
  title       = "Control 72 (Query Ref)"
  description = "Control referencing query.query_16"
  query       = query.query_16
  severity    = "low"
  tags        = local.common_tags
}

control "control_73" {
  title       = "Control 73 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_73' as resource, 'Control 73 passed' as reason"
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
  title       = "Control 78 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_78' as resource, 'Control 78 passed' as reason"
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
  title       = "Control 80 (Query Ref)"
  description = "Control referencing query.query_40"
  query       = query.query_40
  severity    = "low"
  tags        = local.common_tags
}

control "control_81" {
  title       = "Control 81 (Query Ref)"
  description = "Control referencing query.query_43"
  query       = query.query_43
  severity    = "medium"
  tags        = local.common_tags
}

control "control_82" {
  title       = "Control 82 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_82' as resource, 'Control 82 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_83" {
  title       = "Control 83 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_83' as resource, 'Control 83 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_84" {
  title       = "Control 84 (Query Ref)"
  description = "Control referencing query.query_52"
  query       = query.query_52
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
  title       = "Control 87 (Query Ref)"
  description = "Control referencing query.query_61"
  query       = query.query_61
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
  title       = "Control 92 (Query Ref)"
  description = "Control referencing query.query_76"
  query       = query.query_76
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
  title       = "Control 94 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_94' as resource, 'Control 94 passed' as reason"
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
  title       = "Control 96 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_96' as resource, 'Control 96 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_97" {
  title       = "Control 97 (Query Ref)"
  description = "Control referencing query.query_91"
  query       = query.query_91
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
  title       = "Control 103 (Query Ref)"
  description = "Control referencing query.query_9"
  query       = query.query_9
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
  title       = "Control 105 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_105' as resource, 'Control 105 passed' as reason"
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
  title       = "Control 107 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_107' as resource, 'Control 107 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_108" {
  title       = "Control 108 (Query Ref)"
  description = "Control referencing query.query_24"
  query       = query.query_24
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
  title       = "Control 112 (Query Ref)"
  description = "Control referencing query.query_36"
  query       = query.query_36
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
  title       = "Control 114 (Query Ref)"
  description = "Control referencing query.query_42"
  query       = query.query_42
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
  title       = "Control 119 (Query Ref)"
  description = "Control referencing query.query_57"
  query       = query.query_57
  severity    = "critical"
  tags        = local.common_tags
}

control "control_120" {
  title       = "Control 120 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_120' as resource, 'Control 120 passed' as reason"
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
  title       = "Control 122 (Query Ref)"
  description = "Control referencing query.query_66"
  query       = query.query_66
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
  title       = "Control 126 (Query Ref)"
  description = "Control referencing query.query_78"
  query       = query.query_78
  severity    = "high"
  tags        = local.common_tags
}

control "control_127" {
  title       = "Control 127 (Query Ref)"
  description = "Control referencing query.query_81"
  query       = query.query_81
  severity    = "critical"
  tags        = local.common_tags
}

control "control_128" {
  title       = "Control 128 (Query Ref)"
  description = "Control referencing query.query_84"
  query       = query.query_84
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
  title       = "Control 132 (Query Ref)"
  description = "Control referencing query.query_96"
  query       = query.query_96
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
  title       = "Control 135 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_135' as resource, 'Control 135 passed' as reason"
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
  title       = "Control 137 (Query Ref)"
  description = "Control referencing query.query_11"
  query       = query.query_11
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
  title       = "Control 139 (Query Ref)"
  description = "Control referencing query.query_17"
  query       = query.query_17
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
  title       = "Control 141 (Query Ref)"
  description = "Control referencing query.query_23"
  query       = query.query_23
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
  title       = "Control 145 (Query Ref)"
  description = "Control referencing query.query_35"
  query       = query.query_35
  severity    = "medium"
  tags        = local.common_tags
}

control "control_146" {
  title       = "Control 146 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_146' as resource, 'Control 146 passed' as reason"
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
  title       = "Control 148 (Query Ref)"
  description = "Control referencing query.query_44"
  query       = query.query_44
  severity    = "low"
  tags        = local.common_tags
}

control "control_149" {
  title       = "Control 149 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_149' as resource, 'Control 149 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

