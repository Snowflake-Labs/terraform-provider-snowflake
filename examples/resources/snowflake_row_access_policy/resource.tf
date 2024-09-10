resource "snowflake_row_access_policy" "example_row_access_policy" {
  name     = "EXAMPLE_ROW_ACCESS_POLICY"
  database = "EXAMPLE_DB"
  schema   = "EXAMPLE_SCHEMA"
  argument {
    name = "ARG1"
    type = "VARCHAR"
  }
  body    = "case when current_role() in ('ANALYST') then true else false end"
  comment = "comment"
}
