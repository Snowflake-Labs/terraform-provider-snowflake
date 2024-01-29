package resources

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"slices"
)

var grantPrivilegesToShareGrantExactlyOneOfValidation = []string{
	"database_name",
	"schema_name",
	//"function_name",
	"table_name",
	"all_tables_in_schema",
	"tag_name",
	"view_name",
}

var grantPrivilegesToShareSchema = map[string]*schema.Schema{
	"share_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the share on which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"privileges": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "The privileges to grant on the share.",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"database_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "TODO",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"schema_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "TODO",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	//	TODO(SNOW-1021686): Because function identifier contains arguments which are not supported right now
	//"function_name": {
	//	Type:        schema.TypeString,
	//	Optional:    true,
	//	ForceNew:    true,
	//	Description: "TODO",
	//	ValidateDiagFunc: IsValidIdentifier[sdk.FunctionIdentifier](),
	//	ExactlyOneOf: grantPrivilegesToShareGrantExactlyOneOfValidation,
	//},
	"table_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "TODO",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"all_tables_in_schema": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "TODO",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"tag_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "TODO",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"view_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "TODO",
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
		if err := d.Set("share_name", id.ShareName.FullyQualifiedName()); err != nil {
			return nil, err
		}
		if err := d.Set("privileges", id.Privileges); err != nil {
			return nil, err
		}

		switch id.Kind {
		case OnDatabaseShareGrantKind:
			if err := d.Set("database_name", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnSchemaShareGrantKind:
			if err := d.Set("schema_name", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnFunctionShareGrantKind:
			if err := d.Set("function_name", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnTableShareGrantKind:
			if err := d.Set("table_name", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnAllTablesInSchemaShareGrantKind:
			if err := d.Set("all_tables_in_schema", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnTagShareGrantKind:
			if err := d.Set("tag_name", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnViewShareGrantKind:
			if err := d.Set("view_name", id.Identifier.FullyQualifiedName()); err != nil {
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

	err := client.Grants.GrantPrivilegeToShare(ctx, getObjectPrivilegesFromSchema(d), getShareGrantOn(d), id.ShareName)
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

	err = client.Grants.RevokePrivilegeFromShare(ctx, getObjectPrivilegesFromSchema(d), getShareGrantOn(d), id.ShareName)
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

	// TODO: Read for
	var privileges []string

	logging.DebugLogger.Printf("[DEBUG] Filtering grants to be set on account: count = %d", len(grants))
	for _, grant := range grants {
		if grant.GrantedTo != sdk.ObjectTypeShare {
			continue
		}
		// Only consider privileges that are already present in the ID, so we
		// don't delete privileges managed by other resources.
		if !slices.Contains(id.Privileges, grant.Privilege) {
			continue
		}
		if grant.GranteeName.Name() == id.ShareName.Name() { // TODO: id.ShareName should be outside resource identifier (forgot the name)
			if grantedOn == grant.GrantedOn {
				privileges = append(privileges, grant.Privilege)
			}
		}
	}

	logging.DebugLogger.Printf("[DEBUG] Setting privileges: %v", privileges)
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
	id.ShareName = sdk.NewAccountObjectIdentifier(d.Get("share_name").(string))
	id.Privileges = expandStringList(d.Get("privileges").(*schema.Set).List())

	databaseName, databaseNameOk := d.GetOk("database_name")
	schemaName, schemaNameOk := d.GetOk("schema_name")
	tableName, tableNameOk := d.GetOk("table_name")
	allTablesInSchema, allTablesInSchemaOk := d.GetOk("all_tables_in_schema")
	tagName, tagNameOk := d.GetOk("tag_name")
	viewName, viewNameOk := d.GetOk("view_name")

	switch {
	case databaseNameOk:
		id.Kind = OnDatabaseShareGrantKind
		id.Identifier = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(databaseName.(string))
	case schemaNameOk:
		id.Kind = OnSchemaShareGrantKind
		id.Identifier = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schemaName.(string))
	case tableNameOk:
		id.Kind = OnTableShareGrantKind
		id.Identifier = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(tableName.(string))
	case allTablesInSchemaOk:
		id.Kind = OnAllTablesInSchemaShareGrantKind
		id.Identifier = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(allTablesInSchema.(string))
	case tagNameOk:
		id.Kind = OnTagShareGrantKind
		id.Identifier = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(tagName.(string))
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

	databaseName, databaseNameOk := d.GetOk("database_name")
	schemaName, schemaNameOk := d.GetOk("schema_name")
	tableName, tableNameOk := d.GetOk("table_name")
	allTablesInSchema, allTablesInSchemaOk := d.GetOk("all_tables_in_schema")
	tagName, tagNameOk := d.GetOk("tag_name")
	viewName, viewNameOk := d.GetOk("view_name")

	switch {
	case len(databaseName.(string)) > 0 && databaseNameOk:
		grantOn.Database = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(databaseName.(string))
	case len(schemaName.(string)) > 0 && schemaNameOk:
		grantOn.Schema = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schemaName.(string))
	case len(tableName.(string)) > 0 && tableNameOk:
		grantOn.Table = &sdk.OnTable{
			Name: sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(tableName.(string)),
		}
	case len(allTablesInSchema.(string)) > 0 && allTablesInSchemaOk:
		grantOn.Table = &sdk.OnTable{
			AllInSchema: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(allTablesInSchema.(string)),
		}
	case len(tagName.(string)) > 0 && tagNameOk:
		grantOn.Tag = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(tagName.(string))
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
	return opts, "", nil
}
