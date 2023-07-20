provider "snowflake" {
  account                = "..." # required if not using profile. Can also be set via SNOWFLAKE_ACCOUNT env var
  username               = "..." # required if not using profile or token. Can also be set via SNOWFLAKE_USER env var
  password               = "..."
  oauth_access_token     = "..."
  private_key_path       = "..."
  private_key            = "..."
  private_key_passphrase = "..."
  oauth_refresh_token    = "..."
  oauth_client_id        = "..."
  oauth_client_secret    = "..."
  oauth_endpoint         = "..."
  oauth_redirect_url     = "..."

  // optional
  region    = "..." # required if using legacy format for account identifier
  role      = "..."
  host      = "..."
  warehouse = "..."
  session_params = {
    query_tag = "..."
  }
}


provider "snowflake" {
  profile = "securityadmin"
}
