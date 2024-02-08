query "total_employee" {
  sql = <<-EOQ
    SELECT COUNT(*) AS "Total Employees", CONCAT($1::text, ' ', $2::text, ' ', $3::text) as "Params" FROM employee;
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