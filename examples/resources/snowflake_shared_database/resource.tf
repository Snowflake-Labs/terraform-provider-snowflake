# 1. Preparing database to share
resource "snowflake_share" "test" {
  provider = primary_account # notice the provider fields
  name     = "share_name"
  accounts = ["<secondary_account_organization_name>.<secondary_account_name>"]
}

resource "snowflake_database" "test" {
  provider = primary_account
  name     = "shared_database"
}

resource "snowflake_grant_privileges_to_share" "test" {
  provider    = primary_account
  to_share    = snowflake_share.test.name
  privileges  = ["USAGE"]
  on_database = snowflake_database.test.name
}

# 2. Creating shared database
resource "snowflake_shared_database" "test" {
  provider                     = secondary_account
  depends_on                   = [snowflake_grant_privileges_to_share.test]
  name                         = snowflake_database.test.name # shared database should have the same as the "imported" one
  from_share                   = "<primary_account_organization_name>.<primary_account_name>.${snowflake_share.test.name}"
  is_transient                 = false
  external_volume              = "external_volume_name"
  catalog                      = "catalog_name"
  replace_invalid_characters   = false
  default_ddl_collation        = "en_US"
  storage_serialization_policy = "OPTIMIZED"
  log_level                    = "OFF"
  trace_level                  = "OFF"
  comment                      = "A shared database"
}
