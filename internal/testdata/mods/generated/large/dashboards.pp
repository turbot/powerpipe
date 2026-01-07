
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

dashboard "dashboard_10" {
  title = "Dashboard 10"
  description = "Test dashboard 10 for performance benchmarking"

  tags = {
    category = "test"
    index    = "10"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_10.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_10.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_10.sql
    }
  }
}

dashboard "dashboard_11" {
  title = "Dashboard 11"
  description = "Test dashboard 11 for performance benchmarking"

  tags = {
    category = "test"
    index    = "11"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_11.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_11.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_11.sql
    }
  }
}

dashboard "dashboard_12" {
  title = "Dashboard 12"
  description = "Test dashboard 12 for performance benchmarking"

  tags = {
    category = "test"
    index    = "12"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_12.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_12.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_12.sql
    }
  }
}

dashboard "dashboard_13" {
  title = "Dashboard 13"
  description = "Test dashboard 13 for performance benchmarking"

  tags = {
    category = "test"
    index    = "13"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_13.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_13.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_13.sql
    }
  }
}

dashboard "dashboard_14" {
  title = "Dashboard 14"
  description = "Test dashboard 14 for performance benchmarking"

  tags = {
    category = "test"
    index    = "14"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_14.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_14.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_14.sql
    }
  }
}

dashboard "dashboard_15" {
  title = "Dashboard 15"
  description = "Test dashboard 15 for performance benchmarking"

  tags = {
    category = "test"
    index    = "15"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_15.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_15.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_15.sql
    }
  }
}

dashboard "dashboard_16" {
  title = "Dashboard 16"
  description = "Test dashboard 16 for performance benchmarking"

  tags = {
    category = "test"
    index    = "16"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_16.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_16.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_16.sql
    }
  }
}

dashboard "dashboard_17" {
  title = "Dashboard 17"
  description = "Test dashboard 17 for performance benchmarking"

  tags = {
    category = "test"
    index    = "17"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_17.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_17.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_17.sql
    }
  }
}

dashboard "dashboard_18" {
  title = "Dashboard 18"
  description = "Test dashboard 18 for performance benchmarking"

  tags = {
    category = "test"
    index    = "18"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_18.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_18.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_18.sql
    }
  }
}

dashboard "dashboard_19" {
  title = "Dashboard 19"
  description = "Test dashboard 19 for performance benchmarking"

  tags = {
    category = "test"
    index    = "19"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_19.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_19.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_19.sql
    }
  }
}

dashboard "dashboard_20" {
  title = "Dashboard 20"
  description = "Test dashboard 20 for performance benchmarking"

  tags = {
    category = "test"
    index    = "20"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_20.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_20.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_20.sql
    }
  }
}

dashboard "dashboard_21" {
  title = "Dashboard 21"
  description = "Test dashboard 21 for performance benchmarking"

  tags = {
    category = "test"
    index    = "21"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_21.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_21.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_21.sql
    }
  }
}

dashboard "dashboard_22" {
  title = "Dashboard 22"
  description = "Test dashboard 22 for performance benchmarking"

  tags = {
    category = "test"
    index    = "22"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_22.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_22.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_22.sql
    }
  }
}

dashboard "dashboard_23" {
  title = "Dashboard 23"
  description = "Test dashboard 23 for performance benchmarking"

  tags = {
    category = "test"
    index    = "23"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_23.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_23.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_23.sql
    }
  }
}

dashboard "dashboard_24" {
  title = "Dashboard 24"
  description = "Test dashboard 24 for performance benchmarking"

  tags = {
    category = "test"
    index    = "24"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_24.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_24.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_24.sql
    }
  }
}

dashboard "dashboard_25" {
  title = "Dashboard 25"
  description = "Test dashboard 25 for performance benchmarking"

  tags = {
    category = "test"
    index    = "25"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_25.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_25.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_25.sql
    }
  }
}

dashboard "dashboard_26" {
  title = "Dashboard 26"
  description = "Test dashboard 26 for performance benchmarking"

  tags = {
    category = "test"
    index    = "26"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_26.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_26.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_26.sql
    }
  }
}

dashboard "dashboard_27" {
  title = "Dashboard 27"
  description = "Test dashboard 27 for performance benchmarking"

  tags = {
    category = "test"
    index    = "27"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_27.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_27.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_27.sql
    }
  }
}

dashboard "dashboard_28" {
  title = "Dashboard 28"
  description = "Test dashboard 28 for performance benchmarking"

  tags = {
    category = "test"
    index    = "28"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_28.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_28.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_28.sql
    }
  }
}

dashboard "dashboard_29" {
  title = "Dashboard 29"
  description = "Test dashboard 29 for performance benchmarking"

  tags = {
    category = "test"
    index    = "29"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_29.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_29.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_29.sql
    }
  }
}

dashboard "dashboard_30" {
  title = "Dashboard 30"
  description = "Test dashboard 30 for performance benchmarking"

  tags = {
    category = "test"
    index    = "30"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_30.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_30.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_30.sql
    }
  }
}

dashboard "dashboard_31" {
  title = "Dashboard 31"
  description = "Test dashboard 31 for performance benchmarking"

  tags = {
    category = "test"
    index    = "31"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_31.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_31.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_31.sql
    }
  }
}

dashboard "dashboard_32" {
  title = "Dashboard 32"
  description = "Test dashboard 32 for performance benchmarking"

  tags = {
    category = "test"
    index    = "32"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_32.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_32.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_32.sql
    }
  }
}

