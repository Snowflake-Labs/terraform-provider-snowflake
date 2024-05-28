package resources_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

type grantType int

const (
	normal grantType = iota
	onFuture
	onAll
)

func TestGetPropertyAsPointer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
	}, map[string]interface{}{
		"integer": 123,
		"string":  "some string",
		"boolean": true,
		"invalid": true,
	})

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Equal(t, true, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))
}

// queriedAccountRolePrivilegesEqualTo will check if all the privileges specified in the argument are granted in Snowflake.
func queriedPrivilegesEqualTo(query func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error), privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		grants, err := query(client, ctx)
		if err != nil {
			return err
		}
		for _, grant := range grants {
			if (grant.GrantTo == sdk.ObjectTypeDatabaseRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole) && grant.Privilege == "USAGE" {
				continue
			}
			if !slices.Contains(privileges, grant.Privilege) {
				return fmt.Errorf("grant not expected, grant: %v, not in %v", grants, privileges)
			}
		}

		return nil
	}
}

// queriedAccountRolePrivilegesContainAtLeast will check if all the privileges specified in the argument are granted in Snowflake.
// Any additional grants will be ignored.
func queriedPrivilegesContainAtLeast(query func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error), roleName sdk.ObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		grants, err := query(client, ctx)
		if err != nil {
			return err
		}
		var grantedPrivileges []string
		for _, grant := range grants {
			grantedPrivileges = append(grantedPrivileges, grant.Privilege)
		}
		notAllPrivilegesInGrantedPrivileges := slices.ContainsFunc(privileges, func(privilege string) bool {
			return !slices.Contains(grantedPrivileges, privilege)
		})
		if notAllPrivilegesInGrantedPrivileges {
			return fmt.Errorf("not every privilege from the list: %v was found in grant privileges: %v, for role name: %s", privileges, grantedPrivileges, roleName.FullyQualifiedName())
		}

		return nil
	}
}

func TestGetFirstNestedObjectByKey(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"int_property": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type: schema.TypeInt,
					},
				},
			},
		},
		"string_property": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type: schema.TypeString,
					},
				},
			},
		},
		"list_property": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type: schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"multiple_list_properties": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type: schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"list": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"empty list": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"not_property": {
			Type: schema.TypeString,
		},
	}, map[string]any{
		"int_property": []any{
			map[string]any{
				"value": 123,
			},
		},
		"string_property": []any{
			map[string]any{
				"value": "some string",
			},
		},
		"list":       []any{"one"},
		"empty_list": []any{},
		"list_property": []any{
			map[string]any{
				"value": []any{"one", "two", "three"},
			},
		},
		"multiple_list_properties": []any{
			map[string]any{
				"value": []any{"one", "two", "three"},
			},
			map[string]any{
				"value": []any{"one", "two", "three"},
			},
		},
		"not_property": "not a property",
	})

	intValue, err := resources.GetPropertyOfFirstNestedObjectByKey[int](d, "int_property", "value")
	assert.NoError(t, err)
	assert.Equal(t, 123, *intValue)

	stringValue, err := resources.GetPropertyOfFirstNestedObjectByKey[string](d, "string_property", "value")
	assert.NoError(t, err)
	assert.Equal(t, "some string", *stringValue)

	listValue, err := resources.GetPropertyOfFirstNestedObjectByKey[[]any](d, "list_property", "value")
	assert.NoError(t, err)
	assert.Equal(t, []any{"one", "two", "three"}, *listValue)

	_, err = resources.GetPropertyOfFirstNestedObjectByKey[any](d, "non_existing_property_key", "non_existing_value_key")
	assert.ErrorContains(t, err, "nested property non_existing_property_key not found")

	_, err = resources.GetPropertyOfFirstNestedObjectByKey[any](d, "not_property", "value")
	assert.ErrorContains(t, err, "nested property not_property is not an array or has incorrect number of values: 0, expected: 1")

	_, err = resources.GetPropertyOfFirstNestedObjectByKey[any](d, "empty_list", "value")
	assert.ErrorContains(t, err, "nested property empty_list not found") // Empty list is a default value, so it's treated as "not set"

	_, err = resources.GetPropertyOfFirstNestedObjectByKey[any](d, "multiple_list_properties", "value")
	assert.ErrorContains(t, err, "nested property multiple_list_properties is not an array or has incorrect number of values: 2, expected: 1")

	_, err = resources.GetPropertyOfFirstNestedObjectByKey[any](d, "list", "value")
	assert.ErrorContains(t, err, "nested property list is not of type map[string]any, got: string")

	_, err = resources.GetPropertyOfFirstNestedObjectByKey[any](d, "int_property", "non_existing_value_key")
	assert.ErrorContains(t, err, "nested value key non_existing_value_key couldn't be found in the nested property map int_property")

	_, err = resources.GetPropertyOfFirstNestedObjectByKey[int](d, "string_property", "value")
	assert.ErrorContains(t, err, "nested property string_property.value is not of type int, got: string")
}
