resource "snowflake_tag" "test" {
  name     = var.tag_name
  database = var.database
  schema   = var.schema
}

resource "snowflake_tag_association" "test" {
  object_identifiers = [var.account_fully_qualified_name]
  object_type        = "ACCOUNT"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = "v1"
}
