package config_test

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/stretchr/testify/require"
)

func Test_FromModelPoc(t *testing.T) {
	t.Run("test basic", func(t *testing.T) {
		someModel := Some("test", "Some Name")
		expectedOutput := strings.TrimPrefix(`
"resource" "snowflake_share" "test" {
  "name" = "Some Name"
}
`, "\n")
		result := config.FromModelPoc(t, someModel)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test full", func(t *testing.T) {
		someModel := Some("test", "Some Name").
			WithComment("Some Comment").
			WithStringList("a", "b", "a").
			WithStringSet("a", "b", "c").
			WithObjectList(
				Item{IntField: 1, StringField: "first item"},
				Item{IntField: 2, StringField: "second item"},
			).WithDependsOn("some_other_resource.some_name", "other_resource.some_other_name", "third_resource.third_name")
		expectedOutput := strings.TrimPrefix(`
"resource" "snowflake_share" "test" {
  "comment" = "Some Comment"
  "name" = "Some Name"
  "string_list" = ["a", "b", "a"]
  "string_set" = ["a", "b", "c"]
  "object_list" = {
    "int_field" = 1
    "string_field" = "first item"
  }
  "object_list" = {
    "int_field" = 2
    "string_field" = "second item"
  }
  "depends_on" = [some_other_resource.some_name, other_resource.some_other_name, third_resource.third_name]
}
`, "\n")

		result := config.FromModelPoc(t, someModel)

		require.Equal(t, expectedOutput, result)
	})
}
