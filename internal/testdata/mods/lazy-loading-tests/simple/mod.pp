mod "lazy_simple" {
  title       = "Simple Lazy Loading Test"
  description = "Basic mod with clear, simple resources for baseline lazy loading testing"
}

locals {
  common_tags = {
    test     = "true"
    category = "simple"
  }
}
