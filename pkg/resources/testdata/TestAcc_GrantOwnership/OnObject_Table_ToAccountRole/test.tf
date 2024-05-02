resource "snowflake_role" "test" {
  name = var.account_role_name
}

resource "snowflake_database" "test" {
  name = var.database_name
}

resource "snowflake_schema" "test" {
  name     = var.schema_name
  database = snowflake_database.test.name
}

resource "snowflake_table" "test" {
  name     = var.table_name
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    object_type = "TABLE"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_table.test.name}\""
  }
}
