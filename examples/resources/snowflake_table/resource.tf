resource "snowflake_schema" "schema" {
  database            = "database"
  name                = "schema"
  data_retention_days = 1
}

resource "snowflake_table" "table" {
  database            = snowflake_schema.schema.database
  schema              = snowflake_schema.schema.name
  name                = "table"
  comment             = "A table."
  cluster_by          = ["to_date(DATE)"]
  data_retention_days = snowflake_schema.schema.data_retention_days
  change_tracking     = false

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

  column {
    name    = "extra"
    type    = "VARIANT"
    comment = "extra data"
  }

  primary_key {
    name = "my_key"
    keys = ["data"]
  }
}
