resource "snowflake_tag" "test" {
  name           = var.tag_name
  database       = var.database
  schema         = var.schema
  allowed_values = ["finance", "hr"]
  comment        = "Terraform acceptance test"
}

resource "snowflake_tag_association" "test" {
  object_identifier {
    name = var.database
  }
  object_type = "DATABASE"
  tag_id      = snowflake_tag.test.id
  tag_value   = "finance"
}
