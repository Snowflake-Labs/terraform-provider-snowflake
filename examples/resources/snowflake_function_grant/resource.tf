resource snowflake_function_grant grant {
  database_name   = "db"
  schema_name     = "schema"
  function_name  = "function"

  arguments   = [
    {
      "name": "a",
      "type": "array"
    },
    {
      "name": "b",
      "type": "string"
    }
  ]
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
