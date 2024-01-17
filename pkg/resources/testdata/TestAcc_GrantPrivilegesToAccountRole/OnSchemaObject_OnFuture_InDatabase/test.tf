resource "snowflake_grant_privileges_to_account_role" "test" {
  role_name = var.name
  privileges         = var.privileges
  with_grant_option  = var.with_grant_option

  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_database        = var.database
    }
  }
}
