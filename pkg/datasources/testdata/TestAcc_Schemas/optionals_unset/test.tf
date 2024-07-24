resource "snowflake_schema" "test" {
  name         = var.name
  database     = var.database
  comment      = var.comment
  is_transient = true
  is_managed   = true
}

data "snowflake_schemas" "test" {
  with_describe   = false
  with_parameters = false
  depends_on      = [snowflake_schema.test]
  like            = var.name
  starts_with     = var.name
  limit {
    rows = 1
  }
}
