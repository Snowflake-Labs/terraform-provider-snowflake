package resources

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"slices"
	"strings"
)

var grantPrivilegesToDatabaseRoleSchema = map[string]*schema.Schema{
	"database_role_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the database role to which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
	},
	"privileges": {
		Type:             schema.TypeSet,
		Optional:         true,
		Description:      "The privileges to grant on the database role.",
		ValidateDiagFunc: doesNotContainOwnershipGrant(),
		ExactlyOneOf: []string{
			"privileges",
			"all_privileges",
		},
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"all_privileges": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Grant all privileges on the database role.",
		ExactlyOneOf: []string{
			"privileges",
			"all_privileges",
		},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		ForceNew:    true,
		Description: "Specifies whether the grantee can grant the privileges to other users.",
	},
	"on_database": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the database on which privileges will be granted. If the identifier is not fully qualified (in the form of <db_name>.≤database_role_name>), the command looks for the database role in the current database for the session. All privileges are limited to the database that contains the database role, as well as other objects in the same database.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"on_database",
			"on_schema",
			"on_schema_object",
		},
	},
	"on_schema": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the schema on which privileges will be granted.",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"on_database",
			"on_schema",
			"on_schema_object",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"schema_name": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "The fully qualified name of the schema.",
					ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
					ExactlyOneOf: []string{
						"on_schema.0.schema_name",
						"on_schema.0.all_schemas_in_database",
						"on_schema.0.future_schemas_in_database",
					},
				},
				"all_schemas_in_database": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "The fully qualified name of the database.",
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
					ExactlyOneOf: []string{
						"on_schema.0.schema_name",
						"on_schema.0.all_schemas_in_database",
						"on_schema.0.future_schemas_in_database",
					},
				},
				"future_schemas_in_database": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "The fully qualified name of the database.",
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
					ExactlyOneOf: []string{
						"on_schema.0.schema_name",
						"on_schema.0.all_schemas_in_database",
						"on_schema.0.future_schemas_in_database",
					},
				},
			},
		},
	},
	"on_schema_object": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the schema object on which privileges will be granted.",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"on_database",
			"on_schema",
			"on_schema_object",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_type": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "The object type of the schema object on which privileges will be granted. Valid values are: ALERT | DYNAMIC TABLE | EVENT TABLE | FILE FORMAT | FUNCTION | PROCEDURE | SECRET | SEQUENCE | PIPE | MASKING POLICY | PASSWORD POLICY | ROW ACCESS POLICY | SESSION POLICY | TAG | STAGE | STREAM | TABLE | EXTERNAL TABLE | TASK | VIEW | MATERIALIZED VIEW | NETWORK RULE | PACKAGES POLICY | ICEBERG TABLE",
					RequiredWith: []string{
						"on_schema_object.0.object_name",
					},
					ConflictsWith: []string{
						"on_schema_object.0.all",
						"on_schema_object.0.future",
					},
					ValidateFunc: validation.StringInSlice([]string{
						"ALERT",
						"DYNAMIC TABLE",
						"EVENT TABLE",
						"FILE FORMAT",
						"FUNCTION",
						"PROCEDURE",
						"SECRET",
						"SEQUENCE",
						"PIPE",
						"MASKING POLICY",
						"PASSWORD POLICY",
						"ROW ACCESS POLICY",
						"SESSION POLICY",
						"TAG",
						"STAGE",
						"STREAM",
						"TABLE",
						"EXTERNAL TABLE",
						"TASK",
						"VIEW",
						"MATERIALIZED VIEW",
						"NETWORK RULE",
						"PACKAGES POLICY",
						"ICEBERG TABLE",
					}, true),
				},
				"object_name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "The fully qualified name of the object on which privileges will be granted.",
					RequiredWith: []string{
						"on_schema_object.0.object_type",
					},
					ExactlyOneOf: []string{
						"on_schema_object.0.object_name",
						"on_schema_object.0.all",
						"on_schema_object.0.future",
					},
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
				},
				"all": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantPrivilegesOnDatabaseRoleBulkOperationSchema,
					},
					ConflictsWith: []string{
						"on_schema_object.0.object_type",
					},
					ExactlyOneOf: []string{
						"on_schema_object.0.object_name",
						"on_schema_object.0.all",
						"on_schema_object.0.future",
					},
				},
				"future": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on future objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantPrivilegesOnDatabaseRoleBulkOperationSchema,
					},
					ConflictsWith: []string{
						"on_schema_object.0.object_type",
					},
					ExactlyOneOf: []string{
						"on_schema_object.0.object_name",
						"on_schema_object.0.all",
						"on_schema_object.0.future",
					},
				},
			},
		},
	},
}

