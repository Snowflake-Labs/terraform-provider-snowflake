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

var grantDatabaseRoleSchema = map[string]*schema.Schema{
	"database_role_name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The fully qualified name of the database role which will be granted to share or parent role.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
	},
	"parent_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The fully qualified name of the parent account role which will create a parent-child relationship between the roles.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"parent_role_name",
			"parent_database_role_name",
			"share_name",
		},
	},
	"parent_database_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The fully qualified name of the parent database role which will create a parent-child relationship between the roles.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		ExactlyOneOf: []string{
			"parent_role_name",
			"parent_database_role_name",
			"share_name",
		},
	},
	"share_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The fully qualified name of the share on which privileges will be granted.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"parent_role_name",
			"parent_database_role_name",
			"share_name",
		},
	},
}

func GrantDatabaseRole() *schema.Resource {
	return &schema.Resource{
		Create: CreateGrantDatabaseRole,
		Read:   ReadGrantDatabaseRole,
		Delete: DeleteGrantDatabaseRole,
		Schema: grantDatabaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), helpers.IDDelimiter)
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid ID specified: %v, expected <database_role_name>|<object_type>|<target_identifier>", d.Id())
				}
				if err := d.Set("database_role_name", parts[0]); err != nil {
					return nil, err
				}
				switch parts[1] {
				case "ROLE":
					if err := d.Set("parent_role_name", strings.Trim(parts[2], "\"")); err != nil {
						return nil, err
					}
				case "DATABASE ROLE":
					if err := d.Set("parent_database_role_name", parts[2]); err != nil {
						return nil, err
					}
				case "SHARE":
					if err := d.Set("share_name", parts[2]); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid object type specified: %v, expected ROLE, DATABASE ROLE, or SHARE", parts[1])
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

// CreateGrantDatabaseRole implements schema.CreateFunc.
func CreateGrantDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	databaseRoleName := d.Get("database_role_name").(string)
	databaseRoleIdentifier := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(databaseRoleName)
	// format of snowflakeResourceID is <database_role_identifier>|<object type>|<parent_role_name>
	var snowflakeResourceID string
	if parentRoleName, ok := d.GetOk("parent_role_name"); ok && parentRoleName.(string) != "" {
		parentRoleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parentRoleName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(databaseRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleIdentifier.FullyQualifiedName())
		req := sdk.NewGrantDatabaseRoleRequest(databaseRoleIdentifier).WithAccountRole(parentRoleIdentifier)
		if err := client.DatabaseRoles.Grant(ctx, req); err != nil {
			return err
		}
	} else if parentDatabaseRoleName, ok := d.GetOk("parent_database_role_name"); ok && parentDatabaseRoleName.(string) != "" {
		parentRoleIdentifier := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parentDatabaseRoleName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(databaseRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeDatabaseRole.String(), parentRoleIdentifier.FullyQualifiedName())
		req := sdk.NewGrantDatabaseRoleRequest(databaseRoleIdentifier).WithDatabaseRole(parentRoleIdentifier)
		if err := client.DatabaseRoles.Grant(ctx, req); err != nil {
			return err
		}
	} else if shareName, ok := d.GetOk("share_name"); ok && shareName.(string) != "" {
		shareIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(shareName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(databaseRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeShare.String(), shareIdentifier.FullyQualifiedName())
		req := sdk.NewGrantDatabaseRoleToShareRequest(databaseRoleIdentifier, shareIdentifier)
		if err := client.DatabaseRoles.GrantToShare(ctx, req); err != nil {
			return err
		}
	}
	d.SetId(snowflakeResourceID)
	return ReadGrantDatabaseRole(d, meta)
}

// ReadGrantDatabaseRole implements schema.ReadFunc.
func ReadGrantDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	databaseRoleName := parts[0]
	databaseRoleIdentifier := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(databaseRoleName)
	objectType := parts[1]
	targetIdentifier := parts[2]
	ctx := context.Background()
	grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
		Of: &sdk.ShowGrantsOf{
			DatabaseRole: databaseRoleIdentifier,
		},
	})
	if err != nil {
		log.Printf("[DEBUG] database role (%s) not found", databaseRoleIdentifier.FullyQualifiedName())
		d.SetId("")
		return nil
	}

	var found bool
	for _, grant := range grants {
		if grant.GrantedTo == sdk.ObjectType(objectType) {
			if grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeShare {
				if grant.GranteeName.FullyQualifiedName() == targetIdentifier {
					found = true
					break
				}
			} else {
				/*
					note that grantee_name is not saved as a valid identifier in the
					SHOW GRANTS OF DATABASE ROLE <database_role_name> command
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
		log.Printf("[DEBUG] database role grant (%s) not found", d.Id())
		d.SetId("")
	}

	return nil
}

// DeleteGrantDatabaseRole implements schema.DeleteFunc.
func DeleteGrantDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	parts := strings.Split(d.Id(), "|")
	id := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(parts[0])
	objectType := parts[1]
	granteeName := parts[2]
	ctx := context.Background()
	switch objectType {
	case "ROLE":
		if err := client.DatabaseRoles.Revoke(ctx, sdk.NewRevokeDatabaseRoleRequest(id).WithAccountRole(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(granteeName))); err != nil {
			return err
		}
	case "DATABASE ROLE":
		if err := client.DatabaseRoles.Revoke(ctx, sdk.NewRevokeDatabaseRoleRequest(id).WithDatabaseRole(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(granteeName))); err != nil {
			return err
		}
	case "SHARE":
		if err := client.DatabaseRoles.RevokeFromShare(ctx, sdk.NewRevokeDatabaseRoleFromShareRequest(id, sdk.NewAccountObjectIdentifierFromFullyQualifiedName(granteeName))); err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}
