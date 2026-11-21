mod "server_reload" {
  title    = "Server reload"
  database = connection.sqlite.server_reload
}

variable "who" {
  type    = string
  default = "first"
}

control "who" {
  title = "Who"
  sql   = "select 'who' as resource, 'ok' as status, 'who is ' || $1 as reason"
  args  = [var.who]
}

benchmark "who" {
  title    = "Who"
  children = [control.who]
}
