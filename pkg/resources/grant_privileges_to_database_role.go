package resources

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"slices"
	"strings"
)

// TODO:
//	- resolve todos
// 	- remove logs
// 	- make error messages consistent
//	- write documentation (document on_all, always_apply etc.)
//	- test import

var grantPrivilegesToDatabaseRoleSchema = map[string]*schema.Schema{
	"database_role_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the database role to which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
	},
	"privileges": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The privileges to grant on the database role.",
		ExactlyOneOf: []string{
			"privileges",
			"all_privileges",
		},
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: isNotOwnershipGrant(),
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
		Description: "If specified, allows the recipient role to grant the privileges to other roles.",
	},
	"always_apply": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the resource will always produce a “plan” and on “apply” it will re-grant defined privileges. It is supposed to be used only in “grant privileges on all X’s in database / schema Y” or “grant all privileges to X” scenarios to make sure that every new object in a given database / schema is granted by the account role and every new privilege is granted to the database role. Important note: this flag is not compliant with the Terraform assumptions of the config being eventually convergent (producing an empty plan).",
		// TODO: conflicts with
		//AtLeastOneOf: []string{
		//	"all_privileges",
		//	"on_schema.0.all_schemas_in_database",
		//	"on_schema_object.0.all",
		//},
	},
	"always_apply_trigger": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
		// TODO: Fix desc
		Description: "This field should not be set and its main purpose is to achieve the functionality described by always_apply field. This is value will be flipped to the opposite value on every terraform apply, thus creating a new plan that will re-apply grants.",
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

func isNotOwnershipGrant() func(value any, path cty.Path) diag.Diagnostics {
	return func(value any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if privilege, ok := value.(string); ok && strings.ToUpper(privilege) == "OWNERSHIP" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unsupported privilege 'OWNERSHIP'",
				// TODO: Change when a new resource for granting ownership will be available
				Detail:        "Granting ownership is only allowed in dedicated resources (snowflake_user_ownership_grant, snowflake_role_ownership_grant)",
				AttributePath: nil,
			})
		}
		return diags
	}
}

func GrantPrivilegesToDatabaseRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateGrantPrivilegesToDatabaseRole,
		ReadContext:   ReadGrantPrivilegesToDatabaseRole,
		DeleteContext: DeleteGrantPrivilegesToDatabaseRole,
		UpdateContext: UpdateGrantPrivilegesToDatabaseRole,

		Schema: grantPrivilegesToDatabaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportGrantPrivilegesToDatabaseRole,
		},
	}
}

func ImportGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := ParseGrantPrivilegesToDatabaseRoleId(d.Id())
	if err != nil {
		return nil, err
	}
	if err := d.Set("database_role_name", id.DatabaseRoleName); err != nil {
		return nil, err
	}
	if err := d.Set("with_grant_option", id.WithGrantOption); err != nil {
		return nil, err
	}
	if err := d.Set("always_apply", id.AlwaysApply); err != nil {
		return nil, err
	}
	if err := d.Set("all_privileges", id.AllPrivileges); err != nil {
		return nil, err
	}
	if err := d.Set("privileges", id.Privileges); err != nil {
		return nil, err
	}

	switch id.Kind {
	case OnDatabaseDatabaseRoleGrantKind:
		if err := d.Set("on_database", id.Data.(*OnDatabaseGrantData).DatabaseName.FullyQualifiedName()); err != nil {
			return nil, err
		}
	case OnSchemaDatabaseRoleGrantKind:
		data := id.Data.(*OnSchemaGrantData)
		var onSchema map[string]any

		switch data.Kind {
		case OnSchemaSchemaGrantKind:
			onSchema["schema_name"] = data.SchemaName.FullyQualifiedName()
		case OnAllSchemasInDatabaseSchemaGrantKind:
			onSchema["all_schemas_in_database"] = data.DatabaseName.FullyQualifiedName()
		case OnFutureSchemasInDatabaseSchemaGrantKind:
			onSchema["future_schemas_in_database"] = data.DatabaseName.FullyQualifiedName()
		}

		if err := d.Set("on_schema", []any{onSchema}); err != nil {
			return nil, err
		}
	case OnSchemaObjectDatabaseRoleGrantKind:
		data := id.Data.(*OnSchemaObjectGrantData)
		var onSchemaObject map[string]any

		switch data.Kind {
		case OnObjectSchemaObjectGrantKind:
			onSchemaObject["object_type"] = data.Object.ObjectType.String()
			onSchemaObject["object_name"] = data.Object.Name.FullyQualifiedName()
		case OnAllSchemaObjectGrantKind:
			var onAll map[string]any

			onAll["object_name_plural"] = data.OnAllOrFuture.ObjectNamePlural.String()
			switch data.OnAllOrFuture.Kind {
			case InDatabaseBulkOperationGrantKind:
				onAll["in_database"] = data.OnAllOrFuture.Database.FullyQualifiedName()
			case InSchemaBulkOperationGrantKind:
				onAll["in_schema"] = data.OnAllOrFuture.Schema.FullyQualifiedName()
			}

			onSchemaObject["all"] = []any{onAll}
		case OnFutureSchemaObjectGrantKind:
			var onFuture map[string]any

			onFuture["object_name_plural"] = data.OnAllOrFuture.ObjectNamePlural.String()
			switch data.OnAllOrFuture.Kind {
			case InDatabaseBulkOperationGrantKind:
				onFuture["in_database"] = data.OnAllOrFuture.Database.FullyQualifiedName()
			case InSchemaBulkOperationGrantKind:
				onFuture["in_schema"] = data.OnAllOrFuture.Schema.FullyQualifiedName()
			}

			onSchemaObject["future"] = []any{onFuture}
		}

		if err := d.Set("on_schema_object", []any{onSchemaObject}); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func CreateGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := createGrantPrivilegesToDatabaseRoleIdFromSchema(d)
	err := client.Grants.GrantPrivilegesToDatabaseRole(
		ctx,
		getDatabaseRolePrivilegesFromSchema(d),
		getDatabaseRoleGrantOn(d),
		sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(d.Get("database_role_name").(string)),
		&sdk.GrantPrivilegesToDatabaseRoleOptions{
			WithGrantOption: sdk.Bool(d.Get("with_grant_option").(bool)),
		},
	)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "An error occurred when granting privileges to database role",
			Detail:   fmt.Sprintf("Id: %s\nError: %s", id.DatabaseRoleName, err.Error()),
		})
	}

	d.SetId(id.String())

	return ReadGrantPrivilegesToDatabaseRole(ctx, d, meta)
}

func UpdateGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id, err := ParseGrantPrivilegesToDatabaseRoleId(d.Id())
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse internal identifier",
			Detail:   fmt.Sprintf("Id: %s\nErr: %s", d.Id(), err.Error()), // TODO: link to the documentation (?). It should describe how the identifier looks.
		})
	}

	if d.HasChange("privileges") {
		before, after := d.GetChange("privileges")
		privilegesBeforeChange := expandStringList(before.(*schema.Set).List())
		privilegesAfterChange := expandStringList(after.(*schema.Set).List())

		var privilegesToAdd, privilegesToRemove []string

		for _, privilegeBeforeChange := range privilegesBeforeChange {
			if !slices.Contains(privilegesAfterChange, privilegeBeforeChange) {
				privilegesToRemove = append(privilegesToRemove, privilegeBeforeChange)
			}
		}

		for _, privilegeAfterChange := range privilegesAfterChange {
			if !slices.Contains(privilegesBeforeChange, privilegeAfterChange) {
				privilegesToAdd = append(privilegesToAdd, privilegeAfterChange)
			}
		}

		log.Println("PRIVILEGES TO CHANGE:", privilegesToAdd, privilegesToRemove)
		grantOn := getDatabaseRoleGrantOn(d)

		if len(privilegesToAdd) > 0 {
			err = client.Grants.GrantPrivilegesToDatabaseRole(
				ctx,
				getDatabaseRolePrivileges(
					false,
					privilegesToAdd,
					id.Kind == OnDatabaseDatabaseRoleGrantKind,
					id.Kind == OnSchemaDatabaseRoleGrantKind,
					id.Kind == OnSchemaObjectDatabaseRoleGrantKind,
				),
				grantOn,
				id.DatabaseRoleName,
				new(sdk.GrantPrivilegesToDatabaseRoleOptions),
			)
			if err != nil {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to grant added privileges",
					Detail:   fmt.Sprintf("Id: %s\nErr: %s", d.Id(), err.Error()), // TODO: link to the documentation (?). It should describe how the identifier looks.
				})
			}
		}

		if len(privilegesToRemove) > 0 {
			err = client.Grants.RevokePrivilegesFromDatabaseRole(
				ctx,
				getDatabaseRolePrivileges(
					false,
					privilegesToRemove,
					id.Kind == OnDatabaseDatabaseRoleGrantKind,
					id.Kind == OnSchemaDatabaseRoleGrantKind,
					id.Kind == OnSchemaObjectDatabaseRoleGrantKind,
				),
				grantOn,
				id.DatabaseRoleName,
				new(sdk.RevokePrivilegesFromDatabaseRoleOptions),
			)
			if err != nil {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to revoke removed privileges",
					Detail:   fmt.Sprintf("Id: %s\nErr: %s", d.Id(), err.Error()), // TODO: link to the documentation (?). It should describe how the identifier looks.
				})
			}
		}

		id.Privileges = privilegesAfterChange
	}

	log.Println("UPDATE always_apply has change:", d.HasChange("always_apply"))
	if d.HasChange("always_apply") {
		log.Println("UPDATE always_apply get:", d.Get("always_apply"))
		id.AlwaysApply = d.Get("always_apply").(bool)
	}

	if id.AlwaysApply {
		log.Println("UPDATE applying grants")
		err := client.Grants.GrantPrivilegesToDatabaseRole(
			ctx,
			getDatabaseRolePrivilegesFromSchema(d),
			getDatabaseRoleGrantOn(d),
			id.DatabaseRoleName,
			&sdk.GrantPrivilegesToDatabaseRoleOptions{
				WithGrantOption: &id.WithGrantOption,
			},
		)
		if err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when granting privileges to database role",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", id.DatabaseRoleName, err.Error()),
			})
		}
	}

	d.SetId(id.String())

	return ReadGrantPrivilegesToDatabaseRole(ctx, d, meta)
}

func DeleteGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	id, err := ParseGrantPrivilegesToDatabaseRoleId(d.Id())
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse internal identifier",
			Detail:   fmt.Sprintf("Id: %s\nErr: %s", d.Id(), err.Error()), // TODO: link to the documentation (?). It should describe how the identifier looks.
		})
	}

	err = client.Grants.RevokePrivilegesFromDatabaseRole(
		ctx,
		getDatabaseRolePrivilegesFromSchema(d),
		getDatabaseRoleGrantOn(d),
		id.DatabaseRoleName,
		&sdk.RevokePrivilegesFromDatabaseRoleOptions{},
	)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "An error occurred when revoking privileges from database role",
			Detail:   fmt.Sprintf("Id: %s\nError: %s", id.DatabaseRoleName, err.Error()),
		})
	}

	d.SetId("")

	return diags
}

func ReadGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Println("READ BEGINS")
	var diags diag.Diagnostics

	id, err := ParseGrantPrivilegesToDatabaseRoleId(d.Id())
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse internal identifier",
			Detail:   fmt.Sprintf("Id: %s\nErr: %s", d.Id(), err.Error()), // TODO: link to the documentation (?). It should describe how the identifier looks.
		})
	}

	log.Println("READ id.AlwaysApply:", id.AlwaysApply)
	if id.AlwaysApply {
		// Change the value of always_apply_trigger to produce a plan
		triggerId, _ := uuid.GenerateUUID() // TODO handle error
		log.Printf("READ applying triggerId: %s, was: %s\n", triggerId, d.Get("always_apply_trigger"))
		if err := d.Set("always_apply_trigger", triggerId); err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error setting always_apply_trigger for database role",
				Detail:   fmt.Sprintf("Id: %s\nErr: %s", d.Id(), err.Error()), // TODO: link to the documentation (?). It should describe how the identifier looks.
			})
		}
	}

	if id.AllPrivileges {
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Show with all_privileges option is skipped for now.", // TODO: Details
			Detail:   "<TODO_LINK>",                                         // TODO: link to the design decisions doc
		})
	}

	opts, grantedOn, diagnostics := prepareShowGrantsRequest(id)
	if len(diagnostics) != 0 {
		return append(diags, diagnostics...)
	}

	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	grants, err := client.Grants.Show(ctx, opts)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to retrieve grants",
			Detail:   fmt.Sprintf("Id: %s\nErr: %s", d.Id(), err.Error()), // TODO: link to the documentation (?). It should describe how the identifier looks.
		})
	}

	var privileges []string

	// TODO: Refactor - check if correct with new conventions
	// TODO: Compare privileges
	for _, grant := range grants {
		log.Println("GRANT:", grant)
		// Accept only DATABASE ROLEs
		if grant.GrantTo != sdk.ObjectTypeDatabaseRole && grant.GrantedTo != sdk.ObjectTypeDatabaseRole {
			continue
		}
		// TODO: What about all_privileges, right now we cannot assure that the list of privileges is correct
		// Only consider privileges that are already present in the ID so we
		// don't delete privileges managed by other resources.
		if !slices.Contains(id.Privileges, grant.Privilege) {
			continue
		}
		// TODO: What about GranteeName with database roles is it fully qualified or not ? if yes, refactor GranteeName.
		if id.WithGrantOption == grant.GrantOption && id.DatabaseRoleName.Name() == grant.GranteeName.Name() {
			// future grants do not have grantedBy, only current grants do. If grantedby
			// is an empty string it means the grant could not have been created by terraform
			if (opts.Future == nil || *opts.Future == false) && grant.GrantedBy.Name() == "" {
				continue
			}
			// grant_on is for future grants, granted_on is for current grants. They function the same way though in a test for matching the object type
			if grantedOn == grant.GrantedOn || grantedOn == grant.GrantOn {
				privileges = append(privileges, grant.Privilege)
			}
		}
	}

	log.Println("PRIVILEGES:", id.Privileges, expandStringList(d.Get("privileges").(*schema.Set).List()), privileges)
	if err := d.Set("privileges", privileges); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error setting privileges for database role",
			Detail:   fmt.Sprintf("Id: %s\nErr: %s", d.Id(), err.Error()), // TODO: link to the documentation (?). It should describe how the identifier looks.
		})
	}

	return diags
}

