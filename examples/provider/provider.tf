provider snowflake {
  // required
  username = "..."
  account  = "..."
  region   = "..."

  // optional, at exactly one must be set
  password           = "..."
  oauth_access_token = "..."
  private_key_path   = "..."

  // optional
  role = "..."
}
