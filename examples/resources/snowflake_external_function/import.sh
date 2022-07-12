# format is database name | schema name | external function name | <list of function arg types, separated with '-'>
terraform import snowflake_external_function.example 'dbName|schemaName|externalFunctionName|varchar-varchar-varchar'
