# format is database name | schema name | external table name | privilege | true/false for with_grant_option
terraform import snowflake_external_table_grant.example 'dbName|schemaName|externalTableName|SELECT|false'
