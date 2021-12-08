resource "snowflake_masking_policy" "example_masking_policy" {
  name               = "EXAMPLE_MASKING_POLICY"
  database           = "EXAMPLE_DB"
  schema             = "EXAMPLE_SCHEMA"
  value_data_type    = "string"
  masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
  return_data_type   = "string"
}
