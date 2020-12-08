# format is database name | schema name | view name | privilege | true/false for with_grant_option
terraform import snowflake_view_grant.example 'dbName|schemaName|viewName|USAGE|false'