dashboard "dashboard_33" {
  title = "Dashboard 33"
  description = "Test dashboard 33 for performance benchmarking"

  tags = {
    category = "test"
    index    = "33"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_33.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_33.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_33.sql
    }
  }
}

dashboard "dashboard_34" {
  title = "Dashboard 34"
  description = "Test dashboard 34 for performance benchmarking"

  tags = {
    category = "test"
    index    = "34"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_34.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_34.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_34.sql
    }
  }
}

dashboard "dashboard_35" {
  title = "Dashboard 35"
  description = "Test dashboard 35 for performance benchmarking"

  tags = {
    category = "test"
    index    = "35"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_35.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_35.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_35.sql
    }
  }
}

dashboard "dashboard_36" {
  title = "Dashboard 36"
  description = "Test dashboard 36 for performance benchmarking"

  tags = {
    category = "test"
    index    = "36"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_36.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_36.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_36.sql
    }
  }
}

dashboard "dashboard_37" {
  title = "Dashboard 37"
  description = "Test dashboard 37 for performance benchmarking"

  tags = {
    category = "test"
    index    = "37"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_37.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_37.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_37.sql
    }
  }
}

dashboard "dashboard_38" {
  title = "Dashboard 38"
  description = "Test dashboard 38 for performance benchmarking"

  tags = {
    category = "test"
    index    = "38"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_38.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_38.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_38.sql
    }
  }
}

dashboard "dashboard_39" {
  title = "Dashboard 39"
  description = "Test dashboard 39 for performance benchmarking"

  tags = {
    category = "test"
    index    = "39"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_39.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_39.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_39.sql
    }
  }
}

dashboard "dashboard_40" {
  title = "Dashboard 40"
  description = "Test dashboard 40 for performance benchmarking"

  tags = {
    category = "test"
    index    = "40"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_40.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_40.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_40.sql
    }
  }
}

dashboard "dashboard_41" {
  title = "Dashboard 41"
  description = "Test dashboard 41 for performance benchmarking"

  tags = {
    category = "test"
    index    = "41"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_41.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_41.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_41.sql
    }
  }
}

dashboard "dashboard_42" {
  title = "Dashboard 42"
  description = "Test dashboard 42 for performance benchmarking"

  tags = {
    category = "test"
    index    = "42"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_42.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_42.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_42.sql
    }
  }
}

dashboard "dashboard_43" {
  title = "Dashboard 43"
  description = "Test dashboard 43 for performance benchmarking"

  tags = {
    category = "test"
    index    = "43"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_43.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_43.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_43.sql
    }
  }
}

dashboard "dashboard_44" {
  title = "Dashboard 44"
  description = "Test dashboard 44 for performance benchmarking"

  tags = {
    category = "test"
    index    = "44"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_44.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_44.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_44.sql
    }
  }
}

dashboard "dashboard_45" {
  title = "Dashboard 45"
  description = "Test dashboard 45 for performance benchmarking"

  tags = {
    category = "test"
    index    = "45"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_45.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_45.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_45.sql
    }
  }
}

dashboard "dashboard_46" {
  title = "Dashboard 46"
  description = "Test dashboard 46 for performance benchmarking"

  tags = {
    category = "test"
    index    = "46"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_46.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_46.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_46.sql
    }
  }
}

dashboard "dashboard_47" {
  title = "Dashboard 47"
  description = "Test dashboard 47 for performance benchmarking"

  tags = {
    category = "test"
    index    = "47"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_47.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_47.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_47.sql
    }
  }
}

dashboard "dashboard_48" {
  title = "Dashboard 48"
  description = "Test dashboard 48 for performance benchmarking"

  tags = {
    category = "test"
    index    = "48"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_48.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_48.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_48.sql
    }
  }
}

dashboard "dashboard_49" {
  title = "Dashboard 49"
  description = "Test dashboard 49 for performance benchmarking"

  tags = {
    category = "test"
    index    = "49"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_49.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_49.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_49.sql
    }
  }
}

dashboard "dashboard_50" {
  title = "Dashboard 50"
  description = "Test dashboard 50 for performance benchmarking"

  tags = {
    category = "test"
    index    = "50"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_50.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_50.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_50.sql
    }
  }
}

dashboard "dashboard_51" {
  title = "Dashboard 51"
  description = "Test dashboard 51 for performance benchmarking"

  tags = {
    category = "test"
    index    = "51"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_51.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_51.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_51.sql
    }
  }
}

dashboard "dashboard_52" {
  title = "Dashboard 52"
  description = "Test dashboard 52 for performance benchmarking"

  tags = {
    category = "test"
    index    = "52"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_52.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_52.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_52.sql
    }
  }
}

dashboard "dashboard_53" {
  title = "Dashboard 53"
  description = "Test dashboard 53 for performance benchmarking"

  tags = {
    category = "test"
    index    = "53"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_53.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_53.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_53.sql
    }
  }
}

dashboard "dashboard_54" {
  title = "Dashboard 54"
  description = "Test dashboard 54 for performance benchmarking"

  tags = {
    category = "test"
    index    = "54"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_54.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_54.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_54.sql
    }
  }
}

dashboard "dashboard_55" {
  title = "Dashboard 55"
  description = "Test dashboard 55 for performance benchmarking"

  tags = {
    category = "test"
    index    = "55"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_55.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_55.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_55.sql
    }
  }
}

dashboard "dashboard_56" {
  title = "Dashboard 56"
  description = "Test dashboard 56 for performance benchmarking"

  tags = {
    category = "test"
    index    = "56"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_56.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_56.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_56.sql
    }
  }
}

