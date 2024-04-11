package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var grantOwnershipSchema = map[string]*schema.Schema{
	"account_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the account role to which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"account_role_name",
			"database_role_name",
		},
	},
	"database_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the database role to which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		ExactlyOneOf: []string{
			"account_role_name",
			"database_role_name",
		},
	},
	"outbound_privileges": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies whether to remove or transfer all existing outbound privileges on the object when ownership is transferred to a new role. Available options are: REVOKE for removing existing privileges and COPY to transfer them with ownership. For more information head over to [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#optional-parameters).",
		ValidateFunc: validation.StringInSlice([]string{
			"COPY",
			"REVOKE",
		}, true),
	},
	"on": {
		Type:        schema.TypeList,
		Required:    true,
		ForceNew:    true,
		Description: "Configures which object(s) should transfer their ownership to the specified role.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_type": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: fmt.Sprintf("Specifies the type of object on which you are transferring ownership. Available values are: %s", strings.Join(sdk.ValidGrantOwnershipObjectTypesString, " | ")),
					RequiredWith: []string{
						"on.0.object_name",
					},
					ValidateFunc: validation.StringInSlice(sdk.ValidGrantOwnershipObjectTypesString, true),
				},
				"object_name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Specifies the identifier for the object on which you are transferring ownership.",
					RequiredWith: []string{
						"on.0.object_type",
					},
					ExactlyOneOf: []string{
						"on.0.object_name",
						"on.0.all",
						"on.0.future",
					},
				},
				"all": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantOwnershipBulkOperationSchema("all"),
					},
					ExactlyOneOf: []string{
						"on.0.object_name",
						"on.0.all",
						"on.0.future",
					},
				},
				"future": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantOwnershipBulkOperationSchema("future"),
					},
					ExactlyOneOf: []string{
						"on.0.object_name",
						"on.0.all",
						"on.0.future",
					},
				},
			},
		},
	},
}

func grantOwnershipBulkOperationSchema(branchName string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"object_type_plural": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  fmt.Sprintf("Specifies the type of object in plural form on which you are transferring ownership. Available values are: %s. For more information head over to [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#required-parameters).", strings.Join(sdk.ValidGrantOwnershipPluralObjectTypesString, " | ")),
			ValidateFunc: validation.StringInSlice(sdk.ValidGrantOwnershipPluralObjectTypesString, true),
		},
		"in_database": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			Description:      "The fully qualified name of the database.",
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			ExactlyOneOf: []string{
				fmt.Sprintf("on.0.%s.0.in_database", branchName),
				fmt.Sprintf("on.0.%s.0.in_schema", branchName),
			},
		},
		"in_schema": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			Description:      "The fully qualified name of the schema.",
			ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
			ExactlyOneOf: []string{
				fmt.Sprintf("on.0.%s.0.in_database", branchName),
				fmt.Sprintf("on.0.%s.0.in_schema", branchName),
			},
		},
	}
}

func GrantOwnership() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateGrantOwnership,
		// There's no Update, because every field is marked as ForceNew
		DeleteContext: DeleteGrantOwnership,
		ReadContext:   ReadGrantOwnership,

		Schema: grantOwnershipSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportGrantOwnership(),
		},
	}
}

func ImportGrantOwnership() schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
		logging.DebugLogger.Printf("[DEBUG] Entering import grant privileges to account role")
		id, err := ParseGrantOwnershipId(d.Id())
		if err != nil {
			return nil, err
		}
		logging.DebugLogger.Printf("[DEBUG] Imported identifier: %s", id.String())

		switch id.GrantOwnershipTargetRoleKind {
		case ToAccountGrantOwnershipTargetRoleKind:
			if err := d.Set("account_role_name", id.AccountRoleName.Name()); err != nil {
				return nil, err
			}
		case ToDatabaseGrantOwnershipTargetRoleKind:
			if err := d.Set("database_role_name", id.DatabaseRoleName.FullyQualifiedName()); err != nil {
				return nil, err
			}
		}

		if id.OutboundPrivilegesBehavior != nil {
			if err := d.Set("outbound_privileges", *id.OutboundPrivilegesBehavior); err != nil {
				return nil, err
			}
		}

		switch id.Kind {
		case OnObjectGrantOwnershipKind:
			data := id.Data.(*OnObjectGrantOwnershipData)

			onObject := make(map[string]any)
			onObject["object_type"] = data.ObjectType.String()
			if objectName, ok := any(data.ObjectName).(sdk.AccountObjectIdentifier); ok {
				onObject["object_name"] = objectName.Name()
			} else {
				onObject["object_name"] = data.ObjectName.FullyQualifiedName()
			}

			if err := d.Set("on", []any{onObject}); err != nil {
				return nil, err
			}
		case OnAllGrantOwnershipKind, OnFutureGrantOwnershipKind:
			data := id.Data.(*BulkOperationGrantData)

			on := make(map[string]any)
			onAllOrFuture := make(map[string]any)
			onAllOrFuture["object_type_plural"] = data.ObjectNamePlural.String()
			switch data.Kind {
			case InDatabaseBulkOperationGrantKind:
				onAllOrFuture["in_database"] = data.Database.Name()
			case InSchemaBulkOperationGrantKind:
				onAllOrFuture["in_schema"] = data.Schema.FullyQualifiedName()
			}

			switch id.Kind {
			case OnAllGrantOwnershipKind:
				on["all"] = []any{onAllOrFuture}
			case OnFutureGrantOwnershipKind:
				on["future"] = []any{onAllOrFuture}
			}

			if err := d.Set("on", []any{on}); err != nil {
				return nil, err
			}
		}

		return []*schema.ResourceData{d}, nil
	}
}

func CreateGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := createGrantOwnershipIdFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	logging.DebugLogger.Printf("[DEBUG] created identifier from schema: %s", id.String())

	grantOn, err := getOwnershipGrantOn(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Grants.GrantOwnership(
		ctx,
		*grantOn,
		getOwnershipGrantTo(d),
		getOwnershipGrantOpts(id),
	)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred during grant ownership",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", id.String(), err),
			},
		}
	}

	logging.DebugLogger.Printf("[DEBUG] Setting identifier to %s", id.String())
	d.SetId(id.String())

	return ReadGrantOwnership(ctx, d, meta)
}

func DeleteGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := ParseGrantOwnershipId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	grantOn, err := getOwnershipGrantOn(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if grantOn.Future != nil {
		// TODO (SNOW-1182623): Still waiting for the response on the behavior/SQL syntax we should use here
		log.Printf("[WARN] Unsupported operation, please revoke ownership transfer manually")
	} else {
		accountRoleName, err := client.ContextFunctions.CurrentRole(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		err = client.Grants.GrantOwnership( // TODO: Should we always set outbound privileges to COPY in delete operation or set it to the config value?
			ctx,
			*grantOn,
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(sdk.NewAccountObjectIdentifier(accountRoleName)),
			},
			getOwnershipGrantOpts(id),
		)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "An error occurred when transferring ownership back to the original role",
					Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
				},
			}
		}
	}

	d.SetId("")

	return nil
}

func ReadGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := ParseGrantOwnershipId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	opts, grantedOn := prepareShowGrantsRequestForGrantOwnership(id)
	if opts == nil {
		return nil
	}

	client := meta.(*provider.Context).Client

	grants, err := client.Grants.Show(ctx, opts)
	if err != nil {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to retrieve grants. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	ownershipFound := false

	for _, grant := range grants {
		if grant.Privilege != "OWNERSHIP" {
			continue
		}

		// Future grants do not have grantedBy, only current grants do.
		// If grantedby is an empty string, it means terraform could not have created the grant
		if (opts.Future == nil || !*opts.Future) && grant.GrantedBy.Name() == "" {
			continue
		}

		// grant_on is for future grants, granted_on is for current grants.
		// They function the same way though in a test for matching the object type
		if grantedOn != grant.GrantedOn && grantedOn != grant.GrantOn {
			continue
		}

		switch id.GrantOwnershipTargetRoleKind {
		case ToAccountGrantOwnershipTargetRoleKind:
			if grant.GranteeName.Name() == id.AccountRoleName.Name() {
				ownershipFound = true
			}
		case ToDatabaseGrantOwnershipTargetRoleKind:
			if grant.GranteeName.Name() == id.DatabaseRoleName.Name() {
				ownershipFound = true
			}
		}
	}

	if !ownershipFound {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Couldn't find OWNERSHIP privilege on the target object. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Id: %s", d.Id()),
			},
		}
	}

	return nil
}

