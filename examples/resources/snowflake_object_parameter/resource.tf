resource "snowflake_database" "d" {
  name = "TEST_DB"
}

resource "snowflake_object_parameter" "o" {
  key         = "SUSPEND_TASK_AFTER_NUM_FAILURES"
  value       = "33"
  object_type = "DATABASE"
  object_identifier {
    name = snowflake_database.d.name
  }
}

resource "snowflake_schema" "s" {
  name     = "TEST_SCHEMA"
  database = snowflake_database.d.name
}

resource "snowflake_object_parameter" "o2" {
  key         = "USER_TASK_TIMEOUT_MS"
  value       = "500"
  object_type = "SCHEMA"
  object_identifier {
    database = snowflake_database.d.name
    name     = snowflake_schema.s.name
  }
}

resource "snowflake_table" "t" {
  name     = "TEST_TABLE"
  database = snowflake_database.d.name
  schema   = snowflake_schema.s.name
  column {
    name = "id"
    type = "NUMBER"
  }
}

resource "snowflake_object_parameter" "o3" {
  key         = "DATA_RETENTION_TIME_IN_DAYS"
  value       = "89"
  object_type = "TABLE"
  object_identifier {
    database = snowflake_database.d.name
    schema   = snowflake_schema.s.name
    name     = snowflake_table.t.name
  }
}

// Setting object parameter at account level
resource "snowflake_object_parameter" "o4" {
  key        = "DATA_RETENTION_TIME_IN_DAYS"
  value      = "89"
  on_account = true
}
