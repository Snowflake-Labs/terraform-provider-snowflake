# format is database name | schema name | stream name | privilege | true/false for with_grant_option
terraform import snowflake_stream_grant.example 'dbName|schemaName|streamName|SELECT|false'
