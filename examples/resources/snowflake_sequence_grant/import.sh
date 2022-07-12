# format is database name | schema name | sequence name | privilege | true/false for with_grant_option
terraform import snowflake_sequence_grant.example 'dbName|schemaName|sequenceName|USAGE|false'
