package resources

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validGrantOwnershipObjectTypes = []sdk.ObjectType{
	sdk.ObjectTypeAggregationPolicy,
	sdk.ObjectTypeAlert,
	sdk.ObjectTypeAuthenticationPolicy,
	sdk.ObjectTypeComputePool,
	sdk.ObjectTypeDatabase,
	sdk.ObjectTypeDatabaseRole,
	sdk.ObjectTypeDynamicTable,
	sdk.ObjectTypeEventTable,
	sdk.ObjectTypeExternalTable,
	sdk.ObjectTypeExternalVolume,
	sdk.ObjectTypeFailoverGroup,
	sdk.ObjectTypeFileFormat,
	sdk.ObjectTypeFunction,
	sdk.ObjectTypeHybridTable,
	sdk.ObjectTypeIcebergTable,
	sdk.ObjectTypeImageRepository,
	sdk.ObjectTypeIntegration,
	sdk.ObjectTypeMaterializedView,
	sdk.ObjectTypeNetworkPolicy,
	sdk.ObjectTypeNetworkRule,
	sdk.ObjectTypePackagesPolicy,
	sdk.ObjectTypePipe,
	sdk.ObjectTypeProcedure,
	sdk.ObjectTypeMaskingPolicy,
	sdk.ObjectTypePasswordPolicy,
	sdk.ObjectTypeProjectionPolicy,
	sdk.ObjectTypeRole,
	sdk.ObjectTypeRowAccessPolicy,
	sdk.ObjectTypeSchema,
	sdk.ObjectTypeSessionPolicy,
	sdk.ObjectTypeSecret,
	sdk.ObjectTypeSequence,
	sdk.ObjectTypeStage,
	sdk.ObjectTypeStream,
	sdk.ObjectTypeTable,
	sdk.ObjectTypeTag,
	sdk.ObjectTypeTask,
	sdk.ObjectTypeUser,
	sdk.ObjectTypeView,
	sdk.ObjectTypeWarehouse,
}

var validGrantOwnershipObjectTypesString []string
var validGrantOwnershipPluralObjectTypesString []string

func init() {
	validGrantOwnershipObjectTypesString = make([]string, len(validGrantOwnershipObjectTypes))
	validGrantOwnershipPluralObjectTypesString = make([]string, len(validGrantOwnershipObjectTypes))

	for i, objectType := range validGrantOwnershipObjectTypes {
		validGrantOwnershipObjectTypesString[i] = objectType.String()
		validGrantOwnershipObjectTypesString[i] = objectType.Plural().String()
	}
}

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
		Description: "Specifies whether to remove or transfer all existing outbound privileges on the object when ownership is transferred to a new role. Available options are: REVOKE for removing existing privileges and COPY to transfer them with ownership.",
		ValidateFunc: validation.StringInSlice([]string{
			"COPY",
			"REVOKE",
		}, true),
	},
	"on": {
		Type:        schema.TypeList,
		Required:    true,
		ForceNew:    true,
		Description: "TODO",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"on.0.object_name",
			"on.0.all",
			"on.0.future",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_type": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Specifies the type of object on which you are transferring ownership. Available values are: AGGREGATION POLICY | ALERT | AUTHENTICATION POLICY | COMPUTE POOL | DATABASE | DATABASE ROLE | DYNAMIC TABLE | EVENT TABLE | EXTERNAL TABLE | EXTERNAL VOLUME | FAILOVER GROUP | FILE FORMAT | FUNCTION | HYBRID TABLE | ICEBERG TABLE | IMAGE REPOSITORY | INTEGRATION | MATERIALIZED VIEW | NETWORK POLICY | NETWORK RULE | PACKAGES POLICY | PIPE | PROCEDURE | MASKING POLICY | PASSWORD POLICY | PROJECTION POLICY | REPLICATION GROUP | ROLE | ROW ACCESS POLICY | SCHEMA | SESSION POLICY | SECRET | SEQUENCE | STAGE | STREAM | TABLE | TAG | TASK | USER | VIEW | WAREHOUSE",
					RequiredWith: []string{
						"on.0.object_name",
					},
					ValidateFunc: validation.StringInSlice(validGrantOwnershipObjectTypesString, true),
				},
				"object_name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Specifies the identifier for the object on which you are transferring ownership.",
					RequiredWith: []string{
						"on.0.object_type",
					},
				},
				"all": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantOwnershipBulkOperationSchema,
					},
					ExactlyOneOf: []string{
						"on.0.all.0.in_database",
						"on.0.all.0.in_schema",
					},
				},
				"future": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantOwnershipBulkOperationSchema,
					},
					ExactlyOneOf: []string{
						"on.0.future.0.in_database",
						"on.0.future.0.in_schema",
					},
				},
			},
		},
	},
}

