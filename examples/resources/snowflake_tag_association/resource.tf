resource "snowflake_database" "test" {
  name = "database"
}

resource "snowflake_schema" "test" {
  name     = "schema"
  database = snowflake_database.test.name
}

resource "snowflake_tag" "test" {
  name           = "cost_center"
  database       = snowflake_database.test.name
  schema         = snowflake_schema.test.name
  allowed_values = ["finance", "engineering"]
}

resource "snowflake_tag_association" "db_association" {
  object_identifiers = [snowflake_database.test.fully_qualified_name]
  object_type        = "DATABASE"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = "finance"
}

resource "snowflake_table" "test" {
  database = snowflake_database.test.name
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
  object_identifiers = [snowflake_table.test.fully_qualified_name]
  object_type        = "TABLE"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = "engineering"
}

resource "snowflake_tag_association" "column_association" {
  # For now, column fully qualified names have to be constructed manually.
  object_identifiers = [format("%s.\"column1\"", snowflake_table.test.fully_qualified_name)]
  object_type        = "COLUMN"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = "engineering"
}

resource "snowflake_tag_association" "account_association" {
  object_identifiers = ["\"ORGANIZATION_NAME\".\"ACCOUNT_NAME\""]
  object_type        = "ACCOUNT"
  tag_id             = snowflake_tag.test.fully_qualified_name
  tag_value          = "engineering"
}
