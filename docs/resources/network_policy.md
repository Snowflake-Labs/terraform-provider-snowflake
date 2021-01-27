---
page_title: "snowflake_network_policy Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_network_policy`



## Example Usage

```terraform
resource snowflake_network_policy policy {
  name    = "policy"
  comment = "A policy."

  allowed_ip_list = ["192.168.0.100/24"]
  blocked_ip_list = ["192.168.0.101"]
}
```

## Schema

### Required

- **allowed_ip_list** (Set of String, Required) Specifies one or more IPv4 addresses (CIDR notation) that are allowed access to your Snowflake account
- **name** (String, Required) Specifies the identifier for the network policy; must be unique for the account in which the network policy is created.

### Optional

- **blocked_ip_list** (Set of String, Optional) Specifies one or more IPv4 addresses (CIDR notation) that are denied access to your Snowflake account<br><br>**Do not** add `0.0.0.0/0` to `blocked_ip_list`
- **comment** (String, Optional) Specifies a comment for the network policy.
- **id** (String, Optional) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_network_policy.example policyname
```
