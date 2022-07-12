# format is database name | schema name | | privilege | true/false for with_grant_option
terraform import snowflake_schema_grant.example 'databaseName|schemaName||MONITOR|false'
