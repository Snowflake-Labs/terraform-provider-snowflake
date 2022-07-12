# format is database name | schema name | file format name | privilege | true/false for with_grant_option
terraform import snowflake_file_format_grant.example 'dbName|schemaName|fileFormatName|USAGE|false'
