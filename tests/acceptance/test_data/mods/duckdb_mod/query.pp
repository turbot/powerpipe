query "total_employee" {
  sql = <<-EOQ
    SELECT COUNT(*) AS "Total Employees" FROM employee;
  EOQ
}

query "json_casting" {
  sql = <<-EOQ
    SELECT preferences::JSON FROM employee;
  EOQ
}

query "extensions" {
  sql = <<-EOQ
    select extension_name, installed from duckdb_extensions();
  EOQ
}

query "params_only" {
  sql = <<-EOQ
    SELECT CONCAT($1::text, ' ', $2::text, ' ', $3::text) as "Params" FROM employee;
  EOQ
  param "p1"{
    description = "First parameter"
    default = "default_parameter_1"
  }
  param "p2"{
    description = "Second parameter"
    default = "default_parameter_2"
  }
  param "p3"{
    description = "Third parameter"
    default = "default_parameter_3"
  }
}

dashboard "testing_card_blocks" {
  title = "Testing card blocks"

  container {
    card "card1" {
      query = query.total_employee
      width = 2
    }

    card "card2" {
      query = query.params_only
      width = 2
    }
  }
}
