# Handling account null comment

This test shows that the problem from [this issue](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3402)
is now handled by the provider. Because of the limitations in the [terraform plugin testing framework](https://github.com/hashicorp/terraform-plugin-testing)
we cannot create account externally and then import that account in the first step of the test. This can only be tested manually.

## Snowflake setup

Before running Terraform tests you have to create an account we would like to import.
Run the following script to create an account:
```snowflake
CREATE ACCOUNT TESTING_ACCOUNT
    ADMIN_NAME = '<admin_name>'
    ADMIN_PASSWORD = '<password>'
    ADMIN_USER_TYPE = SERVICE
    EMAIL = '<email>'
    EDITION = STANDARD
    COMMENT = NULL;
```

## Test steps 

1. Build the provider by running `make install-tf`.
2. Copy the Terraform code from `main.tf` and run `terraform init`.
3. Run `terraform import snowflake_account.test_account '<organization_name>.<account_name>'`.
4. Right now, you should get the same error as in [this issue](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3402).
5. Modify your `~/.terraformrc`, so that it looks like the following configuration:
```terraform
provider_installation {
  dev_overrides {
      "Snowflake-Labs/snowflake" = "<path to .terraform.d/plugins>" # should be logged by `make install-tf` command
  }
  direct {}
}
```
6. Run `terraform init -upgrade` to make sure the overridden plugin is used (you will be notified by warning logged by Terraform CLI)
7. Run `terraform import snowflake_account.test_account '<organization_name>.<account_name>'`.
8. The import should be passing now. Run `terraform plan` to make sure the Read operation is also passing.

## Test cleanup

To clean up the test either run `terraform apply -auto-approve -destroy` or in a case where import didn't work
run the following Snowflake script:
```snowflake
DROP ACCOUNT TESTING_ACCOUNT GRACE_PERIOD_IN_DAYS = 3;
```
