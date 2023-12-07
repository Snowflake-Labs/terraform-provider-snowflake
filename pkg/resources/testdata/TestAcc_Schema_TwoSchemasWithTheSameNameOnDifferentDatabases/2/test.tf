resource "snowflake_schema" "test" {
  name = var.name
  database = var.database
}

resource "snowflake_database" "test" {
  name = var.new_database
}

resource "snowflake_schema" "test_2" {
  name = var.name
  database = snowflake_database.test.name
}
