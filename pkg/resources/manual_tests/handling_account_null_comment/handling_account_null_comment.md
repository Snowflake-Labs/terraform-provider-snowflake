# Handling account null comment

This test shows that the problem from [this issue](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3402)
is now handled by the provider. The issue occurs when importing an account that has `null` comment.
Because of the limitations in the [terraform plugin testing framework](https://github.com/hashicorp/terraform-plugin-testing)
we cannot create account externally and then import that account in the first step of the test. This can only be tested manually.

## Snowflake setup

Before running Terraform tests you have to create an account we would like to import.
Run the following script to create an account:
```snowflake
CREATE ACCOUNT TESTING_ACCOUNT
    ADMIN_NAME = '<admin_name>' -- TODO: Replace
    ADMIN_PASSWORD = '<password>' -- TODO: Replace
    ADMIN_USER_TYPE = SERVICE
    EMAIL = '<email>' -- TODO: Replace
    EDITION = STANDARD
    COMMENT = NULL;
```

## Test steps 

In this test, we'll make a use of building the provider locally and overriding it in the `~/.terraformrc`.
For more details on that, please visit our [advanced debugging guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/CONTRIBUTING.md#advanced-debugging).

1. Copy the Terraform code from `main.tf` and initialize the project by running `terraform init`.
2. Run `terraform import snowflake_account.test_account '<organization_name>.<account_name>'`.
3. Right now, you should get the same error as in [this issue](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3402).
4. Modify your `~/.terraformrc` to use the locally built provider
5. Run `terraform init -upgrade` to make sure the overridden plugin is used (you will be notified by warning logged by Terraform CLI)
6. Run `terraform import snowflake_account.test_account '<organization_name>.<account_name>'`.
7. The import should be passing now. Run `terraform plan` to make sure the Read operation is also passing.

## Test cleanup

To clean up the test either run `terraform apply -auto-approve -destroy` or in a case where import didn't work
run the following Snowflake script:
```snowflake
DROP ACCOUNT TESTING_ACCOUNT GRACE_PERIOD_IN_DAYS = 3;
```
