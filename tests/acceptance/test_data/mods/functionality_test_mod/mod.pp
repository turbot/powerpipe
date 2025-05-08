mod "functionality_test_mod"{
  title = "Functionality test mod"
  description = "This is a simple mod used for testing different steampipe features and funtionalities."
  database = var.database
}

variable "database" {
  type    = connection.postgres
  default = connection.postgres.my_connection
}