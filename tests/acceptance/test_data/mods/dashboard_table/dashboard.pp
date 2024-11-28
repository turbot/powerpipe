dashboard "testing_card_blocks" {
  title = "Testing card blocks"

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
        primary_key = true
      }

      column "name" {
        display = "all"
      }
    }
  }
}