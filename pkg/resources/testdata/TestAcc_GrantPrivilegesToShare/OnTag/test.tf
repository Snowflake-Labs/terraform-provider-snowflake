resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  name     = var.schema
  database = snowflake_database.test.name
}

resource "snowflake_tag" "test" {
  name     = var.tag_name
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
}

resource "snowflake_share" "test" {
  name = var.share_name
}

resource "snowflake_grant_privileges_to_share" "test_setup" {
  share_name    = snowflake_share.test.name
  privileges    = ["USAGE"]
  database_name = snowflake_database.test.name
}

resource "snowflake_grant_privileges_to_share" "test" {
  share_name = snowflake_share.test.name
  privileges = var.privileges
  tag_name   = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_tag.test.name}\""
}