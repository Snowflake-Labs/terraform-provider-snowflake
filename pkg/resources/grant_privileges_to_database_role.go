package resources

import (
	"context"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"slices"
	"strings"
)

var grantPrivilegesToDatabaseRoleSchema = map[string]*schema.Schema{
	"database_role_name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The fully qualified name of the database role to which privileges will be granted.",
	},
	"privileges": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:         true,
		ValidateDiagFunc: doesNotContainOwnershipGrant(),
		ExactlyOneOf: []string{
			"privileges",
			"all_privileges",
		},
	},
	"all_privileges": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		ExactlyOneOf: []string{
			"privileges",
			"all_privileges",
		},
	},
	"with_grant_option": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"database_name": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  false,
	},
	"on_schema": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"schema_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"all_schemas_in_database": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"future_schemas_in_database": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	},
	"on_schema_object": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"object_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"all": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: grantPrivilegesOnDatabaseRoleBulkOperationSchema,
					},
				},
				"future": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: grantPrivilegesOnDatabaseRoleBulkOperationSchema,
					},
				},
			},
		},
	},
}

var grantPrivilegesOnDatabaseRoleBulkOperationSchema = map[string]*schema.Schema{
	"object_type_plural": {
		Type:     schema.TypeString,
		Required: true,
	},
	"in_database": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"in_schema": {
		Type:     schema.TypeString,
		Optional: true,
	},
}

func doesNotContainOwnershipGrant() func(value any, path cty.Path) diag.Diagnostics {
	return func(value any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if privileges, ok := value.([]string); ok {
			if slices.ContainsFunc(privileges, func(privilege string) bool {
				return strings.ToUpper(privilege) == "OWNERSHIP"
			}) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unsupported privilege type 'OWNERSHIP'.",
					// TODO: Change when a new resource for granting ownership will be available
					Detail:        "Granting ownership is only allowed in dedicated resources (snowflake_user_ownership_grant, snowflake_role_ownership_grant)",
					AttributePath: nil,
				})
			}
		}
		return diags
	}
}

func GrantPrivilegesToDatabaseRole() *schema.Resource {
	return &schema.Resource{
		Create: CreateGrantPrivilegesToRole,
		Read:   ReadGrantPrivilegesToRole,
		Delete: DeleteGrantPrivilegesToRole,
		Update: UpdateGrantPrivilegesToRole,

		Schema: grantPrivilegesToRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				return nil, nil
			},
		},
	}
}