func prepareShowGrantsRequest(id GrantPrivilegesToDatabaseRoleId) (*sdk.ShowGrantOptions, sdk.ObjectType, diag.Diagnostics) {
	opts := new(sdk.ShowGrantOptions)
	var grantedOn sdk.ObjectType
	var diags diag.Diagnostics

	switch id.Kind {
	case OnDatabaseDatabaseRoleGrantKind:
		grantedOn = sdk.ObjectTypeDatabase
		data := id.Data.(*OnDatabaseGrantData)
		opts.On = &sdk.ShowGrantsOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypeDatabase,
				Name:       data.DatabaseName,
			},
		}
	case OnSchemaDatabaseRoleGrantKind:
		grantedOn = sdk.ObjectTypeSchema
		data := id.Data.(*OnSchemaGrantData)

		switch data.Kind {
		case OnSchemaSchemaGrantKind:
			opts.On = &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypeSchema,
					Name:       data.SchemaName,
				},
			}
		case OnAllSchemasInDatabaseSchemaGrantKind:
			// TODO: Document
			return nil, "", append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Skipping",
				Detail:   "TODO",
			})
		case OnFutureSchemasInDatabaseSchemaGrantKind:
			// TODO: show future on database (collisions with other on future triggers and over fetching is ok ?)
			opts.Future = sdk.Bool(true)
			opts.In = &sdk.ShowGrantsIn{
				Database: data.DatabaseName,
			}
		}
	case OnSchemaObjectDatabaseRoleGrantKind:
		data := id.Data.(*OnSchemaObjectGrantData)

		switch data.Kind {
		case OnObjectSchemaObjectGrantKind:
			grantedOn = data.Object.ObjectType
			opts.On = &sdk.ShowGrantsOn{
				Object: data.Object,
			}
		case OnAllSchemaObjectGrantKind:
			// TODO: Document
			return nil, "", append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Skipping",
				Detail:   "TODO",
			})
		case OnFutureSchemaObjectGrantKind:
			grantedOn = data.OnAllOrFuture.ObjectNamePlural.Singular()
			opts.Future = sdk.Bool(true)

			switch data.OnAllOrFuture.Kind {
			case InDatabaseBulkOperationGrantKind:
				opts.In = &sdk.ShowGrantsIn{
					Database: data.OnAllOrFuture.Database,
				}
			case InSchemaBulkOperationGrantKind:
				opts.In = &sdk.ShowGrantsIn{
					Schema: data.OnAllOrFuture.Schema,
				}
			}
		}
	}

	return opts, grantedOn, diags
}

