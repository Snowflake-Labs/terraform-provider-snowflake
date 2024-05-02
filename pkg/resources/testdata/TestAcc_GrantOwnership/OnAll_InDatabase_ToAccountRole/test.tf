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

resource "snowflake_table" "test2" {
  name     = var.second_table_name
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_grant_ownership" "test" {
  depends_on        = [snowflake_table.test, snowflake_table.test2]
  account_role_name = snowflake_role.test.name
  on {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.test.name
    }
  }
}