var grantPrivilegesOnDatabaseRoleBulkOperationSchema = map[string]*schema.Schema{
	"object_type_plural": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The plural object type of the schema object on which privileges will be granted. Valid values are: ALERTS | DYNAMIC TABLES | EVENT TABLES | FILE FORMATS | FUNCTIONS | PROCEDURES | SECRETS | SEQUENCES | PIPES | MASKING POLICIES | PASSWORD POLICIES | ROW ACCESS POLICIES | SESSION POLICIES | TAGS | STAGES | STREAMS | TABLES | EXTERNAL TABLES | TASKS | VIEWS | MATERIALIZED VIEWS | NETWORK RULES | PACKAGES POLICIES | ICEBERG TABLES",
		ValidateFunc: validation.StringInSlice([]string{
			"ALERTS",
			"DYNAMIC TABLES",
			"EVENT TABLES",
			"FILE FORMATS",
			"FUNCTIONS",
			"PROCEDURES",
			"SECRETS",
			"SEQUENCES",
			"PIPES",
			"MASKING POLICIES",
			"PASSWORD POLICIES",
			"ROW ACCESS POLICIES",
			"SESSION POLICIES",
			"TAGS",
			"STAGES",
			"STREAMS",
			"TABLES",
			"EXTERNAL TABLES",
			"TASKS",
			"VIEWS",
			"MATERIALIZED VIEWS",
			"NETWORK RULES",
			"PACKAGES POLICIES",
			"ICEBERG TABLES",
		}, true),
	},
	"in_database": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"in_schema": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
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
		CreateContext: CreateGrantPrivilegesToDatabaseRole,
		ReadContext:   ReadGrantPrivilegesToDatabaseRole,
		//Delete: DeleteGrantPrivilegesToRole,
		//Update: UpdateGrantPrivilegesToRole,

		Schema: grantPrivilegesToDatabaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				return nil, nil
			},
		},
	}
}

func CreateGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	databaseRoleName := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(d.Get("database_role_name").(string))
	err := client.Grants.GrantPrivilegesToDatabaseRole(
		ctx,
		getDatabaseRolePrivileges(d, getDatabaseRoleGrantOn(d)),
		getDatabaseRoleGrantOn(d),
		databaseRoleName,
		&sdk.GrantPrivilegesToDatabaseRoleOptions{
			WithGrantOption: GetPropertyAsPointer[bool](d, "with_grant_option"),
		},
	)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("An error occurred when granting privileges to database role (%s)", databaseRoleName),
			Detail:   fmt.Sprintf("Error: %s", err.Error()),
		})
		return diags
	}

	// TODO: Identifier d.SetId()

	return ReadGrantPrivilegesToDatabaseRole(ctx, d, meta)
}

func ReadGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	opts := &sdk.ShowGrantOptions{
		Future: nil,
		On:     nil,
		To:     nil,
		Of:     nil,
		In:     nil,
	}

	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	client.Grants.Show(ctx, opts)

	return diags
}

