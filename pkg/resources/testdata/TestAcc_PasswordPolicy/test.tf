resource "snowflake_password_policy" "pa" {
  name       = var.name
  database   = var.database
  schema     = var.schema
  min_length = var.min_length
  max_length = var.max_length
  comment    = var.comment
  or_replace = true
}