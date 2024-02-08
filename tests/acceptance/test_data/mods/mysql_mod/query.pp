query "total_employee" {
  sql = <<-EOQ
    SELECT 1;
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
    // SELECT COUNT(*) AS "Total Employees", CONCAT(? , ' ', ? , ' ', ?) as "Params" FROM employee;
    // SELECT COUNT(*) AS "Total Employees", CONCAT(@p1, ' ', @p2, ' ', @p3) as "Params" FROM employee;