resource "snowflake_account_role" "test" {
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
  database = snowflake_database.test.name
  name     = var.table_name
  schema   = snowflake_schema.test.name

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_materialized_view" "test" {
  database  = snowflake_database.test.name
  name      = var.materialized_view_name
  schema    = snowflake_schema.test.name
  statement = "select * from \"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_table.test.name}\""
  warehouse = var.warehouse_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "MATERIALIZED VIEW"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_materialized_view.test.name}\""
  }
}
