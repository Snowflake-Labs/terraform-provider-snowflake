---
page_title: "snowflake_managed_account Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_managed_account`





## Schema

### Required

- **admin_name** (String, Required) Identifier, as well as login name, for the initial user in the managed account. This user serves as the account administrator for the account.
- **admin_password** (String, Required) Password for the initial user in the managed account.
- **name** (String, Required) Identifier for the managed account; must be unique for your account.

### Optional

- **comment** (String, Optional) Specifies a comment for the managed account.
- **id** (String, Optional) The ID of this resource.
- **type** (String, Optional) Specifies the type of managed account.

### Read-only

- **cloud** (String, Read-only) Cloud in which the managed account is located.
- **created_on** (String, Read-only) Date and time when the managed account was created.
- **locator** (String, Read-only) Display name of the managed account.
- **region** (String, Read-only) Snowflake Region in which the managed account is located.
- **url** (String, Read-only) URL for accessing the managed account, particularly through the web interface.


