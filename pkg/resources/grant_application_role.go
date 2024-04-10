package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantApplicationRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the application role to grant.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
	},
	"parent_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The fully qualified name of the account role on which application role will be granted.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"parent_role_name",
			"application_name",
		},
	},
	"application_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The fully qualified name of the application on which application role will be granted.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"parent_role_name",
			"application_name",
		},
	},
}

func GrantApplicationRole() *schema.Resource {
	return &schema.Resource{
		Create: CreateGrantApplicationRole,
		Read:   ReadGrantApplicationRole,
		Delete: DeleteGrantApplicationRole,
		Schema: grantApplicationRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), helpers.IDDelimiter)
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid ID specified: %v, expected <name>|<object_type>|<target_identifier>", d.Id())
				}
				if err := d.Set("name", parts[0]); err != nil {
					return nil, err
				}
				switch parts[1] {
				case "ROLE":
					if err := d.Set("parent_role_name", strings.Trim(parts[2], "\"")); err != nil {
						return nil, err
					}
				case "APPLICATION":
					if err := d.Set("application_name", strings.Trim(parts[2], "\"")); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid object type specified: %v, expected ROLE, APPLICATION", parts[1])
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func CreateGrantApplicationRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	name := d.Get("name").(string)
	applicationRoleIdentifier := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(name)
	// format of snowflakeResourceID is <application_role_identifier>|<object type>|<parent_role_name>
	var snowflakeResourceID string
	if parentRoleName, ok := d.GetOk("parent_role_name"); ok && parentRoleName.(string) != "" {
		parentRoleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parentRoleName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(applicationRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleIdentifier.FullyQualifiedName())
		req := sdk.NewGrantApplicationRoleRequest(applicationRoleIdentifier).WithTo(*sdk.NewKindOfRoleRequest().WithRoleName(&parentRoleIdentifier))
		if err := client.ApplicationRoles.Grant(ctx, req); err != nil {
			return err
		}
	} else if applicationName, ok := d.GetOk("application_name"); ok && applicationName.(string) != "" {
		applicationIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(applicationName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(applicationRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeApplication.String(), applicationIdentifier.FullyQualifiedName())
		req := sdk.NewGrantApplicationRoleRequest(applicationRoleIdentifier).WithTo(*sdk.NewKindOfRoleRequest().WithApplicationName(&applicationIdentifier))
		if err := client.ApplicationRoles.Grant(ctx, req); err != nil {
			return err
		}
	}
	d.SetId(snowflakeResourceID)
	return ReadGrantApplicationRole(d, meta)
}

func ReadGrantApplicationRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	applicationRoleName := parts[0]
	applicationRoleIdentifier := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(applicationRoleName)
	objectType := parts[1]
	targetIdentifier := parts[2]
	ctx := context.Background()
	grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
		Of: &sdk.ShowGrantsOf{
			ApplicationRole: applicationRoleIdentifier,
		},
	})
	if err != nil {
		log.Printf("[DEBUG] application role (%s) not found", applicationRoleIdentifier.FullyQualifiedName())
		d.SetId("")
		return nil
	}

	var found bool
	for _, grant := range grants {
		if grant.GrantedTo == sdk.ObjectType(objectType) {
			if grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeApplication {
				if grant.GranteeName.FullyQualifiedName() == targetIdentifier {
					found = true
					break
				}
			} else {
				/*
					note that grantee_name is not saved as a valid identifier in the
					SHOW GRANTS OF APPLICATION ROLE <application_role_name> command
					for example, "ABC"."test_parent_role" is saved as ABC."test_parent_role"
					or "ABC"."test_parent_role" is saved as ABC.test_parent_role
					and our internal mapper thereby fails to parse it correctly, returning "ABC."test_parent_role"
					so this funny string replacement is needed to make it work
				*/
				s := grant.GranteeName.FullyQualifiedName()
				if !strings.Contains(s, "\"") {
					parts := strings.Split(s, ".")
					s = sdk.NewDatabaseObjectIdentifier(parts[0], parts[1]).FullyQualifiedName()
				} else {
					parts := strings.Split(s, "\".\"")
					if len(parts) < 2 {
						parts = strings.Split(s, "\".")
						if len(parts) < 2 {
							parts = strings.Split(s, ".\"")
							if len(parts) < 2 {
								s = strings.Trim(s, "\"")
								parts = strings.Split(s, ".")
							}
						}
					}
					s = sdk.NewDatabaseObjectIdentifier(parts[0], parts[1]).FullyQualifiedName()
				}
				if s == targetIdentifier {
					found = true
					break
				}
			}
		}
	}
	if !found {
		log.Printf("[DEBUG] application role grant (%s) not found", d.Id())
		d.SetId("")
	}

	return nil
}

func DeleteGrantApplicationRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	parts := strings.Split(d.Id(), "|")
	id := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[0])
	objectType := parts[1]
	granteeName := parts[2]
	ctx := context.Background()
	switch objectType {
	case "ROLE":
		applicationRoleName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(granteeName)
		if err := client.ApplicationRoles.Revoke(ctx, sdk.NewRevokeApplicationRoleRequest(id).WithFrom(*sdk.NewKindOfRoleRequest().WithRoleName(&applicationRoleName))); err != nil {
			return err
		}
	case "APPLICATION":
		applicationName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(granteeName)
		if err := client.ApplicationRoles.Revoke(ctx, sdk.NewRevokeApplicationRoleRequest(id).WithFrom(*sdk.NewKindOfRoleRequest().WithApplicationName(&applicationName))); err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}
