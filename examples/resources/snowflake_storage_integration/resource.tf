locals {
  snowflake = {
    role_name = "snowflake"
  }
}

data "aws_caller_identity" "current" {}

resource "aws_iam_role" "snowflake" {
  name               = local.snowflake.role_name
  assume_role_policy = data.aws_iam_policy_document.assume_snowflake.json
}

data "aws_iam_policy_document" "assume_snowflake" {
  statement {
    principals {
      type        = "AWS"
      identifiers = [snowflake_storage_integration.integration.storage_aws_iam_user_arn]
    }
    actions = ["sts:AssumeRole"]
    condition {
      test     = "StringEquals"
      values   = [snowflake_storage_integration.integration.storage_aws_external_id]
      variable = "sts:ExternalId"
    }
  }
}

resource "snowflake_storage_integration" "integration" {
  name    = "storage_integration"
  comment = "A storage integration."
  type    = "EXTERNAL_STAGE"

  enabled = true

  storage_allowed_locations = ["*"]

  storage_provider     = "S3"
  storage_aws_role_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${local.snowflake.role_name}"
}
