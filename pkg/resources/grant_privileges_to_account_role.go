package resources

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var grantPrivilegesToAccountRoleSchema = map[string]*schema.Schema{
	"account_role_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the account role to which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	// According to docs https://docs.snowflake.com/en/user-guide/data-exchange-marketplace-privileges#usage-notes IMPORTED PRIVILEGES
	// will be returned as USAGE in SHOW GRANTS command. In addition, USAGE itself is a valid privilege, but both cannot be set at the
	// same time (IMPORTED PRIVILEGES can only be granted to the database created from SHARE and USAGE in every other case).
	// To handle both cases, additional logic was added in read operation where IMPORTED PRIVILEGES is replaced with USAGE.
	"privileges": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The privileges to grant on the account role.",
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
		Description: "Grant all privileges on the account role.",
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
	"always_apply": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the resource will always produce a “plan” and on “apply” it will re-grant defined privileges. It is supposed to be used only in “grant privileges on all X’s in database / schema Y” or “grant all privileges to X” scenarios to make sure that every new object in a given database / schema is granted by the account role and every new privilege is granted to the database role. Important note: this flag is not compliant with the Terraform assumptions of the config being eventually convergent (producing an empty plan).",
	},
	"always_apply_trigger": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "This is a helper field and should not be set. Its main purpose is to help to achieve the functionality described by the always_apply field.",
	},
	"on_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		ForceNew:    true,
		Description: "If true, the privileges will be granted on the account.",
		ExactlyOneOf: []string{
			"on_account",
			"on_account_object",
			"on_schema",
			"on_schema_object",
		},
	},
	"on_account_object": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the account object on which privileges will be granted ",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"on_account",
			"on_account_object",
			"on_schema",
			"on_schema_object",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_type": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "The object type of the account object on which privileges will be granted. Valid values are: USER | RESOURCE MONITOR | WAREHOUSE | COMPUTE POOL | DATABASE | INTEGRATION | FAILOVER GROUP | REPLICATION GROUP | EXTERNAL VOLUME",
					ValidateFunc: validation.StringInSlice([]string{
						"USER",
						"RESOURCE MONITOR",
						"WAREHOUSE",
						"COMPUTE POOL",
						"DATABASE",
						"INTEGRATION",
						"FAILOVER GROUP",
						"REPLICATION GROUP",
						"EXTERNAL VOLUME",
					}, true),
				},
				"object_name": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         true,
					Description:      "The fully qualified name of the object on which privileges will be granted.",
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
			},
		},
	},
	"on_schema": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the schema on which privileges will be granted.",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"on_account",
			"on_account_object",
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
			"on_account",
			"on_account_object",
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
					ValidateDiagFunc: ValidGrantedObjectType(),
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
						Schema: grantPrivilegesOnAccountRoleBulkOperationSchema,
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
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantPrivilegesOnAccountRoleBulkOperationSchema,
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

var grantPrivilegesOnAccountRoleBulkOperationSchema = map[string]*schema.Schema{
	"object_type_plural": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The plural object type of the schema object on which privileges will be granted. Valid values are: ALERTS | DYNAMIC TABLES | EVENT TABLES | FILE FORMATS | FUNCTIONS | PROCEDURES | SECRETS | SEQUENCES | PIPES | MASKING POLICIES | PASSWORD POLICIES | ROW ACCESS POLICIES | SESSION POLICIES | TAGS | STAGES | STREAMS | TABLES | EXTERNAL TABLES | TASKS | VIEWS | MATERIALIZED VIEWS | NETWORK RULES | PACKAGES POLICIES | ICEBERG TABLES",
		ValidateDiagFunc: ValidGrantedPluralObjectType(),
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

func GrantPrivilegesToAccountRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateGrantPrivilegesToAccountRole,
		UpdateContext: UpdateGrantPrivilegesToAccountRole,
		DeleteContext: DeleteGrantPrivilegesToAccountRole,
		ReadContext:   ReadGrantPrivilegesToAccountRole,

		Schema: grantPrivilegesToAccountRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportGrantPrivilegesToAccountRole(),
		},
	}
}

