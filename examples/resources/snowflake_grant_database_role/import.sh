# format is database_role_name (string) | object_type (ROLE|DATABASE ROLE|SHARE) | grantee_name (string)
terraform import snowflake_grant_database_role.example '"ABC"."test_db_role"|ROLE|"test_parent_role"'
