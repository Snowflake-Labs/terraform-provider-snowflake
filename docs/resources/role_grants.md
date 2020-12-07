---
page_title: "snowflake_role_grants Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_role_grants`



## Example Usage

```terraform
resource "snowflake_role" "role" {
  name    = "rking_test_role"
  comment = "for testing"
}

resource "snowflake_user" "user" {
  name    = "rking_test_user"
  comment = "for testing"
}

resource "snowflake_user" "user2" {
  name    = "rking_test_user2"
  comment = "for testing"
}

resource "snowflake_role" "other_role" {
  name = "rking_test_role2"
}

resource "snowflake_role_grants" "grants" {
  name = "foo"

  role_name = "${snowflake_role.role.name}"

  roles = [
    "${snowflake_role.other_role.name}",
  ]

  users = [
    "${snowflake_user.user.name}",
    "${snowflake_user.user2.name}",
  ]
}
```

## Schema

### Required

- **role_name** (String, Required) The name of the role we are granting.

### Optional

- **id** (String, Optional) The ID of this resource.
- **roles** (Set of String, Optional) Grants role to this specified role.
- **users** (Set of String, Optional) Grants role to this specified user.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_role_grants.example rolename
```
