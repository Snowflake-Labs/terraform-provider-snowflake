resource "snowflake_database" "test" {
  name = var.database
}

resource "snowflake_schema" "test" {
  database = snowflake_database.test.name
  name = var.schema
}

resource "snowflake_table" "test" {
  database        = snowflake_database.test.name
  schema          = snowflake_schema.test.name
  name            = var.table
  change_tracking = true
  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }

  column {
    name = "SOME_TEXT"
    type = "VARCHAR"
  }
}

resource "snowflake_cortex_search_service" "test" {
  depends_on = [snowflake_table.test]

  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name       = var.name
  on         = var.on
  target_lag = "2 minutes"
  warehouse  = var.warehouse
  query      = var.query
  comment    = var.comment
}

data "snowflake_cortex_search_services" "test" {
  like = snowflake_cortex_search_service.test.name
  in {
    database = snowflake_cortex_search_service.test.database
  }
  starts_with = snowflake_cortex_search_service.test.name
  limit {
    rows = 1
  }
}
