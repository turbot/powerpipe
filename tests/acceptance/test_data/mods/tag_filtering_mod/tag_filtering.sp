benchmark "tag_filtering_benchmark" {
  title       = "Tag filtering benchmark"
  description = "Contains controls with a mix of deprecated tags to validate tag filtering."

  # All controls participate so tag filters can slice them.
  children = [
    control.deprecated_true,
    control.deprecated_false,
    control.no_deprecated_tag,
    control.other_tag_only,
  ]
}

control "deprecated_true" {
  title       = "Control with deprecated=true"
  description = "Should be included when filtering deprecated=true, excluded otherwise."
  query       = query.always_ok

  tags = {
    deprecated = "true"
  }
}

control "deprecated_false" {
  title       = "Control with deprecated=false"
  description = "Should be included for deprecated!=true."
  query       = query.always_ok

  tags = {
    deprecated = "false"
  }
}

control "no_deprecated_tag" {
  title       = "Control without deprecated tag"
  description = "Lacks deprecated tag; should be included when filtering deprecated!=true."
  query       = query.always_ok
}

control "other_tag_only" {
  title       = "Control with other tag only"
  description = "Has env tag but no deprecated tag; should still be included in deprecated!=true filters."
  query       = query.always_ok

  tags = {
    env = "qa"
  }
}

query "always_ok" {
  title       = "Always OK query"
  description = "Returns a single OK status row."
  sql         = "select 'ok' as status, 'dummy' as resource, 'acceptance' as reason"
}
