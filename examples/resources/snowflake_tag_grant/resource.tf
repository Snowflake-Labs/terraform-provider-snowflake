resource "snowflake_tag_grant" "example" {
  database_name = "database"
  schema_name   = "schema"
  tag_name      = "tag"
  roles         = ["TEST_ROLE"]
  privilege     = "OWNERSHIP"

}
