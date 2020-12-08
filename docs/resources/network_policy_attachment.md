---
page_title: "snowflake_network_policy_attachment Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_network_policy_attachment`



## Example Usage

```terraform
resource snowflake_network_policy_attachment attach {
  network_policy_name = "policy"
  set_for_account     = false
  users = ["user1", "user2"]
}
```

## Schema

### Required

- **network_policy_name** (String, Required) Specifies the identifier for the network policy; must be unique for the account in which the network policy is created.

### Optional

- **id** (String, Optional) The ID of this resource.
- **set_for_account** (Boolean, Optional) Specifies whether the network policy should be applied globally to your Snowflake account<br><br>**Note:** The Snowflake user running `terraform apply` must be on an IP address allowed by the network policy to set that policy globally on the Snowflake account.<br><br>Additionally, a Snowflake account can only have one network policy set globally at any given time. This resource does not enforce one-policy-per-account, it is the user's responsibility to enforce this. If multiple network policy resources have `set_for_account: true`, the final policy set on the account will be non-deterministic.
- **users** (Set of String, Optional) Specifies which users the network policy should be attached to

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_network_policy_attachment.example attachment_policyname
```
