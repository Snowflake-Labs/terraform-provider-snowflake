resource snowflake_procedure_grant grant {
  database_name   = "db"
  schema_name     = "schema"
  procedure_name  = "procedure"

  arguments {
    name = "a"
    type = "array"
  }
  arguments {
    name = "b"
    type = "string"
  }
  return_type = "string"

  privilege = "select"
  roles = [
    "role1",
    "role2",
  ]

  shares = [
    "share1",
    "share2",
  ]

  on_future         = false
  with_grant_option = false
}
