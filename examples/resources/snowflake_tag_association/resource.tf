resource "snowflake_database" "database" {
  name    = "db1"
}

resource "snowflake_schema" "schema" {
  name     = "schema1"
  database = snowflake_database.database.name
}

resource "snowflake_tag" "tag" {
  name           = "cost_center"
  database       = snowflake_database.database.name
  schema         = snowflake_schema.schema.name
  allowed_values = ["finance", "engineering"]
}

resource "snowflake_tag_association" "association" {
  object_name     = snowflake_database.database.name
  object_type     = "DATABASE"
  tag_id          = snowflake_tag.tag.id
  tag_value       = "finance"
  skip_validation = true
}