// TODO(SNOW-1229218): Make sdk.ObjectType + string objectName to sdk.ObjectIdentifier mapping available in the sdk (for all object types).
func getOnObjectIdentifier(objectType sdk.ObjectType, objectName string) (sdk.ObjectIdentifier, error) {
	identifier, err := helpers.DecodeSnowflakeParameterID(objectName)
	if err != nil {
		return nil, err
	}

	switch objectType {
	case sdk.ObjectTypeComputePool,
		sdk.ObjectTypeDatabase,
		sdk.ObjectTypeExternalVolume,
		sdk.ObjectTypeFailoverGroup,
		sdk.ObjectTypeIntegration,
		sdk.ObjectTypeNetworkPolicy,
		sdk.ObjectTypeReplicationGroup,
		sdk.ObjectTypeRole,
		sdk.ObjectTypeUser,
		sdk.ObjectTypeWarehouse:
		return sdk.NewAccountObjectIdentifier(objectName), nil
	case sdk.ObjectTypeDatabaseRole,
		sdk.ObjectTypeSchema:
		if _, ok := identifier.(sdk.DatabaseObjectIdentifier); !ok {
			return nil, sdk.NewError(fmt.Sprintf("invalid object_name %s, expected database object identifier", objectName))
		}
	case sdk.ObjectTypeAggregationPolicy,
		sdk.ObjectTypeAlert,
		sdk.ObjectTypeAuthenticationPolicy,
		sdk.ObjectTypeDataMetricFunction,
		sdk.ObjectTypeDynamicTable,
		sdk.ObjectTypeEventTable,
		sdk.ObjectTypeExternalTable,
		sdk.ObjectTypeFileFormat,
		sdk.ObjectTypeFunction,
		sdk.ObjectTypeGitRepository,
		sdk.ObjectTypeHybridTable,
		sdk.ObjectTypeIcebergTable,
		sdk.ObjectTypeImageRepository,
		sdk.ObjectTypeMaterializedView,
		sdk.ObjectTypeNetworkRule,
		sdk.ObjectTypePackagesPolicy,
		sdk.ObjectTypePipe,
		sdk.ObjectTypeProcedure,
		sdk.ObjectTypeMaskingPolicy,
		sdk.ObjectTypePasswordPolicy,
		sdk.ObjectTypeProjectionPolicy,
		sdk.ObjectTypeRowAccessPolicy,
		sdk.ObjectTypeSessionPolicy,
		sdk.ObjectTypeSecret,
		sdk.ObjectTypeSequence,
		sdk.ObjectTypeStage,
		sdk.ObjectTypeStream,
		sdk.ObjectTypeTable,
		sdk.ObjectTypeTag,
		sdk.ObjectTypeTask,
		sdk.ObjectTypeView:
		if _, ok := identifier.(sdk.SchemaObjectIdentifier); !ok {
			return nil, sdk.NewError(fmt.Sprintf("invalid object_name %s, expected schema object identifier", objectName))
		}
	default:
		return nil, sdk.NewError(fmt.Sprintf("object_type %s is not supported, please create a feature request for the provider if given object_type should be supported", objectType))
	}

	return identifier, nil
}

func getOwnershipGrantOn(d *schema.ResourceData) (*sdk.OwnershipGrantOn, error) {
	ownershipGrantOn := new(sdk.OwnershipGrantOn)

	on := d.Get("on").([]any)[0].(map[string]any)
	onObjectType := on["object_type"].(string)
	onObjectName := on["object_name"].(string)
	onAll := on["all"].([]any)
	onFuture := on["future"].([]any)

	switch {
	case len(onObjectType) > 0 && len(onObjectName) > 0:
		objectType := sdk.ObjectType(strings.ToUpper(onObjectType))
		objectName, err := getOnObjectIdentifier(objectType, onObjectName)
		if err != nil {
			return nil, err
		}
		ownershipGrantOn.Object = &sdk.Object{
			ObjectType: objectType,
			Name:       objectName,
		}
	case len(onAll) > 0:
		ownershipGrantOn.All = getGrantOnSchemaObjectIn(onAll[0].(map[string]any))
	case len(onFuture) > 0:
		ownershipGrantOn.Future = getGrantOnSchemaObjectIn(onFuture[0].(map[string]any))
	}

	return ownershipGrantOn, nil
}

func getOwnershipGrantTo(d *schema.ResourceData) sdk.OwnershipGrantTo {
	var ownershipGrantTo sdk.OwnershipGrantTo

	if accountRoleName, ok := d.GetOk("account_role_name"); ok {
		ownershipGrantTo.AccountRoleName = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(accountRoleName.(string)))
	}

	if databaseRoleName, ok := d.GetOk("database_role_name"); ok {
		ownershipGrantTo.DatabaseRoleName = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(databaseRoleName.(string)))
	}

	return ownershipGrantTo
}

func getOwnershipGrantOpts(id *GrantOwnershipId) *sdk.GrantOwnershipOptions {
	opts := new(sdk.GrantOwnershipOptions)

	if id != nil && id.OutboundPrivilegesBehavior != nil {
		outboundPrivileges := id.OutboundPrivilegesBehavior.ToOwnershipCurrentGrantsOutboundPrivileges()
		if outboundPrivileges != nil {
			opts.CurrentGrants = &sdk.OwnershipCurrentGrants{
				OutboundPrivileges: *outboundPrivileges,
			}
		}
	}

	return opts
}

