resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_database" "test" {
  name = var.database_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}

provider "snowflake" {
  profile = "default"
  alias   = "secondary"
  role    = snowflake_account_role.test.name
}

resource "snowflake_schema" "test" {
  provider   = snowflake.secondary
  depends_on = [snowflake_grant_ownership.test]
  database   = snowflake_database.test.name
  name       = var.schema_name
}
