package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantPrivilegesToShareGrantExactlyOneOfValidation = []string{
	"on_database",
	"on_schema",
	// TODO(SNOW-990811): "function_name",
	"on_table",
	"on_all_tables_in_schema",
	"on_tag",
	"on_view",
}

var grantPrivilegesToShareSchema = map[string]*schema.Schema{
	"to_share": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the share on which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"privileges": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "The privileges to grant on the share. See available list of privileges: https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-share#syntax",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"on_database": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the database on which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_schema": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the schema on which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	//	TODO(SNOW-1021686): Because function identifier contains arguments which are not supported right now
	// "function_name": {
	//	Type:        schema.TypeString,
	//	Optional:    true,
	//	ForceNew:    true,
	//	Description: "The fully qualified name of the function on which privileges will be granted.",
	//	ValidateDiagFunc: IsValidIdentifier[sdk.FunctionIdentifier](),
	//	ExactlyOneOf: grantPrivilegesToShareGrantExactlyOneOfValidation,
	// },
	"on_table": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the table on which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_all_tables_in_schema": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified identifier for the schema for which the specified privilege will be granted for all tables.",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_tag": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the tag on which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_view": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the view on which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
}

func GrantPrivilegesToShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateGrantPrivilegesToShare,
		UpdateContext: UpdateGrantPrivilegesToShare,
		DeleteContext: DeleteGrantPrivilegesToShare,
		ReadContext:   ReadGrantPrivilegesToShare,

		Schema: grantPrivilegesToShareSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportGrantPrivilegesToShare(),
		},
	}
}

func ImportGrantPrivilegesToShare() func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		id, err := ParseGrantPrivilegesToShareId(d.Id())
		if err != nil {
			return nil, err
		}
		if err := d.Set("to_share", id.ShareName.Name()); err != nil {
			return nil, err
		}
		if err := d.Set("privileges", id.Privileges); err != nil {
			return nil, err
		}

		switch id.Kind {
		case OnDatabaseShareGrantKind:
			if err := d.Set("on_database", id.Identifier.Name()); err != nil {
				return nil, err
			}
		case OnSchemaShareGrantKind:
			if err := d.Set("on_schema", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		// TODO(SNOW-990811) case OnFunctionShareGrantKind:
		//	if err := d.Set("function_name", id.Identifier.FullyQualifiedName()); err != nil {
		//		return nil, err
		//	}
		case OnTableShareGrantKind:
			if err := d.Set("on_table", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnAllTablesInSchemaShareGrantKind:
			if err := d.Set("on_all_tables_in_schema", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnTagShareGrantKind:
			if err := d.Set("on_tag", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnViewShareGrantKind:
			if err := d.Set("on_view", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		}

		return []*schema.ResourceData{d}, nil
	}
}

func CreateGrantPrivilegesToShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id := createGrantPrivilegesToShareIdFromSchema(d)
	log.Printf("[DEBUG] created identifier from schema: %s", id.String())

	err := client.Grants.GrantPrivilegeToShare(ctx, getObjectPrivilegesFromSchema(d), getShareGrantOn(d), sdk.NewAccountObjectIdentifier(id.ShareName.Name()))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when granting privileges to share",
				Detail:   fmt.Sprintf("Id: %s\nShare name: %s\nError: %s", id.String(), id.ShareName, err.Error()),
			},
		}
	}

	d.SetId(id.String())

	return ReadGrantPrivilegesToShare(ctx, d, meta)
}

func UpdateGrantPrivilegesToShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id, err := ParseGrantPrivilegesToShareId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	if d.HasChange("privileges") {
		oldPrivileges, newPrivileges := d.GetChange("privileges")
		privilegesBeforeChange := expandStringList(oldPrivileges.(*schema.Set).List())
		privilegesAfterChange := expandStringList(newPrivileges.(*schema.Set).List())

		var privilegesToAdd, privilegesToRemove []sdk.ObjectPrivilege

		for _, privilegeBeforeChange := range privilegesBeforeChange {
			if !slices.Contains(privilegesAfterChange, privilegeBeforeChange) {
				privilegesToRemove = append(privilegesToRemove, sdk.ObjectPrivilege(privilegeBeforeChange))
			}
		}

		for _, privilegeAfterChange := range privilegesAfterChange {
			if !slices.Contains(privilegesBeforeChange, privilegeAfterChange) {
				privilegesToAdd = append(privilegesToAdd, sdk.ObjectPrivilege(privilegeAfterChange))
			}
		}

		grantOn := getShareGrantOn(d)

		if len(privilegesToAdd) > 0 {
			err = client.Grants.GrantPrivilegeToShare(
				ctx,
				privilegesToAdd,
				grantOn,
				sdk.NewAccountObjectIdentifier(id.ShareName.Name()),
			)
			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to grant added privileges",
						Detail:   fmt.Sprintf("Id: %s\nPrivileges to add: %v\nError: %s", d.Id(), privilegesToAdd, err.Error()),
					},
				}
			}
		}

		if len(privilegesToRemove) > 0 {
			err = client.Grants.RevokePrivilegeFromShare(
				ctx,
				privilegesToRemove,
				grantOn,
				sdk.NewAccountObjectIdentifier(id.ShareName.Name()),
			)
			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to revoke removed privileges",
						Detail:   fmt.Sprintf("Id: %s\nPrivileges to remove: %v\nError: %s", d.Id(), privilegesToRemove, err.Error()),
					},
				}
			}
		}

		id.Privileges = privilegesAfterChange
		d.SetId(id.String())
	}

	return ReadGrantPrivilegesToShare(ctx, d, meta)
}

func DeleteGrantPrivilegesToShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id, err := ParseGrantPrivilegesToShareId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	err = client.Grants.RevokePrivilegeFromShare(ctx, getObjectPrivilegesFromSchema(d), getShareGrantOn(d), sdk.NewAccountObjectIdentifier(id.ShareName.Name()))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when revoking privileges from share",
				Detail:   fmt.Sprintf("Id: %s\nShare name: %s\nError: %s", d.Id(), id.ShareName.FullyQualifiedName(), err.Error()),
			},
		}
	}

	d.SetId("")

	return nil
}

func ReadGrantPrivilegesToShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := ParseGrantPrivilegesToShareId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	opts, grantedOn, diags := prepareShowGrantsRequestForShare(id)
	if len(diags) != 0 {
		return diags
	}

	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	grants, err := client.Grants.Show(ctx, opts)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve grants",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	var privileges []string
	for _, grant := range grants {
		if grant.GrantedTo != sdk.ObjectTypeShare {
			continue
		}
		// Only consider privileges that are already present in the ID, so we
		// don't delete privileges managed by other resources.
		if !slices.Contains(id.Privileges, grant.Privilege) {
			continue
		}
		if grant.GranteeName.Name() == id.ShareName.Name() {
			if grantedOn == grant.GrantedOn {
				privileges = append(privileges, grant.Privilege)
			}
		}
	}

	// REFERENCE_USAGE is a special pseudo-privilege that you can grant or revoke,
	// but it won't show up when querying privileges (not returned by show grants ... query).
	// That's why we have to check it manually outside the loop and append it whenever it's specified in the configuration.
	if slices.Contains(id.Privileges, sdk.ObjectPrivilegeReferenceUsage.String()) {
		privileges = append(privileges, sdk.ObjectPrivilegeReferenceUsage.String())
	}

	if err := d.Set("privileges", privileges); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error setting privileges for account role",
				Detail:   fmt.Sprintf("Id: %s\nPrivileges: %v\nError: %s", d.Id(), privileges, err.Error()),
			},
		}
	}

	return nil
}

func createGrantPrivilegesToShareIdFromSchema(d *schema.ResourceData) *GrantPrivilegesToShareId {
	id := new(GrantPrivilegesToShareId)
	id.ShareName = sdk.NewAccountObjectIdentifier(d.Get("to_share").(string))
	id.Privileges = expandStringList(d.Get("privileges").(*schema.Set).List())

	databaseName, databaseNameOk := d.GetOk("on_database")
	schemaName, schemaNameOk := d.GetOk("on_schema")
	// TODO(SNOW-990811) functionName, functionNameOk := d.GetOk("function_name")
	tableName, tableNameOk := d.GetOk("on_table")
	allTablesInSchema, allTablesInSchemaOk := d.GetOk("on_all_tables_in_schema")
	tagName, tagNameOk := d.GetOk("on_tag")
	viewName, viewNameOk := d.GetOk("on_view")

	switch {
	case databaseNameOk:
		id.Kind = OnDatabaseShareGrantKind
		id.Identifier = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(databaseName.(string))
	case schemaNameOk:
		id.Kind = OnSchemaShareGrantKind
		id.Identifier = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schemaName.(string))
	// TODO(SNOW-990811) case functionNameOk:
	//	id.Kind = OnFunctionShareGrantKind
	//	id.Identifier = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(functionName.(string))
	case tableNameOk:
		id.Kind = OnTableShareGrantKind
		id.Identifier = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(tableName.(string))
	case allTablesInSchemaOk:
		id.Kind = OnAllTablesInSchemaShareGrantKind
		id.Identifier = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(allTablesInSchema.(string))
	case tagNameOk:
		id.Kind = OnTagShareGrantKind
		id.Identifier = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(tagName.(string))
	case viewNameOk:
		id.Kind = OnViewShareGrantKind
		id.Identifier = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(viewName.(string))
	}

	return id
}

