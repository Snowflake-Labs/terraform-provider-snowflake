---
page_title: "Unassigning Policies"
subcategory: ""
description: |-

---
# Unassigning policies

For some objects, like network policies, Snowflake [docs](https://docs.snowflake.com/en/sql-reference/sql/drop-network-policy#usage-notes) suggest that a network policy cannot be dropped successfully if it is currently assigned to another object. Currently, the provider does not unassign such objects automatically.

Before dropping the resource:
- if the objects the policy is assigned to are managed in Terraform, follow the example below
- if they are not managed in Terraform, list them with `SELECT * from table(information_schema.policy_references(policy_name=>'<string>'));` and unassign them manually with `ALTER ...`

## Example

When you have a configuration like
```terraform
resource "snowflake_network_policy" "example" {
  name = "network_policy_name"
}

resource "snowflake_oauth_integration_for_custom_clients" "example" {
  name               = "integration"
  oauth_client_type  = "CONFIDENTIAL"
  oauth_redirect_uri = "https://example.com"
  blocked_roles_list = ["ACCOUNTADMIN", "SECURITYADMIN"]
  network_policy                   = snowflake_network_policy.example.fully_qualified_name
}
```

and try removing the network policy, Terraform fails with
```
│ Error deleting network policy EXAMPLE, err = 001492 (42601): SQL compilation error:
│ Cannot perform Drop operation on network policy EXAMPLE. The policy is attached to INTEGRATION with name EXAMPLE. Unset the network policy from INTEGRATION and try the
│ Drop operation again.
```

In order to remove the policy correctly, first adjust the configuration to
```terraform
resource "snowflake_network_policy" "example" {
  name = "network_policy_name"
}

resource "snowflake_oauth_integration_for_custom_clients" "example" {
  name               = "integration"
  oauth_client_type  = "CONFIDENTIAL"
  oauth_redirect_uri = "https://example.com"
  blocked_roles_list = ["ACCOUNTADMIN", "SECURITYADMIN"]
}
```

Note that the network policy has been unassigned. Now, run `terraform apply`. This should cause the policy to be unassigned. Now, adjust the configuration once again to
```terraform
resource "snowflake_oauth_integration_for_custom_clients" "example" {
  name               = "integration"
  oauth_client_type  = "CONFIDENTIAL"
  oauth_redirect_uri = "https://example.com"
  blocked_roles_list = ["ACCOUNTADMIN", "SECURITYADMIN"]
}
```

Now the network policy should be removed successfully.

This behavior will be fixed in the provider in the future.
