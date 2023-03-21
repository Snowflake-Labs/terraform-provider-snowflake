package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validTagPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeApply,
)

var tagGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the tag on which to grant privileges.",
		ForceNew:    true,
	},
	"tag_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the tag on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the tag.",
		Default:      "APPLY",
		ValidateFunc: validation.StringInSlice(validTagPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the tag on which to grant privileges.",
		ForceNew:    true,
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
}

// TagGrant returns a pointer to the resource representing a tag grant.
func TagGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateTagGrant,
			Read:   ReadTagGrant,
			Update: UpdateTagGrant,
			Delete: DeleteTagGrant,

			Schema: tagGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					grantID, err := ParseTagGrantID(d.Id())
					if err != nil {
						return nil, err
					}
					if err := d.Set("tag_name", grantID.ObjectName); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", grantID.SchemaName); err != nil {
						return nil, err
					}
					if err := d.Set("database_name", grantID.DatabaseName); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", grantID.Privilege); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
						return nil, err
					}
					if err := d.Set("roles", grantID.Roles); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validTagPrivileges,
	}
}

// CreateTagGrant implements schema.CreateFunc.
func CreateTagGrant(d *schema.ResourceData, meta interface{}) error {
	tagName := d.Get("tag_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.TagGrant(databaseName, schemaName, tagName)

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewTagGrantID(databaseName, schemaName, tagName, privilege, roles, grantOption)
	d.SetId(grantID.String())

	return ReadTagGrant(d, meta)
}

// ReadTagGrant implements schema.ReadFunc.
func ReadTagGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseTagGrantID(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}

	if err := d.Set("tag_name", grantID.ObjectName); err != nil {
		return err
	}

	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}

	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	builder := snowflake.TagGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

	return readGenericGrant(d, meta, tagGrantSchema, builder, false, validTagPrivileges)
}

// UpdateTagGrant implements schema.UpdateFunc.
func UpdateTagGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update is roles. if nothing changed,
	// nothing to update and we're done.
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd, rolesToRevoke := changeDiff(d, "roles")

	grantID, err := ParseTagGrantID(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.TagGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		rolesToRevoke,
		nil,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		grantID.WithGrantOption,
		rolesToAdd,
		nil,
	); err != nil {
		return err
	}

	return ReadTagGrant(d, meta)
}

// DeleteTagGrant implements schema.DeleteFunc.
func DeleteTagGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseTagGrantID(d.Id())
	if err != nil {
		return err
	}

	builder := snowflake.TagGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)

	return deleteGenericGrant(d, meta, builder)
}

type TagGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewTagGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, withGrantOption bool) *TagGrantID {
	return &TagGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *TagGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func ParseTagGrantID(s string) (*TagGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		withGrantOption := false
		roles := []string{}
		if len(idParts) == 6 {
			withGrantOption = idParts[5] == "true"
			roles = helpers.SplitStringToSlice(idParts[4], ",")
		} else {
			withGrantOption = idParts[4] == "true"
		}
		return &TagGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			ObjectName:      idParts[2],
			Privilege:       idParts[3],
			Roles:           roles,
			WithGrantOption: withGrantOption,
			IsOldID:         true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 6 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 6 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 6", len(idParts))
	}
	return &TagGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		WithGrantOption: idParts[4] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[5], ","),
		IsOldID:         false,
	}, nil
}
