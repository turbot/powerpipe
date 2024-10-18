query "pipes_workspace_query" {
  sql = <<-EOQ
    select account_aliases from all_aws.aws_account;
  EOQ
}
