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

resource "snowflake_external_table" "delta_external_table" {
  database     = "db"
  schema       = "schema"
  name         = "delta_external_table"
  comment      = "External table with Delta format"
  file_format  = "TYPE = PARQUET"
  table_format = "DELTA"
  location     = "@stage/delta_table"
  refresh_on_create = false
  auto_refresh      = false

  column {
    name = "id"
    type = "int"
    as   = "value:id::int"
  }

  column {
    name = "data"
    type = "text"
    as   = "value:data::string"
  }
}
