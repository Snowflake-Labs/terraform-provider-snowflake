resource "snowflake_tag" "test" {
  name           = var.tag_name
  database       = var.database
  schema         = var.schema
  allowed_values = []
}

resource "snowflake_tag_association" "test" {
  object_identifiers = [var.schema_fully_qualified_name]

  object_type = "SCHEMA"
  tag_id      = snowflake_tag.test.fully_qualified_name
  tag_value   = "TAG_VALUE"
}
