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
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_3.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_4.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_5.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_2.sql
        }
      }
    }
  }
}

dashboard "dashboard_3" {
  title       = "Dashboard 3"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_4.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_3.sql
    }
  }
}

dashboard "dashboard_4" {
  title       = "Dashboard 4"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_5.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_6.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_4.sql
      }
    }
  }
}

dashboard "dashboard_5" {
  title       = "Dashboard 5"
  description = "Dashboard with 3 levels of container nesting"
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

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_8.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_5.sql
        }
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
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_9.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_10.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_11.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_8.sql
        }
      }
    }
  }
}

dashboard "dashboard_9" {
  title       = "Dashboard 9"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_10.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_9.sql
    }
  }
}

dashboard "dashboard_10" {
  title       = "Dashboard 10"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_11.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_12.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_10.sql
      }
    }
  }
}

dashboard "dashboard_11" {
  title       = "Dashboard 11"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_12.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_13.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_14.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_11.sql
        }
      }
    }
  }
}

dashboard "dashboard_12" {
  title       = "Dashboard 12"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_13.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_12.sql
    }
  }
}

dashboard "dashboard_13" {
  title       = "Dashboard 13"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_14.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_15.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_13.sql
      }
    }
  }
}

dashboard "dashboard_14" {
  title       = "Dashboard 14"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_15.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_16.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_17.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_14.sql
        }
      }
    }
  }
}

dashboard "dashboard_15" {
  title       = "Dashboard 15"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_16.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_15.sql
    }
  }
}

dashboard "dashboard_16" {
  title       = "Dashboard 16"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_17.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_18.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_16.sql
      }
    }
  }
}

dashboard "dashboard_17" {
  title       = "Dashboard 17"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_18.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_19.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_20.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_17.sql
        }
      }
    }
  }
}

dashboard "dashboard_18" {
  title       = "Dashboard 18"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_19.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_18.sql
    }
  }
}

dashboard "dashboard_19" {
  title       = "Dashboard 19"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_20.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_21.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_19.sql
      }
    }
  }
}

dashboard "dashboard_20" {
  title       = "Dashboard 20"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_21.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_22.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_23.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_20.sql
        }
      }
    }
  }
}

dashboard "dashboard_21" {
  title       = "Dashboard 21"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_22.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_21.sql
    }
  }
}

dashboard "dashboard_22" {
  title       = "Dashboard 22"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_23.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_24.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_22.sql
      }
    }
  }
}

dashboard "dashboard_23" {
  title       = "Dashboard 23"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_24.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_25.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_26.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_23.sql
        }
      }
    }
  }
}

dashboard "dashboard_24" {
  title       = "Dashboard 24"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_25.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_24.sql
    }
  }
}

dashboard "dashboard_25" {
  title       = "Dashboard 25"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_26.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_27.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_25.sql
      }
    }
  }
}

dashboard "dashboard_26" {
  title       = "Dashboard 26"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_27.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_28.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_29.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_26.sql
        }
      }
    }
  }
}

dashboard "dashboard_27" {
  title       = "Dashboard 27"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_28.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_27.sql
    }
  }
}

dashboard "dashboard_28" {
  title       = "Dashboard 28"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_29.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_30.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_28.sql
      }
    }
  }
}

dashboard "dashboard_29" {
  title       = "Dashboard 29"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_30.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_31.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_32.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_29.sql
        }
      }
    }
  }
}

dashboard "dashboard_30" {
  title       = "Dashboard 30"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_31.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_30.sql
    }
  }
}

dashboard "dashboard_31" {
  title       = "Dashboard 31"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_32.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_33.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_31.sql
      }
    }
  }
}

dashboard "dashboard_32" {
  title       = "Dashboard 32"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_33.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_34.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_35.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_32.sql
        }
      }
    }
  }
}

dashboard "dashboard_33" {
  title       = "Dashboard 33"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_34.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_33.sql
    }
  }
}

dashboard "dashboard_34" {
  title       = "Dashboard 34"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_35.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_36.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_34.sql
      }
    }
  }
}

dashboard "dashboard_35" {
  title       = "Dashboard 35"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_36.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_37.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_38.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_35.sql
        }
      }
    }
  }
}

dashboard "dashboard_36" {
  title       = "Dashboard 36"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_37.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_36.sql
    }
  }
}

dashboard "dashboard_37" {
  title       = "Dashboard 37"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_38.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_39.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_37.sql
      }
    }
  }
}

dashboard "dashboard_38" {
  title       = "Dashboard 38"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_39.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_40.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_41.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_38.sql
        }
      }
    }
  }
}

dashboard "dashboard_39" {
  title       = "Dashboard 39"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_40.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_39.sql
    }
  }
}

dashboard "dashboard_40" {
  title       = "Dashboard 40"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_41.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_42.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_40.sql
      }
    }
  }
}

dashboard "dashboard_41" {
  title       = "Dashboard 41"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_42.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_43.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_44.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_41.sql
        }
      }
    }
  }
}

dashboard "dashboard_42" {
  title       = "Dashboard 42"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_43.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_42.sql
    }
  }
}

dashboard "dashboard_43" {
  title       = "Dashboard 43"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_44.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_45.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_43.sql
      }
    }
  }
}

dashboard "dashboard_44" {
  title       = "Dashboard 44"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_45.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_46.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_47.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_44.sql
        }
      }
    }
  }
}

dashboard "dashboard_45" {
  title       = "Dashboard 45"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_46.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_45.sql
    }
  }
}

dashboard "dashboard_46" {
  title       = "Dashboard 46"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_47.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_48.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_46.sql
      }
    }
  }
}

dashboard "dashboard_47" {
  title       = "Dashboard 47"
  description = "Dashboard with 3 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_48.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_49.sql
      }

      container {
        title = "Level 3 Container"

        card {
          title = "Level 3 Card"
          width = 4
          sql   = query.query_50.sql
        }

        table {
          title = "Data Table"
          sql   = query.query_47.sql
        }
      }
    }
  }
}

dashboard "dashboard_48" {
  title       = "Dashboard 48"
  description = "Dashboard with 1 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_49.sql
    }

    table {
      title = "Data Table"
      sql   = query.query_48.sql
    }
  }
}

dashboard "dashboard_49" {
  title       = "Dashboard 49"
  description = "Dashboard with 2 levels of container nesting"
  tags        = local.common_tags

  container {
    title = "Level 1 Container"

    card {
      title = "Level 1 Card"
      width = 4
      sql   = query.query_50.sql
    }

    container {
      title = "Level 2 Container"

      card {
        title = "Level 2 Card"
        width = 4
        sql   = query.query_51.sql
      }

      table {
        title = "Data Table"
        sql   = query.query_49.sql
      }
    }
  }
}

