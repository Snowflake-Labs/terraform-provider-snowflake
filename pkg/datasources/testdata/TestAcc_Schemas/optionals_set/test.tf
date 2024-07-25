resource "snowflake_schema" "test" {
  name         = var.name
  database     = var.database
  comment      = var.comment
  is_transient = true
  is_managed   = true
}

resource "snowflake_table" "test" {
  database = var.database
  schema   = snowflake_schema.test.name
  name     = "table"

  column {
    name = "id"
    type = "int"
  }
}

data "snowflake_schemas" "test" {
  depends_on  = [snowflake_table.test]
  like        = var.name
  starts_with = var.name
  limit {
    rows = 1
  }
}
