# Default provider for most resources
provider "snowflake" {
  role = "SYSADMIN"
}

# Alternative provider with masking_admin role
provider "snowflake" {
  alias = "masking"
  role  = "MASKING_ADMIN"
}

resource "snowflake_masking_policy" "policy" {
  provider = snowflake.masking # Create masking policy with masking_admin role

  name               = "EXAMPLE_MASKING_POLICY"
  database           = "EXAMPLE_DB"
  schema             = "EXAMPLE_SCHEMA"
  value_data_type    = "VARCHAR"
  masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
  return_data_type   = "VARCHAR"
}

# Table is created by the default provider
resource "snowflake_table" "table" {
  database = "EXAMPLE_DB"
  schema   = "EXAMPLE_SCHEMA"
  name     = "table"

  column {
    name = "secret"
    type = "VARCHAR(16777216)"
  }

  lifecycle {
    # Masking policy is managed by a standalone resource and shouldn't be changed by the table resource.
    ignore_changes = [column[0].masking_policy]
  }
}

resource "snowflake_table_column_masking_policy_application" "application" {
  provider = snowflake.masking # Apply masking policy with masking_admin role

  table          = snowflake_table.table.fully_qualified_name
  column         = "secret"
  masking_policy = snowflake_masking_policy.policy.fully_qualified_name
}
