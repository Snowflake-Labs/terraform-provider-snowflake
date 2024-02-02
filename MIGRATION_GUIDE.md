# Migration guide

This document is meant to help you migrate your Terraform config to the new newest version. In migration guides, we will only 
describe deprecations or breaking changes and help you to change your configuration to keep the same (or similar) behavior
across different versions.

## vX.XX.X -> v0.85.0

### Migration from old grant resources to new ones

In recent changes, we introduced new grant resources intended to replace old grant solutions. Here's some of the useful
information that may help during the migration of the grant resources. Here's an example of the `snowflake_database_grant` to `snowflake_grant_privileges_to_account_role` migration.
The migration can be done in two ways. Either you can remove old grant resources and replace them with new ones or perform
more complicated migration, but without revoking any grant (no downtime migration). We'll focus on the second one as the first approach
is pretty straight forward. As an example we'll take `snowflake_database_grant` that grants one privilege to two roles:

```terraform
resource "snowflake_database_grant" "old_resource" {
  depends_on = [ snowflake_database.test, snowflake_role.a, snowflake_role.b ]
  database_name = snowflake_database.test.name
  privilege = "USAGE"
  roles = [ snowflake_role.a.name, snowflake_role.b.name ]
}
```

#### 1. terraform list

Run `terraform state list` to search the grants you're looking for (for larger configurations it's best to filter the results), 
for example, `terraform state list | grep "snowflake_database_grant"`.

#### 2. terraform rm

Now choose which one you would to migrate next and remove it from the state with `terraform state rm <resource_address>`. 
In our example, `terraform state rm snowflake_database_grant.old_resource`. After running the command, you can remove the resource from the configuration 
(removing the state will "detach" it from the resource block, so after removing it, the Terraform won't try to revoke USAGE from our roles).

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
  role_name  = each.key
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

#### 3.3.4. Thoughts on configuration generation

Config generation may be a good solution for a few reasons, but it also comes with limitations:
- Manual review/fixing
    - Half of the values could be removed because they're the same as the default values
- No `for_each` capabilities
    - You cannot specify `for_each` in the import block like in the second approach which promotes incremental migration
    - Generated configurations can't use `for_each` which results in much more configuration code
- No resource reference
    - As you can see `account_role_name` and `object_name` are plain values, but the values most likely should be referenced by other resources' names.

[Hashicorp documentation reference on generating configuration limitations](https://developer.hashicorp.com/terraform/language/import/generating-configuration)

## v0.84.0 ➞ v0.85.0

### snowflake_notification_integration resource changes
#### *(behavior change)* notification_provider
`notification_provider` becomes required and has three possible values `AZURE_STORAGE_QUEUE`, `AWS_SNS`, and `GCP_PUBSUB`.
It is still possible to set it to `AWS_SQS` but because there is no underlying SQL, so it will result in an error.
Attributes `aws_sqs_arn` and `aws_sqs_role_arn` will be ignored.
Computed attributes `aws_sqs_external_id` and `aws_sqs_iam_user_arn` won't be updated.

#### *(behavior change)* force new for multiple attributes
Force new was added for the following attributes (because no usable SQL alter statements for them):
- `azure_storage_queue_primary_uri`
- `azure_tenant_id`
- `gcp_pubsub_subscription_name`
- `gcp_pubsub_topic_name`

#### *(deprecation)* direction
`direction` parameter is deprecated because it is added automatically on the SDK level.

#### *(deprecation)* type
`type` parameter is deprecated because it is added automatically on the SDK level (and basically it's always `QUEUE`).

## v0.73.0 ➞ v0.74.0
### Provider configuration changes

In this change we have done a provider refactor to make it more complete and customizable by supporting more options that 
were already available in Golang Snowflake driver. This lead to several attributes being added and a few deprecated.
We will focus on the deprecated ones and show you how to adapt your current configuration to the new changes.

#### *(rename)* username ➞ user

```terraform
provider "snowflake" {
  # before
  username = "username"
  
  # after
  user = "username"
}
```

#### *(structural change)* OAuth API

```terraform
provider "snowflake" {
  # before
  browser_auth        = false
  oauth_access_token  = "<access_token>"
  oauth_refresh_token = "<refresh_token>"
  oauth_client_id     = "<client_id>"
  oauth_client_secret = "<client_secret>"
  oauth_endpoint      = "<endpoint>"
  oauth_redirect_url  = "<redirect_uri>"

  # after
  authenticator = "ExternalBrowser"
  token         = "<access_token>"
  token_accessor {
    refresh_token   = "<refresh_token>"
    client_id       = "<client_id>"
    client_secret   = "<client_secret>"
    token_endpoint  = "<endpoint>"
    redirect_uri    = "<redirect_uri>"
  }
}
```

#### *(remove redundant information)* region

Specifying a region is a legacy thing and according to https://docs.snowflake.com/en/user-guide/admin-account-identifier
you can specify a region as a part of account parameter. Specifying account parameter with the region is also considered legacy,
but with this approach it will be easier to convert only your account identifier to the new preferred way of specifying account identifier.

```terraform
provider "snowflake" {
  # before
  region = "<cloud_region_id>"

  # after
  account = "<account_locator>.<cloud_region_id>"
}
```

#### *(todo)* private key path

```terraform
provider "snowflake" {
  # before
  private_key_path = "<filepath>"

  # after
  private_key = file("<filepath>")
}
```

#### *(rename)* session_params ➞ params

```terraform
provider "snowflake" {
  # before
  session_params = {}
  
  # after
  params = {}
}
```

#### *(behavior change)* authenticator (JWT)

Before the change `authenticator` parameter did not have to be set for private key authentication and was deduced by the provider. The change is a result of the introduced configuration alignment with an underlying [gosnowflake driver](https://github.com/snowflakedb/gosnowflake). The authentication type is required there, and it defaults to user+password one. From this version, set `authenticator` to `JWT` explicitly.
