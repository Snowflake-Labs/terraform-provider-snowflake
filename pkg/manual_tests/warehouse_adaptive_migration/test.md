# Manual test: single-apply migration to `snowflake_warehouse_adaptive`

Verifies [#5201](https://github.com/snowflakedb/terraform-provider-snowflake/issues/5201) against a locally built provider.

The configs live in `steps/` and are copied to `main.tf` one at a time. They are kept in a subdirectory on purpose: Terraform reads every `.tf` file in its working directory, so leaving them alongside `main.tf` would make it parse all of them at once and fail on duplicate `terraform` blocks and duplicate resources. It does not descend into subdirectories, so `steps/` is invisible to it.

## 0. Build the provider and set up the dev override

From the repository root:

```sh
make build-local
```

This produces `./terraform-provider-snowflake` at the repo root.

Add the override to `~/.terraformrc`, pointing at the **directory** containing the binary (not the binary itself):

```hcl
provider_installation {

  dev_overrides {
    "registry.terraform.io/snowflakedb/snowflake" = "<path_to_repository_root>"
  }

  direct {}
}
```

> Remember to remove or comment out this block afterwards - it silently overrides the provider for every Terraform run on your machine.

The configs declare `required_version = ">= 1.7"`, since `removed` blocks need Terraform 1.7 or higher (`import` blocks need 1.5). Credentials come from the provider's own defaults - the `default` profile in `~/.snowflake/config`, or `SNOWFLAKE_*` environment variables. Set `SNOWFLAKE_PROFILE` if you want a different profile.

Do **not** run `terraform init`. With a dev override, Terraform resolves the provider from the local path and prints a warning on every command. Every step below is run from this directory:

```sh
cd pkg/manual_tests/warehouse_adaptive_migration
```

`main.tf`, `terraform.tfstate*`, `.terraform/`, and `apply.log` are all working files created while running this test - do not commit them.

## 1. Create the standard warehouse

```sh
cp steps/step1_initial.tf main.tf
terraform apply
```

Expected: creates `ADAPTIVE_MIGRATION_TEST_WH` as type `STANDARD`, size `MEDIUM`.

Record the baseline in Snowflake - `created_on` is what proves later that the warehouse was never recreated:

```sql
SHOW WAREHOUSES LIKE 'ADAPTIVE_MIGRATION_TEST_WH';
```

## 2. Migrate in a single apply

```sh
cp steps/step2_migration.tf main.tf
terraform plan
```

Check the plan **before** applying:

- `snowflake_warehouse_adaptive.example` must be **imported and updated** - `~ update`, never `-/+ destroy and then create`
- `warehouse_type` changes `STANDARD` -> `ADAPTIVE`
- `snowflake_warehouse.example` is *forgotten*, not destroyed. Terraform words this as "will no longer be managed by Terraform, but will not be destroyed"

Then:

```sh
terraform apply
```

Verify in Snowflake:

```sql
SHOW WAREHOUSES LIKE 'ADAPTIVE_MIGRATION_TEST_WH';
-- type          = ADAPTIVE
-- state         = ENABLED   (not STARTED / SUSPENDED)
-- size, min_cluster_count, max_cluster_count, scaling_policy, auto_suspend = NULL
-- created_on    = unchanged from step 1  <- the warehouse was altered, not recreated
-- comment       = "created as a standard warehouse"
```

Also confirm the apply log is clean: there should be **no** `Command not supported on an Adaptive Warehouse` error. Before the `AlterWithSuspend` fix, a successful apply still logged that. To see the provider's SQL and log lines:

```sh
TF_LOG=DEBUG terraform apply 2>&1 | tee apply.log
grep -i "not supported on an Adaptive Warehouse" apply.log   # expect no matches
```

And that the warehouse actually works:

```sql
USE WAREHOUSE ADAPTIVE_MIGRATION_TEST_WH;
SELECT COUNT(*) FROM TABLE(GENERATOR(ROWCOUNT => 100000));
```

## 3. Drop the migration blocks

```sh
cp steps/step3_cleanup.tf main.tf
terraform plan
```

Expected: `No changes.` An empty plan here is the real proof that the migration converged - the standard resource is gone from state and the adaptive resource matches Snowflake.

## 4. Negative check: interactive warehouse is still rejected

Create an interactive warehouse by hand:

```sql
CREATE INTERACTIVE WAREHOUSE ADAPTIVE_MIGRATION_TEST_INTERACTIVE_WH;
```

Add an import block for it to `main.tf`:

```terraform
resource "snowflake_warehouse_adaptive" "interactive_attempt" {
  name = "ADAPTIVE_MIGRATION_TEST_INTERACTIVE_WH"
}

import {
  to = snowflake_warehouse_adaptive.interactive_attempt
  id = "ADAPTIVE_MIGRATION_TEST_INTERACTIVE_WH"
}
```

`terraform plan` must fail with:

```
warehouse "ADAPTIVE_MIGRATION_TEST_INTERACTIVE_WH" is an interactive warehouse and cannot be converted to ADAPTIVE; use snowflake_warehouse_interactive instead
```

Then remove those two blocks again.

## 5. Clean up

```sh
cp steps/step3_cleanup.tf main.tf
terraform destroy
```

```sql
DROP WAREHOUSE IF EXISTS ADAPTIVE_MIGRATION_TEST_INTERACTIVE_WH;
SHOW WAREHOUSES LIKE 'ADAPTIVE_MIGRATION_TEST%';
```

Finally, remove the `dev_overrides` block from `~/.terraformrc`.
