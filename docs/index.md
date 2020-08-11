# Snowflake Provider

This is a terraform provider plugin for managing [Snowflake](http://snowflakedb.com) accounts.
Coverage is focused on part of Snowflake related to access control.

## Example

```hcl
provider snowflake {
  // required
  username = "..."
  account  = "..."
  region   = "..."

  // optional, at exactly one must be set
  password           = "..."
  oauth_access_token = "..."
  private_key_path   = "..."

  // optional
  role  = "..."
}
```

## Authentication

The Snowflake provider support multiple ways to authenticate:

* Password
* OAuth Access Token
* Browser Auth
* Private Key

In all cases account, username, and region are required.

### Keypair Authentication Environment Variables

You should generate the public and private keys and set up environment variables.

```shell

cd ~/.ssh
openssl genrsa -out snowflake_key 4096
openssl rsa -in snowflake_key -pubout -out snowflake_key.pub
```

To export the variables into your provider:

```shell
export SNOWFLAKE_USER="..."
export SNOWFLAKE_PRIVATE_KEY_PATH="~/.ssh/snowflake_key"
```

### OAuth Access Token

If you have an OAuth access token, export these credentials as environment variables:

```shell
export SNOWFLAKE_USER='...'
export SNOWFLAKE_OAUTH_ACCESS_TOKEN='...'
```

Note that once this access token expires, you'll need to request a new one through an external application.

### Username and Password Environment Variables

If you choose to use Username and Password Authentication, export these credentials:

```shell
export SNOWFLAKE_USER='...'
export SNOWFLAKE_PASSWORD='...'
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Snowflake
 `provider` block:

* `account` - (required) The name of the Snowflake account. Can also come from the
  `SNOWFLAKE_ACCOUNT` environment variable.
* `username` - (required) Username for username+password authentication. Can come from the
  `SNOWFLAKE_PASSWORD` environment variable.
* `region` - (required) [Snowflake region](https://docs.snowflake.com/en/user-guide/intro-regions.html) to use. Can be source from the `SNOWFLAKE_REGION` environment variable.
* `password` - (optional) Password for username+password auth. Cannot be used with `browser_auth` or
  `private_key_path`. Can be source from `SNOWFLAKE_PASSWORD` environment variable.
* `oauth_access_token` - (optional) Token for use with OAuth. Generating the token is left to other
  tools. Cannot be used with `browser_auth`, `private_key_path` or `password`. Can be source from
  `SNOWFLAKE_OAUTH_ACCESS_TOKEN` environment variable.
* `private_key_path` - (optional) Path to a private key for using keypair authentication.. Cannot be
  used with `browser_auth`, `oauth_access_token` or `password`. Can be source from
  `SNOWFLAKE_PRIVATE_KEY_PATH` environment variable.
* `role` - (optional) Snowflake role to use for operations. If left unset, default role for user
  will be used. Can come from the `SNOWFLAKE_ROLE` environment variable.
