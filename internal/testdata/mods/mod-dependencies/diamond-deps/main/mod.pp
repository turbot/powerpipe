mod "diamond_main" {
  title       = "Diamond Main"
  description = "Main mod for diamond dependency tests (main -> left, right; left -> shared; right -> shared)"

  require {
    mod "github.com/test/left" {
      version = "v1.0.0"
    }
    mod "github.com/test/right" {
      version = "v1.0.0"
    }
  }
}

query "main_query" {
  title = "Main Query"
  sql   = "SELECT 'diamond_main' as source"
}

# Control using shared dependency through left
control "via_left" {
  title = "Via Left"
  query = shared.query.shared_query
}

# Control using shared dependency through right
control "via_right" {
  title = "Via Right"
  query = shared.query.shared_query
}

# Benchmark including controls from all paths
benchmark "diamond_benchmark" {
  title = "Diamond Benchmark"
  children = [
    control.via_left,
    control.via_right,
    left.control.left_control,
    right.control.right_control
  ]
}
