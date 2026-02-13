locals {
  common_tags = {
    service     = var.service_name
    environment = var.environment
  }
}

dashboard "dashboard_with_variable_tags" {
  title       = "Dashboard With Variable Tags"
  description = "A dashboard using variable tags"

  tags = local.common_tags

  container {
    title = "Overview"

    card {
      width = 2
      sql   = "SELECT 1 as value"
    }
  }
}

dashboard "dashboard_with_merge_tags" {
  title       = "Dashboard With Merge Tags"
  description = "A dashboard using merge() for tags"

  tags = merge(local.common_tags, {
    additional = "extra_tag"
  })

  container {
    title = "Details"

    card {
      width = 2
      sql   = "SELECT 2 as value"
    }
  }
}

dashboard "dashboard_with_literal_tags" {
  title       = "Dashboard With Literal Tags"
  description = "A dashboard using literal tags"

  tags = {
    type     = "literal"
    category = "test"
  }

  container {
    title = "Info"

    card {
      width = 2
      sql   = "SELECT 3 as value"
    }
  }
}
