resource "snowflake_function" "f" {
  database        = var.database
  schema          = var.schema
  name            = var.name
  return_type     = "VARCHAR"
  return_behavior = "IMMUTABLE"
  statement       = "SELECT PARAM"

  arguments {
    name = "PARAM"
    type = "VARCHAR"
  }
}
