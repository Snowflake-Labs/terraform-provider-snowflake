package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantRoleSchema = map[string]*schema.Schema{
	"role_name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The fully qualified name of the role which will be granted to the user or parent role.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"user_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The fully qualified name of the user on which specified role will be granted.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"user_name",
			"parent_role_name",
		},
	},
	"parent_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The fully qualified name of the parent role which will create a parent-child relationship between the roles.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"user_name",
			"parent_role_name",
		},
	},
}

func GrantRole() *schema.Resource {
	return &schema.Resource{
		Create: CreateGrantRole,
		Read:   ReadGrantRole,
		Delete: DeleteGrantRole,
		Schema: grantRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), helpers.IDDelimiter)
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid ID specified: %v, expected <role_name>|<grantee_object_type>|<grantee_identifier>", d.Id())
				}
				if err := d.Set("role_name", parts[0]); err != nil {
					return nil, err
				}
				switch parts[1] {
				case "ROLE":
					if err := d.Set("parent_role_name", parts[2]); err != nil {
						return nil, err
					}
				case "USER":
					if err := d.Set("usere_name", parts[2]); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid object type specified: %v, expected ROLE or USER", parts[1])
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

// CreateGrantRole implements schema.CreateFunc.
func CreateGrantRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	roleName := d.Get("role_name").(string)
	roleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(roleName)
	// format of snowflakeResourceID is <role_identifier>|<object type>|<target_identifier>
	var snowflakeResourceID string
	if parentRoleName, ok := d.GetOk("parent_role_name"); ok && parentRoleName.(string) != "" {
		parentRoleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parentRoleName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(roleIdentifier.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleIdentifier.FullyQualifiedName())
		req := sdk.NewGrantRoleRequest(roleIdentifier, sdk.GrantRole{
			Role: &parentRoleIdentifier,
		})
		if err := client.Roles.Grant(ctx, req); err != nil {
			return err
		}
	} else if userName, ok := d.GetOk("user_name"); ok && userName.(string) != "" {
		userIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(userName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(roleIdentifier.FullyQualifiedName(), sdk.ObjectTypeUser.String(), userIdentifier.FullyQualifiedName())
		req := sdk.NewGrantRoleRequest(roleIdentifier, sdk.GrantRole{
			User: &userIdentifier,
		})
		if err := client.Roles.Grant(ctx, req); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid role grant specified: %v", d)
	}
	d.SetId(snowflakeResourceID)
	return ReadGrantRole(d, meta)
}

func ReadGrantRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	roleName := parts[0]
	roleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(roleName)
	objectType := parts[1]
	targetIdentifier := parts[2]
	ctx := context.Background()
	grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
		Of: &sdk.ShowGrantsOf{
			Role: roleIdentifier,
		},
	})
	if err != nil {
		log.Printf("[DEBUG] role (%s) not found", roleIdentifier.FullyQualifiedName())
		d.SetId("")
		return nil
	}

	var found bool
	for _, grant := range grants {
		if grant.GrantedTo == sdk.ObjectType(objectType) {
			if grant.GranteeName.FullyQualifiedName() == targetIdentifier {
				found = true
				break
			}
		}
	}
	if !found {
		log.Printf("[DEBUG] role grant (%s) not found", d.Id())
		d.SetId("")
	}

	return nil
}

func DeleteGrantRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[0])
	objectType := parts[1]
	granteeName := parts[2]
	ctx := context.Background()
	granteeIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(granteeName)
	switch objectType {
	case "ROLE":
		if err := client.Roles.Revoke(ctx, sdk.NewRevokeRoleRequest(id, sdk.RevokeRole{Role: &granteeIdentifier})); err != nil {
			return err
		}
	case "USER":
		if err := client.Roles.Revoke(ctx, sdk.NewRevokeRoleRequest(id, sdk.RevokeRole{User: &granteeIdentifier})); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid object type specified: %v, expected ROLE or USER", objectType)
	}
	d.SetId("")
	return nil
}
