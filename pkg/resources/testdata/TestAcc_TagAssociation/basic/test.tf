resource "snowflake_tag" "test" {
  name           = var.tag_name
  database       = var.database
  schema         = var.schema
  allowed_values = ["bar", "foo", "external"]
  comment        = "Terraform acceptance test"
}

resource "snowflake_tag_association" "test" {
  object_identifiers = [var.database_fully_qualified_name]
  object_type        = "DATABASE"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = var.tag_value
}
