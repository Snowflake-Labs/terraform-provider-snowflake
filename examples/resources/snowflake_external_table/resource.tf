resource "snowflake_external_table" "external_table" {
  database    = "db"
  schema      = "schema"
  name        = "external_table"
  comment     = "External table"
  file_format = "TYPE = CSV FIELD_DELIMITER = '|'"

  column {
    name = "id"
    type = "int"
  }

  column {
    name = "data"
    type = "text"
  }
}

# with a location pointing to an existing stage
# name is hardcoded, please see resource documentation for other options
resource "snowflake_external_table" "external_table_with_location" {
  database = "db"
  schema   = "schema"
  name     = "external_table_with_location"
  location = "@MYDB.MYSCHEMA.MYSTAGE"

  column {
    name = "id"
    type = "int"
  }
}