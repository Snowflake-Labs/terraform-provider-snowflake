resource "snowflake_masking_policy" "test" {
 name               = "EXAMPLE_MASKING_POLICY"
  database           = "EXAMPLE_DB"
  schema             = "EXAMPLE_SCHEMA"
	signature {
		column {
			name = "val"
			type = "VARCHAR"
		}
	}
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR"
}