func getObjectPrivilegesFromSchema(d *schema.ResourceData) []sdk.ObjectPrivilege {
	privileges := expandStringList(d.Get("privileges").(*schema.Set).List())
	objectPrivileges := make([]sdk.ObjectPrivilege, len(privileges))
	for i, privilege := range privileges {
		objectPrivileges[i] = sdk.ObjectPrivilege(privilege)
	}
	return objectPrivileges
}

func getShareGrantOn(d *schema.ResourceData) *sdk.ShareGrantOn {
	grantOn := new(sdk.ShareGrantOn)

	databaseName, databaseNameOk := d.GetOk("on_database")
	schemaName, schemaNameOk := d.GetOk("on_schema")
	// TODO(SNOW-990811) functionName, functionNameOk := d.GetOk("on_function")
	tableName, tableNameOk := d.GetOk("on_table")
	allTablesInSchema, allTablesInSchemaOk := d.GetOk("on_all_tables_in_schema")
	tagName, tagNameOk := d.GetOk("on_tag")
	viewName, viewNameOk := d.GetOk("on_view")

	switch {
	case len(databaseName.(string)) > 0 && databaseNameOk:
		grantOn.Database = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(databaseName.(string))
	case len(schemaName.(string)) > 0 && schemaNameOk:
		grantOn.Schema = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schemaName.(string))
	// TODO(SNOW-990811) case len(functionName.(string)) > 0 && functionNameOk:
	//	grantOn.Function = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(functionName.(string))
	case len(tableName.(string)) > 0 && tableNameOk:
		grantOn.Table = &sdk.OnTable{
			Name: sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(tableName.(string)),
		}
	case len(allTablesInSchema.(string)) > 0 && allTablesInSchemaOk:
		grantOn.Table = &sdk.OnTable{
			AllInSchema: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(allTablesInSchema.(string)),
		}
	case len(tagName.(string)) > 0 && tagNameOk:
		grantOn.Tag = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(tagName.(string))
	case len(viewName.(string)) > 0 && viewNameOk:
		grantOn.View = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(viewName.(string))
	}

	return grantOn
}

func prepareShowGrantsRequestForShare(id GrantPrivilegesToShareId) (*sdk.ShowGrantOptions, sdk.ObjectType, diag.Diagnostics) {
	opts := new(sdk.ShowGrantOptions)
	var objectType sdk.ObjectType

	switch id.Kind {
	case OnDatabaseShareGrantKind:
		objectType = sdk.ObjectTypeDatabase
	case OnSchemaShareGrantKind:
		objectType = sdk.ObjectTypeSchema
	case OnTableShareGrantKind:
		objectType = sdk.ObjectTypeTable
	case OnAllTablesInSchemaShareGrantKind:
		return nil, "", diag.Diagnostics{
			diag.Diagnostic{
				// TODO: link to the design decisions doc (SNOW-990811)
				Severity: diag.Warning,
				Summary:  "Show with OnAll option is skipped.",
				Detail:   "See our document on design decisions for grants: <LINK (coming soon)>",
			},
		}
	case OnTagShareGrantKind:
		objectType = sdk.ObjectTypeTag
	case OnViewShareGrantKind:
		objectType = sdk.ObjectTypeView
	}

	opts.On = &sdk.ShowGrantsOn{
		Object: &sdk.Object{
			ObjectType: objectType,
			Name:       id.Identifier,
		},
	}

	return opts, objectType, nil
}