dashboard "dashboard_57" {
  title = "Dashboard 57"
  description = "Test dashboard 57 for performance benchmarking"

  tags = {
    category = "test"
    index    = "57"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_57.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_57.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_57.sql
    }
  }
}

dashboard "dashboard_58" {
  title = "Dashboard 58"
  description = "Test dashboard 58 for performance benchmarking"

  tags = {
    category = "test"
    index    = "58"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_58.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_58.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_58.sql
    }
  }
}

dashboard "dashboard_59" {
  title = "Dashboard 59"
  description = "Test dashboard 59 for performance benchmarking"

  tags = {
    category = "test"
    index    = "59"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_59.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_59.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_59.sql
    }
  }
}

dashboard "dashboard_60" {
  title = "Dashboard 60"
  description = "Test dashboard 60 for performance benchmarking"

  tags = {
    category = "test"
    index    = "60"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_60.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_60.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_60.sql
    }
  }
}

dashboard "dashboard_61" {
  title = "Dashboard 61"
  description = "Test dashboard 61 for performance benchmarking"

  tags = {
    category = "test"
    index    = "61"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_61.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_61.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_61.sql
    }
  }
}

dashboard "dashboard_62" {
  title = "Dashboard 62"
  description = "Test dashboard 62 for performance benchmarking"

  tags = {
    category = "test"
    index    = "62"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_62.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_62.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_62.sql
    }
  }
}

dashboard "dashboard_63" {
  title = "Dashboard 63"
  description = "Test dashboard 63 for performance benchmarking"

  tags = {
    category = "test"
    index    = "63"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_63.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_63.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_63.sql
    }
  }
}

dashboard "dashboard_64" {
  title = "Dashboard 64"
  description = "Test dashboard 64 for performance benchmarking"

  tags = {
    category = "test"
    index    = "64"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_64.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_64.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_64.sql
    }
  }
}

dashboard "dashboard_65" {
  title = "Dashboard 65"
  description = "Test dashboard 65 for performance benchmarking"

  tags = {
    category = "test"
    index    = "65"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_65.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_65.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_65.sql
    }
  }
}

dashboard "dashboard_66" {
  title = "Dashboard 66"
  description = "Test dashboard 66 for performance benchmarking"

  tags = {
    category = "test"
    index    = "66"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_66.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_66.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_66.sql
    }
  }
}

dashboard "dashboard_67" {
  title = "Dashboard 67"
  description = "Test dashboard 67 for performance benchmarking"

  tags = {
    category = "test"
    index    = "67"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_67.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_67.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_67.sql
    }
  }
}

dashboard "dashboard_68" {
  title = "Dashboard 68"
  description = "Test dashboard 68 for performance benchmarking"

  tags = {
    category = "test"
    index    = "68"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_68.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_68.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_68.sql
    }
  }
}

dashboard "dashboard_69" {
  title = "Dashboard 69"
  description = "Test dashboard 69 for performance benchmarking"

  tags = {
    category = "test"
    index    = "69"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_69.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_69.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_69.sql
    }
  }
}

dashboard "dashboard_70" {
  title = "Dashboard 70"
  description = "Test dashboard 70 for performance benchmarking"

  tags = {
    category = "test"
    index    = "70"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_70.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_70.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_70.sql
    }
  }
}

dashboard "dashboard_71" {
  title = "Dashboard 71"
  description = "Test dashboard 71 for performance benchmarking"

  tags = {
    category = "test"
    index    = "71"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_71.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_71.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_71.sql
    }
  }
}

dashboard "dashboard_72" {
  title = "Dashboard 72"
  description = "Test dashboard 72 for performance benchmarking"

  tags = {
    category = "test"
    index    = "72"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_72.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_72.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_72.sql
    }
  }
}

dashboard "dashboard_73" {
  title = "Dashboard 73"
  description = "Test dashboard 73 for performance benchmarking"

  tags = {
    category = "test"
    index    = "73"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_73.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_73.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_73.sql
    }
  }
}

dashboard "dashboard_74" {
  title = "Dashboard 74"
  description = "Test dashboard 74 for performance benchmarking"

  tags = {
    category = "test"
    index    = "74"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_74.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_74.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_74.sql
    }
  }
}

dashboard "dashboard_75" {
  title = "Dashboard 75"
  description = "Test dashboard 75 for performance benchmarking"

  tags = {
    category = "test"
    index    = "75"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_75.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_75.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_75.sql
    }
  }
}

dashboard "dashboard_76" {
  title = "Dashboard 76"
  description = "Test dashboard 76 for performance benchmarking"

  tags = {
    category = "test"
    index    = "76"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_76.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_76.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_76.sql
    }
  }
}

dashboard "dashboard_77" {
  title = "Dashboard 77"
  description = "Test dashboard 77 for performance benchmarking"

  tags = {
    category = "test"
    index    = "77"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_77.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_77.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_77.sql
    }
  }
}

dashboard "dashboard_78" {
  title = "Dashboard 78"
  description = "Test dashboard 78 for performance benchmarking"

  tags = {
    category = "test"
    index    = "78"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_78.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_78.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_78.sql
    }
  }
}

dashboard "dashboard_79" {
  title = "Dashboard 79"
  description = "Test dashboard 79 for performance benchmarking"

  tags = {
    category = "test"
    index    = "79"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_79.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_79.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_79.sql
    }
  }
}

dashboard "dashboard_80" {
  title = "Dashboard 80"
  description = "Test dashboard 80 for performance benchmarking"

  tags = {
    category = "test"
    index    = "80"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_80.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_80.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_80.sql
    }
  }
}

