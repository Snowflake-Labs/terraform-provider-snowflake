# format is <database_name>.<schema_name>.<procedure_name>(<arg types, separated with ','>)
terraform import snowflake_procedure.example 'dbName.schemaName.procedureName(varchar, varchar, varchar)'
