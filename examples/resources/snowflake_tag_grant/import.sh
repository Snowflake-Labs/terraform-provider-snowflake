# format is database name | schema name | tag name | privilege | roles | true/false for with_grant_option
terraform import snowflake_tag_grant.example 'dbName|schemaName|tagName|APPLY|ROLE1,ROLE2|false'
