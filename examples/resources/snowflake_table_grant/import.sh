# format is database name | schema name | table name | privilege | true/false for with_grant_option
terraform import snowflake_table_grant.example 'databaseName|schemaName|tableName|MODIFY|true'
