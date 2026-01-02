# Generated dashboards for lazy loading testing

dashboard "dashboard_0" {
  title       = "Dashboard 0"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_1.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_0.sql
    }
  }
}

dashboard "dashboard_1" {
  title       = "Dashboard 1"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_2.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_3.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_1.sql
      }
    }
  }
}

dashboard "dashboard_2" {
  title       = "Dashboard 2"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_3.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_2.sql
    }
  }
}

dashboard "dashboard_3" {
  title       = "Dashboard 3"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_4.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_5.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_3.sql
      }
    }
  }
}

dashboard "dashboard_4" {
  title       = "Dashboard 4"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_5.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_4.sql
    }
  }
}

dashboard "dashboard_5" {
  title       = "Dashboard 5"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_6.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_7.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_5.sql
      }
    }
  }
}

dashboard "dashboard_6" {
  title       = "Dashboard 6"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_7.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_6.sql
    }
  }
}

dashboard "dashboard_7" {
  title       = "Dashboard 7"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_8.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_9.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_7.sql
      }
    }
  }
}

dashboard "dashboard_8" {
  title       = "Dashboard 8"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_9.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_8.sql
    }
  }
}

dashboard "dashboard_9" {
  title       = "Dashboard 9"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_10.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_11.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_9.sql
      }
    }
  }
}

