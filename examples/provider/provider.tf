terraform {
  required_providers {
    snowflake = {
      source = "Snowflake-Labs/snowflake"
    }
  }
}

# A simple configuration of the provider with a default authentication.
# A default value for `authenticator` is `snowflake`, enabling authentication with `user` and `password`.
provider "snowflake" {
  organization_name = "..." # required if not using profile. Can also be set via SNOWFLAKE_ORGANIZATION_NAME env var
  account_name      = "..." # required if not using profile. Can also be set via SNOWFLAKE_ACCOUNT_NAME env var
  user              = "..." # required if not using profile or token. Can also be set via SNOWFLAKE_USER env var
  password          = "..."

  // optional
  role      = "..."
  host      = "..."
  warehouse = "..."
  params = {
    query_tag = "..."
  }
}

# A simple configuration of the provider with private key authentication.
provider "snowflake" {
  organization_name      = "..." # required if not using profile. Can also be set via SNOWFLAKE_ORGANIZATION_NAME env var
  account_name           = "..." # required if not using profile. Can also be set via SNOWFLAKE_ACCOUNT_NAME env var
  user                   = "..." # required if not using profile or token. Can also be set via SNOWFLAKE_USER env var
  authenticator          = "SNOWFLAKE_JWT"
  private_key            = "-----BEGIN ENCRYPTED PRIVATE KEY-----..."
  private_key_passphrase = "passphrase"
}

# By using the `profile` field, missing fields will be populated from ~/.snowflake/config TOML file
provider "snowflake" {
  profile = "securityadmin"
}
