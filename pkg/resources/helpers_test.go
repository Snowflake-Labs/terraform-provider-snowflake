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

func Test_GetPropertyAsPointer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"second_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"third_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"second_string": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"third_string": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"second_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"third_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}, map[string]interface{}{
		"integer":        123,
		"second_integer": 0,
		"string":         "some string",
		"second_string":  "",
		"boolean":        true,
		"second_boolean": false,
		"invalid":        true,
	})

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Equal(t, true, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[int](d, "second_integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[int](d, "third_integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Nil(t, resources.GetPropertyAsPointer[string](d, "second_integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[string](d, "third_string"))
	assert.Equal(t, true, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "second_boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "third_boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))
}

// TODO [SNOW-1511594]: provide TestResourceDataRaw with working GetRawConfig()
func Test_GetConfigPropertyAsPointerAllowingZeroValue(t *testing.T) {
	t.Skip("TestResourceDataRaw does not set up the ResourceData correctly - GetRawConfig is nil")
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"second_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"third_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"second_string": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"third_string": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"second_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"third_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}, map[string]interface{}{
		"integer":        123,
		"second_integer": 0,
		"string":         "some string",
		"second_string":  "",
		"boolean":        true,
		"second_boolean": false,
		"invalid":        true,
	})

	assert.Equal(t, 123, *resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "integer"))
	assert.Equal(t, 0, *resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "second_integer"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "third_integer"))
	assert.Equal(t, "some string", *resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "string"))
	assert.Equal(t, "", *resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "second_integer"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "third_string"))
	assert.Equal(t, true, *resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "boolean"))
	assert.Equal(t, false, *resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "second_boolean"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "third_boolean"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "invalid"))
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

func TestListDiff(t *testing.T) {
	testCases := []struct {
		Name    string
		Before  []any
		After   []any
		Added   []any
		Removed []any
	}{
		{
			Name:    "no changes",
			Before:  []any{1, 2, 3, 4},
			After:   []any{1, 2, 3, 4},
			Removed: []any{},
			Added:   []any{},
		},
		{
			Name:    "only removed",
			Before:  []any{1, 2, 3, 4},
			After:   []any{},
			Removed: []any{1, 2, 3, 4},
			Added:   []any{},
		},
		{
			Name:    "only added",
			Before:  []any{},
			After:   []any{1, 2, 3, 4},
			Removed: []any{},
			Added:   []any{1, 2, 3, 4},
		},
		{
			Name:    "added repeated items",
			Before:  []any{2},
			After:   []any{1, 2, 1},
			Removed: []any{},
			Added:   []any{1, 1},
		},
		{
			Name:    "removed repeated items",
			Before:  []any{1, 2, 1},
			After:   []any{2},
			Removed: []any{1, 1},
			Added:   []any{},
		},
		{
			Name:    "simple diff: ints",
			Before:  []any{1, 2, 3, 4, 5, 6, 7, 8, 9},
			After:   []any{1, 3, 5, 7, 9, 12, 13, 14},
			Removed: []any{2, 4, 6, 8},
			Added:   []any{12, 13, 14},
		},
		{
			Name:    "simple diff: strings",
			Before:  []any{"one", "two", "three", "four"},
			After:   []any{"five", "two", "four", "six"},
			Removed: []any{"one", "three"},
			Added:   []any{"five", "six"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			added, removed := resources.ListDiff(tc.Before, tc.After)
			assert.Equal(t, tc.Added, added)
			assert.Equal(t, tc.Removed, removed)
		})
	}
}
