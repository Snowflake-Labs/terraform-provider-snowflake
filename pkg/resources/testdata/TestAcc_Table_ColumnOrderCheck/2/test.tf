resource "snowflake_table" "test" {
  database = var.database
  schema   = var.schema
  name     = var.table

  column {
    type = "TIMESTAMP_NTZ(9)"
    name = "column_3"
  }
  column {
    type = "NUMBER(1,0)"
    name = "column_1"
  }
  column {
    type = "VARCHAR(7)"
    name = "column_5"
  }
  column {
    type = "VARCHAR(16777216)"
    name = "column_4"
  }
  column {
    type = "VARIANT"
    name = "column_2"
  }
}
