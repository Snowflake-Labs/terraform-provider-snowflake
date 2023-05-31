# Snowflake Terraform Provider

> ⚠️ **Please note**: If you believe you have found a security issue, _please responsibly disclose_ by contacting us at [team-cloud-foundation-tools-dl@snowflake.com](mailto:team-cloud-foundation-tools-dl@snowflake.com).

----

![.github/workflows/ci.yml](https://github.com/Snowflake-Labs/terraform-provider-snowflake/workflows/.github/workflows/ci.yml/badge.svg)

This is a terraform provider for managing [Snowflake](https://www.snowflake.com/) resources.

## Getting started

> If you're still using the `chanzuckerberg/snowflake` source, see [Upgrading from CZI Provider](./CZI_UPGRADE.md) to upgrade to the current version.

Install the Snowflake Terraform provider by adding a requirement block and a provider block to your Terraform codebase:
```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "~> 0.61"
    }
  }
}

provider "snowflake" {
  account  = "abc12345" # the Snowflake account identifier
  username = "johndoe"
  password = "v3ry$3cr3t"
  role     = "ACCOUNTADMIN"
}
```

For more information on provider configuration see the [provider docs on the Terraform registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs).

Don't forget to run `terraform init` and you're ready to go! 🚀

Start browsing the [registry docs](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs) to find resources and data sources to use.

## Getting Help

Some links that might help you:

- The [introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0) shows how to set up your Snowflake account for Terraform (service user, role, authentication, etc) and how to create your first resources in Terraform.
- The [docs on the Terraform registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest) are a complete reference of all resources and data sources supported and contain more advanced examples.
- The [discussions area](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions) of this repo, we use this forum to discuss new features and changes to the provider.
- **If you are an enterprise customer**, reach out to your account team. This helps us prioritize issues.
- The [issues section](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues) might already have an issue addressing your question.

## Contributing

Cf. [Contributing](./CONTRIBUTING.md).
