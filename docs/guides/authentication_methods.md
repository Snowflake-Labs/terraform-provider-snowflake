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

[//]: # (TODO: https://snowflakecomputing.atlassian.net/browse/SNOW-1791729)

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

Snowflake authenticator is a default value, but it may change over time, so if you want to be explicit about it you can specify the following configuration:

```terraform
provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  password          = "<password>"
  authenticator     = "Snowflake"
}
```

In case of any problems, please go to the [common issues section](#common-issues).

### JWT authenticator flow 

To use JWT authentication, you have to firstly generate key-pairs used by Snowflake.
To correctly generate the necessary keys, follow [this guide](https://docs.snowflake.com/en/user-guide/key-pair-auth#configuring-key-pair-authentication) from the official Snowflake documentation.
After you [set the generated public key](https://docs.snowflake.com/en/user-guide/key-pair-auth#assign-the-public-key-to-a-snowflake-user) to the Terraform user and [verified it](https://docs.snowflake.com/en/user-guide/key-pair-auth#verify-the-user-s-public-key-fingerprint),
you can proceed with the following provider configuration:

```terraform
provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  password          = "<password>"
  authenticator     = "JWT"
  private_key = file("~/.ssh/snowflake_private_key.p8")
}
```

To load the private key you can utilize the built-in [file](https://developer.hashicorp.com/terraform/language/functions/file) function.

#### JWT authenticator flow using passphrase

If you would like to use key-pair utilizing passphrase, you can add it to the configuration like so:

```terraform
provider "snowflake" {
  organization_name      = "<organization_name>"
  account_name           = "<account_name>"
  user                   = "<user_name>"
  password               = "<password>"
  authenticator          = "JWT"
  private_key = file("~/.ssh/snowflake_private_key.p8")
  private_key_passphrase = "<passphrase>"
}
```

In case of any problems, please go to the [common issues section](#common-issues).

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

In case of any problems, please go to the [common issues section](#common-issues).

### Okta authenticator flow

To set up a new Okta account for this flow, follow [this guide](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/pkg/manual_tests/authentication_methods/README.md#okta-authenticator-test).
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

In case of any problems, please go to the [common issues section](#common-issues).

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

### Errors similar to (http: 404): open snowflake connection: 261004 (08004): failed to auth for unknown reason.

This can be caused by missing or incorrect host. When the host field is not set, it's being guessed based on organization name, account name, and other parameters.
For some deployments it will work fine, but for more custom ones, setting the host is necessary to successfully establish the connection.
