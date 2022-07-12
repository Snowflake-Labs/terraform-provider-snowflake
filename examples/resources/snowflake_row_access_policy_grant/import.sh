# format is database name | schema name | row access policy name | privilege | true/false for with_grant_option
terraform import snowflake_row_access_policy_grant.example 'dbName|schemaName|rowAccessPolicyName|SELECT|false'
