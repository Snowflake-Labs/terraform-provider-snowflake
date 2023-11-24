# create and destroy resource
resource "snowflake_unsafe_execute" "test" {
  execute = "CREATE DATABASE ABC"
  revert  = "DROP DATABASE ABC"
}

# create and destroy resource using qualified name
resource "snowflake_unsafe_execute" "test" {
  execute = "CREATE DATABASE \"abc\""
  revert  = "DROP DATABASE \"abc\""
}

# grant and revoke privilege USAGE to ROLE on database
resource "snowflake_unsafe_execute" "test" {
  execute = "GRANT USAGE ON DATABASE ABC TO ROLE XYZ"
  revert  = "REVOKE USAGE ON DATABASE ABC FROM ROLE XYZ"
}

# grant and revoke with for_each
variable "database_grants" {
  type = list(object({
    database_name = string
    role_id       = string
    privileges    = list(string)
  }))
}

resource "snowflake_unsafe_execute" "test" {
  for_each = { for index, db_grant in var.database_grants : index => db_grant }
  execute  = "GRANT ${join(",", each.value.privileges)} ON DATABASE ${each.value.database_name} TO ROLE ${each.value.role_id}"
  revert   = "REVOKE ${join(",", each.value.privileges)} ON DATABASE ${each.value.database_name} FROM ROLE ${each.value.role_id}"
}
