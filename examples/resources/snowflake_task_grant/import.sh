# format is database name | schema name | task name | privilege | true/false for with_grant_option
terraform import snowflake_pipe_grant.example 'dbName|schemaName|taskName|OPERATE|false'
