---
page_title: "Authentication methods"
subcategory: ""
description: |-

---
# Authentication methods

This guide focuses on providing an example on every authentication method available in the provider.
Each method includes steps for setting dependencies, like MFA app, and getting encrypted/unencrypted keys.
For now, we provide examples for the most common use cases.
The rest of the options (Okta, ExternalBrowser, TokenAccessor) are planned to be added later on.

[//]: # (TODO: SNOW-1791729)

## Protecting secret values

When using any of the provided methods, remember to securely store sensitive information.

Here's a list of useful materials on keeping your secrets safe when using Terraform:
- https://developer.hashicorp.com/terraform/cloud-docs/architectural-details/security-model
- https://developer.hashicorp.com/terraform/tutorials/secrets/secrets-vault
- https://developer.hashicorp.com/terraform/language/state/sensitive-data

Read more on Snowflake's password protection: https://docs.snowflake.com/en/user-guide/leaked-password-protection.

## Table of contents

* [Protecting secret values](#protecting-secret-values)
* [Snowflake authenticator flow (login + password)](#snowflake-authenticator-flow-login--password)
* [JWT authenticator flow](#jwt-authenticator-flow-)
  * [JWT authenticator flow using passphrase](#jwt-authenticator-flow-using-passphrase)
* [MFA authenticator flow](#mfa-authenticator-flow)
  * [MFA token caching](#mfa-token-caching)
* [Okta authenticator flow](#okta-authenticator-flow)
* [Common issues](#common-issues)
  * [How can I get my organization name?](#how-can-i-get-my-organization-name)
  * [How can I get my account name?](#how-can-i-get-my-account-name)
  * [Errors similar to (http: 404): open snowflake connection: 261004 (08004): failed to auth for unknown reason.](#errors-similar-to-http-404-open-snowflake-connection-261004-08004-failed-to-auth-for-unknown-reason)

## Authentication flows

### Snowflake authenticator flow (login + password)

Provider setup in this case is pretty straightforward:

```terraform
provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  password          = "<password>"
}
```

Without passing any authenticator, we depend on the underlying Go Snowflake driver and Snowflake itself to fill this field out.
This means that we do not provision the default, and it may change at some point, so if you want to be explicit, you can define Snowflake authenticator like so:

```terraform
provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  password          = "<password>"
  authenticator     = "Snowflake"
}
```

### JWT authenticator flow 

To use JWT authentication, you have to firstly generate key-pairs used by Snowflake.
To correctly generate the necessary keys, follow [this guide](https://docs.snowflake.com/en/user-guide/key-pair-auth#configuring-key-pair-authentication) from the official Snowflake documentation.
After you [set the generated public key](https://docs.snowflake.com/en/user-guide/key-pair-auth#assign-the-public-key-to-a-snowflake-user) to the Terraform user and [verify it](https://docs.snowflake.com/en/user-guide/key-pair-auth#verify-the-user-s-public-key-fingerprint),
you can proceed with the following provider configuration:

```terraform
provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  authenticator     = "JWT"
  private_key       = file("~/.ssh/snowflake_private_key.p8")
}
```

To load the private key you can utilize the built-in [file](https://developer.hashicorp.com/terraform/language/functions/file) function.
If you have any issues with this method, one of the possible root causes could be an additional newline at the end of the file that causes error in the underlying Go Snowflake driver.
If this doesn't help, you can try other methods of supplying this field:
- Filling the key directly by using [multi-string notation](https://developer.hashicorp.com/terraform/language/expressions/strings#heredoc-strings)
- Sourcing it from the environment variable:
```shell
export SNOWFLAKE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----..."
# Alternatively, source from a file.
export SNOWFLAKE_PRIVATE_KEY=$(cat ~/.ssh/snowflake_private_key.p8)

export SNOWFLAKE_PRIVATE_KEY_PASSPHRASE="..."
```
- Using TOML configuration file:
```toml
[default]
private_key = "..."
private_key_passphrase = "..."
```

In case of any other issues, take a look at related topics:
- https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3332#issuecomment-2618957814
- https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3350#issuecomment-2604851052

#### JWT authenticator flow using passphrase

If you would like to use key-pair utilizing passphrase, you can add it to the configuration like so:

```terraform
provider "snowflake" {
  organization_name      = "<organization_name>"
  account_name           = "<account_name>"
  user                   = "<user_name>"
  authenticator          = "JWT"
  private_key            = file("~/.ssh/snowflake_private_key.p8")
  private_key_passphrase = "<passphrase>"
}
```

### MFA authenticator flow

Before being able to log in with MFA method, you have to prepare your Terraform user by following [this guide](https://docs.snowflake.com/en/user-guide/security-mfa) in the official Snowflake documentation.
Once MFA is set up on your Terraform user, you can use one of the following configurations.
Choosing the configuration depends on the preferred confirmation method (push notification or passcode) and the one that is available (not always both options are available).

The configuration that uses push notification:

```terraform
provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  password          = "<password>"
  authenticator     = "UsernamePasswordMFA"
}
```

and the configuration that uses passcode:

```terraform
provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  password          = "<password>"
  authenticator     = "UsernamePasswordMFA"
  passcode          = "000000"
}
```

#### MFA token caching

MFA token caching can help to reduce the number of prompts that must be acknowledged while connecting and authenticating to Snowflake, especially when multiple connection attempts are made within a relatively short time interval.
Follow [this guide](https://docs.snowflake.com/en/user-guide/security-mfa#using-mfa-token-caching-to-minimize-the-number-of-prompts-during-authentication-optional) to enable it.

### Okta authenticator flow

To set up a new Okta account for this flow, follow [this guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/b863d2e79ae6ae021552c4348e3012b8053ede17/pkg/manual_tests/authentication_methods/README.md#okta-authenticator-test).
If you already have an Okta account, skip the first point and follow the next steps.
The guide includes writing the provider configuration in the TOML file, but here's what it should look like fully in HCL:

```terraform
provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  password          = "<password>"
  authenticator     = "Okta"
  okta_url          = "https://dev-123456.okta.com"
}
```

## Common issues

### How can I get my organization name?

If you are logged into account that is in the same organization as Terraform user (or logged in as Terraform user), you can run:
```snowflake
SELECT CURRENT_ORGANIZATION_NAME();
```
The output of this command is your `<organization_name>`.

### How can I get my account name?

If you are logged into as a user that is in the same account as Terraform user (or logged in as Terraform user), you can run:
```snowflake
SELECT CURRENT_ACCOUNT_NAME();
```
The output of this command is your `<account_name>`.

## General recommendations

### Be sure you are passing all the required fields

This point is not only referring to double-checking the fields you are passing, but also to inform you that depending on the account 
you want to log into, a different set of parameters may be required.

Whenever you are on a Snowflake deployment that has different url than the default one:
`<organization_name>-<account_name>.snowflakecomputing.com`, you may encounter errors similar to:

```text
open snowflake connection: 261004 (08004): failed to auth for unknown reason.
```

This error can be raised for a number of reasons, but explicitly specifying the host has effectively prevented such occurrences so far.