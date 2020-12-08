resource snowflake_table table {
  database = "database"
  schema   = "schmea"
  name     = "table"
  comment  = "A table."
  owner    = "me"

  column {
    name = "id"
    type = "int"
  }

  column {
    name = "data"
    type = "text"
  }
}
