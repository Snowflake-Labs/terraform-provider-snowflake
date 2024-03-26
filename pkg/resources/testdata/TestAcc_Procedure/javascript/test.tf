resource "snowflake_procedure" "p" {
  database    = var.database
  schema      = var.schema
  name        = var.name
  language    = "JAVASCRIPT"
  return_type = "VARCHAR"
  comment     = var.comment
  execute_as      = var.execute_as
  statement   = <<EOT
    return "Hi"
  EOT
}
