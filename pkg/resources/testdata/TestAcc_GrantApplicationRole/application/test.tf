resource "snowflake_stage" "stage" {
  name     = var.random_name
  database = var.database_name
  schema   = var.schema_name
}

locals {
  stage_identifier            = "\"${var.database_name}\".\"${var.schema_name}\".\"${snowflake_stage.stage.name}\""
  file_prefix                 = abspath(path.root)
  application_role_identifier = "\"${var.application_name}\".\"app_role_1\""
}

resource "snowflake_unsafe_execute" "put_manifest" {
  execute = "PUT 'file://${local.file_prefix}/manifest.yml' @${local.stage_identifier} AUTO_COMPRESS = FALSE"
  revert  = "SELECT 1"
}


resource "snowflake_unsafe_execute" "put_setup" {
  execute = "PUT 'file://${local.file_prefix}/setup.sql' @${local.stage_identifier} AUTO_COMPRESS = FALSE"
  revert  = "SELECT 1"
}

resource "snowflake_unsafe_execute" "application_package" {
  depends_on = [snowflake_unsafe_execute.put_manifest, snowflake_unsafe_execute.put_setup]
  execute    = "CREATE APPLICATION PACKAGE \"${var.random_name}\""
  revert     = "DROP APPLICATION PACKAGE \"${var.random_name}\" "
}

resource "snowflake_unsafe_execute" "application_version" {
  depends_on = [snowflake_unsafe_execute.application_package]
  execute    = "ALTER APPLICATION PACKAGE \"${var.random_name}\" ADD VERSION v1 USING '@${local.stage_identifier}'"
  revert     = "SELECT 1"
}

resource "snowflake_unsafe_execute" "application" {
  depends_on = [snowflake_unsafe_execute.application_version]
  execute    = "CREATE APPLICATION \"${var.application_name}\" FROM APPLICATION PACKAGE \"${var.random_name}\" USING VERSION v1"
  revert     = "DROP APPLICATION \"${var.application_name}\""
}

resource "snowflake_unsafe_execute" "application2" {
  depends_on = [snowflake_unsafe_execute.application_version]
  execute    = "CREATE APPLICATION \"${var.application_name2}\" FROM APPLICATION PACKAGE \"${var.random_name}\" USING VERSION v1"
  revert     = "DROP APPLICATION \"${var.application_name2}\""
}

resource "snowflake_grant_application_role" "g" {
  depends_on            = [snowflake_unsafe_execute.application2]
  application_role_name = local.application_role_identifier
  application_name      = "\"${var.application_name2}\""
}
