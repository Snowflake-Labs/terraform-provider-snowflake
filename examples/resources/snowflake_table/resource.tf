resource snowflake_table table {
  database   = "database"
  schema     = "schmea"
  name       = "table"
  comment    = "A table."
  cluster_by = ["to_date(DATE)"]
  
  owner      = "me"
  
  column {
    name = "id"
    type = "int"
  }

  column {
    name = "data"
    type = "text"
  }

  column {
    name = "DATE"
    type = "TIMESTAMP_NTZ(9)"
  }
}