dashboard "dashboard_81" {
  title = "Dashboard 81"
  description = "Test dashboard 81 for performance benchmarking"

  tags = {
    category = "test"
    index    = "81"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_81.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_81.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_81.sql
    }
  }
}

dashboard "dashboard_82" {
  title = "Dashboard 82"
  description = "Test dashboard 82 for performance benchmarking"

  tags = {
    category = "test"
    index    = "82"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_82.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_82.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_82.sql
    }
  }
}

dashboard "dashboard_83" {
  title = "Dashboard 83"
  description = "Test dashboard 83 for performance benchmarking"

  tags = {
    category = "test"
    index    = "83"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_83.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_83.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_83.sql
    }
  }
}

dashboard "dashboard_84" {
  title = "Dashboard 84"
  description = "Test dashboard 84 for performance benchmarking"

  tags = {
    category = "test"
    index    = "84"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_84.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_84.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_84.sql
    }
  }
}

dashboard "dashboard_85" {
  title = "Dashboard 85"
  description = "Test dashboard 85 for performance benchmarking"

  tags = {
    category = "test"
    index    = "85"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_85.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_85.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_85.sql
    }
  }
}

dashboard "dashboard_86" {
  title = "Dashboard 86"
  description = "Test dashboard 86 for performance benchmarking"

  tags = {
    category = "test"
    index    = "86"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_86.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_86.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_86.sql
    }
  }
}

dashboard "dashboard_87" {
  title = "Dashboard 87"
  description = "Test dashboard 87 for performance benchmarking"

  tags = {
    category = "test"
    index    = "87"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_87.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_87.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_87.sql
    }
  }
}

dashboard "dashboard_88" {
  title = "Dashboard 88"
  description = "Test dashboard 88 for performance benchmarking"

  tags = {
    category = "test"
    index    = "88"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_88.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_88.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_88.sql
    }
  }
}

dashboard "dashboard_89" {
  title = "Dashboard 89"
  description = "Test dashboard 89 for performance benchmarking"

  tags = {
    category = "test"
    index    = "89"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_89.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_89.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_89.sql
    }
  }
}

dashboard "dashboard_90" {
  title = "Dashboard 90"
  description = "Test dashboard 90 for performance benchmarking"

  tags = {
    category = "test"
    index    = "90"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_90.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_90.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_90.sql
    }
  }
}

dashboard "dashboard_91" {
  title = "Dashboard 91"
  description = "Test dashboard 91 for performance benchmarking"

  tags = {
    category = "test"
    index    = "91"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_91.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_91.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_91.sql
    }
  }
}

dashboard "dashboard_92" {
  title = "Dashboard 92"
  description = "Test dashboard 92 for performance benchmarking"

  tags = {
    category = "test"
    index    = "92"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_92.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_92.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_92.sql
    }
  }
}

dashboard "dashboard_93" {
  title = "Dashboard 93"
  description = "Test dashboard 93 for performance benchmarking"

  tags = {
    category = "test"
    index    = "93"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_93.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_93.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_93.sql
    }
  }
}

dashboard "dashboard_94" {
  title = "Dashboard 94"
  description = "Test dashboard 94 for performance benchmarking"

  tags = {
    category = "test"
    index    = "94"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_94.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_94.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_94.sql
    }
  }
}

dashboard "dashboard_95" {
  title = "Dashboard 95"
  description = "Test dashboard 95 for performance benchmarking"

  tags = {
    category = "test"
    index    = "95"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_95.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_95.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_95.sql
    }
  }
}

dashboard "dashboard_96" {
  title = "Dashboard 96"
  description = "Test dashboard 96 for performance benchmarking"

  tags = {
    category = "test"
    index    = "96"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_96.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_96.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_96.sql
    }
  }
}

dashboard "dashboard_97" {
  title = "Dashboard 97"
  description = "Test dashboard 97 for performance benchmarking"

  tags = {
    category = "test"
    index    = "97"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_97.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_97.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_97.sql
    }
  }
}

dashboard "dashboard_98" {
  title = "Dashboard 98"
  description = "Test dashboard 98 for performance benchmarking"

  tags = {
    category = "test"
    index    = "98"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_98.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_98.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_98.sql
    }
  }
}

dashboard "dashboard_99" {
  title = "Dashboard 99"
  description = "Test dashboard 99 for performance benchmarking"

  tags = {
    category = "test"
    index    = "99"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_99.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_99.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_99.sql
    }
  }
}

dashboard "dashboard_100" {
  title = "Dashboard 100"
  description = "Test dashboard 100 for performance benchmarking"

  tags = {
    category = "test"
    index    = "100"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_100.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_100.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_100.sql
    }
  }
}

dashboard "dashboard_101" {
  title = "Dashboard 101"
  description = "Test dashboard 101 for performance benchmarking"

  tags = {
    category = "test"
    index    = "101"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_101.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_101.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_101.sql
    }
  }
}

dashboard "dashboard_102" {
  title = "Dashboard 102"
  description = "Test dashboard 102 for performance benchmarking"

  tags = {
    category = "test"
    index    = "102"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_102.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_102.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_102.sql
    }
  }
}

dashboard "dashboard_103" {
  title = "Dashboard 103"
  description = "Test dashboard 103 for performance benchmarking"

  tags = {
    category = "test"
    index    = "103"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_103.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_103.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_103.sql
    }
  }
}

dashboard "dashboard_104" {
  title = "Dashboard 104"
  description = "Test dashboard 104 for performance benchmarking"

  tags = {
    category = "test"
    index    = "104"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_104.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_104.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_104.sql
    }
  }
}

