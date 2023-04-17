# format is database_name|schema_name|masking_policy_name|privilege|with_grant_option|roles
terraform import snowflake_masking_policy_grant.example "dbName|schemaName|maskingPolicyName|USAGE|false|role1,role2"

