resource "snowflake_tag" "test" {
  name           = var.tag_name
  database       = var.database
  schema         = var.schema
  allowed_values = []
}

resource "snowflake_tag_association" "test" {
  object_identifier {
    database = var.database
    name     = var.schema
  }

  object_type = "SCHEMA"
  tag_id      = snowflake_tag.test.id
  tag_value   = "TAG_VALUE"
}