dashboard "dashboard_105" {
  title = "Dashboard 105"
  description = "Test dashboard 105 for performance benchmarking"

  tags = {
    category = "test"
    index    = "105"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_105.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_105.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_105.sql
    }
  }
}

dashboard "dashboard_106" {
  title = "Dashboard 106"
  description = "Test dashboard 106 for performance benchmarking"

  tags = {
    category = "test"
    index    = "106"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_106.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_106.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_106.sql
    }
  }
}

dashboard "dashboard_107" {
  title = "Dashboard 107"
  description = "Test dashboard 107 for performance benchmarking"

  tags = {
    category = "test"
    index    = "107"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_107.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_107.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_107.sql
    }
  }
}

dashboard "dashboard_108" {
  title = "Dashboard 108"
  description = "Test dashboard 108 for performance benchmarking"

  tags = {
    category = "test"
    index    = "108"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_108.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_108.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_108.sql
    }
  }
}

dashboard "dashboard_109" {
  title = "Dashboard 109"
  description = "Test dashboard 109 for performance benchmarking"

  tags = {
    category = "test"
    index    = "109"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_109.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_109.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_109.sql
    }
  }
}

dashboard "dashboard_110" {
  title = "Dashboard 110"
  description = "Test dashboard 110 for performance benchmarking"

  tags = {
    category = "test"
    index    = "110"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_110.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_110.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_110.sql
    }
  }
}

dashboard "dashboard_111" {
  title = "Dashboard 111"
  description = "Test dashboard 111 for performance benchmarking"

  tags = {
    category = "test"
    index    = "111"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_111.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_111.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_111.sql
    }
  }
}

dashboard "dashboard_112" {
  title = "Dashboard 112"
  description = "Test dashboard 112 for performance benchmarking"

  tags = {
    category = "test"
    index    = "112"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_112.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_112.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_112.sql
    }
  }
}

dashboard "dashboard_113" {
  title = "Dashboard 113"
  description = "Test dashboard 113 for performance benchmarking"

  tags = {
    category = "test"
    index    = "113"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_113.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_113.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_113.sql
    }
  }
}

dashboard "dashboard_114" {
  title = "Dashboard 114"
  description = "Test dashboard 114 for performance benchmarking"

  tags = {
    category = "test"
    index    = "114"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_114.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_114.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_114.sql
    }
  }
}

dashboard "dashboard_115" {
  title = "Dashboard 115"
  description = "Test dashboard 115 for performance benchmarking"

  tags = {
    category = "test"
    index    = "115"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_115.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_115.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_115.sql
    }
  }
}

dashboard "dashboard_116" {
  title = "Dashboard 116"
  description = "Test dashboard 116 for performance benchmarking"

  tags = {
    category = "test"
    index    = "116"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_116.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_116.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_116.sql
    }
  }
}

dashboard "dashboard_117" {
  title = "Dashboard 117"
  description = "Test dashboard 117 for performance benchmarking"

  tags = {
    category = "test"
    index    = "117"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_117.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_117.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_117.sql
    }
  }
}

dashboard "dashboard_118" {
  title = "Dashboard 118"
  description = "Test dashboard 118 for performance benchmarking"

  tags = {
    category = "test"
    index    = "118"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_118.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_118.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_118.sql
    }
  }
}

dashboard "dashboard_119" {
  title = "Dashboard 119"
  description = "Test dashboard 119 for performance benchmarking"

  tags = {
    category = "test"
    index    = "119"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_119.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_119.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_119.sql
    }
  }
}

dashboard "dashboard_120" {
  title = "Dashboard 120"
  description = "Test dashboard 120 for performance benchmarking"

  tags = {
    category = "test"
    index    = "120"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_120.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_120.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_120.sql
    }
  }
}

dashboard "dashboard_121" {
  title = "Dashboard 121"
  description = "Test dashboard 121 for performance benchmarking"

  tags = {
    category = "test"
    index    = "121"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_121.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_121.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_121.sql
    }
  }
}

dashboard "dashboard_122" {
  title = "Dashboard 122"
  description = "Test dashboard 122 for performance benchmarking"

  tags = {
    category = "test"
    index    = "122"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_122.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_122.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_122.sql
    }
  }
}

dashboard "dashboard_123" {
  title = "Dashboard 123"
  description = "Test dashboard 123 for performance benchmarking"

  tags = {
    category = "test"
    index    = "123"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_123.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_123.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_123.sql
    }
  }
}

dashboard "dashboard_124" {
  title = "Dashboard 124"
  description = "Test dashboard 124 for performance benchmarking"

  tags = {
    category = "test"
    index    = "124"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_124.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_124.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_124.sql
    }
  }
}

dashboard "dashboard_125" {
  title = "Dashboard 125"
  description = "Test dashboard 125 for performance benchmarking"

  tags = {
    category = "test"
    index    = "125"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_125.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_125.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_125.sql
    }
  }
}

dashboard "dashboard_126" {
  title = "Dashboard 126"
  description = "Test dashboard 126 for performance benchmarking"

  tags = {
    category = "test"
    index    = "126"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_126.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_126.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_126.sql
    }
  }
}

dashboard "dashboard_127" {
  title = "Dashboard 127"
  description = "Test dashboard 127 for performance benchmarking"

  tags = {
    category = "test"
    index    = "127"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_127.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_127.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_127.sql
    }
  }
}

dashboard "dashboard_128" {
  title = "Dashboard 128"
  description = "Test dashboard 128 for performance benchmarking"

  tags = {
    category = "test"
    index    = "128"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_128.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_128.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_128.sql
    }
  }
}

