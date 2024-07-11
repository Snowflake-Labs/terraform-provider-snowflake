# Simple usage
data "snowflake_roles" "simple" {
}

output "simple_output" {
  value = data.snowflake_roles.simple.roles
}

# Filtering (like)
data "snowflake_roles" "like" {
  like = "role-name"
}

output "like_output" {
  value = data.snowflake_roles.like.roles
}

# Filtering (in class)
data "snowflake_roles" "in_class" {
  in_class = "SNOWFLAKE.CORE.BUDGET"
}

output "in_class_output" {
  value = data.snowflake_roles.in_class.roles
}

# Ensure the number of roles is equal to at least one element (with the use of postcondition)
data "snowflake_roles" "assert_with_postcondition" {
  like = "role-name-%"
  lifecycle {
    postcondition {
      condition     = length(self.roles) > 0
      error_message = "there should be at least one role"
    }
  }
}

# Ensure the number of roles is equal to at exactly one element (with the use of check block)
check "role_check" {
  data "snowflake_roles" "assert_with_check_block" {
    like = "role-name"
  }

  assert {
    condition     = length(data.snowflake_roles.assert_with_check_block.roles) == 1
    error_message = "Roles filtered by '${data.snowflake_roles.assert_with_check_block.like}' returned ${length(data.snowflake_roles.assert_with_check_block.roles)} roles where one was expected"
  }
}
