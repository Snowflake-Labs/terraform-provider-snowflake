resource "snowflake_schema" "test" {
  database = var.database
  name     = var.schema
  data_retention_days = var.schema_data_retention_time
}

resource "snowflake_table" "test" {
  depends_on = [snowflake_schema.test]
  database = var.database
  schema     = var.schema
  name     = var.table
  data_retention_time_in_days = var.table_data_retention_time

  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}
