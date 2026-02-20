mod "right" {
  title   = "Right Branch"
  version = "1.0.0"

  require {
    mod "github.com/test/shared" {
      version = "v1.0.0"
    }
  }
}

query "right_query" {
  title = "Right Query"
  sql   = "SELECT 'right' as source"
}

control "right_control" {
  title = "Right Control"
  query = shared.query.shared_query
}
