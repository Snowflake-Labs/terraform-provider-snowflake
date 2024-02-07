resource "snowflake_row_access_policy" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  signature = {
    V = "BOOLEAN",
    X = "TIMESTAMP_NTZ"
  }
  row_access_expression = "case when current_role() in ('ANALYST') then false else true end"
  comment               = "Terraform acceptance test - changed comment"
}
