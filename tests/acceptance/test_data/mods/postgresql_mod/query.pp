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

// query bigint
query "total" {
  sql = <<-EOQ
    SELECT COUNT(*) as total FROM employees;
  EOQ
}

// query string
query "name" {
  sql = <<-EOQ
    SELECT name FROM employees;
  EOQ
}

// query small int
query "age" {
  sql = <<-EOQ
    SELECT age FROM employees;
  EOQ
}

// query float
query "salary" {
  sql = <<-EOQ
    SELECT salary FROM employees;
  EOQ
}

query "search_path" {
  sql = <<-EOQ
    SHOW search_path;
  EOQ
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