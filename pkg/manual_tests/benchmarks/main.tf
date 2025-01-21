module "schema" {
  source         = "./schema"
  resource_count = 1
}

module "warehouse" {
  source         = "./warehouse"
  resource_count = 0
}

module "task" {
  source         = "./task"
  resource_count = 0
}

provider "snowflake" {
  profile = "secondary_test_account"
}

terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "= 1.0.1"
    }
  }
}
