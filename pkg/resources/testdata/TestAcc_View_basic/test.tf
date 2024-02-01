resource "snowflake_view" "test" {
  name        = var.name
  comment     = var.comment
  database    = var.database
  schema      = var.schema
  is_secure   = var.is_secure
  or_replace  = var.or_replace
  copy_grants = var.copy_grants
  statement   = var.statement
}
