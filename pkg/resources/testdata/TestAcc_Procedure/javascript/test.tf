resource "snowflake_procedure" "p" {
  database    = var.database
  schema      = var.schema
  name        = var.name
  language    = "JAVASCRIPT"
  return_type = "VARCHAR"
  comment     = var.comment
  statement   = <<EOT
    return "Hi"
  EOT
}
