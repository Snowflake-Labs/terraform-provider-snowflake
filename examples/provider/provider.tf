provider "snowflake" {
  // required
  username = "..."
  account  = "..."
  region   = "..."

  // optional, at exactly one must be set
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
  role = "..."
  host = "..."
}
