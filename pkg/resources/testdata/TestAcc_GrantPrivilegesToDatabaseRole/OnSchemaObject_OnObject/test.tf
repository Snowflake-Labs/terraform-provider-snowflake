resource "snowflake_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.table_name

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "\"${var.database}\".\"${var.name}\""
  privileges         = var.privileges
  with_grant_option  = var.with_grant_option

  on_schema_object {
    object_type = "TABLE"
    object_name = snowflake_table.test.fully_qualified_name
  }
}
