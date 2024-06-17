resource "snowflake_database" "test" {
  name = "things"
}

resource "snowflake_schema" "test_schema" {
  name     = "things"
  database = snowflake_database.test.name
}

resource "snowflake_sequence" "test_sequence" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test_schema.name
  name     = "thing_counter"
}