dashboard "dashboard_129" {
  title = "Dashboard 129"
  description = "Test dashboard 129 for performance benchmarking"

  tags = {
    category = "test"
    index    = "129"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_129.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_129.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_129.sql
    }
  }
}

dashboard "dashboard_130" {
  title = "Dashboard 130"
  description = "Test dashboard 130 for performance benchmarking"

  tags = {
    category = "test"
    index    = "130"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_130.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_130.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_130.sql
    }
  }
}

dashboard "dashboard_131" {
  title = "Dashboard 131"
  description = "Test dashboard 131 for performance benchmarking"

  tags = {
    category = "test"
    index    = "131"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_131.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_131.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_131.sql
    }
  }
}

dashboard "dashboard_132" {
  title = "Dashboard 132"
  description = "Test dashboard 132 for performance benchmarking"

  tags = {
    category = "test"
    index    = "132"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_132.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_132.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_132.sql
    }
  }
}

dashboard "dashboard_133" {
  title = "Dashboard 133"
  description = "Test dashboard 133 for performance benchmarking"

  tags = {
    category = "test"
    index    = "133"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_133.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_133.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_133.sql
    }
  }
}

dashboard "dashboard_134" {
  title = "Dashboard 134"
  description = "Test dashboard 134 for performance benchmarking"

  tags = {
    category = "test"
    index    = "134"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_134.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_134.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_134.sql
    }
  }
}

dashboard "dashboard_135" {
  title = "Dashboard 135"
  description = "Test dashboard 135 for performance benchmarking"

  tags = {
    category = "test"
    index    = "135"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_135.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_135.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_135.sql
    }
  }
}

dashboard "dashboard_136" {
  title = "Dashboard 136"
  description = "Test dashboard 136 for performance benchmarking"

  tags = {
    category = "test"
    index    = "136"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_136.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_136.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_136.sql
    }
  }
}

dashboard "dashboard_137" {
  title = "Dashboard 137"
  description = "Test dashboard 137 for performance benchmarking"

  tags = {
    category = "test"
    index    = "137"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_137.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_137.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_137.sql
    }
  }
}

dashboard "dashboard_138" {
  title = "Dashboard 138"
  description = "Test dashboard 138 for performance benchmarking"

  tags = {
    category = "test"
    index    = "138"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_138.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_138.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_138.sql
    }
  }
}

dashboard "dashboard_139" {
  title = "Dashboard 139"
  description = "Test dashboard 139 for performance benchmarking"

  tags = {
    category = "test"
    index    = "139"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_139.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_139.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_139.sql
    }
  }
}

dashboard "dashboard_140" {
  title = "Dashboard 140"
  description = "Test dashboard 140 for performance benchmarking"

  tags = {
    category = "test"
    index    = "140"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_140.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_140.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_140.sql
    }
  }
}

dashboard "dashboard_141" {
  title = "Dashboard 141"
  description = "Test dashboard 141 for performance benchmarking"

  tags = {
    category = "test"
    index    = "141"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_141.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_141.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_141.sql
    }
  }
}

dashboard "dashboard_142" {
  title = "Dashboard 142"
  description = "Test dashboard 142 for performance benchmarking"

  tags = {
    category = "test"
    index    = "142"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_142.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_142.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_142.sql
    }
  }
}

dashboard "dashboard_143" {
  title = "Dashboard 143"
  description = "Test dashboard 143 for performance benchmarking"

  tags = {
    category = "test"
    index    = "143"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_143.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_143.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_143.sql
    }
  }
}

dashboard "dashboard_144" {
  title = "Dashboard 144"
  description = "Test dashboard 144 for performance benchmarking"

  tags = {
    category = "test"
    index    = "144"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_144.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_144.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_144.sql
    }
  }
}

dashboard "dashboard_145" {
  title = "Dashboard 145"
  description = "Test dashboard 145 for performance benchmarking"

  tags = {
    category = "test"
    index    = "145"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_145.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_145.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_145.sql
    }
  }
}

dashboard "dashboard_146" {
  title = "Dashboard 146"
  description = "Test dashboard 146 for performance benchmarking"

  tags = {
    category = "test"
    index    = "146"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_146.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_146.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_146.sql
    }
  }
}

dashboard "dashboard_147" {
  title = "Dashboard 147"
  description = "Test dashboard 147 for performance benchmarking"

  tags = {
    category = "test"
    index    = "147"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_147.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_147.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_147.sql
    }
  }
}

dashboard "dashboard_148" {
  title = "Dashboard 148"
  description = "Test dashboard 148 for performance benchmarking"

  tags = {
    category = "test"
    index    = "148"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_148.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_148.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_148.sql
    }
  }
}

dashboard "dashboard_149" {
  title = "Dashboard 149"
  description = "Test dashboard 149 for performance benchmarking"

  tags = {
    category = "test"
    index    = "149"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_149.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_149.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_149.sql
    }
  }
}

dashboard "dashboard_150" {
  title = "Dashboard 150"
  description = "Test dashboard 150 for performance benchmarking"

  tags = {
    category = "test"
    index    = "150"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_150.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_150.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_150.sql
    }
  }
}

dashboard "dashboard_151" {
  title = "Dashboard 151"
  description = "Test dashboard 151 for performance benchmarking"

  tags = {
    category = "test"
    index    = "151"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_151.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_151.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_151.sql
    }
  }
}

