# format is database name | schema name | stage name | privilege | true/false for with_grant_option
terraform import snowflake_stage_grant.example 'databaseName|schemaName|stageName|USAGE|true'
