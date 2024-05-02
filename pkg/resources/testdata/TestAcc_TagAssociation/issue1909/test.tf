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
    name = var.column_name
    type = "VARIANT"
  }
}
resource "snowflake_table" "test2" {
  name     = var.table_name2
  database = var.database
  schema   = var.schema
  column {
    name = var.column_name
    type = "VARIANT"
  }
}

resource "snowflake_tag_association" "test" {
  object_identifier {
    database = var.database
    schema   = var.schema
    name     = "${snowflake_table.test.name}.${snowflake_table.test.column[0].name}"
  }
  object_identifier {
    database = var.database
    schema   = var.schema
    name     = "${snowflake_table.test2.name}.${snowflake_table.test2.column[0].name}"
  }
  object_type = "COLUMN"
  tag_id      = snowflake_tag.test.id
  tag_value   = "v1"
}
