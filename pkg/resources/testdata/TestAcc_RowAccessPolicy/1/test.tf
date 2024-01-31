resource "snowflake_row_access_policy" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  signature = {
    N = "VARCHAR"
    V = "VARCHAR",
  }
  row_access_expression = "case when current_role() in ('ANALYST') then true else false end"
  comment               = "Terraform acceptance test"
}
