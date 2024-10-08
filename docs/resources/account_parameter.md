---
page_title: "snowflake_account_parameter Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# snowflake_account_parameter (Resource)



## Example Usage

```terraform
resource "snowflake_account_parameter" "p" {
  key   = "ALLOW_ID_TOKEN"
  value = "true"
}

resource "snowflake_account_parameter" "p2" {
  key   = "CLIENT_ENCRYPTION_KEY_SIZE"
  value = "256"
}
```

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/guides/identifiers#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) Name of account parameter. Valid values are those in [account parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#account-parameters).
- `value` (String) Value of account parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_account_parameter.p <parameter_name>
```
