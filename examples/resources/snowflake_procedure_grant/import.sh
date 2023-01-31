# format is database name | schema name | procedure signature | privilege | true/false for with_grant_option
terraform import snowflake_procedure_grant.example 'dbName|schemaName|procedureName(ARG1TYPE,ARG2TYPE)|USAGE|false'