dashboard "dashboard_152" {
  title = "Dashboard 152"
  description = "Test dashboard 152 for performance benchmarking"

  tags = {
    category = "test"
    index    = "152"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_152.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_152.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_152.sql
    }
  }
}

dashboard "dashboard_153" {
  title = "Dashboard 153"
  description = "Test dashboard 153 for performance benchmarking"

  tags = {
    category = "test"
    index    = "153"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_153.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_153.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_153.sql
    }
  }
}

dashboard "dashboard_154" {
  title = "Dashboard 154"
  description = "Test dashboard 154 for performance benchmarking"

  tags = {
    category = "test"
    index    = "154"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_154.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_154.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_154.sql
    }
  }
}

dashboard "dashboard_155" {
  title = "Dashboard 155"
  description = "Test dashboard 155 for performance benchmarking"

  tags = {
    category = "test"
    index    = "155"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_155.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_155.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_155.sql
    }
  }
}

dashboard "dashboard_156" {
  title = "Dashboard 156"
  description = "Test dashboard 156 for performance benchmarking"

  tags = {
    category = "test"
    index    = "156"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_156.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_156.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_156.sql
    }
  }
}

dashboard "dashboard_157" {
  title = "Dashboard 157"
  description = "Test dashboard 157 for performance benchmarking"

  tags = {
    category = "test"
    index    = "157"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_157.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_157.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_157.sql
    }
  }
}

dashboard "dashboard_158" {
  title = "Dashboard 158"
  description = "Test dashboard 158 for performance benchmarking"

  tags = {
    category = "test"
    index    = "158"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_158.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_158.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_158.sql
    }
  }
}

dashboard "dashboard_159" {
  title = "Dashboard 159"
  description = "Test dashboard 159 for performance benchmarking"

  tags = {
    category = "test"
    index    = "159"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_159.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_159.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_159.sql
    }
  }
}

dashboard "dashboard_160" {
  title = "Dashboard 160"
  description = "Test dashboard 160 for performance benchmarking"

  tags = {
    category = "test"
    index    = "160"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_160.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_160.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_160.sql
    }
  }
}

dashboard "dashboard_161" {
  title = "Dashboard 161"
  description = "Test dashboard 161 for performance benchmarking"

  tags = {
    category = "test"
    index    = "161"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_161.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_161.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_161.sql
    }
  }
}

dashboard "dashboard_162" {
  title = "Dashboard 162"
  description = "Test dashboard 162 for performance benchmarking"

  tags = {
    category = "test"
    index    = "162"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_162.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_162.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_162.sql
    }
  }
}

dashboard "dashboard_163" {
  title = "Dashboard 163"
  description = "Test dashboard 163 for performance benchmarking"

  tags = {
    category = "test"
    index    = "163"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_163.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_163.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_163.sql
    }
  }
}

dashboard "dashboard_164" {
  title = "Dashboard 164"
  description = "Test dashboard 164 for performance benchmarking"

  tags = {
    category = "test"
    index    = "164"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_164.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_164.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_164.sql
    }
  }
}

dashboard "dashboard_165" {
  title = "Dashboard 165"
  description = "Test dashboard 165 for performance benchmarking"

  tags = {
    category = "test"
    index    = "165"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_165.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_165.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_165.sql
    }
  }
}

dashboard "dashboard_166" {
  title = "Dashboard 166"
  description = "Test dashboard 166 for performance benchmarking"

  tags = {
    category = "test"
    index    = "166"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_166.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_166.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_166.sql
    }
  }
}

dashboard "dashboard_167" {
  title = "Dashboard 167"
  description = "Test dashboard 167 for performance benchmarking"

  tags = {
    category = "test"
    index    = "167"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_167.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_167.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_167.sql
    }
  }
}

dashboard "dashboard_168" {
  title = "Dashboard 168"
  description = "Test dashboard 168 for performance benchmarking"

  tags = {
    category = "test"
    index    = "168"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_168.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_168.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_168.sql
    }
  }
}

dashboard "dashboard_169" {
  title = "Dashboard 169"
  description = "Test dashboard 169 for performance benchmarking"

  tags = {
    category = "test"
    index    = "169"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_169.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_169.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_169.sql
    }
  }
}

dashboard "dashboard_170" {
  title = "Dashboard 170"
  description = "Test dashboard 170 for performance benchmarking"

  tags = {
    category = "test"
    index    = "170"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_170.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_170.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_170.sql
    }
  }
}

dashboard "dashboard_171" {
  title = "Dashboard 171"
  description = "Test dashboard 171 for performance benchmarking"

  tags = {
    category = "test"
    index    = "171"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_171.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_171.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_171.sql
    }
  }
}

dashboard "dashboard_172" {
  title = "Dashboard 172"
  description = "Test dashboard 172 for performance benchmarking"

  tags = {
    category = "test"
    index    = "172"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_172.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_172.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_172.sql
    }
  }
}

dashboard "dashboard_173" {
  title = "Dashboard 173"
  description = "Test dashboard 173 for performance benchmarking"

  tags = {
    category = "test"
    index    = "173"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_173.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_173.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_173.sql
    }
  }
}

dashboard "dashboard_174" {
  title = "Dashboard 174"
  description = "Test dashboard 174 for performance benchmarking"

  tags = {
    category = "test"
    index    = "174"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_174.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_174.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_174.sql
    }
  }
}

dashboard "dashboard_175" {
  title = "Dashboard 175"
  description = "Test dashboard 175 for performance benchmarking"

  tags = {
    category = "test"
    index    = "175"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_175.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_175.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_175.sql
    }
  }
}

