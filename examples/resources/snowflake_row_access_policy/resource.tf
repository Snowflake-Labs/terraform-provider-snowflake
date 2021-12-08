resource "snowflake_row_access_policy" "example_row_access_policy" {
  name               = "EXAMPLE_ROW_ACCESS_POLICY"
  database           = "EXAMPLE_DB"
  schema             = "EXAMPLE_SCHEMA"
  signature          = {
    A = "VARCHAR",
    B = "VARCHAR"
  }
  row_access_expression = "case when current_role() in ('ANALYST') then true else false end"
}
