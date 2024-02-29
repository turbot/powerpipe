query "generic_query_with_dimensions" {
  description = "parameterized query to simulate control results, with rows conataining all possible statuses(with extra dimensions)"
  sql = <<-EOQ
select num as id, 
    case 
        when (num<=$1) then 'ok' 
        when (num>$1 and num<=$1+$2) then 'alarm'
        when (num>$1+$2 and num<=$1+$2+$3) then 'error' 
        when (num>$1+$2+$3 and num<=$1+$2+$3+$4) then 'skip' 
        when (num>$1+$2+$3+$4 and num<=$1+$2+$3+$4+$5) then 'info' 
    end status, 
    'steampipe' as resource, 
    case 
        when (num<=$1) then 'Resource satisfies condition' 
        when (num>$1 and num<=$1+$2) then 'Resource does not satisfy condition' 
        when (num>$1+$2 and num<=$1+$2+$3) then 'Resource has some error' 
        when (num>$1+$2+$3 and num<=$1+$2+$3+$4) then 'Resource is skipped' 
        when (num>$1+$2+$3+$4 and num<=$1+$2+$3+$4+$5) then 'Information' 
    end reason,
'0.1.0' as version,
'xyz' as module
from generate_series(1, ($1::int+$2::int+$3::int+$4::int+$5::int)) num
EOQ
  param "number_of_ok" {
    description = "Number of resources in OK"
    default = 0
  }
  param "number_of_alarm" {
    description = "Number of resources in ALARM"
    default = 0
  }
  param "number_of_error" {
    description = "Number of resources in ERROR"
    default = 0
  }
  param "number_of_skip" {
    description = "Number of resources in SKIP"
    default = 0
  }
  param "number_of_info" {
    description = "Number of resources in INFO"
    default = 0
  }
}