dashboard "dashboard_176" {
  title = "Dashboard 176"
  description = "Test dashboard 176 for performance benchmarking"

  tags = {
    category = "test"
    index    = "176"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_176.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_176.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_176.sql
    }
  }
}

dashboard "dashboard_177" {
  title = "Dashboard 177"
  description = "Test dashboard 177 for performance benchmarking"

  tags = {
    category = "test"
    index    = "177"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_177.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_177.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_177.sql
    }
  }
}

dashboard "dashboard_178" {
  title = "Dashboard 178"
  description = "Test dashboard 178 for performance benchmarking"

  tags = {
    category = "test"
    index    = "178"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_178.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_178.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_178.sql
    }
  }
}

dashboard "dashboard_179" {
  title = "Dashboard 179"
  description = "Test dashboard 179 for performance benchmarking"

  tags = {
    category = "test"
    index    = "179"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_179.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_179.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_179.sql
    }
  }
}

dashboard "dashboard_180" {
  title = "Dashboard 180"
  description = "Test dashboard 180 for performance benchmarking"

  tags = {
    category = "test"
    index    = "180"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_180.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_180.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_180.sql
    }
  }
}

dashboard "dashboard_181" {
  title = "Dashboard 181"
  description = "Test dashboard 181 for performance benchmarking"

  tags = {
    category = "test"
    index    = "181"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_181.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_181.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_181.sql
    }
  }
}

dashboard "dashboard_182" {
  title = "Dashboard 182"
  description = "Test dashboard 182 for performance benchmarking"

  tags = {
    category = "test"
    index    = "182"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_182.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_182.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_182.sql
    }
  }
}

dashboard "dashboard_183" {
  title = "Dashboard 183"
  description = "Test dashboard 183 for performance benchmarking"

  tags = {
    category = "test"
    index    = "183"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_183.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_183.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_183.sql
    }
  }
}

dashboard "dashboard_184" {
  title = "Dashboard 184"
  description = "Test dashboard 184 for performance benchmarking"

  tags = {
    category = "test"
    index    = "184"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_184.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_184.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_184.sql
    }
  }
}

dashboard "dashboard_185" {
  title = "Dashboard 185"
  description = "Test dashboard 185 for performance benchmarking"

  tags = {
    category = "test"
    index    = "185"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_185.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_185.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_185.sql
    }
  }
}

dashboard "dashboard_186" {
  title = "Dashboard 186"
  description = "Test dashboard 186 for performance benchmarking"

  tags = {
    category = "test"
    index    = "186"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_186.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_186.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_186.sql
    }
  }
}

dashboard "dashboard_187" {
  title = "Dashboard 187"
  description = "Test dashboard 187 for performance benchmarking"

  tags = {
    category = "test"
    index    = "187"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_187.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_187.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_187.sql
    }
  }
}

dashboard "dashboard_188" {
  title = "Dashboard 188"
  description = "Test dashboard 188 for performance benchmarking"

  tags = {
    category = "test"
    index    = "188"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_188.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_188.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_188.sql
    }
  }
}

dashboard "dashboard_189" {
  title = "Dashboard 189"
  description = "Test dashboard 189 for performance benchmarking"

  tags = {
    category = "test"
    index    = "189"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_189.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_189.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_189.sql
    }
  }
}

dashboard "dashboard_190" {
  title = "Dashboard 190"
  description = "Test dashboard 190 for performance benchmarking"

  tags = {
    category = "test"
    index    = "190"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_190.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_190.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_190.sql
    }
  }
}

dashboard "dashboard_191" {
  title = "Dashboard 191"
  description = "Test dashboard 191 for performance benchmarking"

  tags = {
    category = "test"
    index    = "191"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_191.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_191.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_191.sql
    }
  }
}

dashboard "dashboard_192" {
  title = "Dashboard 192"
  description = "Test dashboard 192 for performance benchmarking"

  tags = {
    category = "test"
    index    = "192"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_192.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_192.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_192.sql
    }
  }
}

dashboard "dashboard_193" {
  title = "Dashboard 193"
  description = "Test dashboard 193 for performance benchmarking"

  tags = {
    category = "test"
    index    = "193"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_193.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_193.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_193.sql
    }
  }
}

dashboard "dashboard_194" {
  title = "Dashboard 194"
  description = "Test dashboard 194 for performance benchmarking"

  tags = {
    category = "test"
    index    = "194"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_194.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_194.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_194.sql
    }
  }
}

dashboard "dashboard_195" {
  title = "Dashboard 195"
  description = "Test dashboard 195 for performance benchmarking"

  tags = {
    category = "test"
    index    = "195"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_195.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_195.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_195.sql
    }
  }
}

dashboard "dashboard_196" {
  title = "Dashboard 196"
  description = "Test dashboard 196 for performance benchmarking"

  tags = {
    category = "test"
    index    = "196"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_196.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_196.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_196.sql
    }
  }
}

dashboard "dashboard_197" {
  title = "Dashboard 197"
  description = "Test dashboard 197 for performance benchmarking"

  tags = {
    category = "test"
    index    = "197"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_197.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_197.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_197.sql
    }
  }
}

dashboard "dashboard_198" {
  title = "Dashboard 198"
  description = "Test dashboard 198 for performance benchmarking"

  tags = {
    category = "test"
    index    = "198"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_198.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_198.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_198.sql
    }
  }
}

dashboard "dashboard_199" {
  title = "Dashboard 199"
  description = "Test dashboard 199 for performance benchmarking"

  tags = {
    category = "test"
    index    = "199"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_199.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_199.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_199.sql
    }
  }
}
