terraform {
  required_version = ">= 1.3.6"
  required_providers {
    snowflake = {
      source  = "snowflake-labs/snowflake"
      version = "0.94.1"
    }
  }
}

resource "snowflake_user" "user1" {
  name = "example_user1"
}

resource "snowflake_user" "user2" {
  name = "example_user2"
}

resource "snowflake_user" "user3" {
  name = "example_user3"
}

resource "snowflake_account_role" "role1" {
  name = "example_role1"
}

resource "snowflake_account_role" "role2" {
  name = "example_role2"
}

resource "snowflake_account_role" "role3" {
  name = "example_role3"
}

locals {
  yaml_roles = yamldecode(file("${path.module}/config.yaml"))
  grant_to_user = distinct(flatten([
    for k, v in local.yaml_roles.roles : v.grant_to.user == null ? [] : [
    for u in v.grant_to.user : {
      role    = k
      to_user = u
    }
  ]]))
  grant_to_role = distinct(flatten([
    for k, v in local.yaml_roles.roles : v.grant_to.role == null ? [] : [
    for r in v.grant_to.role : {
      role    = k
      to_role = r
    }
  ]]))
}

output "grant_to_user_output" {
  value = local.grant_to_user
}

output "grant_to_role_output" {
  value = local.grant_to_role
}

resource "snowflake_grant_account_role" "user_grants" {
  depends_on = [
    snowflake_user.user1, snowflake_user.user2, snowflake_user.user3, snowflake_account_role.role1,
    snowflake_account_role.role2, snowflake_account_role.role3
  ]
  for_each  = {for entry in local.grant_to_user : "${entry.role}.${entry.to_user}" => entry}
  role_name = each.value.role
  user_name = each.value.to_user
}

resource "snowflake_grant_account_role" "role_grants" {
  depends_on = [
    snowflake_user.user1, snowflake_user.user2, snowflake_user.user3, snowflake_account_role.role1,
    snowflake_account_role.role2, snowflake_account_role.role3
  ]
  for_each         = {for entry in local.grant_to_role : "${entry.role}.${entry.to_role}" => entry}
  role_name        = each.value.role
  parent_role_name = each.value.to_role
}
