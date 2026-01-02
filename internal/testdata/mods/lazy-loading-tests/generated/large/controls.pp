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
  title       = "Control 2 (Query Ref)"
  description = "Control referencing query.query_6"
  query       = query.query_6
  severity    = "high"
  tags        = local.common_tags
}

control "control_3" {
  title       = "Control 3 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_3' as resource, 'Control 3 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_4" {
  title       = "Control 4 (Query Ref)"
  description = "Control referencing query.query_12"
  query       = query.query_12
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
  title       = "Control 7 (Query Ref)"
  description = "Control referencing query.query_21"
  query       = query.query_21
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
  title       = "Control 11 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_11' as resource, 'Control 11 passed' as reason"
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
  title       = "Control 13 (Query Ref)"
  description = "Control referencing query.query_39"
  query       = query.query_39
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
  title       = "Control 23 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_23' as resource, 'Control 23 passed' as reason"
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
  title       = "Control 26 (Query Ref)"
  description = "Control referencing query.query_78"
  query       = query.query_78
  severity    = "high"
  tags        = local.common_tags
}

control "control_27" {
  title       = "Control 27 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_27' as resource, 'Control 27 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_28" {
  title       = "Control 28 (Query Ref)"
  description = "Control referencing query.query_84"
  query       = query.query_84
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
  title       = "Control 31 (Query Ref)"
  description = "Control referencing query.query_93"
  query       = query.query_93
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
  description = "Control referencing query.query_105"
  query       = query.query_105
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
  title       = "Control 37 (Query Ref)"
  description = "Control referencing query.query_111"
  query       = query.query_111
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
  description = "Control referencing query.query_117"
  query       = query.query_117
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
  description = "Control referencing query.query_126"
  query       = query.query_126
  severity    = "high"
  tags        = local.common_tags
}

control "control_43" {
  title       = "Control 43 (Query Ref)"
  description = "Control referencing query.query_129"
  query       = query.query_129
  severity    = "critical"
  tags        = local.common_tags
}

control "control_44" {
  title       = "Control 44 (Query Ref)"
  description = "Control referencing query.query_132"
  query       = query.query_132
  severity    = "low"
  tags        = local.common_tags
}

control "control_45" {
  title       = "Control 45 (Query Ref)"
  description = "Control referencing query.query_135"
  query       = query.query_135
  severity    = "medium"
  tags        = local.common_tags
}

control "control_46" {
  title       = "Control 46 (Query Ref)"
  description = "Control referencing query.query_138"
  query       = query.query_138
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
  title       = "Control 50 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_50' as resource, 'Control 50 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_51" {
  title       = "Control 51 (Query Ref)"
  description = "Control referencing query.query_153"
  query       = query.query_153
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
  title       = "Control 55 (Query Ref)"
  description = "Control referencing query.query_165"
  query       = query.query_165
  severity    = "critical"
  tags        = local.common_tags
}

control "control_56" {
  title       = "Control 56 (Query Ref)"
  description = "Control referencing query.query_168"
  query       = query.query_168
  severity    = "low"
  tags        = local.common_tags
}

control "control_57" {
  title       = "Control 57 (Query Ref)"
  description = "Control referencing query.query_171"
  query       = query.query_171
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
  title       = "Control 60 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_60' as resource, 'Control 60 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_61" {
  title       = "Control 61 (Query Ref)"
  description = "Control referencing query.query_183"
  query       = query.query_183
  severity    = "medium"
  tags        = local.common_tags
}

control "control_62" {
  title       = "Control 62 (Query Ref)"
  description = "Control referencing query.query_186"
  query       = query.query_186
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
  description = "Control referencing query.query_192"
  query       = query.query_192
  severity    = "low"
  tags        = local.common_tags
}

control "control_65" {
  title       = "Control 65 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_65' as resource, 'Control 65 passed' as reason"
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
  title       = "Control 67 (Query Ref)"
  description = "Control referencing query.query_201"
  query       = query.query_201
  severity    = "critical"
  tags        = local.common_tags
}

control "control_68" {
  title       = "Control 68 (Query Ref)"
  description = "Control referencing query.query_204"
  query       = query.query_204
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
  description = "Control referencing query.query_225"
  query       = query.query_225
  severity    = "critical"
  tags        = local.common_tags
}

control "control_76" {
  title       = "Control 76 (Query Ref)"
  description = "Control referencing query.query_228"
  query       = query.query_228
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
  title       = "Control 79 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_79' as resource, 'Control 79 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_80" {
  title       = "Control 80 (Query Ref)"
  description = "Control referencing query.query_240"
  query       = query.query_240
  severity    = "low"
  tags        = local.common_tags
}

control "control_81" {
  title       = "Control 81 (Query Ref)"
  description = "Control referencing query.query_243"
  query       = query.query_243
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
  title       = "Control 83 (Query Ref)"
  description = "Control referencing query.query_249"
  query       = query.query_249
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
  title       = "Control 85 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_85' as resource, 'Control 85 passed' as reason"
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
  description = "Control referencing query.query_261"
  query       = query.query_261
  severity    = "critical"
  tags        = local.common_tags
}

control "control_88" {
  title       = "Control 88 (Query Ref)"
  description = "Control referencing query.query_264"
  query       = query.query_264
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
  title       = "Control 91 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_91' as resource, 'Control 91 passed' as reason"
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
  title       = "Control 93 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_93' as resource, 'Control 93 passed' as reason"
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
  title       = "Control 95 (Query Ref)"
  description = "Control referencing query.query_285"
  query       = query.query_285
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
  description = "Control referencing query.query_291"
  query       = query.query_291
  severity    = "medium"
  tags        = local.common_tags
}

control "control_98" {
  title       = "Control 98 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_98' as resource, 'Control 98 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_99" {
  title       = "Control 99 (Query Ref)"
  description = "Control referencing query.query_297"
  query       = query.query_297
  severity    = "critical"
  tags        = local.common_tags
}

control "control_100" {
  title       = "Control 100 (Query Ref)"
  description = "Control referencing query.query_300"
  query       = query.query_300
  severity    = "low"
  tags        = local.common_tags
}

control "control_101" {
  title       = "Control 101 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_101' as resource, 'Control 101 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_102" {
  title       = "Control 102 (Query Ref)"
  description = "Control referencing query.query_306"
  query       = query.query_306
  severity    = "high"
  tags        = local.common_tags
}

control "control_103" {
  title       = "Control 103 (Query Ref)"
  description = "Control referencing query.query_309"
  query       = query.query_309
  severity    = "critical"
  tags        = local.common_tags
}

control "control_104" {
  title       = "Control 104 (Query Ref)"
  description = "Control referencing query.query_312"
  query       = query.query_312
  severity    = "low"
  tags        = local.common_tags
}

control "control_105" {
  title       = "Control 105 (Query Ref)"
  description = "Control referencing query.query_315"
  query       = query.query_315
  severity    = "medium"
  tags        = local.common_tags
}

control "control_106" {
  title       = "Control 106 (Query Ref)"
  description = "Control referencing query.query_318"
  query       = query.query_318
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
  description = "Control referencing query.query_324"
  query       = query.query_324
  severity    = "low"
  tags        = local.common_tags
}

control "control_109" {
  title       = "Control 109 (Query Ref)"
  description = "Control referencing query.query_327"
  query       = query.query_327
  severity    = "medium"
  tags        = local.common_tags
}

control "control_110" {
  title       = "Control 110 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_110' as resource, 'Control 110 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_111" {
  title       = "Control 111 (Query Ref)"
  description = "Control referencing query.query_333"
  query       = query.query_333
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
  title       = "Control 114 (Query Ref)"
  description = "Control referencing query.query_342"
  query       = query.query_342
  severity    = "high"
  tags        = local.common_tags
}

control "control_115" {
  title       = "Control 115 (Query Ref)"
  description = "Control referencing query.query_345"
  query       = query.query_345
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
  title       = "Control 117 (Query Ref)"
  description = "Control referencing query.query_351"
  query       = query.query_351
  severity    = "medium"
  tags        = local.common_tags
}

control "control_118" {
  title       = "Control 118 (Query Ref)"
  description = "Control referencing query.query_354"
  query       = query.query_354
  severity    = "high"
  tags        = local.common_tags
}

control "control_119" {
  title       = "Control 119 (Query Ref)"
  description = "Control referencing query.query_357"
  query       = query.query_357
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
  description = "Control referencing query.query_366"
  query       = query.query_366
  severity    = "high"
  tags        = local.common_tags
}

control "control_123" {
  title       = "Control 123 (Query Ref)"
  description = "Control referencing query.query_369"
  query       = query.query_369
  severity    = "critical"
  tags        = local.common_tags
}

control "control_124" {
  title       = "Control 124 (Query Ref)"
  description = "Control referencing query.query_372"
  query       = query.query_372
  severity    = "low"
  tags        = local.common_tags
}

control "control_125" {
  title       = "Control 125 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_125' as resource, 'Control 125 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_126" {
  title       = "Control 126 (Query Ref)"
  description = "Control referencing query.query_378"
  query       = query.query_378
  severity    = "high"
  tags        = local.common_tags
}

control "control_127" {
  title       = "Control 127 (Query Ref)"
  description = "Control referencing query.query_381"
  query       = query.query_381
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
  title       = "Control 129 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_129' as resource, 'Control 129 passed' as reason"
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
  title       = "Control 131 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_131' as resource, 'Control 131 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_132" {
  title       = "Control 132 (Query Ref)"
  description = "Control referencing query.query_396"
  query       = query.query_396
  severity    = "low"
  tags        = local.common_tags
}

control "control_133" {
  title       = "Control 133 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_133' as resource, 'Control 133 passed' as reason"
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
  title       = "Control 137 (Query Ref)"
  description = "Control referencing query.query_11"
  query       = query.query_11
  severity    = "medium"
  tags        = local.common_tags
}

control "control_138" {
  title       = "Control 138 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_138' as resource, 'Control 138 passed' as reason"
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
  title       = "Control 140 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_140' as resource, 'Control 140 passed' as reason"
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
  title       = "Control 142 (Query Ref)"
  description = "Control referencing query.query_26"
  query       = query.query_26
  severity    = "high"
  tags        = local.common_tags
}

control "control_143" {
  title       = "Control 143 (Query Ref)"
  description = "Control referencing query.query_29"
  query       = query.query_29
  severity    = "critical"
  tags        = local.common_tags
}

control "control_144" {
  title       = "Control 144 (Query Ref)"
  description = "Control referencing query.query_32"
  query       = query.query_32
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
  title       = "Control 148 (Query Ref)"
  description = "Control referencing query.query_44"
  query       = query.query_44
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

control "control_150" {
  title       = "Control 150 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_150' as resource, 'Control 150 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_151" {
  title       = "Control 151 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_151' as resource, 'Control 151 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_152" {
  title       = "Control 152 (Query Ref)"
  description = "Control referencing query.query_56"
  query       = query.query_56
  severity    = "low"
  tags        = local.common_tags
}

control "control_153" {
  title       = "Control 153 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_153' as resource, 'Control 153 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_154" {
  title       = "Control 154 (Query Ref)"
  description = "Control referencing query.query_62"
  query       = query.query_62
  severity    = "high"
  tags        = local.common_tags
}

control "control_155" {
  title       = "Control 155 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_155' as resource, 'Control 155 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_156" {
  title       = "Control 156 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_156' as resource, 'Control 156 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_157" {
  title       = "Control 157 (Query Ref)"
  description = "Control referencing query.query_71"
  query       = query.query_71
  severity    = "medium"
  tags        = local.common_tags
}

control "control_158" {
  title       = "Control 158 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_158' as resource, 'Control 158 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_159" {
  title       = "Control 159 (Query Ref)"
  description = "Control referencing query.query_77"
  query       = query.query_77
  severity    = "critical"
  tags        = local.common_tags
}

control "control_160" {
  title       = "Control 160 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_160' as resource, 'Control 160 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_161" {
  title       = "Control 161 (Query Ref)"
  description = "Control referencing query.query_83"
  query       = query.query_83
  severity    = "medium"
  tags        = local.common_tags
}

control "control_162" {
  title       = "Control 162 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_162' as resource, 'Control 162 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_163" {
  title       = "Control 163 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_163' as resource, 'Control 163 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_164" {
  title       = "Control 164 (Query Ref)"
  description = "Control referencing query.query_92"
  query       = query.query_92
  severity    = "low"
  tags        = local.common_tags
}

control "control_165" {
  title       = "Control 165 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_165' as resource, 'Control 165 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_166" {
  title       = "Control 166 (Query Ref)"
  description = "Control referencing query.query_98"
  query       = query.query_98
  severity    = "high"
  tags        = local.common_tags
}

control "control_167" {
  title       = "Control 167 (Query Ref)"
  description = "Control referencing query.query_101"
  query       = query.query_101
  severity    = "critical"
  tags        = local.common_tags
}

control "control_168" {
  title       = "Control 168 (Query Ref)"
  description = "Control referencing query.query_104"
  query       = query.query_104
  severity    = "low"
  tags        = local.common_tags
}

control "control_169" {
  title       = "Control 169 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_169' as resource, 'Control 169 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_170" {
  title       = "Control 170 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_170' as resource, 'Control 170 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_171" {
  title       = "Control 171 (Query Ref)"
  description = "Control referencing query.query_113"
  query       = query.query_113
  severity    = "critical"
  tags        = local.common_tags
}

control "control_172" {
  title       = "Control 172 (Query Ref)"
  description = "Control referencing query.query_116"
  query       = query.query_116
  severity    = "low"
  tags        = local.common_tags
}

control "control_173" {
  title       = "Control 173 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_173' as resource, 'Control 173 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_174" {
  title       = "Control 174 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_174' as resource, 'Control 174 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_175" {
  title       = "Control 175 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_175' as resource, 'Control 175 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_176" {
  title       = "Control 176 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_176' as resource, 'Control 176 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_177" {
  title       = "Control 177 (Query Ref)"
  description = "Control referencing query.query_131"
  query       = query.query_131
  severity    = "medium"
  tags        = local.common_tags
}

control "control_178" {
  title       = "Control 178 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_178' as resource, 'Control 178 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_179" {
  title       = "Control 179 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_179' as resource, 'Control 179 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_180" {
  title       = "Control 180 (Query Ref)"
  description = "Control referencing query.query_140"
  query       = query.query_140
  severity    = "low"
  tags        = local.common_tags
}

control "control_181" {
  title       = "Control 181 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_181' as resource, 'Control 181 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_182" {
  title       = "Control 182 (Query Ref)"
  description = "Control referencing query.query_146"
  query       = query.query_146
  severity    = "high"
  tags        = local.common_tags
}

control "control_183" {
  title       = "Control 183 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_183' as resource, 'Control 183 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_184" {
  title       = "Control 184 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_184' as resource, 'Control 184 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_185" {
  title       = "Control 185 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_185' as resource, 'Control 185 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_186" {
  title       = "Control 186 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_186' as resource, 'Control 186 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_187" {
  title       = "Control 187 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_187' as resource, 'Control 187 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_188" {
  title       = "Control 188 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_188' as resource, 'Control 188 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_189" {
  title       = "Control 189 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_189' as resource, 'Control 189 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_190" {
  title       = "Control 190 (Query Ref)"
  description = "Control referencing query.query_170"
  query       = query.query_170
  severity    = "high"
  tags        = local.common_tags
}

control "control_191" {
  title       = "Control 191 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_191' as resource, 'Control 191 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_192" {
  title       = "Control 192 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_192' as resource, 'Control 192 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_193" {
  title       = "Control 193 (Query Ref)"
  description = "Control referencing query.query_179"
  query       = query.query_179
  severity    = "medium"
  tags        = local.common_tags
}

control "control_194" {
  title       = "Control 194 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_194' as resource, 'Control 194 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_195" {
  title       = "Control 195 (Query Ref)"
  description = "Control referencing query.query_185"
  query       = query.query_185
  severity    = "critical"
  tags        = local.common_tags
}

control "control_196" {
  title       = "Control 196 (Query Ref)"
  description = "Control referencing query.query_188"
  query       = query.query_188
  severity    = "low"
  tags        = local.common_tags
}

control "control_197" {
  title       = "Control 197 (Query Ref)"
  description = "Control referencing query.query_191"
  query       = query.query_191
  severity    = "medium"
  tags        = local.common_tags
}

control "control_198" {
  title       = "Control 198 (Query Ref)"
  description = "Control referencing query.query_194"
  query       = query.query_194
  severity    = "high"
  tags        = local.common_tags
}

control "control_199" {
  title       = "Control 199 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_199' as resource, 'Control 199 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_200" {
  title       = "Control 200 (Query Ref)"
  description = "Control referencing query.query_200"
  query       = query.query_200
  severity    = "low"
  tags        = local.common_tags
}

control "control_201" {
  title       = "Control 201 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_201' as resource, 'Control 201 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_202" {
  title       = "Control 202 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_202' as resource, 'Control 202 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_203" {
  title       = "Control 203 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_203' as resource, 'Control 203 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_204" {
  title       = "Control 204 (Query Ref)"
  description = "Control referencing query.query_212"
  query       = query.query_212
  severity    = "low"
  tags        = local.common_tags
}

control "control_205" {
  title       = "Control 205 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_205' as resource, 'Control 205 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_206" {
  title       = "Control 206 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_206' as resource, 'Control 206 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_207" {
  title       = "Control 207 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_207' as resource, 'Control 207 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_208" {
  title       = "Control 208 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_208' as resource, 'Control 208 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_209" {
  title       = "Control 209 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_209' as resource, 'Control 209 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_210" {
  title       = "Control 210 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_210' as resource, 'Control 210 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_211" {
  title       = "Control 211 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_211' as resource, 'Control 211 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_212" {
  title       = "Control 212 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_212' as resource, 'Control 212 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_213" {
  title       = "Control 213 (Query Ref)"
  description = "Control referencing query.query_239"
  query       = query.query_239
  severity    = "medium"
  tags        = local.common_tags
}

control "control_214" {
  title       = "Control 214 (Query Ref)"
  description = "Control referencing query.query_242"
  query       = query.query_242
  severity    = "high"
  tags        = local.common_tags
}

control "control_215" {
  title       = "Control 215 (Query Ref)"
  description = "Control referencing query.query_245"
  query       = query.query_245
  severity    = "critical"
  tags        = local.common_tags
}

control "control_216" {
  title       = "Control 216 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_216' as resource, 'Control 216 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_217" {
  title       = "Control 217 (Query Ref)"
  description = "Control referencing query.query_251"
  query       = query.query_251
  severity    = "medium"
  tags        = local.common_tags
}

control "control_218" {
  title       = "Control 218 (Query Ref)"
  description = "Control referencing query.query_254"
  query       = query.query_254
  severity    = "high"
  tags        = local.common_tags
}

control "control_219" {
  title       = "Control 219 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_219' as resource, 'Control 219 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_220" {
  title       = "Control 220 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_220' as resource, 'Control 220 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_221" {
  title       = "Control 221 (Query Ref)"
  description = "Control referencing query.query_263"
  query       = query.query_263
  severity    = "medium"
  tags        = local.common_tags
}

control "control_222" {
  title       = "Control 222 (Query Ref)"
  description = "Control referencing query.query_266"
  query       = query.query_266
  severity    = "high"
  tags        = local.common_tags
}

control "control_223" {
  title       = "Control 223 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_223' as resource, 'Control 223 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_224" {
  title       = "Control 224 (Query Ref)"
  description = "Control referencing query.query_272"
  query       = query.query_272
  severity    = "low"
  tags        = local.common_tags
}

control "control_225" {
  title       = "Control 225 (Query Ref)"
  description = "Control referencing query.query_275"
  query       = query.query_275
  severity    = "medium"
  tags        = local.common_tags
}

control "control_226" {
  title       = "Control 226 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_226' as resource, 'Control 226 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_227" {
  title       = "Control 227 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_227' as resource, 'Control 227 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_228" {
  title       = "Control 228 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_228' as resource, 'Control 228 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_229" {
  title       = "Control 229 (Query Ref)"
  description = "Control referencing query.query_287"
  query       = query.query_287
  severity    = "medium"
  tags        = local.common_tags
}

control "control_230" {
  title       = "Control 230 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_230' as resource, 'Control 230 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_231" {
  title       = "Control 231 (Query Ref)"
  description = "Control referencing query.query_293"
  query       = query.query_293
  severity    = "critical"
  tags        = local.common_tags
}

control "control_232" {
  title       = "Control 232 (Query Ref)"
  description = "Control referencing query.query_296"
  query       = query.query_296
  severity    = "low"
  tags        = local.common_tags
}

control "control_233" {
  title       = "Control 233 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_233' as resource, 'Control 233 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_234" {
  title       = "Control 234 (Query Ref)"
  description = "Control referencing query.query_302"
  query       = query.query_302
  severity    = "high"
  tags        = local.common_tags
}

control "control_235" {
  title       = "Control 235 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_235' as resource, 'Control 235 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_236" {
  title       = "Control 236 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_236' as resource, 'Control 236 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_237" {
  title       = "Control 237 (Query Ref)"
  description = "Control referencing query.query_311"
  query       = query.query_311
  severity    = "medium"
  tags        = local.common_tags
}

control "control_238" {
  title       = "Control 238 (Query Ref)"
  description = "Control referencing query.query_314"
  query       = query.query_314
  severity    = "high"
  tags        = local.common_tags
}

control "control_239" {
  title       = "Control 239 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_239' as resource, 'Control 239 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_240" {
  title       = "Control 240 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_240' as resource, 'Control 240 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_241" {
  title       = "Control 241 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_241' as resource, 'Control 241 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_242" {
  title       = "Control 242 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_242' as resource, 'Control 242 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_243" {
  title       = "Control 243 (Query Ref)"
  description = "Control referencing query.query_329"
  query       = query.query_329
  severity    = "critical"
  tags        = local.common_tags
}

control "control_244" {
  title       = "Control 244 (Query Ref)"
  description = "Control referencing query.query_332"
  query       = query.query_332
  severity    = "low"
  tags        = local.common_tags
}

control "control_245" {
  title       = "Control 245 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_245' as resource, 'Control 245 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_246" {
  title       = "Control 246 (Query Ref)"
  description = "Control referencing query.query_338"
  query       = query.query_338
  severity    = "high"
  tags        = local.common_tags
}

control "control_247" {
  title       = "Control 247 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_247' as resource, 'Control 247 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_248" {
  title       = "Control 248 (Query Ref)"
  description = "Control referencing query.query_344"
  query       = query.query_344
  severity    = "low"
  tags        = local.common_tags
}

control "control_249" {
  title       = "Control 249 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_249' as resource, 'Control 249 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_250" {
  title       = "Control 250 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_250' as resource, 'Control 250 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_251" {
  title       = "Control 251 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_251' as resource, 'Control 251 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_252" {
  title       = "Control 252 (Query Ref)"
  description = "Control referencing query.query_356"
  query       = query.query_356
  severity    = "low"
  tags        = local.common_tags
}

control "control_253" {
  title       = "Control 253 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_253' as resource, 'Control 253 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_254" {
  title       = "Control 254 (Query Ref)"
  description = "Control referencing query.query_362"
  query       = query.query_362
  severity    = "high"
  tags        = local.common_tags
}

control "control_255" {
  title       = "Control 255 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_255' as resource, 'Control 255 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_256" {
  title       = "Control 256 (Query Ref)"
  description = "Control referencing query.query_368"
  query       = query.query_368
  severity    = "low"
  tags        = local.common_tags
}

control "control_257" {
  title       = "Control 257 (Query Ref)"
  description = "Control referencing query.query_371"
  query       = query.query_371
  severity    = "medium"
  tags        = local.common_tags
}

control "control_258" {
  title       = "Control 258 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_258' as resource, 'Control 258 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_259" {
  title       = "Control 259 (Query Ref)"
  description = "Control referencing query.query_377"
  query       = query.query_377
  severity    = "critical"
  tags        = local.common_tags
}

control "control_260" {
  title       = "Control 260 (Query Ref)"
  description = "Control referencing query.query_380"
  query       = query.query_380
  severity    = "low"
  tags        = local.common_tags
}

control "control_261" {
  title       = "Control 261 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_261' as resource, 'Control 261 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_262" {
  title       = "Control 262 (Query Ref)"
  description = "Control referencing query.query_386"
  query       = query.query_386
  severity    = "high"
  tags        = local.common_tags
}

control "control_263" {
  title       = "Control 263 (Query Ref)"
  description = "Control referencing query.query_389"
  query       = query.query_389
  severity    = "critical"
  tags        = local.common_tags
}

control "control_264" {
  title       = "Control 264 (Query Ref)"
  description = "Control referencing query.query_392"
  query       = query.query_392
  severity    = "low"
  tags        = local.common_tags
}

control "control_265" {
  title       = "Control 265 (Query Ref)"
  description = "Control referencing query.query_395"
  query       = query.query_395
  severity    = "medium"
  tags        = local.common_tags
}

control "control_266" {
  title       = "Control 266 (Query Ref)"
  description = "Control referencing query.query_398"
  query       = query.query_398
  severity    = "high"
  tags        = local.common_tags
}

control "control_267" {
  title       = "Control 267 (Query Ref)"
  description = "Control referencing query.query_1"
  query       = query.query_1
  severity    = "critical"
  tags        = local.common_tags
}

control "control_268" {
  title       = "Control 268 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_268' as resource, 'Control 268 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_269" {
  title       = "Control 269 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_269' as resource, 'Control 269 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_270" {
  title       = "Control 270 (Query Ref)"
  description = "Control referencing query.query_10"
  query       = query.query_10
  severity    = "high"
  tags        = local.common_tags
}

control "control_271" {
  title       = "Control 271 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_271' as resource, 'Control 271 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_272" {
  title       = "Control 272 (Query Ref)"
  description = "Control referencing query.query_16"
  query       = query.query_16
  severity    = "low"
  tags        = local.common_tags
}

control "control_273" {
  title       = "Control 273 (Query Ref)"
  description = "Control referencing query.query_19"
  query       = query.query_19
  severity    = "medium"
  tags        = local.common_tags
}

control "control_274" {
  title       = "Control 274 (Query Ref)"
  description = "Control referencing query.query_22"
  query       = query.query_22
  severity    = "high"
  tags        = local.common_tags
}

control "control_275" {
  title       = "Control 275 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_275' as resource, 'Control 275 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_276" {
  title       = "Control 276 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_276' as resource, 'Control 276 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_277" {
  title       = "Control 277 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_277' as resource, 'Control 277 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_278" {
  title       = "Control 278 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_278' as resource, 'Control 278 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_279" {
  title       = "Control 279 (Query Ref)"
  description = "Control referencing query.query_37"
  query       = query.query_37
  severity    = "critical"
  tags        = local.common_tags
}

control "control_280" {
  title       = "Control 280 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_280' as resource, 'Control 280 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_281" {
  title       = "Control 281 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_281' as resource, 'Control 281 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_282" {
  title       = "Control 282 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_282' as resource, 'Control 282 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_283" {
  title       = "Control 283 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_283' as resource, 'Control 283 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_284" {
  title       = "Control 284 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_284' as resource, 'Control 284 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_285" {
  title       = "Control 285 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_285' as resource, 'Control 285 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_286" {
  title       = "Control 286 (Query Ref)"
  description = "Control referencing query.query_58"
  query       = query.query_58
  severity    = "high"
  tags        = local.common_tags
}

control "control_287" {
  title       = "Control 287 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_287' as resource, 'Control 287 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_288" {
  title       = "Control 288 (Query Ref)"
  description = "Control referencing query.query_64"
  query       = query.query_64
  severity    = "low"
  tags        = local.common_tags
}

control "control_289" {
  title       = "Control 289 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_289' as resource, 'Control 289 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_290" {
  title       = "Control 290 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_290' as resource, 'Control 290 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_291" {
  title       = "Control 291 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_291' as resource, 'Control 291 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_292" {
  title       = "Control 292 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_292' as resource, 'Control 292 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_293" {
  title       = "Control 293 (Query Ref)"
  description = "Control referencing query.query_79"
  query       = query.query_79
  severity    = "medium"
  tags        = local.common_tags
}

control "control_294" {
  title       = "Control 294 (Query Ref)"
  description = "Control referencing query.query_82"
  query       = query.query_82
  severity    = "high"
  tags        = local.common_tags
}

control "control_295" {
  title       = "Control 295 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_295' as resource, 'Control 295 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_296" {
  title       = "Control 296 (Query Ref)"
  description = "Control referencing query.query_88"
  query       = query.query_88
  severity    = "low"
  tags        = local.common_tags
}

control "control_297" {
  title       = "Control 297 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_297' as resource, 'Control 297 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_298" {
  title       = "Control 298 (Query Ref)"
  description = "Control referencing query.query_94"
  query       = query.query_94
  severity    = "high"
  tags        = local.common_tags
}

control "control_299" {
  title       = "Control 299 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_299' as resource, 'Control 299 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_300" {
  title       = "Control 300 (Query Ref)"
  description = "Control referencing query.query_100"
  query       = query.query_100
  severity    = "low"
  tags        = local.common_tags
}

control "control_301" {
  title       = "Control 301 (Query Ref)"
  description = "Control referencing query.query_103"
  query       = query.query_103
  severity    = "medium"
  tags        = local.common_tags
}

control "control_302" {
  title       = "Control 302 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_302' as resource, 'Control 302 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_303" {
  title       = "Control 303 (Query Ref)"
  description = "Control referencing query.query_109"
  query       = query.query_109
  severity    = "critical"
  tags        = local.common_tags
}

control "control_304" {
  title       = "Control 304 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_304' as resource, 'Control 304 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_305" {
  title       = "Control 305 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_305' as resource, 'Control 305 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_306" {
  title       = "Control 306 (Query Ref)"
  description = "Control referencing query.query_118"
  query       = query.query_118
  severity    = "high"
  tags        = local.common_tags
}

control "control_307" {
  title       = "Control 307 (Query Ref)"
  description = "Control referencing query.query_121"
  query       = query.query_121
  severity    = "critical"
  tags        = local.common_tags
}

control "control_308" {
  title       = "Control 308 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_308' as resource, 'Control 308 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_309" {
  title       = "Control 309 (Query Ref)"
  description = "Control referencing query.query_127"
  query       = query.query_127
  severity    = "medium"
  tags        = local.common_tags
}

control "control_310" {
  title       = "Control 310 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_310' as resource, 'Control 310 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_311" {
  title       = "Control 311 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_311' as resource, 'Control 311 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_312" {
  title       = "Control 312 (Query Ref)"
  description = "Control referencing query.query_136"
  query       = query.query_136
  severity    = "low"
  tags        = local.common_tags
}

control "control_313" {
  title       = "Control 313 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_313' as resource, 'Control 313 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_314" {
  title       = "Control 314 (Query Ref)"
  description = "Control referencing query.query_142"
  query       = query.query_142
  severity    = "high"
  tags        = local.common_tags
}

control "control_315" {
  title       = "Control 315 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_315' as resource, 'Control 315 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_316" {
  title       = "Control 316 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_316' as resource, 'Control 316 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_317" {
  title       = "Control 317 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_317' as resource, 'Control 317 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_318" {
  title       = "Control 318 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_318' as resource, 'Control 318 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_319" {
  title       = "Control 319 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_319' as resource, 'Control 319 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_320" {
  title       = "Control 320 (Query Ref)"
  description = "Control referencing query.query_160"
  query       = query.query_160
  severity    = "low"
  tags        = local.common_tags
}

control "control_321" {
  title       = "Control 321 (Query Ref)"
  description = "Control referencing query.query_163"
  query       = query.query_163
  severity    = "medium"
  tags        = local.common_tags
}

control "control_322" {
  title       = "Control 322 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_322' as resource, 'Control 322 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_323" {
  title       = "Control 323 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_323' as resource, 'Control 323 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_324" {
  title       = "Control 324 (Query Ref)"
  description = "Control referencing query.query_172"
  query       = query.query_172
  severity    = "low"
  tags        = local.common_tags
}

control "control_325" {
  title       = "Control 325 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_325' as resource, 'Control 325 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_326" {
  title       = "Control 326 (Query Ref)"
  description = "Control referencing query.query_178"
  query       = query.query_178
  severity    = "high"
  tags        = local.common_tags
}

control "control_327" {
  title       = "Control 327 (Query Ref)"
  description = "Control referencing query.query_181"
  query       = query.query_181
  severity    = "critical"
  tags        = local.common_tags
}

control "control_328" {
  title       = "Control 328 (Query Ref)"
  description = "Control referencing query.query_184"
  query       = query.query_184
  severity    = "low"
  tags        = local.common_tags
}

control "control_329" {
  title       = "Control 329 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_329' as resource, 'Control 329 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_330" {
  title       = "Control 330 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_330' as resource, 'Control 330 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_331" {
  title       = "Control 331 (Query Ref)"
  description = "Control referencing query.query_193"
  query       = query.query_193
  severity    = "critical"
  tags        = local.common_tags
}

control "control_332" {
  title       = "Control 332 (Query Ref)"
  description = "Control referencing query.query_196"
  query       = query.query_196
  severity    = "low"
  tags        = local.common_tags
}

control "control_333" {
  title       = "Control 333 (Query Ref)"
  description = "Control referencing query.query_199"
  query       = query.query_199
  severity    = "medium"
  tags        = local.common_tags
}

control "control_334" {
  title       = "Control 334 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_334' as resource, 'Control 334 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_335" {
  title       = "Control 335 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_335' as resource, 'Control 335 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_336" {
  title       = "Control 336 (Query Ref)"
  description = "Control referencing query.query_208"
  query       = query.query_208
  severity    = "low"
  tags        = local.common_tags
}

control "control_337" {
  title       = "Control 337 (Query Ref)"
  description = "Control referencing query.query_211"
  query       = query.query_211
  severity    = "medium"
  tags        = local.common_tags
}

control "control_338" {
  title       = "Control 338 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_338' as resource, 'Control 338 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_339" {
  title       = "Control 339 (Query Ref)"
  description = "Control referencing query.query_217"
  query       = query.query_217
  severity    = "critical"
  tags        = local.common_tags
}

control "control_340" {
  title       = "Control 340 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_340' as resource, 'Control 340 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_341" {
  title       = "Control 341 (Query Ref)"
  description = "Control referencing query.query_223"
  query       = query.query_223
  severity    = "medium"
  tags        = local.common_tags
}

control "control_342" {
  title       = "Control 342 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_342' as resource, 'Control 342 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_343" {
  title       = "Control 343 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_343' as resource, 'Control 343 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_344" {
  title       = "Control 344 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_344' as resource, 'Control 344 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_345" {
  title       = "Control 345 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_345' as resource, 'Control 345 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_346" {
  title       = "Control 346 (Query Ref)"
  description = "Control referencing query.query_238"
  query       = query.query_238
  severity    = "high"
  tags        = local.common_tags
}

control "control_347" {
  title       = "Control 347 (Query Ref)"
  description = "Control referencing query.query_241"
  query       = query.query_241
  severity    = "critical"
  tags        = local.common_tags
}

control "control_348" {
  title       = "Control 348 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_348' as resource, 'Control 348 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_349" {
  title       = "Control 349 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_349' as resource, 'Control 349 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_350" {
  title       = "Control 350 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_350' as resource, 'Control 350 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_351" {
  title       = "Control 351 (Query Ref)"
  description = "Control referencing query.query_253"
  query       = query.query_253
  severity    = "critical"
  tags        = local.common_tags
}

control "control_352" {
  title       = "Control 352 (Query Ref)"
  description = "Control referencing query.query_256"
  query       = query.query_256
  severity    = "low"
  tags        = local.common_tags
}

control "control_353" {
  title       = "Control 353 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_353' as resource, 'Control 353 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_354" {
  title       = "Control 354 (Query Ref)"
  description = "Control referencing query.query_262"
  query       = query.query_262
  severity    = "high"
  tags        = local.common_tags
}

control "control_355" {
  title       = "Control 355 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_355' as resource, 'Control 355 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_356" {
  title       = "Control 356 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_356' as resource, 'Control 356 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_357" {
  title       = "Control 357 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_357' as resource, 'Control 357 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_358" {
  title       = "Control 358 (Query Ref)"
  description = "Control referencing query.query_274"
  query       = query.query_274
  severity    = "high"
  tags        = local.common_tags
}

control "control_359" {
  title       = "Control 359 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_359' as resource, 'Control 359 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_360" {
  title       = "Control 360 (Query Ref)"
  description = "Control referencing query.query_280"
  query       = query.query_280
  severity    = "low"
  tags        = local.common_tags
}

control "control_361" {
  title       = "Control 361 (Query Ref)"
  description = "Control referencing query.query_283"
  query       = query.query_283
  severity    = "medium"
  tags        = local.common_tags
}

control "control_362" {
  title       = "Control 362 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_362' as resource, 'Control 362 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_363" {
  title       = "Control 363 (Query Ref)"
  description = "Control referencing query.query_289"
  query       = query.query_289
  severity    = "critical"
  tags        = local.common_tags
}

control "control_364" {
  title       = "Control 364 (Query Ref)"
  description = "Control referencing query.query_292"
  query       = query.query_292
  severity    = "low"
  tags        = local.common_tags
}

control "control_365" {
  title       = "Control 365 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_365' as resource, 'Control 365 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_366" {
  title       = "Control 366 (Query Ref)"
  description = "Control referencing query.query_298"
  query       = query.query_298
  severity    = "high"
  tags        = local.common_tags
}

control "control_367" {
  title       = "Control 367 (Query Ref)"
  description = "Control referencing query.query_301"
  query       = query.query_301
  severity    = "critical"
  tags        = local.common_tags
}

control "control_368" {
  title       = "Control 368 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_368' as resource, 'Control 368 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_369" {
  title       = "Control 369 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_369' as resource, 'Control 369 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_370" {
  title       = "Control 370 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_370' as resource, 'Control 370 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_371" {
  title       = "Control 371 (Query Ref)"
  description = "Control referencing query.query_313"
  query       = query.query_313
  severity    = "critical"
  tags        = local.common_tags
}

control "control_372" {
  title       = "Control 372 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_372' as resource, 'Control 372 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_373" {
  title       = "Control 373 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_373' as resource, 'Control 373 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_374" {
  title       = "Control 374 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_374' as resource, 'Control 374 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_375" {
  title       = "Control 375 (Query Ref)"
  description = "Control referencing query.query_325"
  query       = query.query_325
  severity    = "critical"
  tags        = local.common_tags
}

control "control_376" {
  title       = "Control 376 (Query Ref)"
  description = "Control referencing query.query_328"
  query       = query.query_328
  severity    = "low"
  tags        = local.common_tags
}

control "control_377" {
  title       = "Control 377 (Query Ref)"
  description = "Control referencing query.query_331"
  query       = query.query_331
  severity    = "medium"
  tags        = local.common_tags
}

control "control_378" {
  title       = "Control 378 (Query Ref)"
  description = "Control referencing query.query_334"
  query       = query.query_334
  severity    = "high"
  tags        = local.common_tags
}

control "control_379" {
  title       = "Control 379 (Query Ref)"
  description = "Control referencing query.query_337"
  query       = query.query_337
  severity    = "critical"
  tags        = local.common_tags
}

control "control_380" {
  title       = "Control 380 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_380' as resource, 'Control 380 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_381" {
  title       = "Control 381 (Query Ref)"
  description = "Control referencing query.query_343"
  query       = query.query_343
  severity    = "medium"
  tags        = local.common_tags
}

control "control_382" {
  title       = "Control 382 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_382' as resource, 'Control 382 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_383" {
  title       = "Control 383 (Query Ref)"
  description = "Control referencing query.query_349"
  query       = query.query_349
  severity    = "critical"
  tags        = local.common_tags
}

control "control_384" {
  title       = "Control 384 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_384' as resource, 'Control 384 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_385" {
  title       = "Control 385 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_385' as resource, 'Control 385 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_386" {
  title       = "Control 386 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_386' as resource, 'Control 386 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_387" {
  title       = "Control 387 (Query Ref)"
  description = "Control referencing query.query_361"
  query       = query.query_361
  severity    = "critical"
  tags        = local.common_tags
}

control "control_388" {
  title       = "Control 388 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_388' as resource, 'Control 388 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_389" {
  title       = "Control 389 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_389' as resource, 'Control 389 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_390" {
  title       = "Control 390 (Query Ref)"
  description = "Control referencing query.query_370"
  query       = query.query_370
  severity    = "high"
  tags        = local.common_tags
}

control "control_391" {
  title       = "Control 391 (Query Ref)"
  description = "Control referencing query.query_373"
  query       = query.query_373
  severity    = "critical"
  tags        = local.common_tags
}

control "control_392" {
  title       = "Control 392 (Query Ref)"
  description = "Control referencing query.query_376"
  query       = query.query_376
  severity    = "low"
  tags        = local.common_tags
}

control "control_393" {
  title       = "Control 393 (Query Ref)"
  description = "Control referencing query.query_379"
  query       = query.query_379
  severity    = "medium"
  tags        = local.common_tags
}

control "control_394" {
  title       = "Control 394 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_394' as resource, 'Control 394 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_395" {
  title       = "Control 395 (Query Ref)"
  description = "Control referencing query.query_385"
  query       = query.query_385
  severity    = "critical"
  tags        = local.common_tags
}

control "control_396" {
  title       = "Control 396 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_396' as resource, 'Control 396 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_397" {
  title       = "Control 397 (Query Ref)"
  description = "Control referencing query.query_391"
  query       = query.query_391
  severity    = "medium"
  tags        = local.common_tags
}

control "control_398" {
  title       = "Control 398 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_398' as resource, 'Control 398 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_399" {
  title       = "Control 399 (Query Ref)"
  description = "Control referencing query.query_397"
  query       = query.query_397
  severity    = "critical"
  tags        = local.common_tags
}

control "control_400" {
  title       = "Control 400 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_400' as resource, 'Control 400 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_401" {
  title       = "Control 401 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_401' as resource, 'Control 401 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_402" {
  title       = "Control 402 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_402' as resource, 'Control 402 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_403" {
  title       = "Control 403 (Query Ref)"
  description = "Control referencing query.query_9"
  query       = query.query_9
  severity    = "critical"
  tags        = local.common_tags
}

control "control_404" {
  title       = "Control 404 (Query Ref)"
  description = "Control referencing query.query_12"
  query       = query.query_12
  severity    = "low"
  tags        = local.common_tags
}

control "control_405" {
  title       = "Control 405 (Query Ref)"
  description = "Control referencing query.query_15"
  query       = query.query_15
  severity    = "medium"
  tags        = local.common_tags
}

control "control_406" {
  title       = "Control 406 (Query Ref)"
  description = "Control referencing query.query_18"
  query       = query.query_18
  severity    = "high"
  tags        = local.common_tags
}

control "control_407" {
  title       = "Control 407 (Query Ref)"
  description = "Control referencing query.query_21"
  query       = query.query_21
  severity    = "critical"
  tags        = local.common_tags
}

control "control_408" {
  title       = "Control 408 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_408' as resource, 'Control 408 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_409" {
  title       = "Control 409 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_409' as resource, 'Control 409 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_410" {
  title       = "Control 410 (Query Ref)"
  description = "Control referencing query.query_30"
  query       = query.query_30
  severity    = "high"
  tags        = local.common_tags
}

control "control_411" {
  title       = "Control 411 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_411' as resource, 'Control 411 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_412" {
  title       = "Control 412 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_412' as resource, 'Control 412 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_413" {
  title       = "Control 413 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_413' as resource, 'Control 413 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_414" {
  title       = "Control 414 (Query Ref)"
  description = "Control referencing query.query_42"
  query       = query.query_42
  severity    = "high"
  tags        = local.common_tags
}

control "control_415" {
  title       = "Control 415 (Query Ref)"
  description = "Control referencing query.query_45"
  query       = query.query_45
  severity    = "critical"
  tags        = local.common_tags
}

control "control_416" {
  title       = "Control 416 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_416' as resource, 'Control 416 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_417" {
  title       = "Control 417 (Query Ref)"
  description = "Control referencing query.query_51"
  query       = query.query_51
  severity    = "medium"
  tags        = local.common_tags
}

control "control_418" {
  title       = "Control 418 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_418' as resource, 'Control 418 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_419" {
  title       = "Control 419 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_419' as resource, 'Control 419 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_420" {
  title       = "Control 420 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_420' as resource, 'Control 420 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_421" {
  title       = "Control 421 (Query Ref)"
  description = "Control referencing query.query_63"
  query       = query.query_63
  severity    = "medium"
  tags        = local.common_tags
}

control "control_422" {
  title       = "Control 422 (Query Ref)"
  description = "Control referencing query.query_66"
  query       = query.query_66
  severity    = "high"
  tags        = local.common_tags
}

control "control_423" {
  title       = "Control 423 (Query Ref)"
  description = "Control referencing query.query_69"
  query       = query.query_69
  severity    = "critical"
  tags        = local.common_tags
}

control "control_424" {
  title       = "Control 424 (Query Ref)"
  description = "Control referencing query.query_72"
  query       = query.query_72
  severity    = "low"
  tags        = local.common_tags
}

control "control_425" {
  title       = "Control 425 (Query Ref)"
  description = "Control referencing query.query_75"
  query       = query.query_75
  severity    = "medium"
  tags        = local.common_tags
}

control "control_426" {
  title       = "Control 426 (Query Ref)"
  description = "Control referencing query.query_78"
  query       = query.query_78
  severity    = "high"
  tags        = local.common_tags
}

control "control_427" {
  title       = "Control 427 (Query Ref)"
  description = "Control referencing query.query_81"
  query       = query.query_81
  severity    = "critical"
  tags        = local.common_tags
}

control "control_428" {
  title       = "Control 428 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_428' as resource, 'Control 428 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_429" {
  title       = "Control 429 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_429' as resource, 'Control 429 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_430" {
  title       = "Control 430 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_430' as resource, 'Control 430 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_431" {
  title       = "Control 431 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_431' as resource, 'Control 431 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_432" {
  title       = "Control 432 (Query Ref)"
  description = "Control referencing query.query_96"
  query       = query.query_96
  severity    = "low"
  tags        = local.common_tags
}

control "control_433" {
  title       = "Control 433 (Query Ref)"
  description = "Control referencing query.query_99"
  query       = query.query_99
  severity    = "medium"
  tags        = local.common_tags
}

control "control_434" {
  title       = "Control 434 (Query Ref)"
  description = "Control referencing query.query_102"
  query       = query.query_102
  severity    = "high"
  tags        = local.common_tags
}

control "control_435" {
  title       = "Control 435 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_435' as resource, 'Control 435 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_436" {
  title       = "Control 436 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_436' as resource, 'Control 436 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_437" {
  title       = "Control 437 (Query Ref)"
  description = "Control referencing query.query_111"
  query       = query.query_111
  severity    = "medium"
  tags        = local.common_tags
}

control "control_438" {
  title       = "Control 438 (Query Ref)"
  description = "Control referencing query.query_114"
  query       = query.query_114
  severity    = "high"
  tags        = local.common_tags
}

control "control_439" {
  title       = "Control 439 (Query Ref)"
  description = "Control referencing query.query_117"
  query       = query.query_117
  severity    = "critical"
  tags        = local.common_tags
}

control "control_440" {
  title       = "Control 440 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_440' as resource, 'Control 440 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_441" {
  title       = "Control 441 (Query Ref)"
  description = "Control referencing query.query_123"
  query       = query.query_123
  severity    = "medium"
  tags        = local.common_tags
}

control "control_442" {
  title       = "Control 442 (Query Ref)"
  description = "Control referencing query.query_126"
  query       = query.query_126
  severity    = "high"
  tags        = local.common_tags
}

control "control_443" {
  title       = "Control 443 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_443' as resource, 'Control 443 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_444" {
  title       = "Control 444 (Query Ref)"
  description = "Control referencing query.query_132"
  query       = query.query_132
  severity    = "low"
  tags        = local.common_tags
}

control "control_445" {
  title       = "Control 445 (Query Ref)"
  description = "Control referencing query.query_135"
  query       = query.query_135
  severity    = "medium"
  tags        = local.common_tags
}

control "control_446" {
  title       = "Control 446 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_446' as resource, 'Control 446 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_447" {
  title       = "Control 447 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_447' as resource, 'Control 447 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_448" {
  title       = "Control 448 (Query Ref)"
  description = "Control referencing query.query_144"
  query       = query.query_144
  severity    = "low"
  tags        = local.common_tags
}

control "control_449" {
  title       = "Control 449 (Query Ref)"
  description = "Control referencing query.query_147"
  query       = query.query_147
  severity    = "medium"
  tags        = local.common_tags
}

control "control_450" {
  title       = "Control 450 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_450' as resource, 'Control 450 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_451" {
  title       = "Control 451 (Query Ref)"
  description = "Control referencing query.query_153"
  query       = query.query_153
  severity    = "critical"
  tags        = local.common_tags
}

control "control_452" {
  title       = "Control 452 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_452' as resource, 'Control 452 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_453" {
  title       = "Control 453 (Query Ref)"
  description = "Control referencing query.query_159"
  query       = query.query_159
  severity    = "medium"
  tags        = local.common_tags
}

control "control_454" {
  title       = "Control 454 (Query Ref)"
  description = "Control referencing query.query_162"
  query       = query.query_162
  severity    = "high"
  tags        = local.common_tags
}

control "control_455" {
  title       = "Control 455 (Query Ref)"
  description = "Control referencing query.query_165"
  query       = query.query_165
  severity    = "critical"
  tags        = local.common_tags
}

control "control_456" {
  title       = "Control 456 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_456' as resource, 'Control 456 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_457" {
  title       = "Control 457 (Query Ref)"
  description = "Control referencing query.query_171"
  query       = query.query_171
  severity    = "medium"
  tags        = local.common_tags
}

control "control_458" {
  title       = "Control 458 (Query Ref)"
  description = "Control referencing query.query_174"
  query       = query.query_174
  severity    = "high"
  tags        = local.common_tags
}

control "control_459" {
  title       = "Control 459 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_459' as resource, 'Control 459 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_460" {
  title       = "Control 460 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_460' as resource, 'Control 460 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_461" {
  title       = "Control 461 (Query Ref)"
  description = "Control referencing query.query_183"
  query       = query.query_183
  severity    = "medium"
  tags        = local.common_tags
}

control "control_462" {
  title       = "Control 462 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_462' as resource, 'Control 462 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_463" {
  title       = "Control 463 (Query Ref)"
  description = "Control referencing query.query_189"
  query       = query.query_189
  severity    = "critical"
  tags        = local.common_tags
}

control "control_464" {
  title       = "Control 464 (Query Ref)"
  description = "Control referencing query.query_192"
  query       = query.query_192
  severity    = "low"
  tags        = local.common_tags
}

control "control_465" {
  title       = "Control 465 (Query Ref)"
  description = "Control referencing query.query_195"
  query       = query.query_195
  severity    = "medium"
  tags        = local.common_tags
}

control "control_466" {
  title       = "Control 466 (Query Ref)"
  description = "Control referencing query.query_198"
  query       = query.query_198
  severity    = "high"
  tags        = local.common_tags
}

control "control_467" {
  title       = "Control 467 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_467' as resource, 'Control 467 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_468" {
  title       = "Control 468 (Query Ref)"
  description = "Control referencing query.query_204"
  query       = query.query_204
  severity    = "low"
  tags        = local.common_tags
}

control "control_469" {
  title       = "Control 469 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_469' as resource, 'Control 469 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_470" {
  title       = "Control 470 (Query Ref)"
  description = "Control referencing query.query_210"
  query       = query.query_210
  severity    = "high"
  tags        = local.common_tags
}

control "control_471" {
  title       = "Control 471 (Query Ref)"
  description = "Control referencing query.query_213"
  query       = query.query_213
  severity    = "critical"
  tags        = local.common_tags
}

control "control_472" {
  title       = "Control 472 (Query Ref)"
  description = "Control referencing query.query_216"
  query       = query.query_216
  severity    = "low"
  tags        = local.common_tags
}

control "control_473" {
  title       = "Control 473 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_473' as resource, 'Control 473 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_474" {
  title       = "Control 474 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_474' as resource, 'Control 474 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_475" {
  title       = "Control 475 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_475' as resource, 'Control 475 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_476" {
  title       = "Control 476 (Query Ref)"
  description = "Control referencing query.query_228"
  query       = query.query_228
  severity    = "low"
  tags        = local.common_tags
}

control "control_477" {
  title       = "Control 477 (Query Ref)"
  description = "Control referencing query.query_231"
  query       = query.query_231
  severity    = "medium"
  tags        = local.common_tags
}

control "control_478" {
  title       = "Control 478 (Query Ref)"
  description = "Control referencing query.query_234"
  query       = query.query_234
  severity    = "high"
  tags        = local.common_tags
}

control "control_479" {
  title       = "Control 479 (Query Ref)"
  description = "Control referencing query.query_237"
  query       = query.query_237
  severity    = "critical"
  tags        = local.common_tags
}

control "control_480" {
  title       = "Control 480 (Query Ref)"
  description = "Control referencing query.query_240"
  query       = query.query_240
  severity    = "low"
  tags        = local.common_tags
}

control "control_481" {
  title       = "Control 481 (Query Ref)"
  description = "Control referencing query.query_243"
  query       = query.query_243
  severity    = "medium"
  tags        = local.common_tags
}

control "control_482" {
  title       = "Control 482 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_482' as resource, 'Control 482 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_483" {
  title       = "Control 483 (Query Ref)"
  description = "Control referencing query.query_249"
  query       = query.query_249
  severity    = "critical"
  tags        = local.common_tags
}

control "control_484" {
  title       = "Control 484 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_484' as resource, 'Control 484 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_485" {
  title       = "Control 485 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_485' as resource, 'Control 485 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_486" {
  title       = "Control 486 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_486' as resource, 'Control 486 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_487" {
  title       = "Control 487 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_487' as resource, 'Control 487 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_488" {
  title       = "Control 488 (Query Ref)"
  description = "Control referencing query.query_264"
  query       = query.query_264
  severity    = "low"
  tags        = local.common_tags
}

control "control_489" {
  title       = "Control 489 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_489' as resource, 'Control 489 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_490" {
  title       = "Control 490 (Query Ref)"
  description = "Control referencing query.query_270"
  query       = query.query_270
  severity    = "high"
  tags        = local.common_tags
}

control "control_491" {
  title       = "Control 491 (Query Ref)"
  description = "Control referencing query.query_273"
  query       = query.query_273
  severity    = "critical"
  tags        = local.common_tags
}

control "control_492" {
  title       = "Control 492 (Query Ref)"
  description = "Control referencing query.query_276"
  query       = query.query_276
  severity    = "low"
  tags        = local.common_tags
}

control "control_493" {
  title       = "Control 493 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_493' as resource, 'Control 493 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_494" {
  title       = "Control 494 (Query Ref)"
  description = "Control referencing query.query_282"
  query       = query.query_282
  severity    = "high"
  tags        = local.common_tags
}

control "control_495" {
  title       = "Control 495 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_495' as resource, 'Control 495 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_496" {
  title       = "Control 496 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_496' as resource, 'Control 496 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_497" {
  title       = "Control 497 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_497' as resource, 'Control 497 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_498" {
  title       = "Control 498 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_498' as resource, 'Control 498 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_499" {
  title       = "Control 499 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_499' as resource, 'Control 499 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_500" {
  title       = "Control 500 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_500' as resource, 'Control 500 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_501" {
  title       = "Control 501 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_501' as resource, 'Control 501 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_502" {
  title       = "Control 502 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_502' as resource, 'Control 502 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_503" {
  title       = "Control 503 (Query Ref)"
  description = "Control referencing query.query_309"
  query       = query.query_309
  severity    = "critical"
  tags        = local.common_tags
}

control "control_504" {
  title       = "Control 504 (Query Ref)"
  description = "Control referencing query.query_312"
  query       = query.query_312
  severity    = "low"
  tags        = local.common_tags
}

control "control_505" {
  title       = "Control 505 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_505' as resource, 'Control 505 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_506" {
  title       = "Control 506 (Query Ref)"
  description = "Control referencing query.query_318"
  query       = query.query_318
  severity    = "high"
  tags        = local.common_tags
}

control "control_507" {
  title       = "Control 507 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_507' as resource, 'Control 507 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_508" {
  title       = "Control 508 (Query Ref)"
  description = "Control referencing query.query_324"
  query       = query.query_324
  severity    = "low"
  tags        = local.common_tags
}

control "control_509" {
  title       = "Control 509 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_509' as resource, 'Control 509 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_510" {
  title       = "Control 510 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_510' as resource, 'Control 510 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_511" {
  title       = "Control 511 (Query Ref)"
  description = "Control referencing query.query_333"
  query       = query.query_333
  severity    = "critical"
  tags        = local.common_tags
}

control "control_512" {
  title       = "Control 512 (Query Ref)"
  description = "Control referencing query.query_336"
  query       = query.query_336
  severity    = "low"
  tags        = local.common_tags
}

control "control_513" {
  title       = "Control 513 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_513' as resource, 'Control 513 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_514" {
  title       = "Control 514 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_514' as resource, 'Control 514 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_515" {
  title       = "Control 515 (Query Ref)"
  description = "Control referencing query.query_345"
  query       = query.query_345
  severity    = "critical"
  tags        = local.common_tags
}

control "control_516" {
  title       = "Control 516 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_516' as resource, 'Control 516 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_517" {
  title       = "Control 517 (Query Ref)"
  description = "Control referencing query.query_351"
  query       = query.query_351
  severity    = "medium"
  tags        = local.common_tags
}

control "control_518" {
  title       = "Control 518 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_518' as resource, 'Control 518 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_519" {
  title       = "Control 519 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_519' as resource, 'Control 519 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_520" {
  title       = "Control 520 (Query Ref)"
  description = "Control referencing query.query_360"
  query       = query.query_360
  severity    = "low"
  tags        = local.common_tags
}

control "control_521" {
  title       = "Control 521 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_521' as resource, 'Control 521 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_522" {
  title       = "Control 522 (Query Ref)"
  description = "Control referencing query.query_366"
  query       = query.query_366
  severity    = "high"
  tags        = local.common_tags
}

control "control_523" {
  title       = "Control 523 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_523' as resource, 'Control 523 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_524" {
  title       = "Control 524 (Query Ref)"
  description = "Control referencing query.query_372"
  query       = query.query_372
  severity    = "low"
  tags        = local.common_tags
}

control "control_525" {
  title       = "Control 525 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_525' as resource, 'Control 525 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_526" {
  title       = "Control 526 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_526' as resource, 'Control 526 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_527" {
  title       = "Control 527 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_527' as resource, 'Control 527 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_528" {
  title       = "Control 528 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_528' as resource, 'Control 528 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_529" {
  title       = "Control 529 (Query Ref)"
  description = "Control referencing query.query_387"
  query       = query.query_387
  severity    = "medium"
  tags        = local.common_tags
}

control "control_530" {
  title       = "Control 530 (Query Ref)"
  description = "Control referencing query.query_390"
  query       = query.query_390
  severity    = "high"
  tags        = local.common_tags
}

control "control_531" {
  title       = "Control 531 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_531' as resource, 'Control 531 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_532" {
  title       = "Control 532 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_532' as resource, 'Control 532 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_533" {
  title       = "Control 533 (Query Ref)"
  description = "Control referencing query.query_399"
  query       = query.query_399
  severity    = "medium"
  tags        = local.common_tags
}

control "control_534" {
  title       = "Control 534 (Query Ref)"
  description = "Control referencing query.query_2"
  query       = query.query_2
  severity    = "high"
  tags        = local.common_tags
}

control "control_535" {
  title       = "Control 535 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_535' as resource, 'Control 535 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_536" {
  title       = "Control 536 (Query Ref)"
  description = "Control referencing query.query_8"
  query       = query.query_8
  severity    = "low"
  tags        = local.common_tags
}

control "control_537" {
  title       = "Control 537 (Query Ref)"
  description = "Control referencing query.query_11"
  query       = query.query_11
  severity    = "medium"
  tags        = local.common_tags
}

control "control_538" {
  title       = "Control 538 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_538' as resource, 'Control 538 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_539" {
  title       = "Control 539 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_539' as resource, 'Control 539 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_540" {
  title       = "Control 540 (Query Ref)"
  description = "Control referencing query.query_20"
  query       = query.query_20
  severity    = "low"
  tags        = local.common_tags
}

control "control_541" {
  title       = "Control 541 (Query Ref)"
  description = "Control referencing query.query_23"
  query       = query.query_23
  severity    = "medium"
  tags        = local.common_tags
}

control "control_542" {
  title       = "Control 542 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_542' as resource, 'Control 542 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_543" {
  title       = "Control 543 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_543' as resource, 'Control 543 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_544" {
  title       = "Control 544 (Query Ref)"
  description = "Control referencing query.query_32"
  query       = query.query_32
  severity    = "low"
  tags        = local.common_tags
}

control "control_545" {
  title       = "Control 545 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_545' as resource, 'Control 545 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_546" {
  title       = "Control 546 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_546' as resource, 'Control 546 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_547" {
  title       = "Control 547 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_547' as resource, 'Control 547 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_548" {
  title       = "Control 548 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_548' as resource, 'Control 548 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_549" {
  title       = "Control 549 (Query Ref)"
  description = "Control referencing query.query_47"
  query       = query.query_47
  severity    = "medium"
  tags        = local.common_tags
}

control "control_550" {
  title       = "Control 550 (Query Ref)"
  description = "Control referencing query.query_50"
  query       = query.query_50
  severity    = "high"
  tags        = local.common_tags
}

control "control_551" {
  title       = "Control 551 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_551' as resource, 'Control 551 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_552" {
  title       = "Control 552 (Query Ref)"
  description = "Control referencing query.query_56"
  query       = query.query_56
  severity    = "low"
  tags        = local.common_tags
}

control "control_553" {
  title       = "Control 553 (Query Ref)"
  description = "Control referencing query.query_59"
  query       = query.query_59
  severity    = "medium"
  tags        = local.common_tags
}

control "control_554" {
  title       = "Control 554 (Query Ref)"
  description = "Control referencing query.query_62"
  query       = query.query_62
  severity    = "high"
  tags        = local.common_tags
}

control "control_555" {
  title       = "Control 555 (Query Ref)"
  description = "Control referencing query.query_65"
  query       = query.query_65
  severity    = "critical"
  tags        = local.common_tags
}

control "control_556" {
  title       = "Control 556 (Query Ref)"
  description = "Control referencing query.query_68"
  query       = query.query_68
  severity    = "low"
  tags        = local.common_tags
}

control "control_557" {
  title       = "Control 557 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_557' as resource, 'Control 557 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_558" {
  title       = "Control 558 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_558' as resource, 'Control 558 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_559" {
  title       = "Control 559 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_559' as resource, 'Control 559 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_560" {
  title       = "Control 560 (Query Ref)"
  description = "Control referencing query.query_80"
  query       = query.query_80
  severity    = "low"
  tags        = local.common_tags
}

control "control_561" {
  title       = "Control 561 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_561' as resource, 'Control 561 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_562" {
  title       = "Control 562 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_562' as resource, 'Control 562 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_563" {
  title       = "Control 563 (Query Ref)"
  description = "Control referencing query.query_89"
  query       = query.query_89
  severity    = "critical"
  tags        = local.common_tags
}

control "control_564" {
  title       = "Control 564 (Query Ref)"
  description = "Control referencing query.query_92"
  query       = query.query_92
  severity    = "low"
  tags        = local.common_tags
}

control "control_565" {
  title       = "Control 565 (Query Ref)"
  description = "Control referencing query.query_95"
  query       = query.query_95
  severity    = "medium"
  tags        = local.common_tags
}

control "control_566" {
  title       = "Control 566 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_566' as resource, 'Control 566 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_567" {
  title       = "Control 567 (Query Ref)"
  description = "Control referencing query.query_101"
  query       = query.query_101
  severity    = "critical"
  tags        = local.common_tags
}

control "control_568" {
  title       = "Control 568 (Query Ref)"
  description = "Control referencing query.query_104"
  query       = query.query_104
  severity    = "low"
  tags        = local.common_tags
}

control "control_569" {
  title       = "Control 569 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_569' as resource, 'Control 569 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_570" {
  title       = "Control 570 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_570' as resource, 'Control 570 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_571" {
  title       = "Control 571 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_571' as resource, 'Control 571 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_572" {
  title       = "Control 572 (Query Ref)"
  description = "Control referencing query.query_116"
  query       = query.query_116
  severity    = "low"
  tags        = local.common_tags
}

control "control_573" {
  title       = "Control 573 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_573' as resource, 'Control 573 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_574" {
  title       = "Control 574 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_574' as resource, 'Control 574 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_575" {
  title       = "Control 575 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_575' as resource, 'Control 575 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_576" {
  title       = "Control 576 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_576' as resource, 'Control 576 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_577" {
  title       = "Control 577 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_577' as resource, 'Control 577 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_578" {
  title       = "Control 578 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_578' as resource, 'Control 578 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_579" {
  title       = "Control 579 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_579' as resource, 'Control 579 passed' as reason"
  severity    = "critical"
  tags        = local.common_tags
}

control "control_580" {
  title       = "Control 580 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_580' as resource, 'Control 580 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_581" {
  title       = "Control 581 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_581' as resource, 'Control 581 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_582" {
  title       = "Control 582 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_582' as resource, 'Control 582 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_583" {
  title       = "Control 583 (Query Ref)"
  description = "Control referencing query.query_149"
  query       = query.query_149
  severity    = "critical"
  tags        = local.common_tags
}

control "control_584" {
  title       = "Control 584 (Query Ref)"
  description = "Control referencing query.query_152"
  query       = query.query_152
  severity    = "low"
  tags        = local.common_tags
}

control "control_585" {
  title       = "Control 585 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_585' as resource, 'Control 585 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_586" {
  title       = "Control 586 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_586' as resource, 'Control 586 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_587" {
  title       = "Control 587 (Query Ref)"
  description = "Control referencing query.query_161"
  query       = query.query_161
  severity    = "critical"
  tags        = local.common_tags
}

control "control_588" {
  title       = "Control 588 (Query Ref)"
  description = "Control referencing query.query_164"
  query       = query.query_164
  severity    = "low"
  tags        = local.common_tags
}

control "control_589" {
  title       = "Control 589 (Query Ref)"
  description = "Control referencing query.query_167"
  query       = query.query_167
  severity    = "medium"
  tags        = local.common_tags
}

control "control_590" {
  title       = "Control 590 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_590' as resource, 'Control 590 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_591" {
  title       = "Control 591 (Query Ref)"
  description = "Control referencing query.query_173"
  query       = query.query_173
  severity    = "critical"
  tags        = local.common_tags
}

control "control_592" {
  title       = "Control 592 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_592' as resource, 'Control 592 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_593" {
  title       = "Control 593 (Query Ref)"
  description = "Control referencing query.query_179"
  query       = query.query_179
  severity    = "medium"
  tags        = local.common_tags
}

control "control_594" {
  title       = "Control 594 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_594' as resource, 'Control 594 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_595" {
  title       = "Control 595 (Query Ref)"
  description = "Control referencing query.query_185"
  query       = query.query_185
  severity    = "critical"
  tags        = local.common_tags
}

control "control_596" {
  title       = "Control 596 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_596' as resource, 'Control 596 passed' as reason"
  severity    = "low"
  tags        = local.common_tags
}

control "control_597" {
  title       = "Control 597 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_597' as resource, 'Control 597 passed' as reason"
  severity    = "medium"
  tags        = local.common_tags
}

control "control_598" {
  title       = "Control 598 (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_598' as resource, 'Control 598 passed' as reason"
  severity    = "high"
  tags        = local.common_tags
}

control "control_599" {
  title       = "Control 599 (Query Ref)"
  description = "Control referencing query.query_197"
  query       = query.query_197
  severity    = "critical"
  tags        = local.common_tags
}

