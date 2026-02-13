mod "lazy_large" {
  title       = "Lazy Loading Large Test"
  description = "Generated mod for lazy loading testing: 200 dashboards, 400 queries, 600 controls, 75 benchmarks"
}

# Common tags for all resources
locals {
  common_tags = {
    generator = "lazy_test"
    preset    = "large"
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
