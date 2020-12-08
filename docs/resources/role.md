---
page_title: "snowflake_role Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_role`



## Example Usage

```terraform
resource snowflake_role role {
  name    = "role1"
  comment = "A role."
}
```

## Schema

### Required

- **name** (String, Required)

### Optional

- **comment** (String, Optional)
- **id** (String, Optional) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_role.example roleName
```
