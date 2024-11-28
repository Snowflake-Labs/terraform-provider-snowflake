resource "snowflake_tag" "test" {
  name     = var.name
  schema   = var.schema
  database = var.database

  allowed_values = var.allowed_values

  comment = var.comment
}

data "snowflake_tags" "test" {
  depends_on = [snowflake_tag.test]

  like = var.name
}