func ImportGrantPrivilegesToAccountRole() func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		logging.DebugLogger.Printf("[DEBUG] Entering import grant privileges to account role")
		id, err := ParseGrantPrivilegesToAccountRoleId(d.Id())
		if err != nil {
			return nil, err
		}
		logging.DebugLogger.Printf("[DEBUG] Imported identifier: %s", id.String())
		if err := d.Set("account_role_name", id.RoleName.FullyQualifiedName()); err != nil {
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
		if err := d.Set("on_account", false); err != nil {
			return nil, err
		}

		switch id.Kind {
		case OnAccountAccountRoleGrantKind:
			if err := d.Set("on_account", true); err != nil {
				return nil, err
			}
		case OnAccountObjectAccountRoleGrantKind:
			data := id.Data.(*OnAccountObjectGrantData)
			onAccountObject := make(map[string]any)
			onAccountObject["object_type"] = data.ObjectType.String()
			onAccountObject["object_name"] = data.ObjectName.FullyQualifiedName()

			if err := d.Set("on_account_object", []any{onAccountObject}); err != nil {
				return nil, err
			}
		case OnSchemaAccountRoleGrantKind:
			data := id.Data.(*OnSchemaGrantData)
			onSchema := make(map[string]any)

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
		case OnSchemaObjectAccountRoleGrantKind:
			data := id.Data.(*OnSchemaObjectGrantData)
			onSchemaObject := make(map[string]any)

			switch data.Kind {
			case OnObjectSchemaObjectGrantKind:
				onSchemaObject["object_type"] = data.Object.ObjectType.String()
				onSchemaObject["object_name"] = data.Object.Name.FullyQualifiedName()
			case OnAllSchemaObjectGrantKind:
				onAll := make(map[string]any)

				onAll["object_type_plural"] = data.OnAllOrFuture.ObjectNamePlural.String()
				switch data.OnAllOrFuture.Kind {
				case InDatabaseBulkOperationGrantKind:
					onAll["in_database"] = data.OnAllOrFuture.Database.FullyQualifiedName()
				case InSchemaBulkOperationGrantKind:
					onAll["in_schema"] = data.OnAllOrFuture.Schema.FullyQualifiedName()
				}

				onSchemaObject["all"] = []any{onAll}
			case OnFutureSchemaObjectGrantKind:
				onFuture := make(map[string]any)

				onFuture["object_type_plural"] = data.OnAllOrFuture.ObjectNamePlural.String()
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
}

func CreateGrantPrivilegesToAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	logging.DebugLogger.Printf("[DEBUG] Entering create grant privileges to account role")
	client := meta.(*provider.Context).Client

	id := createGrantPrivilegesToAccountRoleIdFromSchema(d)
	logging.DebugLogger.Printf("[DEBUG] created identifier from schema: %s", id.String())

	err := client.Grants.GrantPrivilegesToAccountRole(
		ctx,
		getAccountRolePrivilegesFromSchema(d),
		getAccountRoleGrantOn(d),
		sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("account_role_name").(string)),
		&sdk.GrantPrivilegesToAccountRoleOptions{
			WithGrantOption: sdk.Bool(d.Get("with_grant_option").(bool)),
		},
	)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when granting privileges to account role",
				Detail:   fmt.Sprintf("Id: %s\nAccount role name: %s\nError: %s", id.String(), id.RoleName, err),
			},
		}
	}

	logging.DebugLogger.Printf("[DEBUG] Setting identifier to %s", id.String())
	d.SetId(id.String())

	return ReadGrantPrivilegesToAccountRole(ctx, d, meta)
}

func UpdateGrantPrivilegesToAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	logging.DebugLogger.Printf("[DEBUG] Entering update grant privileges to account role")
	client := meta.(*provider.Context).Client

	id, err := ParseGrantPrivilegesToAccountRoleId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}
	logging.DebugLogger.Printf("[DEBUG] Parsed identifier to %s", id.String())

	// handle all_privileges -> privileges change (revoke all privileges)
	if d.HasChange("all_privileges") {
		_, allPrivileges := d.GetChange("all_privileges")

		if !allPrivileges.(bool) {
			logging.DebugLogger.Printf("[DEBUG] Revoking all privileges")
			err = client.Grants.RevokePrivilegesFromAccountRole(ctx, &sdk.AccountRoleGrantPrivileges{
				AllPrivileges: sdk.Bool(true),
			},
				getAccountRoleGrantOn(d),
				id.RoleName,
				new(sdk.RevokePrivilegesFromAccountRoleOptions),
			)

			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to revoke all privileges",
						Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
					},
				}
			}
		}

		id.AllPrivileges = allPrivileges.(bool)
	}

	if d.HasChange("privileges") {
		shouldHandlePrivilegesChange := true

		// Skip if all_privileges was set to true
		if d.HasChange("all_privileges") {
			if _, allPrivileges := d.GetChange("all_privileges"); allPrivileges.(bool) {
				shouldHandlePrivilegesChange = false
				id.Privileges = []string{}
			}
		}

		if shouldHandlePrivilegesChange {
			before, after := d.GetChange("privileges")
			privilegesBeforeChange := expandStringList(before.(*schema.Set).List())
			privilegesAfterChange := expandStringList(after.(*schema.Set).List())

			logging.DebugLogger.Printf("[DEBUG] Changes in privileges. Before: %v, after: %v", privilegesBeforeChange, privilegesAfterChange)

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

			grantOn := getAccountRoleGrantOn(d)

			if len(privilegesToAdd) > 0 {
				logging.DebugLogger.Printf("[DEBUG] Granting privileges: %v", privilegesToAdd)
				err = client.Grants.GrantPrivilegesToAccountRole(
					ctx,
					getAccountRolePrivileges(
						false,
						privilegesToAdd,
						id.Kind == OnAccountAccountRoleGrantKind,
						id.Kind == OnAccountObjectAccountRoleGrantKind,
						id.Kind == OnSchemaAccountRoleGrantKind,
						id.Kind == OnSchemaObjectAccountRoleGrantKind,
					),
					grantOn,
					id.RoleName,
					new(sdk.GrantPrivilegesToAccountRoleOptions),
				)
				if err != nil {
					return diag.Diagnostics{
						diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Failed to grant added privileges",
							Detail:   fmt.Sprintf("Id: %s\nPrivileges to add: %v\nError: %s", d.Id(), privilegesToAdd, err),
						},
					}
				}
			}

			if len(privilegesToRemove) > 0 {
				logging.DebugLogger.Printf("[DEBUG] Revoking privileges: %v", privilegesToRemove)
				err = client.Grants.RevokePrivilegesFromAccountRole(
					ctx,
					getAccountRolePrivileges(
						false,
						privilegesToRemove,
						id.Kind == OnAccountAccountRoleGrantKind,
						id.Kind == OnAccountObjectAccountRoleGrantKind,
						id.Kind == OnSchemaAccountRoleGrantKind,
						id.Kind == OnSchemaObjectAccountRoleGrantKind,
					),
					grantOn,
					id.RoleName,
					new(sdk.RevokePrivilegesFromAccountRoleOptions),
				)
				if err != nil {
					return diag.Diagnostics{
						diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Failed to revoke removed privileges",
							Detail:   fmt.Sprintf("Id: %s\nPrivileges to remove: %v\nError: %s", d.Id(), privilegesToRemove, err),
						},
					}
				}
			}

			id.Privileges = privilegesAfterChange
		}
	}

	// handle privileges -> all_privileges change (grant all privileges)
	if d.HasChange("all_privileges") {
		_, allPrivileges := d.GetChange("all_privileges")

		if allPrivileges.(bool) {
			logging.DebugLogger.Printf("[DEBUG] Granting all privileges")
			err = client.Grants.GrantPrivilegesToAccountRole(ctx, &sdk.AccountRoleGrantPrivileges{
				AllPrivileges: sdk.Bool(true),
			},
				getAccountRoleGrantOn(d),
				id.RoleName,
				new(sdk.GrantPrivilegesToAccountRoleOptions),
			)

			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to grant all privileges",
						Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
					},
				}
			}
		}

		id.AllPrivileges = allPrivileges.(bool)
	}

	if d.HasChange("always_apply") {
		id.AlwaysApply = d.Get("always_apply").(bool)
	}

	if id.AlwaysApply {
		logging.DebugLogger.Printf("[DEBUG] Performing always_apply re-grant")
		err := client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			getAccountRolePrivilegesFromSchema(d),
			getAccountRoleGrantOn(d),
			id.RoleName,
			&sdk.GrantPrivilegesToAccountRoleOptions{
				WithGrantOption: &id.WithGrantOption,
			},
		)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Always apply. An error occurred when granting privileges to account role",
					Detail:   fmt.Sprintf("Id: %s\nAccount role name: %s\nError: %s", d.Id(), id.RoleName, err),
				},
			}
		}
	}

	logging.DebugLogger.Printf("[DEBUG] Setting identifier to %s", id.String())
	d.SetId(id.String())

	return ReadGrantPrivilegesToAccountRole(ctx, d, meta)
}

func DeleteGrantPrivilegesToAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	logging.DebugLogger.Printf("[DEBUG] Entering delete grant privileges to account role")
	client := meta.(*provider.Context).Client

	id, err := ParseGrantPrivilegesToAccountRoleId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}
	logging.DebugLogger.Printf("[DEBUG] Parsed identifier: %s", id.String())

	err = client.Grants.RevokePrivilegesFromAccountRole(
		ctx,
		getAccountRolePrivilegesFromSchema(d),
		getAccountRoleGrantOn(d),
		id.RoleName,
		&sdk.RevokePrivilegesFromAccountRoleOptions{},
	)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when revoking privileges from account role",
				Detail:   fmt.Sprintf("Id: %s\nAccount role name: %s\nError: %s", d.Id(), id.RoleName.FullyQualifiedName(), err),
			},
		}
	}

	d.SetId("")

	return nil
}

func ReadGrantPrivilegesToAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	logging.DebugLogger.Printf("[DEBUG] Entering read grant privileges to role")
	id, err := ParseGrantPrivilegesToAccountRoleId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}
	logging.DebugLogger.Printf("[DEBUG] Parsed identifier: %s", id.String())

	if id.AlwaysApply {
		// The Trigger is a string rather than boolean that would be flipped on every terraform apply
		// because it's easier to think about and not to worry about edge cases that may occur with 1bit values.
		// The only place to have the "flip" is Read operation, because there we can set value and produce a plan
		// that later on will be executed in the Update operation.
		//
		// The following example shows that we can end up with the same value as before, which may lead to empty plans:
		// 1. Create configuration with always_apply = false (let's say trigger will be false by default)
		// 2. terraform apply: Create (Read will update it to false)
		// 3. Update config so that always_apply = true
		// 4. terraform apply: Read (updated trigger to false) -> change is not detected (no plan; no Update)
		triggerId, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to generate UUID",
					Detail:   fmt.Sprintf("Original error: %s", err),
				},
			}
		}

		// Change the value of always_apply_trigger to produce a plan
		if err := d.Set("always_apply_trigger", triggerId); err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error setting always_apply_trigger for database role",
					Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
				},
			}
		}
	}

	if id.AllPrivileges {
		log.Printf("[INFO] Show with all_privileges option is skipped. No changes in privileges in Snowflake will be detected. Consider specyfying all privileges in 'privileges' block.")
		return nil
	}

	opts, grantedOn := prepareShowGrantsRequestForAccountRole(id)
	if opts == nil {
		return nil
	}

	client := meta.(*provider.Context).Client

	logging.DebugLogger.Printf("[DEBUG] About to show grants")
	grants, err := client.Grants.Show(ctx, opts)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve grants",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	actualPrivileges := make([]string, 0)
	expectedPrivileges := make([]string, 0)
	expectedPrivileges = append(expectedPrivileges, id.Privileges...)

	if slices.ContainsFunc(expectedPrivileges, func(s string) bool {
		return strings.ToUpper(s) == sdk.AccountObjectPrivilegeImportedPrivileges.String()
	}) {
		expectedPrivileges = append(expectedPrivileges, sdk.AccountObjectPrivilegeUsage.String())
	}

	logging.DebugLogger.Printf("[DEBUG] Filtering grants to be set on account: count = %d", len(grants))
	for _, grant := range grants {
		// Accept only (account) ROLEs
		if grant.GrantTo != sdk.ObjectTypeRole && grant.GrantedTo != sdk.ObjectTypeRole {
			continue
		}
		// Only consider privileges that are already present in the ID, so we
		// don't delete privileges managed by other resources.
		if !slices.Contains(expectedPrivileges, grant.Privilege) {
			continue
		}
		if grant.GrantOption == id.WithGrantOption && grant.GranteeName.Name() == id.RoleName.Name() {
			// Future grants do not have grantedBy, only current grants do.
			// If grantedby is an empty string, it means terraform could not have created the grant.
			// The same goes for the default SNOWFLAKE database, but we don't want to skip in this case
			if (opts.Future == nil || !*opts.Future) && grant.GrantedBy.Name() == "" && grant.Name.Name() != "SNOWFLAKE" {
				continue
			}

			// grant_on is for future grants, granted_on is for current grants.
			// They function the same way though in a test for matching the object type
			//
			// To `grant privilege on application to a role` the user has to use `object_type = "DATABASE"`.
			// It's because Snowflake treats applications as if they were databases. One exception to the rule is
			// the default application named SNOWFLAKE that could be granted with `object_type = "APPLICATION"`.
			// To make the logic simpler, we do not allow it and `object_type = "DATABASE"` should be used for all applications.
			// TODO When implementing SNOW-991421 see if logic added in SNOW-887897 could be moved to the SDK to simplify the resource implementation.
			if grantedOn == sdk.ObjectTypeDatabase && (sdk.ObjectTypeApplication == grant.GrantedOn || sdk.ObjectTypeApplication == grant.GrantOn) {
				actualPrivileges = append(actualPrivileges, grant.Privilege)
			} else if grantedOn == grant.GrantedOn || grantedOn == grant.GrantOn {
				actualPrivileges = append(actualPrivileges, grant.Privilege)
			}
		}
	}

	usageIndex := slices.IndexFunc(actualPrivileges, func(s string) bool { return strings.ToUpper(s) == sdk.AccountObjectPrivilegeUsage.String() })
	if slices.ContainsFunc(expectedPrivileges, func(s string) bool {
		return strings.ToUpper(s) == sdk.AccountObjectPrivilegeImportedPrivileges.String()
	}) && usageIndex >= 0 {
		actualPrivileges[usageIndex] = sdk.AccountObjectPrivilegeImportedPrivileges.String()
	}

	logging.DebugLogger.Printf("[DEBUG] Setting privileges: %v", actualPrivileges)
	if err := d.Set("privileges", actualPrivileges); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error setting privileges for account role",
				Detail:   fmt.Sprintf("Id: %s\nPrivileges: %v\nError: %s", d.Id(), actualPrivileges, err),
			},
		}
	}

	return nil
}

