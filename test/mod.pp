mod "local" {
  title       = "New mod"
  description = "This is a simple mod used for testing different steampipe features and funtionalities."
  require {
    mod "github.com/turbot/tailpipe-mod-aws-cloudtrail-log-detections" {
      version = "*"
    }
  }
}
query "query1" {
  sql = "select 1"
}