# Dashboard with all panel types

dashboard "main" {
  title       = "Main Dashboard"
  description = "Primary test dashboard with all panel types"
  tags        = local.common_tags

  input "filter_selection" {
    title = "Filter Selection"
    type  = "select"
    width = 4
    sql   = "SELECT 'option1' as label, 'opt1' as value UNION ALL SELECT 'option2', 'opt2'"
  }

  container {
    title = "Overview Section"

    card {
      title = "Total Count"
      width = 3
      sql   = query.simple.sql
    }

    card {
      title = "Status Card"
      width = 3
      type  = "info"
      sql   = "SELECT 'active' as value"
    }

    chart {
      title = "Bar Chart"
      width = 6
      type  = "bar"
      sql   = "SELECT 'a' as category, 10 as count UNION ALL SELECT 'b', 20 UNION ALL SELECT 'c', 15"
    }
  }

  container {
    title = "Details Section"

    table {
      title = "Data Table"
      width = 6
      query = query.for_table
    }

    chart {
      title = "Line Chart"
      width = 6
      type  = "line"
      sql   = "SELECT 1 as x, 10 as y UNION ALL SELECT 2, 15 UNION ALL SELECT 3, 12"
    }
  }

  text {
    value = "This is a **markdown** text block for testing"
    width = 12
  }

  image {
    title = "Test Image"
    src   = "https://example.com/image.png"
    alt   = "Test image"
    width = 4
  }
}

# Dashboard with nested containers
dashboard "nested" {
  title       = "Nested Dashboard"
  description = "Dashboard with deeply nested containers"

  container {
    title = "Level 1"

    container {
      title = "Level 2A"

      card {
        sql = "SELECT 1 as value"
      }

      container {
        title = "Level 3"

        card {
          sql = "SELECT 2 as value"
        }
      }
    }

    container {
      title = "Level 2B"

      card {
        sql = "SELECT 3 as value"
      }
    }
  }
}

# Dashboard with graph visualization
dashboard "with_graph" {
  title = "Graph Dashboard"

  graph {
    title     = "Resource Graph"
    direction = "TB"

    node {
      sql = "SELECT 'node1' as id, 'Node 1' as title"
    }

    node {
      sql = "SELECT 'node2' as id, 'Node 2' as title"
    }

    edge {
      sql = "SELECT 'node1' as from_id, 'node2' as to_id"
    }
  }
}

# Dashboard with flow visualization
dashboard "with_flow" {
  title = "Flow Dashboard"

  flow {
    title = "Process Flow"

    node {
      sql = "SELECT 'step1' as id, 'Step 1' as title"
    }

    node {
      sql = "SELECT 'step2' as id, 'Step 2' as title"
    }

    node {
      sql = "SELECT 'step3' as id, 'Step 3' as title"
    }

    edge {
      sql = "SELECT 'step1' as from_id, 'step2' as to_id UNION ALL SELECT 'step2', 'step3'"
    }
  }
}

# Dashboard with hierarchy visualization
dashboard "with_hierarchy" {
  title = "Hierarchy Dashboard"

  hierarchy {
    title = "Org Hierarchy"

    node {
      sql = "SELECT 'root' as id, 'Root' as title"
    }

    node {
      sql = "SELECT 'child1' as id, 'Child 1' as title, 'root' as parent_id"
    }

    node {
      sql = "SELECT 'child2' as id, 'Child 2' as title, 'root' as parent_id"
    }
  }
}

# Dashboard with categories
dashboard "with_categories" {
  title = "Category Dashboard"

  graph {
    title = "Categorized Graph"

    category "active" {
      title = "Active Resources"
      color = "green"
      icon  = "check"
    }

    category "inactive" {
      title = "Inactive Resources"
      color = "red"
      icon  = "x"
    }

    node {
      sql = "SELECT 'node1' as id, 'Active Node' as title, 'active' as category"
    }

    node {
      sql = "SELECT 'node2' as id, 'Inactive Node' as title, 'inactive' as category"
    }

    edge {
      sql = "SELECT 'node1' as from_id, 'node2' as to_id"
    }
  }
}

# Dashboard with global input reference
dashboard "with_inputs" {
  title = "Input Dashboard"

  input "local_input" {
    title = "Local Input"
    type  = "text"
    width = 6
  }

  card {
    sql = "SELECT self.input.local_input.value as value"
  }
}

# Simple dashboard for basic tests
dashboard "simple" {
  title = "Simple Dashboard"

  card {
    sql = "SELECT 42 as value"
  }
}
