resource "snowflake_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.table_name

  column {
    name = "id"
    type = "int"
  }
}

resource "snowflake_grant_privileges_to_database_role" "test" {
  database_role_name = "\"${var.database}\".\"${var.name}\""
  privileges = var.privileges
  with_grant_option = var.with_grant_option

  on_schema_object {
    object_type = "TABLE"
    object_name = "\"${var.database}\".\"${var.schema}\".\"${var.table_name}\""
  }
}
