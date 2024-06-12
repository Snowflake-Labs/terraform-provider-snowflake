resource "snowflake_database" "database" {
  name = "database"
}

resource "snowflake_schema" "schema" {
  name     = "schema"
  database = snowflake_database.database.name
}

resource "snowflake_tag" "tag" {
  name           = "cost_center"
  database       = snowflake_database.database.name
  schema         = snowflake_schema.schema.name
  allowed_values = ["finance", "engineering"]
}
