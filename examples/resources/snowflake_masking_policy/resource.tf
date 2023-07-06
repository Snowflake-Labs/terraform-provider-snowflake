resource "snowflake_masking_policy" "test" {
  name     = "EXAMPLE_MASKING_POLICY"
  database = "EXAMPLE_DB"
  schema   = "EXAMPLE_SCHEMA"
  signature {
    column {
      name = "val"
      type = "VARCHAR"
    }
  }
  masking_expression = <<-EOF
    case 
      when current_role() in ('ROLE_A') then 
        val 
      when is_role_in_session( 'ROLE_B' ) then 
        'ABC123'
      else
        '******'
    end
  EOF

  return_data_type = "VARCHAR"
}
