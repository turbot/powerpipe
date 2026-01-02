mod "left" {
  title   = "Left Branch"
  version = "1.0.0"

  require {
    mod "github.com/test/shared" {
      version = "v1.0.0"
    }
  }
}

query "left_query" {
  title = "Left Query"
  sql   = "SELECT 'left' as source"
}

control "left_control" {
  title = "Left Control"
  query = shared.query.shared_query
}
