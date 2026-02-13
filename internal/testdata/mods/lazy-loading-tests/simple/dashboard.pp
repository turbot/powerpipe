# Simple dashboard for baseline lazy loading tests

dashboard "simple" {
  title       = "Simple Dashboard"
  description = "Basic dashboard with cards for testing lazy loading"
  tags        = local.common_tags

  card {
    title = "Count"
    width = 4
    sql   = query.simple_count.sql
  }

  card {
    title = "Status"
    width = 4
    sql   = query.simple_status.sql
  }

  card {
    title = "Inline"
    width = 4
    sql   = "SELECT 'ok' as value"
  }
}
