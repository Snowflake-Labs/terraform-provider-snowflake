# format is database name | schema name | stored procedure name | <list of arg types, separated with '-'>
terraform import snowflake_procedure.example 'dbName|schemaName|procedureName|varchar-varchar-varchar'
