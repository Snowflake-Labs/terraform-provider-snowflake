resource "snowflake_standard_database" "test" {
  name = "database"
}

resource "snowflake_schema" "test" {
  name     = "schema"
  database = snowflake_standard_database.test.name
}

resource "snowflake_tag" "test" {
  name           = "cost_center"
  database       = snowflake_standard_database.test.name
  schema         = snowflake_schema.test.name
  allowed_values = ["finance", "engineering"]
}

resource "snowflake_tag_association" "db_association" {
  object_identifier {
    name = snowflake_standard_database.test.name
  }
  object_type = "DATABASE"
  tag_id      = snowflake_tag.test.id
  tag_value   = "finance"
}

resource "snowflake_table" "test" {
  database = snowflake_standard_database.test.name
  schema   = snowflake_schema.test.name
  name     = "TABLE_NAME"
  comment  = "Terraform example table"
  column {
    name = "column1"
    type = "VARIANT"
  }
  column {
    name = "column2"
    type = "VARCHAR(16)"
  }
}

resource "snowflake_tag_association" "table_association" {
  object_identifier {
    name     = snowflake_table.test.name
    database = snowflake_standard_database.test.name
    schema   = snowflake_schema.test.name
  }
  object_type = "TABLE"
  tag_id      = snowflake_tag.test.id
  tag_value   = "engineering"
}

resource "snowflake_tag_association" "column_association" {
  object_identifier {
    name     = "${snowflake_table.test.name}.column_name"
    database = snowflake_standard_database.test.name
    schema   = snowflake_schema.test.name
  }
  object_type = "COLUMN"
  tag_id      = snowflake_tag.test.id
  tag_value   = "engineering"
}
