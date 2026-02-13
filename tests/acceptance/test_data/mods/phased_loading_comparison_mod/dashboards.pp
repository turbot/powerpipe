dashboard "dashboard_with_tags" {
  title       = "Dashboard With Tags"
  description = "A dashboard with various tags for comparison testing"

  tags = {
    service  = "test_service"
    category = "comparison"
    type     = "acceptance_test"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = "SELECT 1 as value"
    }
  }
}

dashboard "dashboard_without_tags" {
  title       = "Dashboard Without Tags"
  description = "A dashboard without tags for comparison testing"

  container {
    title = "Details"

    card {
      width = 2
      sql   = "SELECT 2 as value"
    }
  }
}

dashboard "dashboard_with_empty_tags" {
  title       = "Dashboard With Empty Tags"
  description = "A dashboard with empty tags map"

  tags = {}

  container {
    title = "Info"

    card {
      width = 2
      sql   = "SELECT 3 as value"
    }
  }
}
