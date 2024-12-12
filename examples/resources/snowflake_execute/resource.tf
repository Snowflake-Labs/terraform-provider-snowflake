##################################
### simple use cases
##################################

# create and destroy resource
resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE ABC"
  revert  = "DROP DATABASE ABC"
}

# create and destroy resource using qualified name
resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE \"abc\""
  revert  = "DROP DATABASE \"abc\""
}

# with query
resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE ABC"
  revert  = "DROP DATABASE ABC"
  query   = "SHOW DATABASES LIKE '%ABC%'"
}

##################################
### grants example
##################################

# grant and revoke privilege USAGE to ROLE on database
resource "snowflake_execute" "test" {
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

resource "snowflake_execute" "test" {
  for_each = { for index, db_grant in var.database_grants : index => db_grant }
  execute  = "GRANT ${join(",", each.value.privileges)} ON DATABASE ${each.value.database_name} TO ROLE ${each.value.role_id}"
  revert   = "REVOKE ${join(",", each.value.privileges)} ON DATABASE ${each.value.database_name} FROM ROLE ${each.value.role_id}"
}

##################################
### fixing bad configuration
##################################

# bad revert
# 1 - resource created with a bad revert; it is constructed, revert is not validated before destroy happens
resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE ABC"
  revert  = "SELECT 1"
}

# 2 - fix the revert first; resource won't be recreated
resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE ABC"
  revert  = "DROP DATABASE ABC"
}

# bad query
# 1 - resource will be created; query_results will be empty
resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE ABC"
  revert  = "DROP DATABASE ABC"
  query   = "bad query"
}

# 2 - fix the query; query_results will be calculated; resource won't be recreated
resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE ABC"
  revert  = "DROP DATABASE ABC"
  query   = "SHOW DATABASES LIKE '%ABC%'"
}
