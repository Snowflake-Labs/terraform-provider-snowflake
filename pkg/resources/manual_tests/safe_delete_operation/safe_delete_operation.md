# Safe delete operation

Because of the limitations of terraform plugin testing framework, we cannot test the safe delete operation
automatically.
This test provides a guide how to test the safe delete operation manually.

## Snowflake setup

Before running Terraform tests, you have to create a simple database:

```snowflake
CREATE DATABASE TEST_DATABASE;
```

> Note: we also need a schema, but we can use the default one (`PUBLIC`).

## Terraform configuration

Create a new Terraform configuration file `main.tf` with the following content and initialize a new terraform project by
running `terraform init`:

```terraform
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "=1.1.0"
    }
  }
}

provider "snowflake" {
  driver_tracing = "info"
  preview_features_enabled = ["snowflake_table_resource"]
}

resource "snowflake_table" "test" {
  database = "TEST_DATABASE"
  schema   = "PUBLIC"
  name     = "TEMP_TABLE"
  column {
    name = "id"
    type = "NUMBER"
  }
  column {
    name = "name"
    type = "STRING"
  }
}
```

The configuration will be the same for all the steps.

## Test steps

The test will be split into two parts. In the first part,
we will test the delete operation without the safe delete operation enabled and in the second part we will enable it.

### Part 1. delete operation before the changes

1. Run `terraform apply -auto-approve` to create a new table.
2. Run `terraform apply -destroy` to delete the table (do not confirm the deletion by inputting `yes` into the console).
3. Run the following SQL statement to remove the public schema:

```snowflake
DROP SCHEMA TEST_DATABASE.PUBLIC;
```

4. Confirm the deletion by inputting `yes` into the console.

You should see the following error: `Error: [errors.go:22] object does not exist or not authorized` and by checking
the Terraform state by running `terraform state list`, you should see that the table is still there.
To bring back the initial state to run the second part run the following SQL statement:

```snowflake
CREATE SCHEMA TEST_DATABASE.PUBLIC;
```

and remove the table from the state by running `terraform state rm snowflake_table.test`.

### Part 2. delete operation after the changes

> Note: When the changes to the provider are released, you can change the version in the `main.tf` file and run
`terraform init -upgrade` to use the new version.
> If the changes are not yet released, proceed with local setup.

1. Build the project locally by running `make local-build` in the project root, and then follow-up with
   `make install-tf`.
2. Modify your `~/.terraformrc` to use the locally built provider, it should look similar to the following config:

```hcl
provider_installation {
  dev_overrides {
    "snowflakedb/snowflake" = "/<user_path>/.terraform.d/plugins" # TODO: Replace user_path
  }
}
```

3. Run the same steps as in the first part.
4. After following the same steps you should end up with successfully deleted table resource. You can verify that by
   running `terraform state list`.
