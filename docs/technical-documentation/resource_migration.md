
# Resource migration

Here's a guide on how to migrate deprecated resources to their new counter-parts.
The migration can be done in two ways. Either you can remove old grant resources and replace them with new ones or perform
more complicated migration, but without revoking any grant (no downtime migration). We'll focus on the second one as the first approach
is pretty straight forward. As an example we'll take `snowflake_database_grant` to `snowflake_grant_privileges_to_account_role` migration with one privilege granted to two roles:

```terraform
resource "snowflake_database_grant" "old_resource" {
  depends_on = [ snowflake_database.test, snowflake_role.a, snowflake_role.b ]
  database_name = snowflake_database.test.name
  privilege = "USAGE"
  roles = [ snowflake_role.a.name, snowflake_role.b.name ]
}
```

> **Important note:** **Always** save your state, before any state manipulation, so in case of failed migration, you will be able to recover from having incorrect state files.

#### 1. terraform list

First, we need to list all the grant resources that will need to be migrated.
We can do that by running the `terraform state list` command.

> **Tip:** for larger configurations, it's best to filter the results using the grep command. For example: `terraform state list | grep "snowflake_database_grant"`.

#### 2. terraform rm

Now choose which one you would to migrate next and remove it from the state, so when you remove the old resource,
no grant will be revoked. More specifically, the Terraform Delete operation won't be run for removed resource.
It will detach the resource block in your configuration from the actual Snowflake resource.
You can remove resources from the state with the `terraform state rm <resource_address>` command.
In our example, `terraform state rm snowflake_database_grant.old_resource`. After running the command, you can remove the resource from the configuration
(again, removing the state will "detach" the resource block from the Snowflake resource. That's why after removing it, the Terraform won't try to revoke USAGE from our roles).

#### 3. Two options from here

At this point, we have several options for creating new grant resources that will replace the old ones.
We will cover three options:
- Configuration + Terraform CLI
- Configuration + import block
- Generating the configuration with import block and `terraform plan -generate-config-out`

#### 3.1.1. Write a new grant resource that will be an equivalent of the older one

```terraform
resource "snowflake_grant_privileges_to_account_role" "new_resource" {
  depends_on = [snowflake_database.test, snowflake_role.a, snowflake_role.b]
  for_each   = toset([snowflake_role.a.name, snowflake_role.b.name])
  privileges = ["USAGE"]
  account_role_name  = each.key
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}
```

#### 3.1.2. terraform import

Write the `terraform import` command with the ID so that the resource will be able to parse and fill the state correctly.
You can find import syntax in the documentation for a given resource, [here](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_privileges_to_account_role#import)
is the one for `snowflake_grant_privileges_to_account_role`. In our case, the command will look like this:
```shell 
terraform import 'snowflake_grant_privileges_to_account_role.new_resource["role_a_name"]' 'role_a_name|USAGE|false|false|OnAccountObject|DATABASE|database_name'
terraform import 'snowflake_grant_privileges_to_account_role.new_resource["role_b_name"]' 'role_b_name|USAGE|false|false|OnAccountObject|DATABASE|database_name'
```

[Hashicorp documentation reference on import command](https://developer.hashicorp.com/terraform/cli/commands/import)

#### 3.2.1 Write import block with new resource

This is similar to the first approach, but here we don't have to worry about importing each `for_each`
entry one by one. In the `locals` block, we're defining a map of resource name to ID. Then, we have
to write a new resource similar to the first approach. In the end, we have to define an import block
which will import defined IDs to a specified resource.

```terraform
locals {
  migrations = {
    "${snowflake_role.a.name}" = "\"${snowflake_role.a.name}\"|false|false|USAGE|OnAccountObject|DATABASE|\"${snowflake_database.test.name}\""
    "${snowflake_role.b.name}" = "\"${snowflake_role.b.name}\"|false|false|USAGE|OnAccountObject|DATABASE|\"${snowflake_database.test.name}\""
  }
}

resource "snowflake_grant_privileges_to_account_role" "new_resource" {
  depends_on = [snowflake_database.test, snowflake_role.a, snowflake_role.b]
  for_each   = local.migrations
  privileges = ["USAGE"]
  account_role_name  = "\"${each.key}\""
  on_account_object {
    object_type = "DATABASE"
    object_name = "\"${snowflake_database.test.name}\""
  }
}

import {
  to = snowflake_grant_privileges_to_account_role.new_resource[each.key]
  id = each.value
  for_each = local.migrations
}
```

[Hashicorp documentation reference on import block](https://developer.hashicorp.com/terraform/language/import)

#### 3.2.2 Run terraform plan and apply

After running `terraform plan` you'll see if resources can be imported without any change. If that's the case
and nothing has to be adjusted, then we can perform `terraform apply` to import the state into our new grant resources.

#### 3.3.1. Write import block

Unfortunately, `for_each` cannot be used when generating with import blocks, so we have to define them one by one.
For simplicity, we'll define just one import block (the second one would look the same, only with a different role).

```terraform
import {
  to = snowflake_grant_privileges_to_account_role.new_resource_role_a
  id = "\"${snowflake_role.a.name}\"|false|false|USAGE|OnAccountObject|DATABASE|\"${snowflake_database.test.name}\""
}
```
[Hashicorp documentation reference on import block](https://developer.hashicorp.com/terraform/language/import)

#### 3.3.2. terraform plan -generate-config-out

After specifying the import block run the `terraform plan -generate-config-out=generated.tf` command,
which will scan your configuration files search for import blocks, and put the generated configurations inside the `generated.tf` file.

```terraform
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform
resource "snowflake_grant_privileges_to_account_role" "new_resource_role_a" {
  account_role_name    = "\"test_role_321123123\""
  all_privileges       = false
  always_apply         = false
  always_apply_trigger = null
  on_account           = false
  privileges           = ["USAGE"]
  with_grant_option    = false
  on_account_object {
    object_name = "\"test_database_1231321\""
    object_type = "DATABASE"
  }
}
```

#### 3.3.3. terraform plan and apply

After running `terraform plan` you'll see if there are any changes we have to do before applying our generated configuration.
If no errors are appearing you can run `terraform apply` to import state into generated configurations.

#### 3.3.4. Limitations of Generating Configurations

Config generation may be a good solution for a few reasons, but it also comes with limitations:
- Manual review/fixing
    - Half of the values could be removed because they're the same as the default values
- No `for_each` capabilities
    - You cannot specify `for_each` in the import block like in the second approach which promotes incremental migration
    - Generated configurations can't use `for_each` which results in much more configuration code
- No resource reference
    - As you can see `account_role_name` and `object_name` are plain values, but the values most likely should be referenced by other resources' names.

[Hashicorp documentation reference on limitations of generating configurations](https://developer.hashicorp.com/terraform/language/import/generating-configuration)
