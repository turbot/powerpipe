mod "lazy_cross_refs" {
  title       = "Cross References Lazy Loading Test"
  description = "Test mod with resources referencing each other in complex patterns"
}

# Variables for reference testing
variable "default_region" {
  type        = string
  default     = "us-east-1"
  description = "Default region for cross-reference testing"
}

variable "threshold" {
  type    = number
  default = 100
}

# Locals that reference variables
locals {
  region_prefix = "region:${var.default_region}"

  common_tags = {
    test   = "cross-refs"
    region = var.default_region
  }

  sql_template = "SELECT * FROM data WHERE region = '${var.default_region}'"
}
