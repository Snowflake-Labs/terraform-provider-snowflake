terraform {
  required_providers {
    snowflake = {
      source = "Snowflake-Labs/snowflake"
    }
  }
}

# A simple configuration of the provider.
provider "snowflake" {
  account                = "..." # required if not using profile. Can also be set via SNOWFLAKE_ACCOUNT env var
  user                   = "..." # required if not using profile or token. Can also be set via SNOWFLAKE_USER env var
  password               = "..."
  authenticator          = "..." # required if not using password as auth method
  private_key            = "..."
  private_key_passphrase = "..."

  // optional
  role      = "..."
  host      = "..."
  warehouse = "..."
  params = {
    query_tag = "..."
  }
}

# Use profile field only. In this case, the fields are populated from ~/.snowflake/config TOML file.
provider "snowflake" {
  profile = "securityadmin"
}