func getDatabaseRolePrivileges(d *schema.ResourceData, on *sdk.DatabaseRoleGrantOn) *sdk.DatabaseRoleGrantPrivileges {
	var databaseRoleGrantPrivileges *sdk.DatabaseRoleGrantPrivileges

	if d.Get("all_privileges").(bool) {
		databaseRoleGrantPrivileges.AllPrivileges = sdk.Bool(true)
		return databaseRoleGrantPrivileges
	}

	var privileges []string
	if p, ok := d.GetOk("privileges"); ok {
		privileges = expandStringList(p.(*schema.Set).List())
	}

	switch {
	case on.Database != nil:
		databasePrivileges := make([]sdk.AccountObjectPrivilege, len(privileges))
		for i, privilege := range privileges {
			databasePrivileges[i] = sdk.AccountObjectPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.DatabasePrivileges = databasePrivileges
	case on.Schema != nil:
		schemaPrivileges := make([]sdk.SchemaPrivilege, len(privileges))
		for i, privilege := range privileges {
			schemaPrivileges[i] = sdk.SchemaPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.SchemaPrivileges = schemaPrivileges
	case on.SchemaObject != nil:
		schemaObjectPrivileges := make([]sdk.SchemaObjectPrivilege, len(privileges))
		for i, privilege := range privileges {
			schemaObjectPrivileges[i] = sdk.SchemaObjectPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.SchemaObjectPrivileges = schemaObjectPrivileges
	}

	return databaseRoleGrantPrivileges
}

func getDatabaseRoleGrantOn(d *schema.ResourceData) *sdk.DatabaseRoleGrantOn {
	onDatabase, onDatabaseOk := GetProperty[string](d, "on_database")
	onSchema, onSchemaOk := GetProperty[map[string]any](d, "on_schema")
	onSchemaObject, onSchemaObjectOk := GetProperty[map[string]any](d, "on_schema_object")
	var on *sdk.DatabaseRoleGrantOn

	switch {
	case onDatabaseOk:
		on.Database = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(onDatabase))
	case onSchemaOk:
		var grantOnSchema sdk.GrantOnSchema

		schemaName, schemaNameOk := onSchema["schema_name"]
		allSchemasInDatabase, allSchemasInDatabaseOk := onSchema["all_schemas_in_database"]
		futureSchemasInDatabase, futureSchemasInDatabaseOk := onSchema["future_schemas_in_database"]

		switch {
		case schemaNameOk:
			grantOnSchema.Schema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schemaName.(string)))
		case allSchemasInDatabaseOk:
			grantOnSchema.AllSchemasInDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(allSchemasInDatabase.(string)))
		case futureSchemasInDatabaseOk:
			grantOnSchema.FutureSchemasInDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(futureSchemasInDatabase.(string)))
		}

		on.Schema = &grantOnSchema
	case onSchemaObjectOk:
		var grantOnSchemaObject sdk.GrantOnSchemaObject

		objectType, objectTypeOk := onSchemaObject["object_type"]
		objectName, objectNameOk := onSchemaObject["object_name"]
		all, allOk := onSchemaObject["all"]
		future, futureOk := onSchemaObject["future"]

		switch {
		case objectTypeOk && objectNameOk:
			grantOnSchemaObject.SchemaObject = &sdk.Object{
				ObjectType: sdk.ObjectType(objectType.(string)), // TODO: Should we validate it or just cast it
				Name:       sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(objectName.(string)),
			}
		case allOk:
			grantOnSchemaObject.All = getGrantOnSchemaObjectIn(all.(map[string]any))
		case futureOk:
			grantOnSchemaObject.Future = getGrantOnSchemaObjectIn(future.(map[string]any))
		}

		on.SchemaObject = &grantOnSchemaObject
	}

	return on
}

func getGrantOnSchemaObjectIn(m map[string]any) *sdk.GrantOnSchemaObjectIn {
	grantOnSchemaObjectIn := &sdk.GrantOnSchemaObjectIn{
		PluralObjectType: sdk.PluralObjectType(m["object_type_plural"].(string)),
	}

	if inDatabase, inDatabaseOk := m["in_database"]; inDatabaseOk {
		grantOnSchemaObjectIn.InDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(inDatabase.(string)))
	}

	if inSchema, inSchemaOk := m["in_schema"]; inSchemaOk {
		grantOnSchemaObjectIn.InSchema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(inSchema.(string)))
	}

	return grantOnSchemaObjectIn
}
