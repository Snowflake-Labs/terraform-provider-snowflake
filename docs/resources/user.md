---
page_title: "snowflake_user Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_user`



## Example Usage

```terraform
resource snowflake_user user {
  name         = "Snowflake User"
  login_name   = "snowflake_user"
  comment      = "A user of snowflake."
  password     = "secret"
  disabled     = false
  display_name = "Snowflake User"
  email        = "user@snowflake.example"
  first_name   = "Snowflake"
  last_name    = "User"

  default_warehouse = "warehouse"
  default_role      = "role1"

  rsa_public_key   = "..."
  rsa_public_key_2 = "..."

  must_change_password = false
}
```

## Schema

### Required

- **name** (String, Required) Name of the user. Note that if you do not supply login_name this will be used as login_name. [doc](https://docs.snowflake.net/manuals/sql-reference/sql/create-user.html#required-parameters)

### Optional

- **comment** (String, Optional)
- **default_namespace** (String, Optional) Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login.
- **default_role** (String, Optional) Specifies the role that is active by default for the user’s session upon login.
- **default_warehouse** (String, Optional) Specifies the virtual warehouse that is active by default for the user’s session upon login.
- **disabled** (Boolean, Optional)
- **display_name** (String, Optional) Name displayed for the user in the Snowflake web interface.
- **email** (String, Optional) Email address for the user.
- **first_name** (String, Optional) First name of the user.
- **id** (String, Optional) The ID of this resource.
- **last_name** (String, Optional) Last name of the user.
- **login_name** (String, Optional) The name users use to log in. If not supplied, snowflake will use name instead.
- **must_change_password** (Boolean, Optional) Specifies whether the user is forced to change their password on next login (including their first/initial login) into the system.
- **password** (String, Optional) **WARNING:** this will put the password in the terraform state file. Use carefully.
- **rsa_public_key** (String, Optional) Specifies the user’s RSA public key; used for key-pair authentication. Must be on 1 line without header and trailer.
- **rsa_public_key_2** (String, Optional) Specifies the user’s second RSA public key; used to rotate the public and private keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.

### Read-only

- **has_rsa_public_key** (Boolean, Read-only) Will be true if user as an RSA key set.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_user.example userName
```
