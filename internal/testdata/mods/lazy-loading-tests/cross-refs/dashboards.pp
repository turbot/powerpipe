# Dashboards with nested containers and cross-references

dashboard "main" {
  title       = "Cross Reference Dashboard"
  description = "Dashboard demonstrating cross-reference patterns"
  tags        = local.common_tags

  input "region_select" {
    title = "Region"
    type  = "select"
    width = 4
    sql   = "SELECT 'us-east-1' as label, 'us-east-1' as value UNION ALL SELECT 'us-west-2', 'us-west-2'"
  }

  # Top level container
  container {
    title = "Overview"

    # Card referencing shared query
    card {
      title = "Shared Data"
      width = 4
      sql   = query.shared.sql
    }

    # Card referencing base query
    card {
      title = "Base Data"
      width = 4
      sql   = query.base.sql
    }

    # Inline card
    card {
      title = "Static"
      width = 4
      sql   = "SELECT 'static' as value"
    }
  }

  # Nested containers
  container {
    title = "Nested Section"

    container {
      title = "Level 2"

      table {
        title = "Data Table"
        width = 6
        query = query.shared
      }

      container {
        title = "Level 3"

        chart {
          title = "Chart from Query"
          type  = "bar"
          width = 6
          sql   = query.base.sql
        }
      }
    }

    container {
      title = "Another Level 2"

      card {
        sql = query.uses_local.sql
      }
    }
  }
}

# Dashboard referencing input from another dashboard is not supported,
# but dashboards can share queries
dashboard "secondary" {
  title       = "Secondary Dashboard"
  description = "Another dashboard using the same queries"

  card {
    title = "Shared Query Reuse"
    width = 6
    sql   = query.shared.sql
  }

  card {
    title = "Base Query Reuse"
    width = 6
    sql   = query.base.sql
  }

  table {
    title = "Parameterized Table"
    query = query.parameterized
  }
}
