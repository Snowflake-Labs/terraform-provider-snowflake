resource "snowflake_database" "test" {
  name                        = var.database
  data_retention_time_in_days = var.database_data_retention_time
}

resource "snowflake_schema" "test" {
  database                    = snowflake_database.test.name
  name                        = var.schema
  data_retention_time_in_days = var.schema_data_retention_time
}

resource "snowflake_table" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = var.table

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}