func prepareShowGrantsRequestForGrantOwnership(id *GrantOwnershipId) (*sdk.ShowGrantOptions, sdk.ObjectType) {
	opts := new(sdk.ShowGrantOptions)
	var grantedOn sdk.ObjectType

	switch id.Kind {
	case OnObjectGrantOwnershipKind:
		data := id.Data.(*OnObjectGrantOwnershipData)
		grantedOn = data.ObjectType
		opts.On = &sdk.ShowGrantsOn{
			Object: &sdk.Object{
				ObjectType: data.ObjectType,
				Name:       data.ObjectName,
			},
		}
	case OnAllGrantOwnershipKind: // TODO: discuss if we want to let users do this (lose control over ownership for all objects in x during delete operation - we can also add a flag that would skip delete operation when on_all is set)
		switch data := id.Data.(*BulkOperationGrantData); data.Kind {
		case InDatabaseBulkOperationGrantKind:
			log.Printf("[INFO] Show with on.all option is skipped. No changes in ownership on all %s in database %s in Snowflake will be detected.", data.ObjectNamePlural, data.Database)
		case InSchemaBulkOperationGrantKind:
			log.Printf("[INFO] Show with on.all option is skipped. No changes in ownership on all %s in schema %s in Snowflake will be detected.", data.ObjectNamePlural, data.Schema)
		}
		return nil, ""
	case OnFutureGrantOwnershipKind:
		data := id.Data.(*BulkOperationGrantData)
		grantedOn = data.ObjectNamePlural.Singular()
		opts.Future = sdk.Bool(true)

		switch data.Kind {
		case InDatabaseBulkOperationGrantKind:
			opts.In = &sdk.ShowGrantsIn{
				Database: data.Database,
			}
		case InSchemaBulkOperationGrantKind:
			opts.In = &sdk.ShowGrantsIn{
				Schema: data.Schema,
			}
		}
	}

	return opts, grantedOn
}

func createGrantOwnershipIdFromSchema(d *schema.ResourceData) (*GrantOwnershipId, error) {
	id := new(GrantOwnershipId)
	accountRoleName, accountRoleNameOk := d.GetOk("account_role_name")
	databaseRoleName, databaseRoleNameOk := d.GetOk("database_role_name")

	switch {
	case accountRoleNameOk:
		id.GrantOwnershipTargetRoleKind = ToAccountGrantOwnershipTargetRoleKind
		id.AccountRoleName = sdk.NewAccountObjectIdentifierFromFullyQualifiedName(accountRoleName.(string))
	case databaseRoleNameOk:
		id.GrantOwnershipTargetRoleKind = ToDatabaseGrantOwnershipTargetRoleKind
		id.DatabaseRoleName = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(databaseRoleName.(string))
	}

	outboundPrivileges, outboundPrivilegesOk := d.GetOk("outbound_privileges")
	if outboundPrivilegesOk {
		switch OutboundPrivilegesBehavior(outboundPrivileges.(string)) {
		case CopyOutboundPrivilegesBehavior:
			id.OutboundPrivilegesBehavior = sdk.Pointer(CopyOutboundPrivilegesBehavior)
		case RevokeOutboundPrivilegesBehavior:
			id.OutboundPrivilegesBehavior = sdk.Pointer(RevokeOutboundPrivilegesBehavior)
		}
	}

	grantedOn := d.Get("on").([]any)[0].(map[string]any)
	objectType := grantedOn["object_type"].(string)
	objectName := grantedOn["object_name"].(string)
	all := grantedOn["all"].([]any)
	future := grantedOn["future"].([]any)

	switch {
	case len(objectType) > 0 && len(objectName) > 0:
		id.Kind = OnObjectGrantOwnershipKind
		objectType := sdk.ObjectType(objectType)
		objectName, err := getOnObjectIdentifier(objectType, objectName)
		if err != nil {
			return nil, err
		}
		id.Data = &OnObjectGrantOwnershipData{
			ObjectType: objectType,
			ObjectName: objectName,
		}
	case len(all) > 0:
		id.Kind = OnAllGrantOwnershipKind
		id.Data = getBulkOperationGrantData(getGrantOnSchemaObjectIn(all[0].(map[string]any)))
	case len(future) > 0:
		id.Kind = OnFutureGrantOwnershipKind
		id.Data = getBulkOperationGrantData(getGrantOnSchemaObjectIn(future[0].(map[string]any)))
	}

	return id, nil
}
