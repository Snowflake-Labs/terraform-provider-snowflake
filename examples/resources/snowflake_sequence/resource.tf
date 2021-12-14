resource "snowflake_database" "database" {
	name    = "things"
}

resource "snowflake_schema" "test_schema" {
	name     = "things"
	database = snowflake_database.test_database.name
}

resource "snowflake_sequence" "test_sequence" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "thing_counter"
}
