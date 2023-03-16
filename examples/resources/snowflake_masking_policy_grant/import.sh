# format is database name | schema name | masking policy name | privilege | true/false for with_grant_option
terraform import snowflake_masking_policy_grant.example 'dbName|schemaName|maskingPolicyName|USAGE|false'
