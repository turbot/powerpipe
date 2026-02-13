// Invalid resources mixed with valid ones

// Valid query
query "another_valid_query" {
  sql = "select 'valid' as status"
}

// Control referencing non-existent query
control "control_bad_ref" {
  title = "Control with bad reference"
  query = query.does_not_exist
}

// Another valid control
control "another_valid_control" {
  title = "Another Valid Control"
  sql = "select 'ok' as status, 'resource' as resource"
}
