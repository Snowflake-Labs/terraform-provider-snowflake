# format is <database_name>.<schema_name>.<function_name>(<arg types, separated with ','>)
terraform import snowflake_function.example 'dbName.schemaName.functionName(varchar, varchar, varchar)'
