resource "snowflake_share" "test" {
  name = var.to_share
}

resource "snowflake_database" "test" {
  depends_on = [snowflake_share.test]
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_tag" "test" {
  name     = var.on_tag
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  to_share    = snowflake_share.test.name
  privileges  = ["USAGE"]
  on_database = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_share" "test" {
  to_share   = snowflake_share.test.name
  privileges = var.privileges
  on_tag     = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_tag.test.name}\""
}