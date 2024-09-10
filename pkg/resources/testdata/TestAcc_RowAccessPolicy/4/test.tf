resource "snowflake_row_access_policy" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  argument {
    name = "V"
    type = "BOOLEAN"
  }
  argument {
    name = "X"
    type = "TIMESTAMP_NTZ"
  }
  body    = "case when current_role() in ('ANALYST') then false else true end"
  comment = "Terraform acceptance test - changed comment"
}
