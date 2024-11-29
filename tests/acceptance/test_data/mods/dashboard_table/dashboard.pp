dashboard "testing_tables" {
  title = "Testing table blocks"

  container {
    table {
      title = "Employee table"
      width = 4

      sql   = <<-EOQ
        SELECT
          011 AS id,
          'dwight' AS name,
          'manager' AS role
        UNION ALL
        SELECT
            012 AS id,
            'jim' AS name,
            'developer' AS role
        UNION ALL
        SELECT
            013 AS id,
            'pam' AS name,
            'designer' AS role;
      EOQ

      column "id" {
        display = "all"
        diff_mode = "key"
      }

      column "name" {
        display = "all"
      }
    }
  }
}