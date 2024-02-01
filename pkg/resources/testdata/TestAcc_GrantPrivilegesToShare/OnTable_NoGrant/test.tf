resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_table" "test" {
  name     = var.on_table
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_share" "test" {
  name = var.to_share
}
