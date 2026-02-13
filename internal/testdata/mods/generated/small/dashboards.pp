
dashboard "dashboard_0" {
  title = "Dashboard 0"
  description = "Test dashboard 0 for performance benchmarking"

  tags = {
    category = "test"
    index    = "0"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_0.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_0.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_0.sql
    }
  }
}

dashboard "dashboard_1" {
  title = "Dashboard 1"
  description = "Test dashboard 1 for performance benchmarking"

  tags = {
    category = "test"
    index    = "1"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_1.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_1.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_1.sql
    }
  }
}

dashboard "dashboard_2" {
  title = "Dashboard 2"
  description = "Test dashboard 2 for performance benchmarking"

  tags = {
    category = "test"
    index    = "2"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_2.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_2.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_2.sql
    }
  }
}

dashboard "dashboard_3" {
  title = "Dashboard 3"
  description = "Test dashboard 3 for performance benchmarking"

  tags = {
    category = "test"
    index    = "3"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_3.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_3.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_3.sql
    }
  }
}

dashboard "dashboard_4" {
  title = "Dashboard 4"
  description = "Test dashboard 4 for performance benchmarking"

  tags = {
    category = "test"
    index    = "4"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_4.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_4.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_4.sql
    }
  }
}

dashboard "dashboard_5" {
  title = "Dashboard 5"
  description = "Test dashboard 5 for performance benchmarking"

  tags = {
    category = "test"
    index    = "5"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_5.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_5.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_5.sql
    }
  }
}

dashboard "dashboard_6" {
  title = "Dashboard 6"
  description = "Test dashboard 6 for performance benchmarking"

  tags = {
    category = "test"
    index    = "6"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_6.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_6.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_6.sql
    }
  }
}

dashboard "dashboard_7" {
  title = "Dashboard 7"
  description = "Test dashboard 7 for performance benchmarking"

  tags = {
    category = "test"
    index    = "7"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_7.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_7.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_7.sql
    }
  }
}

dashboard "dashboard_8" {
  title = "Dashboard 8"
  description = "Test dashboard 8 for performance benchmarking"

  tags = {
    category = "test"
    index    = "8"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_8.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_8.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_8.sql
    }
  }
}

dashboard "dashboard_9" {
  title = "Dashboard 9"
  description = "Test dashboard 9 for performance benchmarking"

  tags = {
    category = "test"
    index    = "9"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_9.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_9.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_9.sql
    }
  }
}
