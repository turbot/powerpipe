# Resources with very large SQL heredocs

query "huge_sql" {
  title       = "Huge SQL Query"
  description = "Query with very large SQL statement to test heredoc handling"
  sql         = <<-EOQ
    -- This is a very large SQL query to test heredoc handling
    -- Line 1 of many
    -- Line 2 of many
    -- Line 3 of many
    -- Line 4 of many
    -- Line 5 of many
    SELECT
      t1.id as id_1,
      t1.name as name_1,
      t1.description as description_1,
      t1.created_at as created_at_1,
      t1.updated_at as updated_at_1,
      t1.status as status_1,
      t1.category as category_1,
      t1.priority as priority_1,
      t1.owner as owner_1,
      t1.team as team_1,
      t2.id as id_2,
      t2.name as name_2,
      t2.description as description_2,
      t2.created_at as created_at_2,
      t2.updated_at as updated_at_2,
      t2.status as status_2,
      t2.category as category_2,
      t2.priority as priority_2,
      t2.owner as owner_2,
      t2.team as team_2,
      t3.id as id_3,
      t3.name as name_3,
      t3.description as description_3,
      t3.created_at as created_at_3,
      t3.updated_at as updated_at_3,
      t3.status as status_3,
      t3.category as category_3,
      t3.priority as priority_3,
      t3.owner as owner_3,
      t3.team as team_3,
      CASE
        WHEN t1.status = 'active' AND t2.status = 'active' AND t3.status = 'active' THEN 'all_active'
        WHEN t1.status = 'active' OR t2.status = 'active' OR t3.status = 'active' THEN 'some_active'
        ELSE 'none_active'
      END as combined_status,
      COALESCE(t1.priority, t2.priority, t3.priority, 'low') as effective_priority,
      COUNT(*) OVER (PARTITION BY t1.category) as category_count,
      ROW_NUMBER() OVER (ORDER BY t1.created_at DESC) as row_num,
      DENSE_RANK() OVER (ORDER BY t1.priority) as priority_rank,
      LAG(t1.status) OVER (ORDER BY t1.created_at) as previous_status,
      LEAD(t1.status) OVER (ORDER BY t1.created_at) as next_status,
      FIRST_VALUE(t1.name) OVER (PARTITION BY t1.category ORDER BY t1.created_at) as first_in_category,
      LAST_VALUE(t1.name) OVER (PARTITION BY t1.category ORDER BY t1.created_at) as last_in_category
    FROM
      table_1 t1
      LEFT JOIN table_2 t2 ON t1.id = t2.table_1_id AND t2.status != 'deleted'
      LEFT JOIN table_3 t3 ON t2.id = t3.table_2_id AND t3.status != 'deleted'
      LEFT JOIN table_4 t4 ON t1.id = t4.table_1_id
      LEFT JOIN table_5 t5 ON t4.id = t5.table_4_id
    WHERE
      t1.created_at >= NOW() - INTERVAL '90 days'
      AND t1.status IN ('active', 'pending', 'review')
      AND (
        t1.category = 'important'
        OR t1.priority = 'high'
        OR t2.urgency_level > 5
      )
      AND NOT EXISTS (
        SELECT 1 FROM exclusion_table e
        WHERE e.item_id = t1.id
        AND e.reason = 'manual_exclusion'
      )
      AND t1.owner IN (
        SELECT user_id FROM authorized_users
        WHERE access_level >= 'read'
        AND department IN ('engineering', 'security', 'operations')
      )
    GROUP BY
      t1.id, t1.name, t1.description, t1.created_at, t1.updated_at,
      t1.status, t1.category, t1.priority, t1.owner, t1.team,
      t2.id, t2.name, t2.description, t2.created_at, t2.updated_at,
      t2.status, t2.category, t2.priority, t2.owner, t2.team,
      t3.id, t3.name, t3.description, t3.created_at, t3.updated_at,
      t3.status, t3.category, t3.priority, t3.owner, t3.team
    HAVING
      COUNT(*) > 1
      OR SUM(CASE WHEN t2.status = 'error' THEN 1 ELSE 0 END) > 0
    ORDER BY
      t1.priority DESC,
      t1.created_at DESC,
      t1.name ASC
    LIMIT 1000
    OFFSET 0;
    -- End of large SQL query
  EOQ
}

control "control_with_huge_sql" {
  title       = "Control with Large SQL"
  description = "Control using a very large inline SQL statement"
  sql         = <<-EOQ
    -- Large control SQL
    SELECT
      CASE
        WHEN (
          SELECT COUNT(*)
          FROM resources r
          WHERE r.compliant = true
          AND r.checked_at >= NOW() - INTERVAL '24 hours'
          AND r.resource_type IN ('type_a', 'type_b', 'type_c', 'type_d', 'type_e')
        ) = (
          SELECT COUNT(*)
          FROM resources r
          WHERE r.checked_at >= NOW() - INTERVAL '24 hours'
          AND r.resource_type IN ('type_a', 'type_b', 'type_c', 'type_d', 'type_e')
        )
        THEN 'pass'
        ELSE 'fail'
      END as status,
      'compliance_check' as resource,
      'Checked ' || (
        SELECT COUNT(*) FROM resources
        WHERE checked_at >= NOW() - INTERVAL '24 hours'
      )::text || ' resources' as reason
    FROM
      generate_series(1, 1) as dummy;
  EOQ
}
