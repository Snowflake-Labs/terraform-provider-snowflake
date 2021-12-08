# format is database name | schema name | pipe name | privilege | true/false for with_grant_option
terraform import snowflake_pipe_grant.example 'dbName|schemaName|pipeName|OPERATE|false'
