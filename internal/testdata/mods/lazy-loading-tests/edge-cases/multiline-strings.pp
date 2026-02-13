# Resources with complex string formatting

query "multiline_description" {
  title = "Multiline Description Query"
  description = <<-EOT
    This is a multiline description that spans
    multiple lines and includes:
    - Bullet point 1
    - Bullet point 2
    - Bullet point 3

    And even some code:
    ```sql
    SELECT * FROM table
    ```

    With a final paragraph explaining the purpose
    of this query in detail.
  EOT
  sql = "SELECT 1 as value"
}

query "heredoc_sql" {
  title       = "Heredoc SQL Query"
  description = "SQL using heredoc syntax"
  sql         = <<-EOQ
    SELECT
      id,
      name,
      -- This is a comment inside the SQL
      description,
      /*
       * Multi-line comment
       * spanning multiple lines
       */
      status
    FROM
      resources
    WHERE
      enabled = true
      AND name LIKE '%test%'
    ORDER BY
      created_at DESC;
  EOQ
}

control "control_multiline" {
  title = "Multiline Control"
  description = <<-EOT
    This control checks for:

    1. Proper configuration
    2. Security settings
    3. Compliance status

    Severity: HIGH
    Category: Security
  EOT
  sql = <<-EOQ
    SELECT
      'pass' as status,
      'multiline_resource' as resource,
      'All checks passed' as reason
    WHERE
      1 = 1;
  EOQ
}

dashboard "multiline_dashboard" {
  title = "Multiline Dashboard"
  description = <<-EOT
    Dashboard demonstrating multiline strings in:
    - Titles
    - Descriptions
    - SQL queries
  EOT

  text {
    value = <<-EOT
      # Welcome to the Dashboard

      This dashboard demonstrates **markdown** formatting with:

      1. Headers
      2. Lists
      3. **Bold** and *italic* text

      > A blockquote for emphasis

      ```
      And even code blocks
      ```
    EOT
  }

  card {
    title = "Count"
    sql   = query.multiline_description.sql
  }
}
