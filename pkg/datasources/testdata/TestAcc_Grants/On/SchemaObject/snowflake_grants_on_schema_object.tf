resource "snowflake_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.table

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

data "snowflake_grants" "test" {
  grants_on {
    object_name = snowflake_table.test.fully_qualified_name
    object_type = "TABLE"
  }
}
