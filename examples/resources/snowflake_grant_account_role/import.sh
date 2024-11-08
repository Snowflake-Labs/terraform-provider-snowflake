# format is role_name (string) | grantee_object_type (ROLE|USER) | grantee_name (string)
terraform import snowflake_grant_account_role.example '"test_role"|ROLE|"test_parent_role"'
