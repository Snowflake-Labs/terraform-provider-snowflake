---
page_title: "snowflake_security_integrations Data Source - terraform-provider-snowflake"
subcategory: ""
description: |-
  Datasource used to get details of filtered security integrations. Filtering is aligned with the current possibilities for SHOW SECURITY INTEGRATIONS https://docs.snowflake.com/en/sql-reference/sql/show-integrations query (only like is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection security_integrations.
---

# snowflake_security_integrations (Data Source)

Datasource used to get details of filtered security integrations. Filtering is aligned with the current possibilities for [SHOW SECURITY INTEGRATIONS](https://docs.snowflake.com/en/sql-reference/sql/show-integrations) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `security_integrations`.

## Example Usage

```terraform
# Simple usage
data "snowflake_security_integrations" "simple" {
}

output "simple_output" {
  value = data.snowflake_security_integrations.simple.security_integrations
}

# Filtering (like)
data "snowflake_security_integrations" "like" {
  like = "security-integration-name"
}

output "like_output" {
  value = data.snowflake_security_integrations.like.security_integrations
}

# Filtering by prefix (like)
data "snowflake_security_integrations" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_security_integrations.like_prefix.security_integrations
}

# Without additional data (to limit the number of calls make for every found security integration)
data "snowflake_security_integrations" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SECURITY INTEGRATION for every security integration found and attaches its output to security_integrations.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_security_integrations.only_show.security_integrations
}

# Ensure the number of security_integrations is equal to at least one element (with the use of postcondition)
data "snowflake_security_integrations" "assert_with_postcondition" {
  like = "security-integration-name%"
  lifecycle {
    postcondition {
      condition     = length(self.security_integrations) > 0
      error_message = "there should be at least one security integration"
    }
  }
}

# Ensure the number of security_integrations is equal to at exactly one element (with the use of check block)
check "security_integration_check" {
  data "snowflake_security_integrations" "assert_with_check_block" {
    like = "security-integration-name"
  }

  assert {
    condition     = length(data.snowflake_security_integrations.assert_with_check_block.security_integrations) == 1
    error_message = "security integrations filtered by '${data.snowflake_security_integrations.assert_with_check_block.like}' returned ${length(data.snowflake_security_integrations.assert_with_check_block.security_integrations)} security integrations where one was expected"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `like` (String) Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).
- `with_describe` (Boolean) Runs DESC SECURITY INTEGRATION for each security integration returned by SHOW SECURITY INTEGRATIONS. The output of describe is saved to the description field. By default this value is set to true.

### Read-Only

- `id` (String) The ID of this resource.
- `security_integrations` (List of Object) Holds the aggregated output of all security integrations details queries. (see [below for nested schema](#nestedatt--security_integrations))

<a id="nestedatt--security_integrations"></a>
### Nested Schema for `security_integrations`

Read-Only:

- `describe_output` (List of Object) (see [below for nested schema](#nestedobjatt--security_integrations--describe_output))
- `show_output` (List of Object) (see [below for nested schema](#nestedobjatt--security_integrations--show_output))

<a id="nestedobjatt--security_integrations--describe_output"></a>
### Nested Schema for `security_integrations.describe_output`

Read-Only:

- `todo` (List of Object) (see [below for nested schema](#nestedobjatt--security_integrations--describe_output--todo))

<a id="nestedobjatt--security_integrations--describe_output--todo"></a>
### Nested Schema for `security_integrations.describe_output.todo`

Read-Only:

- `default` (String)
- `name` (String)
- `type` (String)
- `value` (String)



<a id="nestedobjatt--security_integrations--show_output"></a>
### Nested Schema for `security_integrations.show_output`

Read-Only:

- `category` (String)
- `comment` (String)
- `created_on` (String)
- `enabled` (Boolean)
- `integration_type` (String)
- `name` (String)
