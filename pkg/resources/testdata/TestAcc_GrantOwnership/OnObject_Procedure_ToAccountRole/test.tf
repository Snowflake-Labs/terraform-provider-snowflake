resource "snowflake_database" "test" {
  name = var.database_name
}

resource "snowflake_schema" "test" {
  name     = var.schema_name
  database = snowflake_database.test.name
}

resource "snowflake_role" "test" {
  name = var.account_role_name
}

# Without this, DESC PROCEDURE cannot get the procedure body, and the plan will always generate a difference.
resource "snowflake_grant_account_role" "test" {
  role_name = snowflake_role.test.name
  parent_role_name = "SYSADMIN"
}

resource "snowflake_procedure" "with_arguments" {
  name     = var.with_arg_procedure_name
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  language = "JAVASCRIPT"
  arguments {
    name = "ARG1"
    type = "VARCHAR"
  }
  return_type = "VARCHAR"
  statement   = "return ARG1"
}

resource "snowflake_grant_ownership" "with_arguments" {
  account_role_name   = snowflake_role.test.name
  on {
    object_type = "PROCEDURE"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_procedure.with_arguments.name}\"(VARCHAR)"
  }
}

resource "snowflake_procedure" "without_arguments" {
  name        = var.without_arg_procedure_name
  database    = snowflake_database.test.name
  schema      = snowflake_schema.test.name
  language    = "JAVASCRIPT"
  return_type = "VARCHAR"
  statement   = "return 'Hi'"
}

resource "snowflake_grant_ownership" "without_arguments" {
  account_role_name   = snowflake_role.test.name
  on {
    object_type = "PROCEDURE"
    object_name = "\"${snowflake_database.test.name}\".\"${snowflake_schema.test.name}\".\"${snowflake_procedure.without_arguments.name}\"()"
  }
}
