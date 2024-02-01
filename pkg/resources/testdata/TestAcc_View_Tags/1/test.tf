resource "snowflake_tag" "test_tag" {
  name     = "tag1"
  database = var.database
  schema   = var.schema
}

resource "snowflake_tag" "test_tag_2" {
  name     = "tag2"
  database = var.database
  schema   = var.schema
}

resource "snowflake_view" "test" {
  name        = var.name
  database    = var.database
  schema      = var.schema
  is_secure   = false
  or_replace  = false
  copy_grants = false
  statement   = var.statement

  tag {
    name     = snowflake_tag.test_tag.name
    schema   = var.schema
    database = var.database
    value    = "some_value"
  }
}
