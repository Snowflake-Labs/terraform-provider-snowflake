# format is <database_name>.<schema_name>.<external_function_name>(<arg types, separated with ','>)
terraform import snowflake_external_function.example 'dbName.schemaName.externalFunctionName(varchar, varchar, varchar)'
