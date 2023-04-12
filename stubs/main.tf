terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "~> 0.61"
    }
  }
}

variable "snowflake_account" {
  type = string
}

variable "snowflake_username" {
  type = string
}

variable "snowflake_private_key" {
  type = string
}

variable "snowflake_region" {
  type = string
}

provider "snowflake" {
  account     = var.snowflake_account
  username    = var.snowflake_username
  private_key = var.snowflake_private_key
  region      = var.snowflake_region
  role        = "SYSADMIN"
}

resource "snowflake_database" "db" {
  name = "TEST_DB"
}

resource "snowflake_schema" "schema" {
  database = snowflake_database.db.name
  name     = "TEST_SCHEMA"
}
