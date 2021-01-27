resource snowflake_external_table external_table {
  database = "db"
  schema   = "schema"
  name     = "external_table"
  comment  = "External table"

  column {
    name = "id"
    type = "int"
  }

  column {
    name = "data"
    type = "text"
  }
}
