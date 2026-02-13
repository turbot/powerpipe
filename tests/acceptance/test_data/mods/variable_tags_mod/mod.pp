mod "variable_tags_mod" {
  title = "Variable Tags Mod"
  description = "A mod with variable-based tags for testing lazy resolution"
}

variable "service_name" {
  type    = string
  default = "test_service_var"
}

variable "environment" {
  type    = string
  default = "testing"
}
