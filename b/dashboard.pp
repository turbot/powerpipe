

detection_benchmark "db1" {

  title         = "Detection testing"
  children = [detection.d1]

 }


detection "d1"{
         title = "top level detection"
         sql = "select 'r' as reason, 'foo' as resource, 'alarm' as status"
   }
query "query_input" {
  param "new_input" {

  }
  sql = <<-EOQ
    select
      $1 as "column 1",
      'value1' as "column 2"
  EOQ
}