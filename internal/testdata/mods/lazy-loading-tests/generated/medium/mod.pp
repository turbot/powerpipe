mod "lazy_medium" {
  title       = "Lazy Loading Medium Test"
  description = "Generated mod for lazy loading testing: 50 dashboards, 100 queries, 150 controls, 25 benchmarks"
}

# Common tags for all resources
locals {
  common_tags = {
    generator = "lazy_test"
    preset    = "medium"
    test      = "true"
  }

  severity_levels = ["low", "medium", "high", "critical"]
}

# Variables for parameterized resources
variable "test_region" {
  type        = string
  default     = "us-east-1"
  description = "Test region for lazy loading tests"
}

variable "test_threshold" {
  type    = number
  default = 100
}
