resource "snowflake_tag" "test" {
  name     = var.tag_name
  database = var.database
  schema   = var.schema
}

resource "snowflake_table" "test" {
  name     = var.table_name
  database = var.database
  schema   = var.schema
  column {
    name = "test_column"
    type = "VARIANT"
  }
}

resource "snowflake_tag_association" "test" {
  // we need to set the object_identifier to avoid the following error:
  // provider_test.go:17: err: resource snowflake_tag_association: object_identifier: Optional or Required must be set, not both
  // we should consider deprecating object_identifier in favor of object_name
  // https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2534#discussion_r1507570740
  // object_name = "\"${var.database}\".\"${var.schema}\".\"${var.table_name}\""
  object_identifier {
    database = var.database
    schema   = var.schema
    name     = snowflake_table.test.name
  }
  object_type = "TABLE"
  tag_id      = snowflake_tag.test.id
  tag_value   = "v1"
}