var grantOwnershipBulkOperationSchema = map[string]*schema.Schema{
	"object_type_plural": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Specifies the type of object in plural form on which you are transferring ownership. Available values are: AGGREGATION POLICIES | ALERTS | AUTHENTICATION POLICIES | COMPUTE POOLS | DATABASES | DATABASE ROLES | DYNAMIC TABLES | EVENT TABLES | EXTERNAL TABLES | EXTERNAL VOLUMES | FAILOVER GROUPS | FILE FORMATS | FUNCTIONS | HYBRID TABLES | ICEBERG TABLES | IMAGE REPOSITORIES | INTEGRATIONS | MATERIALIZED VIEWS | NETWORK POLICIES | NETWORK RULES | PACKAGES POLICIES | PIPES | PROCEDURES | MASKING POLICIES | PASSWORD POLICIES | PROJECTION POLICIES | REPLICATION GROUPS | ROLES | ROW ACCESS POLICIES | SCHEMAS | SESSION POLICIES | SECRETS | SEQUENCES | STAGES | STREAMS | TABLES | TAGS | TASKS | USERS | VIEWS | WAREHOUSES",
		ValidateFunc: validation.StringInSlice(validGrantOwnershipPluralObjectTypesString, true),
	},
	"in_database": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the database.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"in_schema": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the schema.",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
	},
}

func GrantOwnership() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateGrantOwnership,
		UpdateContext: UpdateGrantOwnership,
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
			if err := d.Set("account_role_name", id.AccountRoleName.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case ToDatabaseGrantOwnershipTargetRoleKind:
			if err := d.Set("database_role_name", id.DatabaseRoleName.FullyQualifiedName()); err != nil {
				return nil, err
			}
		}

		if err := d.Set("outbound_privileges", id.OutboundPrivilegesBehavior); err != nil {
			return nil, err
		}

		switch id.Kind {
		case OnObjectGrantOwnershipKind:
			data := id.Data.(*OnObjectGrantOwnershipData)

			onObject := make(map[string]any)
			onObject["object_type"] = data.ObjectType.String()
			onObject["object_name"] = data.ObjectName.FullyQualifiedName()

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
				onAllOrFuture["in_database"] = data.Database.FullyQualifiedName()
			case InSchemaBulkOperationGrantKind:
				onAllOrFuture["in_schema"] = data.Schema.FullyQualifiedName()
			}

			switch id.Kind {
			case OnAllGrantOwnershipKind:
				on["all"] = onAllOrFuture
			case OnFutureGrantOwnershipKind:
				on["future"] = onAllOrFuture
			}

			if err := d.Set("on", []any{on}); err != nil {
				return nil, err
			}
		}

		return []*schema.ResourceData{d}, nil
	}
}

func CreateGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	//db := meta.(*sql.DB)
	//client := sdk.NewClientFromDB(db)

	d.SetId("some id")

	return ReadGrantOwnership(ctx, d, meta)
}

func UpdateGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	return ReadGrantOwnership(ctx, d, meta)
}

func DeleteGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id, err := ParseGrantOwnershipId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	err = client.Grants.GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{},
		sdk.OwnershipGrantTo{},
		&sdk.GrantOwnershipOptions{},
	)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when transferring ownership back to the original role", // TODO: do we save original role on create ?
				Detail:   fmt.Sprintf("Id: %s\nDatabase role name: %s\nError: %s", d.Id(), id.DatabaseRoleName, err.Error()),
			},
		}
	}

	d.SetId("")

	return nil
}

func ReadGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	return nil
}
