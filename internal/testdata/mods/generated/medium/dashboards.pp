
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
