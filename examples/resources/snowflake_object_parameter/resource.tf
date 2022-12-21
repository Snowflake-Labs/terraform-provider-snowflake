resource "snowflake_database" "d" {
  name = "TEST_DB"
}

resource "snowflake_schema" "s" {
  name     = "TEST_SCHEMA"
  database = snowflake_database.d.name
}

resource "snowflake_object_parameter" "o" {
  key         = "ENABLE_STREAM_TASK_REPLICATION"
  value       = "true"
  object_type = "DATABASE"
  object_name = snowflake_database.d.name
}

resource "snowflake_object_parameter" "o2" {
  key         = "PIPE_EXECUTION_PAUSED"
  value       = "false"
  object_type = "SCHEMA"
  object_name = "${snowflake_database.d.name}.${snowflake_schema.s.name}"
}
