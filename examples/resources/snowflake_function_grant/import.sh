# format is database name | schema name | function signature | privilege | true/false for with_grant_option
terraform import snowflake_function_grant.example 'dbName|schemaName|functionName(ARG1 ARG1TYPE, ARG2 ARG2TYPE):RETURNTYPE|USAGE|false'
