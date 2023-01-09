provider "snowflake" {
  // required
  username = "..."
  account  = "..." # the Snowflake account identifier

  // optional, exactly one must be set
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
}