func getDatabaseRolePrivilegesFromSchema(d *schema.ResourceData) *sdk.DatabaseRoleGrantPrivileges {
	_, onDatabaseOk := d.GetOk("on_database")
	_, onSchemaOk := d.GetOk("on_schema")
	_, onSchemaObjectOk := d.GetOk("on_schema_object")

	return getDatabaseRolePrivileges(
		d.Get("all_privileges").(bool),
		expandStringList(d.Get("privileges").(*schema.Set).List()),
		onDatabaseOk,
		onSchemaOk,
		onSchemaObjectOk,
	)
}

func getDatabaseRolePrivileges(allPrivileges bool, privileges []string, onDatabase bool, onSchema bool, onSchemaObject bool) *sdk.DatabaseRoleGrantPrivileges {
	databaseRoleGrantPrivileges := new(sdk.DatabaseRoleGrantPrivileges)

	if allPrivileges {
		databaseRoleGrantPrivileges.AllPrivileges = sdk.Bool(true)
		return databaseRoleGrantPrivileges
	}

	switch {
	case onDatabase:
		databasePrivileges := make([]sdk.AccountObjectPrivilege, len(privileges))
		for i, privilege := range privileges {
			databasePrivileges[i] = sdk.AccountObjectPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.DatabasePrivileges = databasePrivileges
	case onSchema:
		schemaPrivileges := make([]sdk.SchemaPrivilege, len(privileges))
		for i, privilege := range privileges {
			schemaPrivileges[i] = sdk.SchemaPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.SchemaPrivileges = schemaPrivileges
	case onSchemaObject:
		schemaObjectPrivileges := make([]sdk.SchemaObjectPrivilege, len(privileges))
		for i, privilege := range privileges {
			schemaObjectPrivileges[i] = sdk.SchemaObjectPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.SchemaObjectPrivileges = schemaObjectPrivileges
	}

	return databaseRoleGrantPrivileges
}

func getDatabaseRoleGrantOn(d *schema.ResourceData) *sdk.DatabaseRoleGrantOn {
	onDatabase, onDatabaseOk := d.GetOk("on_database")
	onSchemaBlock, onSchemaOk := d.GetOk("on_schema")
	onSchemaObjectBlock, onSchemaObjectOk := d.GetOk("on_schema_object")
	on := new(sdk.DatabaseRoleGrantOn)

	switch {
	case onDatabaseOk:
		on.Database = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(onDatabase.(string)))
	case onSchemaOk:
		onSchema := onSchemaBlock.([]any)[0].(map[string]any)

		grantOnSchema := new(sdk.GrantOnSchema)

		schemaName := onSchema["schema_name"].(string)
		schemaNameOk := len(schemaName) > 0

		allSchemasInDatabase := onSchema["all_schemas_in_database"].(string)
		allSchemasInDatabaseOk := len(allSchemasInDatabase) > 0

		futureSchemasInDatabase := onSchema["future_schemas_in_database"].(string)
		futureSchemasInDatabaseOk := len(futureSchemasInDatabase) > 0

		switch {
		case schemaNameOk:
			grantOnSchema.Schema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schemaName))
		case allSchemasInDatabaseOk:
			grantOnSchema.AllSchemasInDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(allSchemasInDatabase))
		case futureSchemasInDatabaseOk:
			grantOnSchema.FutureSchemasInDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(futureSchemasInDatabase))
		}

		on.Schema = grantOnSchema
	case onSchemaObjectOk:
		onSchemaObject := onSchemaObjectBlock.([]any)[0].(map[string]any)

		grantOnSchemaObject := new(sdk.GrantOnSchemaObject)

		objectType := onSchemaObject["object_type"].(string)
		objectTypeOk := len(objectType) > 0

		objectName := onSchemaObject["object_name"].(string)
		objectNameOk := len(objectName) > 0

		all := onSchemaObject["all"].([]any)
		allOk := len(all) > 0

		future := onSchemaObject["future"].([]any)
		futureOk := len(future) > 0

		switch {
		case objectTypeOk && objectNameOk:
			grantOnSchemaObject.SchemaObject = &sdk.Object{
				ObjectType: sdk.ObjectType(objectType), // TODO: Should we validate it or just cast it
				Name:       sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(objectName),
			}
		case allOk:
			grantOnSchemaObject.All = getGrantOnSchemaObjectIn(all[0].(map[string]any))
		case futureOk:
			grantOnSchemaObject.Future = getGrantOnSchemaObjectIn(future[0].(map[string]any))
		}

		on.SchemaObject = grantOnSchemaObject
	}

	return on
}

