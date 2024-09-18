# basic resource
resource "snowflake_masking_policy" "test" {
  name     = "EXAMPLE_MASKING_POLICY"
  database = "EXAMPLE_DB"
  schema   = "EXAMPLE_SCHEMA"
  argument {
    name = "ARG1"
    type = "VARCHAR"
  }
  argument {
    name = "ARG2"
    type = "NUMBER"
  }
  argument {
    name = "ARG3"
    type = "TIMESTAMP_NTZ"
  }
  body             = <<-EOF
  case
    when current_role() in ('ROLE_A') then
      ARG1
    when is_role_in_session( 'ROLE_B' ) then
      'ABC123'
    else
      '******'
  end
EOF
  return_data_type = "VARCHAR"
}

# resource with all fields set
resource "snowflake_masking_policy" "test" {
  name     = "EXAMPLE_MASKING_POLICY"
  database = "EXAMPLE_DB"
  schema   = "EXAMPLE_SCHEMA"
  argument {
    name = "ARG1"
    type = "VARCHAR"
  }
  argument {
    name = "ARG2"
    type = "NUMBER"
  }
  argument {
    name = "ARG3"
    type = "TIMESTAMP_NTZ"
  }
  body                  = <<-EOF
  case
    when current_role() in ('ROLE_A') then
      ARG1
    when is_role_in_session( 'ROLE_B' ) then
      'ABC123'
    else
      '******'
  end
EOF
  return_data_type      = "VARCHAR"
  exempt_other_policies = "true"
  comment               = "example masking policy"
}
