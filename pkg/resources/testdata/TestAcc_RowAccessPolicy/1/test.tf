resource "snowflake_row_access_policy" "test" {
  name     = var.name
  database = var.database
  schema   = var.schema
  argument {
    name = "N"
    type = "VARCHAR"
  }
  argument {
    name = "V"
    type = "VARCHAR"
  }
  body    = "case when current_role() in ('ANALYST') then true else false end"
  comment = "Terraform acceptance test"
}