func getGrantOnSchemaObjectIn(allOrFuture map[string]any) *sdk.GrantOnSchemaObjectIn {
	pluralObjectType := sdk.PluralObjectType(allOrFuture["object_type_plural"].(string))
	grantOnSchemaObjectIn := &sdk.GrantOnSchemaObjectIn{
		PluralObjectType: pluralObjectType,
	}

	if inDatabase, ok := allOrFuture["in_database"].(string); ok && len(inDatabase) > 0 {
		grantOnSchemaObjectIn.InDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(inDatabase))
	}

	if inSchema, ok := allOrFuture["in_schema"].(string); ok && len(inSchema) > 0 {
		grantOnSchemaObjectIn.InSchema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(inSchema))
	}

	return grantOnSchemaObjectIn
}

func createGrantPrivilegesToDatabaseRoleIdFromSchema(d *schema.ResourceData) *GrantPrivilegesToDatabaseRoleId {
	id := new(GrantPrivilegesToDatabaseRoleId)
	id.DatabaseRoleName = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(d.Get("database_role_name").(string))
	id.AllPrivileges = d.Get("all_privileges").(bool)
	if p, ok := d.GetOk("privileges"); ok {
		id.Privileges = expandStringList(p.(*schema.Set).List())
	}
	id.WithGrantOption = d.Get("with_grant_option").(bool)

	on := getDatabaseRoleGrantOn(d)
	switch {
	case on.Database != nil:
		id.Kind = OnDatabaseDatabaseRoleGrantKind
		id.Data = &OnDatabaseGrantData{
			DatabaseName: *on.Database,
		}
	case on.Schema != nil:
		onSchemaGrantData := new(OnSchemaGrantData)

		switch {
		case on.Schema.Schema != nil:
			onSchemaGrantData.Kind = OnSchemaSchemaGrantKind
			onSchemaGrantData.SchemaName = on.Schema.Schema
		case on.Schema.AllSchemasInDatabase != nil:
			onSchemaGrantData.Kind = OnAllSchemasInDatabaseSchemaGrantKind
			onSchemaGrantData.DatabaseName = on.Schema.AllSchemasInDatabase
		case on.Schema.FutureSchemasInDatabase != nil:
			onSchemaGrantData.Kind = OnFutureSchemasInDatabaseSchemaGrantKind
			onSchemaGrantData.DatabaseName = on.Schema.FutureSchemasInDatabase
		}

		id.Kind = OnSchemaDatabaseRoleGrantKind
		id.Data = onSchemaGrantData
	case on.SchemaObject != nil:
		onSchemaObjectGrantData := new(OnSchemaObjectGrantData)

		switch {
		case on.SchemaObject.SchemaObject != nil:
			onSchemaObjectGrantData.Kind = OnObjectSchemaObjectGrantKind
			onSchemaObjectGrantData.Object = on.SchemaObject.SchemaObject
		case on.SchemaObject.All != nil:
			onSchemaObjectGrantData.Kind = OnAllSchemaObjectGrantKind
			onSchemaObjectGrantData.OnAllOrFuture = getBulkOperationGrantData(on.SchemaObject.All)
		case on.SchemaObject.Future != nil:
			onSchemaObjectGrantData.Kind = OnFutureSchemaObjectGrantKind
			onSchemaObjectGrantData.OnAllOrFuture = getBulkOperationGrantData(on.SchemaObject.Future)
		}

		id.Kind = OnSchemaObjectDatabaseRoleGrantKind
		id.Data = onSchemaObjectGrantData
	}

	return id
}

func getBulkOperationGrantData(in *sdk.GrantOnSchemaObjectIn) *BulkOperationGrantData {
	bulkOperationGrantData := &BulkOperationGrantData{
		ObjectNamePlural: in.PluralObjectType,
	}

	if in.InDatabase != nil {
		bulkOperationGrantData.Kind = InDatabaseBulkOperationGrantKind
		bulkOperationGrantData.Database = in.InDatabase
	}

	if in.InSchema != nil {
		bulkOperationGrantData.Kind = InSchemaBulkOperationGrantKind
		bulkOperationGrantData.Schema = in.InSchema
	}

	return bulkOperationGrantData
}
