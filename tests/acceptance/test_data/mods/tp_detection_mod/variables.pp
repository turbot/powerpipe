variable "database" {
  type        = connection.tailpipe
  description = "Tailpipe database connection string."
  default     = connection.tailpipe.default
}