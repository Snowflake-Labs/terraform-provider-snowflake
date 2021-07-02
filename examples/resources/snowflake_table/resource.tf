resource snowflake_table table {
  database    = "database"
  schema      = "schema"
  name        = "table"
  comment     = "A table."
  cluster_by  = ["to_date(DATE)"]
  primary_key = ["\"data\""]
  owner       = "me"
  
  column {
    name     = "id"
    type     = "int"
    nullable = true
  }

  column {
    name     = "data"
    type     = "text"
    nullable = false
  }

  column {
    name = "DATE"
    type = "TIMESTAMP_NTZ(9)"
  }
}
