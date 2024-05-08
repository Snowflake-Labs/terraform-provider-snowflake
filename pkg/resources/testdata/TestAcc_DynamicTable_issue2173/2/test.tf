resource "snowflake_schema" "other_schema" {
  database = var.database
  name     = var.other_schema
  comment  = "Other schema"
}

resource "snowflake_table" "t" {
  database        = var.database
  schema          = var.schema
  name            = var.table_name
  change_tracking = true
  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_dynamic_table" "dt" {
  depends_on = [snowflake_table.t]
  name       = var.name
  database   = var.database
  schema     = var.schema
  target_lag {
    maximum_duration = "2 minutes"
  }
  warehouse = var.warehouse
  query     = var.query
  comment   = var.comment
}