func prepareShowGrantsRequestForAccountRole(id GrantPrivilegesToAccountRoleId) (*sdk.ShowGrantOptions, sdk.ObjectType) {
	opts := new(sdk.ShowGrantOptions)
	var grantedOn sdk.ObjectType

	switch id.Kind {
	case OnAccountAccountRoleGrantKind:
		grantedOn = sdk.ObjectTypeAccount
		opts.On = &sdk.ShowGrantsOn{
			Account: sdk.Bool(true),
		}
	case OnAccountObjectAccountRoleGrantKind:
		data := id.Data.(*OnAccountObjectGrantData)
		grantedOn = data.ObjectType
		opts.On = &sdk.ShowGrantsOn{
			Object: &sdk.Object{
				ObjectType: data.ObjectType,
				Name:       data.ObjectName,
			},
		}
	case OnSchemaAccountRoleGrantKind:
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
			log.Printf("[INFO] Show with on_schema.all_schemas_in_database option is skipped. No changes in privileges in Snowflake will be detected.")
			return nil, ""
		case OnFutureSchemasInDatabaseSchemaGrantKind:
			opts.Future = sdk.Bool(true)
			opts.In = &sdk.ShowGrantsIn{
				Database: data.DatabaseName,
			}
		}
	case OnSchemaObjectAccountRoleGrantKind:
		data := id.Data.(*OnSchemaObjectGrantData)

		switch data.Kind {
		case OnObjectSchemaObjectGrantKind:
			grantedOn = data.Object.ObjectType
			opts.On = &sdk.ShowGrantsOn{
				Object: data.Object,
			}
		case OnAllSchemaObjectGrantKind:
			log.Printf("[INFO] Show with on_schema_object.on_all option is skipped. No changes in privileges in Snowflake will be detected.")
			return nil, ""
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

	return opts, grantedOn
}

func getAccountRolePrivilegesFromSchema(d *schema.ResourceData) *sdk.AccountRoleGrantPrivileges {
	_, onAccountOk := d.GetOk("on_account")
	_, onAccountObjectOk := d.GetOk("on_account_object")
	_, onSchemaOk := d.GetOk("on_schema")
	_, onSchemaObjectOk := d.GetOk("on_schema_object")

	return getAccountRolePrivileges(
		d.Get("all_privileges").(bool),
		expandStringList(d.Get("privileges").(*schema.Set).List()),
		onAccountOk,
		onAccountObjectOk,
		onSchemaOk,
		onSchemaObjectOk,
	)
}

func getAccountRolePrivileges(allPrivileges bool, privileges []string, onAccount bool, onAccountObject bool, onSchema bool, onSchemaObject bool) *sdk.AccountRoleGrantPrivileges {
	accountRoleGrantPrivileges := new(sdk.AccountRoleGrantPrivileges)

	if allPrivileges {
		accountRoleGrantPrivileges.AllPrivileges = sdk.Bool(true)
		return accountRoleGrantPrivileges
	}

	switch {
	case onAccount:
		globalPrivileges := make([]sdk.GlobalPrivilege, len(privileges))
		for i, privilege := range privileges {
			globalPrivileges[i] = sdk.GlobalPrivilege(privilege)
		}
		accountRoleGrantPrivileges.GlobalPrivileges = globalPrivileges
	case onAccountObject:
		accountObjectPrivileges := make([]sdk.AccountObjectPrivilege, len(privileges))
		for i, privilege := range privileges {
			accountObjectPrivileges[i] = sdk.AccountObjectPrivilege(privilege)
		}
		accountRoleGrantPrivileges.AccountObjectPrivileges = accountObjectPrivileges
	case onSchema:
		schemaPrivileges := make([]sdk.SchemaPrivilege, len(privileges))
		for i, privilege := range privileges {
			schemaPrivileges[i] = sdk.SchemaPrivilege(privilege)
		}
		accountRoleGrantPrivileges.SchemaPrivileges = schemaPrivileges
	case onSchemaObject:
		schemaObjectPrivileges := make([]sdk.SchemaObjectPrivilege, len(privileges))
		for i, privilege := range privileges {
			schemaObjectPrivileges[i] = sdk.SchemaObjectPrivilege(privilege)
		}
		accountRoleGrantPrivileges.SchemaObjectPrivileges = schemaObjectPrivileges
	}

	return accountRoleGrantPrivileges
}

func getAccountRoleGrantOn(d *schema.ResourceData) *sdk.AccountRoleGrantOn {
	_, onAccountOk := d.GetOk("on_account")
	onAccountObjectBlock, onAccountObjectOk := d.GetOk("on_account_object")
	onSchemaBlock, onSchemaOk := d.GetOk("on_schema")
	onSchemaObjectBlock, onSchemaObjectOk := d.GetOk("on_schema_object")
	on := new(sdk.AccountRoleGrantOn)

	switch {
	case onAccountOk:
		on.Account = sdk.Bool(true)
	case onAccountObjectOk:
		onAccountObject := onAccountObjectBlock.([]any)[0].(map[string]any)

		grantOnAccountObject := new(sdk.GrantOnAccountObject)

		objectType := onAccountObject["object_type"].(string)
		objectName := onAccountObject["object_name"].(string)
		objectIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(objectName)

		switch sdk.ObjectType(objectType) {
		case sdk.ObjectTypeDatabase:
			grantOnAccountObject.Database = &objectIdentifier
		case sdk.ObjectTypeFailoverGroup:
			grantOnAccountObject.FailoverGroup = &objectIdentifier
		case sdk.ObjectTypeIntegration:
			grantOnAccountObject.Integration = &objectIdentifier
		case sdk.ObjectTypeReplicationGroup:
			grantOnAccountObject.ReplicationGroup = &objectIdentifier
		case sdk.ObjectTypeResourceMonitor:
			grantOnAccountObject.ResourceMonitor = &objectIdentifier
		case sdk.ObjectTypeUser:
			grantOnAccountObject.User = &objectIdentifier
		case sdk.ObjectTypeWarehouse:
			grantOnAccountObject.Warehouse = &objectIdentifier
		case sdk.ObjectTypeComputePool:
			grantOnAccountObject.ComputePool = &objectIdentifier
		case sdk.ObjectTypeExternalVolume:
			grantOnAccountObject.ExternalVolume = &objectIdentifier
		}

		on.AccountObject = grantOnAccountObject
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
				ObjectType: sdk.ObjectType(objectType),
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

func createGrantPrivilegesToAccountRoleIdFromSchema(d *schema.ResourceData) *GrantPrivilegesToAccountRoleId {
	id := new(GrantPrivilegesToAccountRoleId)
	id.RoleName = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("account_role_name").(string))
	id.AllPrivileges = d.Get("all_privileges").(bool)
	if p, ok := d.GetOk("privileges"); ok {
		id.Privileges = expandStringList(p.(*schema.Set).List())
	}
	id.WithGrantOption = d.Get("with_grant_option").(bool)

	on := getAccountRoleGrantOn(d)
	switch {
	case on.Account != nil:
		id.Kind = OnAccountAccountRoleGrantKind
		id.Data = new(OnAccountGrantData)
	case on.AccountObject != nil:
		onAccountObjectGrantData := new(OnAccountObjectGrantData)

		switch {
		case on.AccountObject.User != nil:
			onAccountObjectGrantData.ObjectType = sdk.ObjectTypeUser
			onAccountObjectGrantData.ObjectName = *on.AccountObject.User
		case on.AccountObject.ResourceMonitor != nil:
			onAccountObjectGrantData.ObjectType = sdk.ObjectTypeResourceMonitor
			onAccountObjectGrantData.ObjectName = *on.AccountObject.ResourceMonitor
		case on.AccountObject.Warehouse != nil:
			onAccountObjectGrantData.ObjectType = sdk.ObjectTypeWarehouse
			onAccountObjectGrantData.ObjectName = *on.AccountObject.Warehouse
		case on.AccountObject.Database != nil:
			onAccountObjectGrantData.ObjectType = sdk.ObjectTypeDatabase
			onAccountObjectGrantData.ObjectName = *on.AccountObject.Database
		case on.AccountObject.Integration != nil:
			onAccountObjectGrantData.ObjectType = sdk.ObjectTypeIntegration
			onAccountObjectGrantData.ObjectName = *on.AccountObject.Integration
		case on.AccountObject.FailoverGroup != nil:
			onAccountObjectGrantData.ObjectType = sdk.ObjectTypeFailoverGroup
			onAccountObjectGrantData.ObjectName = *on.AccountObject.FailoverGroup
		case on.AccountObject.ReplicationGroup != nil:
			onAccountObjectGrantData.ObjectType = sdk.ObjectTypeReplicationGroup
			onAccountObjectGrantData.ObjectName = *on.AccountObject.ReplicationGroup
		case on.AccountObject.ExternalVolume != nil:
			onAccountObjectGrantData.ObjectType = sdk.ObjectTypeExternalVolume
			onAccountObjectGrantData.ObjectName = *on.AccountObject.ExternalVolume
		}

		id.Kind = OnAccountObjectAccountRoleGrantKind
		id.Data = onAccountObjectGrantData
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

		id.Kind = OnSchemaAccountRoleGrantKind
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

		id.Kind = OnSchemaObjectAccountRoleGrantKind
		id.Data = onSchemaObjectGrantData
	}

	return id
}
