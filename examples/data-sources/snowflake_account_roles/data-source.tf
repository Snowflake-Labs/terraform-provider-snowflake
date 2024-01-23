data "snowflake_account_roles" "all" {
}

data "snowflake_account_roles" "by_pattern" {
  pattern = "some_prefix_%"
}
