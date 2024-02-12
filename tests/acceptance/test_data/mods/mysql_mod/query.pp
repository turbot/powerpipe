query "total_employee" {
  sql = <<-EOQ
    SELECT ? as "Params", ? as "Params2";
  EOQ
  param "p1"{
    description = "First parameter"
    default = "default_parameter_1"
  }
  param "p2"{
    description = "Second parameter"
    default = "default_parameter_2"
  }
}
// query "total_employee" {
//   sql = <<-EOQ
//     SELECT COUNT(*) AS "Total Employees", CONCAT(? , ' ', ? , ' ', ?) as "Params" FROM employee;
//   EOQ
//   param "p1"{
//     description = "First parameter"
//     default = "default_parameter_1"
//   }
//   param "p2"{
//     description = "Second parameter"
//     default = "default_parameter_2"
//   }
//   param "p3"{
//     description = "Third parameter"
//     default = "default_parameter_3"
//   }
// }