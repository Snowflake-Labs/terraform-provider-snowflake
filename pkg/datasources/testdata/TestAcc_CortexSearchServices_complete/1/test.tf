resource "snowflake_table" "t" {
  database        = var.database
  schema          = var.schema
  name            = var.table_name
  change_tracking = true
  column {
    name = "id"
    type = "NUMBER(38,0)"
  }
}

resource "snowflake_cortex_search_service" "css" {
  depends_on = [snowflake_table.t]
  name       = var.name
  on         = var.on
  database   = var.database
  schema     = var.schema
  target_lag = "2 minutes"
  warehouse  = var.warehouse
  query      = var.query
  comment    = var.comment
}

data "snowflake_cortex_search_services" "csss" {
  like {
    pattern = snowflake_cortex_search_service.css.name
  }
  in {
    database = snowflake_cortex_search_service.css.database
  }
  starts_with = snowflake_cortex_search_service.css.name
  limit {
    rows = 1
  }
}
