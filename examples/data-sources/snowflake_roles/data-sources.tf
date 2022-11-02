data "snowflake_roles" "this" {

}

data "snowflake_roles" "ad" {
  pattern = "SYSADMIN"
}
