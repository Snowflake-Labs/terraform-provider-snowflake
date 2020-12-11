# format is database name | schema name | materialized view name | privilege | true/false for with_grant_option
terraform import snowflake_materialized_view_grant.example 'dbName|schemaName|materializedViewName|SELECT|false'
