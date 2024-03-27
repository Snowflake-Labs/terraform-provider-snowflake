package resources

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var grantPrivilegesToRoleSchema = map[string]*schema.Schema{
	// According to docs https://docs.snowflake.com/en/user-guide/data-exchange-marketplace-privileges#usage-notes IMPORTED PRIVILEGES
	// will be returned as USAGE in SHOW GRANTS command. In addition, USAGE itself is a valid privilege, but both cannot be set at the
	// same time (IMPORTED PRIVILEGES can only be granted to the database created from SHARE and USAGE in every other case).
	// To handle both cases, additional logic was added in read operation where IMPORTED PRIVILEGES is replaced with USAGE.
	"privileges": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The privileges to grant on the account role.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		ConflictsWith: []string{
			"all_privileges",
		},
	},
	"all_privileges": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Grant all privileges on the account role.",
		ConflictsWith: []string{
			"privileges",
		},
	},
	"on_account": {
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ForceNew:      true,
		Description:   "If true, the privileges will be granted on the account.",
		ConflictsWith: []string{"on_account_object", "on_schema", "on_schema_object"},
	},
	"on_account_object": {
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		MaxItems:      1,
		ConflictsWith: []string{"on_account", "on_schema", "on_schema_object"},
		Description:   "Specifies the account object on which privileges will be granted ",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The object type of the account object on which privileges will be granted. Valid values are: USER | RESOURCE MONITOR | WAREHOUSE | DATABASE | INTEGRATION | FAILOVER GROUP | REPLICATION GROUP | EXTERNAL VOLUME",
					ValidateFunc: validation.StringInSlice([]string{
						"USER",
						"RESOURCE MONITOR",
						"WAREHOUSE",
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
					Description:      "The fully qualified name of the object on which privileges will be granted.",
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
			},
		},
	},
	"on_schema": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"on_account", "on_account_object", "on_schema_object"},
		Description:   "Specifies the schema on which privileges will be granted.",
		ForceNew:      true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"schema_name": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "The fully qualified name of the schema.",
					ConflictsWith:    []string{"on_schema.0.all_schemas_in_database", "on_schema.0.future_schemas_in_database"},
					ForceNew:         true,
					ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
				},
				"all_schemas_in_database": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "The fully qualified name of the database.",
					ConflictsWith:    []string{"on_schema.0.schema_name", "on_schema.0.future_schemas_in_database"},
					ForceNew:         true,
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
				"future_schemas_in_database": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "The fully qualified name of the database.",
					ConflictsWith:    []string{"on_schema.0.schema_name", "on_schema.0.all_schemas_in_database"},
					ForceNew:         true,
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
			},
		},
	},
	"on_schema_object": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"on_account", "on_account_object", "on_schema"},
		Description:   "Specifies the schema object on which privileges will be granted.",
		ForceNew:      true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_type": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      fmt.Sprintf("The object type of the schema object on which privileges will be granted. Valid values are: %s", strings.Join(sdk.ValidGrantToObjectTypesString, " | ")),
					RequiredWith:     []string{"on_schema_object.0.object_name"},
					ConflictsWith:    []string{"on_schema_object.0.all", "on_schema_object.0.future"},
					ForceNew:         true,
					ValidateDiagFunc: StringInSlice(sdk.ValidGrantToObjectTypesString, true),
				},
				"object_name": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "The fully qualified name of the object on which privileges will be granted.",
					RequiredWith:     []string{"on_schema_object.0.object_type"},
					ConflictsWith:    []string{"on_schema_object.0.all", "on_schema_object.0.future"},
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
				},
				"all": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Configures the privilege to be granted on all objects in eihter a database or schema.",
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"object_type_plural": {
								Type:             schema.TypeString,
								Required:         true,
								Description:      fmt.Sprintf("The plural object type of the schema object on which privileges will be granted. Valid values are: %s", strings.Join(sdk.ValidGrantToPluralObjectTypesString, " | ")),
								ForceNew:         true,
								ValidateDiagFunc: StringInSlice(sdk.ValidGrantToPluralObjectTypesString, true),
							},
							"in_database": {
								Type:             schema.TypeString,
								Optional:         true,
								ForceNew:         true,
								Description:      "The fully qualified name of the database.",
								ConflictsWith:    []string{"on_schema_object.0.all.in_schema"},
								ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
							},
							"in_schema": {
								Type:             schema.TypeString,
								Optional:         true,
								ForceNew:         true,
								Description:      "The fully qualified name of the schema.",
								ConflictsWith:    []string{"on_schema_object.0.all.in_database"},
								ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
							},
						},
					},
				},
				"future": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Configures the privilege to be granted on future objects in eihter a database or schema.",
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"object_type_plural": {
								Type:             schema.TypeString,
								Required:         true,
								Description:      fmt.Sprintf("The plural object type of the schema object on which privileges will be granted. Valid values are: %s", strings.Join(sdk.ValidGrantToFuturePluralObjectTypesString, " | ")),
								ForceNew:         true,
								ValidateDiagFunc: StringInSlice(sdk.ValidGrantToFuturePluralObjectTypesString, true),
							},
							"in_database": {
								Type:             schema.TypeString,
								Optional:         true,
								Description:      "The fully qualified name of the database.",
								ConflictsWith:    []string{"on_schema_object.0.all.in_schema"},
								ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
								ForceNew:         true,
							},
							"in_schema": {
								Type:             schema.TypeString,
								Optional:         true,
								Description:      "The fully qualified name of the schema.",
								ConflictsWith:    []string{"on_schema_object.0.all.in_database"},
								ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
								ForceNew:         true,
							},
						},
					},
				},
			},
		},
	},
	"role_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the role to which privileges will be granted.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the grantee can grant the privileges to other users.",
		Default:     false,
		ForceNew:    true,
	},
}

func GrantPrivilegesToRole() *schema.Resource {
	return &schema.Resource{
		Create:             CreateGrantPrivilegesToRole,
		Read:               ReadGrantPrivilegesToRole,
		Delete:             DeleteGrantPrivilegesToRole,
		Update:             UpdateGrantPrivilegesToRole,
		DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_account_role instead.",

		Schema: grantPrivilegesToRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				resourceID := NewGrantPrivilegesToRoleID(d.Id())
				if err := d.Set("role_name", resourceID.RoleName); err != nil {
					return nil, err
				}
				if err := d.Set("privileges", resourceID.Privileges); err != nil {
					return nil, err
				}
				if err := d.Set("all_privileges", resourceID.AllPrivileges); err != nil {
					return nil, err
				}
				if err := d.Set("with_grant_option", resourceID.WithGrantOption); err != nil {
					return nil, err
				}
				if err := d.Set("on_account", resourceID.OnAccount); err != nil {
					return nil, err
				}
				if resourceID.OnAccountObject {
					if err := d.Set("on_account_object", []map[string]interface{}{{
						"object_type": resourceID.ObjectType,
						"object_name": resourceID.ObjectName,
					}}); err != nil {
						return nil, err
					}
				}
				if resourceID.OnSchema {
					var onSchema []interface{}
					if resourceID.SchemaName != "" {
						onSchema = append(onSchema, map[string]interface{}{
							"schema_name": resourceID.SchemaName,
						})
					}
					if resourceID.All {
						onSchema = append(onSchema, map[string]interface{}{
							"all_schemas_in_database": resourceID.DatabaseName,
						})
					}
					if resourceID.Future {
						onSchema = append(onSchema, map[string]interface{}{
							"future_schemas_in_database": resourceID.DatabaseName,
						})
					}
					if err := d.Set("on_schema", onSchema); err != nil {
						return nil, err
					}
				}

				if resourceID.OnSchemaObject {
					var onSchemaObject []interface{}
					if resourceID.ObjectName != "" {
						onSchemaObject = append(onSchemaObject, map[string]interface{}{
							"object_name": resourceID.ObjectName,
							"object_type": resourceID.ObjectType,
						})
					}
					if resourceID.All {
						all := make([]interface{}, 0)
						m := map[string]interface{}{
							"object_type_plural": resourceID.ObjectTypePlural,
						}

						if resourceID.InSchema {
							m["in_schema"] = resourceID.SchemaName
						}
						if resourceID.InDatabase {
							m["in_database"] = resourceID.DatabaseName
						}
						all = append(all, m)
						onSchemaObject = append(onSchemaObject, map[string]interface{}{
							"all": all,
						})
					}
					if resourceID.Future {
						future := make([]interface{}, 0)
						m := map[string]interface{}{
							"object_type_plural": resourceID.ObjectTypePlural,
						}
						if resourceID.InSchema {
							m["in_schema"] = resourceID.SchemaName
						}
						if resourceID.InDatabase {
							m["in_database"] = resourceID.DatabaseName
						}
						future = append(future, m)
						onSchemaObject = append(onSchemaObject, map[string]interface{}{
							"future": future,
						})
					}
					if err := d.Set("on_schema_object", onSchemaObject); err != nil {
						return nil, err
					}
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

// we need to keep track of literally everything to construct a unique identifier that can be imported
type GrantPrivilegesToRoleID struct {
	RoleName         string
	Privileges       []string
	AllPrivileges    bool
	WithGrantOption  bool
	OnAccount        bool
	OnAccountObject  bool
	OnSchema         bool
	OnSchemaObject   bool
	All              bool
	Future           bool
	ObjectType       string
	ObjectName       string
	ObjectTypePlural string
	InSchema         bool
	SchemaName       string
	InDatabase       bool
	DatabaseName     string
}

func NewGrantPrivilegesToRoleID(id string) GrantPrivilegesToRoleID {
	parts := strings.Split(id, "|")
	privileges := strings.Split(parts[1], ",")
	if len(privileges) == 1 && privileges[0] == "" {
		privileges = []string{}
	}
	return GrantPrivilegesToRoleID{
		RoleName:         parts[0],
		Privileges:       privileges,
		AllPrivileges:    parts[2] == "true",
		WithGrantOption:  parts[3] == "true",
		OnAccount:        parts[4] == "true",
		OnAccountObject:  parts[5] == "true",
		OnSchema:         parts[6] == "true",
		OnSchemaObject:   parts[7] == "true",
		All:              parts[8] == "true",
		Future:           parts[9] == "true",
		ObjectType:       parts[10],
		ObjectName:       parts[11],
		ObjectTypePlural: parts[12],
		InSchema:         parts[13] == "true",
		SchemaName:       parts[14],
		InDatabase:       parts[15] == "true",
		DatabaseName:     parts[16],
	}
}

func (v GrantPrivilegesToRoleID) String() string {
	return helpers.EncodeSnowflakeID(v.RoleName, v.Privileges, v.AllPrivileges, v.WithGrantOption, v.OnAccount, v.OnAccountObject, v.OnSchema, v.OnSchemaObject, v.All, v.Future, v.ObjectType, v.ObjectName, v.ObjectTypePlural, v.InSchema, v.SchemaName, v.InDatabase, v.DatabaseName)
}

func CreateGrantPrivilegesToRole(d *schema.ResourceData, meta interface{}) error {
	logging.DebugLogger.Printf("[DEBUG] Entering create grant privileges to role")
	client := meta.(*provider.Context).Client
	logging.DebugLogger.Printf("[DEBUG] Creating new client from db")
	ctx := context.Background()
	resourceID := &GrantPrivilegesToRoleID{}
	var privileges []string
	if p, ok := d.GetOk("privileges"); ok {
		logging.DebugLogger.Printf("[DEBUG] Building privileges list based on config")
		privileges = expandStringList(p.(*schema.Set).List())
		resourceID.Privileges = privileges
	}
	allPrivileges := d.Get("all_privileges").(bool)
	resourceID.AllPrivileges = allPrivileges
	privilegesToGrant, on, err := configureRoleGrantPrivilegeOptions(d, privileges, allPrivileges, resourceID)
	if err != nil {
		return fmt.Errorf("error configuring account role grant privilege options: %w", err)
	}
	withGrantOption := d.Get("with_grant_option").(bool)
	resourceID.WithGrantOption = withGrantOption
	opts := sdk.GrantPrivilegesToAccountRoleOptions{
		WithGrantOption: sdk.Bool(withGrantOption),
	}
	roleName := d.Get("role_name").(string)
	resourceID.RoleName = roleName
	roleID := sdk.NewAccountObjectIdentifier(roleName)
	logging.DebugLogger.Printf("[DEBUG] About to grant privileges to account role")
	err = client.Grants.GrantPrivilegesToAccountRole(ctx, privilegesToGrant, on, roleID, &opts)
	logging.DebugLogger.Printf("[DEBUG] After granting privileges to account role: err = %v", err)
	if err != nil {
		return fmt.Errorf("error granting privileges to account role: %w", err)
	}

	logging.DebugLogger.Printf("[DEBUG] Setting ID to %s", resourceID.String())
	d.SetId(resourceID.String())
	return ReadGrantPrivilegesToRole(d, meta)
}

func ReadGrantPrivilegesToRole(d *schema.ResourceData, meta interface{}) error {
	logging.DebugLogger.Printf("[DEBUG] Entering read grant privileges to role")
	client := meta.(*provider.Context).Client
	logging.DebugLogger.Printf("[DEBUG] Creating new client from db")
	ctx := context.Background()
	resourceID := NewGrantPrivilegesToRoleID(d.Id())
	roleName := resourceID.RoleName
	allPrivileges := resourceID.AllPrivileges
	if allPrivileges {
		log.Printf("[DEBUG] cannot read ALL PRIVILEGES on grant to role %s because this is not returned by API", roleName)
		return nil // cannot read all privileges because its not something returned by API. We can check only if specific privileges are granted to the role
	}
	var opts sdk.ShowGrantOptions
	var grantOn sdk.ObjectType
	if resourceID.OnAccount {
		logging.DebugLogger.Printf("[DEBUG] Preparing to read privileges: on account")
		grantOn = sdk.ObjectTypeAccount
		opts = sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Account: sdk.Bool(true),
			},
		}
	}

	if resourceID.OnAccountObject {
		logging.DebugLogger.Printf("[DEBUG] Preparing to read privileges: on account object")
		objectType := sdk.ObjectType(resourceID.ObjectType)
		grantOn = objectType
		opts = sdk.ShowGrantOptions{
			On: &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: objectType,
					Name:       sdk.NewAccountObjectIdentifierFromFullyQualifiedName(resourceID.ObjectName),
				},
			},
		}
	}

	if resourceID.OnSchema {
		logging.DebugLogger.Printf("[DEBUG] Preparing to read privileges: on schema")
		grantOn = sdk.ObjectTypeSchema
		if resourceID.SchemaName != "" {
			opts = sdk.ShowGrantOptions{
				On: &sdk.ShowGrantsOn{
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeSchema,
						Name:       sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(resourceID.SchemaName),
					},
				},
			}
		}
		if resourceID.All {
			log.Printf("[DEBUG] cannot read ALL SCHEMAS IN DATABASE on grant to role %s because this is not returned by API", roleName)
			return nil // on_all is not supported by API
		}
		if resourceID.Future {
			opts = sdk.ShowGrantOptions{
				Future: sdk.Bool(true),
				In: &sdk.ShowGrantsIn{
					Database: sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(resourceID.DatabaseName)),
				},
			}
		}
	}

	if resourceID.OnSchemaObject {
		logging.DebugLogger.Printf("[DEBUG] Preparing to read privileges: on schema object")
		if resourceID.ObjectName != "" {
			objectType := sdk.ObjectType(resourceID.ObjectType)
			grantOn = objectType
			opts = sdk.ShowGrantOptions{
				On: &sdk.ShowGrantsOn{
					Object: &sdk.Object{
						ObjectType: objectType,
						Name:       sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(resourceID.ObjectName),
					},
				},
			}
		}

		if resourceID.All {
			return nil // on_all is not supported by API
		}

		if resourceID.Future {
			grantOn = sdk.PluralObjectType(resourceID.ObjectTypePlural).Singular()
			if resourceID.InSchema {
				opts = sdk.ShowGrantOptions{
					Future: sdk.Bool(true),
					In: &sdk.ShowGrantsIn{
						Schema: sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(resourceID.SchemaName)),
					},
				}
			}
			if resourceID.InDatabase {
				opts = sdk.ShowGrantOptions{
					Future: sdk.Bool(true),
					In: &sdk.ShowGrantsIn{
						Database: sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(resourceID.DatabaseName)),
					},
				}
			}
		}
	}

	err := readRoleGrantPrivileges(ctx, client, grantOn, resourceID, &opts, d)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGrantPrivilegesToRole(d *schema.ResourceData, meta interface{}) error {
	logging.DebugLogger.Printf("[DEBUG] Entering update grant privileges to role")
	client := meta.(*provider.Context).Client
	logging.DebugLogger.Printf("[DEBUG] Creating new client from db")
	ctx := context.Background()

	// the only thing that can change is "privileges"
	roleName := d.Get("role_name").(string)
	roleID := sdk.NewAccountObjectIdentifier(roleName)

	logging.DebugLogger.Printf("[DEBUG] Checking if privileges have changed")
	if d.HasChange("privileges") {
		logging.DebugLogger.Printf("[DEBUG] Privileges have changed")
		old, new := d.GetChange("privileges")
		oldPrivileges := expandStringList(old.(*schema.Set).List())
		newPrivileges := expandStringList(new.(*schema.Set).List())

		addPrivileges := []string{}
		removePrivileges := []string{}
		for _, oldPrivilege := range oldPrivileges {
			if !slices.Contains(newPrivileges, oldPrivilege) {
				removePrivileges = append(removePrivileges, oldPrivilege)
			}
		}

		for _, newPrivilege := range newPrivileges {
			if !slices.Contains(oldPrivileges, newPrivilege) {
				addPrivileges = append(addPrivileges, newPrivilege)
			}
		}

		// first add new privileges
		if len(addPrivileges) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Adding new privileges")
			privilegesToGrant, on, err := configureRoleGrantPrivilegeOptions(d, addPrivileges, false, &GrantPrivilegesToRoleID{})
			if err != nil {
				return fmt.Errorf("error configuring account role grant privilege options: %w", err)
			}
			logging.DebugLogger.Printf("[DEBUG] About to grant privileges to account role")
			err = client.Grants.GrantPrivilegesToAccountRole(ctx, privilegesToGrant, on, roleID, nil)
			logging.DebugLogger.Printf("[DEBUG] After granting privileges to account role: err = %v", err)
			if err != nil {
				return fmt.Errorf("error granting privileges to account role: %w", err)
			}
		}

		// then remove old privileges
		if len(removePrivileges) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Removing old privileges")
			privilegesToRevoke, on, err := configureRoleGrantPrivilegeOptions(d, removePrivileges, false, &GrantPrivilegesToRoleID{})
			if err != nil {
				return fmt.Errorf("error configuring account role grant privilege options: %w", err)
			}
			logging.DebugLogger.Printf("[DEBUG] About to revoke privileges from account role")
			err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privilegesToRevoke, on, roleID, nil)
			logging.DebugLogger.Printf("[DEBUG] After revoking privileges from account role: err = %v", err)
			if err != nil {
				return fmt.Errorf("error revoking privileges from account role: %w", err)
			}
		}
		logging.DebugLogger.Printf("[DEBUG] Setting new values")
		resourceID := NewGrantPrivilegesToRoleID(d.Id())
		resourceID.Privileges = newPrivileges
		d.SetId(resourceID.String())
	}
	return ReadGrantPrivilegesToRole(d, meta)
}

func DeleteGrantPrivilegesToRole(d *schema.ResourceData, meta interface{}) error {
	logging.DebugLogger.Printf("[DEBUG] Entering delete grant privileges to role")
	client := meta.(*provider.Context).Client
	logging.DebugLogger.Printf("[DEBUG] Creating new client from db")
	ctx := context.Background()

	roleName := d.Get("role_name").(string)
	roleID := sdk.NewAccountObjectIdentifier(roleName)

	var privileges []string
	if p, ok := d.GetOk("privileges"); ok {
		privileges = expandStringList(p.(*schema.Set).List())
	}
	allPrivileges := d.Get("all_privileges").(bool)
	privilegesToRevoke, on, err := configureRoleGrantPrivilegeOptions(d, privileges, allPrivileges, &GrantPrivilegesToRoleID{})
	if err != nil {
		return fmt.Errorf("error configuring account role grant privilege options: %w", err)
	}
	logging.DebugLogger.Printf("[DEBUG] About to revoke privileges from account role")
	err = client.Grants.RevokePrivilegesFromAccountRole(ctx, privilegesToRevoke, on, roleID, nil)
	logging.DebugLogger.Printf("[DEBUG] After revoking privileges from account role: err = %v", err)
	if err != nil {
		return fmt.Errorf("error revoking privileges from account role: %w", err)
	}
	logging.DebugLogger.Printf("[DEBUG] Cleaning resource id")
	d.SetId("")
	return nil
}

func configureRoleGrantPrivilegeOptions(d *schema.ResourceData, privileges []string, allPrivileges bool, resourceID *GrantPrivilegesToRoleID) (*sdk.AccountRoleGrantPrivileges, *sdk.AccountRoleGrantOn, error) {
	logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options")
	var privilegesToGrant *sdk.AccountRoleGrantPrivileges
	on := sdk.AccountRoleGrantOn{}
	if v, ok := d.GetOk("on_account"); ok && v.(bool) {
		logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: on account")
		on.Account = sdk.Bool(true)
		resourceID.OnAccount = true
		privilegesToGrant = setRolePrivilegeOptions(privileges, allPrivileges, true, false, false, false)
		return privilegesToGrant, &on, nil
	}

	if v, ok := d.GetOk("on_account_object"); ok && len(v.([]interface{})) > 0 {
		logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: on account object")
		on.AccountObject = &sdk.GrantOnAccountObject{}
		resourceID.OnAccountObject = true
		onAccountObject := v.([]interface{})[0].(map[string]interface{})
		objectType := sdk.ObjectType(onAccountObject["object_type"].(string))
		resourceID.ObjectType = objectType.String()
		objectName := onAccountObject["object_name"].(string)
		resourceID.ObjectName = objectName
		objectID := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(objectName)
		switch objectType {
		case sdk.ObjectTypeDatabase:
			on.AccountObject.Database = &objectID
		case sdk.ObjectTypeFailoverGroup:
			on.AccountObject.FailoverGroup = &objectID
		case sdk.ObjectTypeIntegration:
			on.AccountObject.Integration = &objectID
		case sdk.ObjectTypeReplicationGroup:
			on.AccountObject.ReplicationGroup = &objectID
		case sdk.ObjectTypeResourceMonitor:
			on.AccountObject.ResourceMonitor = &objectID
		case sdk.ObjectTypeUser:
			on.AccountObject.User = &objectID
		case sdk.ObjectTypeWarehouse:
			on.AccountObject.Warehouse = &objectID
		case sdk.ObjectTypeExternalVolume:
			on.AccountObject.ExternalVolume = &objectID
		default:
			return nil, nil, fmt.Errorf("invalid object type %s", objectType)
		}
		privilegesToGrant = setRolePrivilegeOptions(privileges, allPrivileges, false, true, false, false)
		return privilegesToGrant, &on, nil
	}

	if v, ok := d.GetOk("on_schema"); ok && len(v.([]interface{})) > 0 {
		logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: on schema")
		onSchema := v.([]interface{})[0].(map[string]interface{})
		on.Schema = &sdk.GrantOnSchema{}
		resourceID.OnSchema = true
		if v, ok := onSchema["schema_name"]; ok && len(v.(string)) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting schema name")
			resourceID.SchemaName = v.(string)
			on.Schema.Schema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(v.(string)))
		}
		if v, ok := onSchema["all_schemas_in_database"]; ok && len(v.(string)) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting all schemas in database")
			resourceID.All = true
			resourceID.InDatabase = true
			resourceID.DatabaseName = v.(string)
			on.Schema.AllSchemasInDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(v.(string)))
		}

		if v, ok := onSchema["future_schemas_in_database"]; ok && len(v.(string)) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting future schemas in database")
			resourceID.Future = true
			resourceID.InDatabase = true
			resourceID.DatabaseName = v.(string)
			on.Schema.FutureSchemasInDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(v.(string)))
		}
		privilegesToGrant = setRolePrivilegeOptions(privileges, allPrivileges, false, false, true, false)
		return privilegesToGrant, &on, nil
	}

	if v, ok := d.GetOk("on_schema_object"); ok && len(v.([]interface{})) > 0 {
		logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: on schema object")
		onSchemaObject := v.([]interface{})[0].(map[string]interface{})
		on.SchemaObject = &sdk.GrantOnSchemaObject{}
		resourceID.OnSchemaObject = true
		if v, ok := onSchemaObject["object_type"]; ok && len(v.(string)) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting schema object type")
			resourceID.ObjectType = v.(string)
			on.SchemaObject.SchemaObject = &sdk.Object{
				ObjectType: sdk.ObjectType(v.(string)),
			}
		}
		if v, ok := onSchemaObject["object_name"]; ok && len(v.(string)) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting schema object name")
			resourceID.ObjectName = v.(string)
			on.SchemaObject.SchemaObject.Name = sdk.Pointer(sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(v.(string)))
		}
		if v, ok := onSchemaObject["all"]; ok && len(v.([]interface{})) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting all")
			all := v.([]interface{})[0].(map[string]interface{})
			on.SchemaObject.All = &sdk.GrantOnSchemaObjectIn{}
			resourceID.All = true
			pluralObjectType := all["object_type_plural"].(string)
			resourceID.ObjectTypePlural = pluralObjectType
			on.SchemaObject.All.PluralObjectType = sdk.PluralObjectType(pluralObjectType)
			if v, ok := all["in_database"]; ok && len(v.(string)) > 0 {
				logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting all in database")
				resourceID.InDatabase = true
				resourceID.DatabaseName = v.(string)
				on.SchemaObject.All.InDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(v.(string)))
			}
			if v, ok := all["in_schema"]; ok && len(v.(string)) > 0 {
				logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting all in schema")
				resourceID.InSchema = true
				resourceID.SchemaName = v.(string)
				on.SchemaObject.All.InSchema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(v.(string)))
			}
		}

		if v, ok := onSchemaObject["future"]; ok && len(v.([]interface{})) > 0 {
			logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting future")
			future := v.([]interface{})[0].(map[string]interface{})
			resourceID.Future = true
			on.SchemaObject.Future = &sdk.GrantOnSchemaObjectIn{}
			pluralObjectType := future["object_type_plural"].(string)
			resourceID.ObjectTypePlural = pluralObjectType
			on.SchemaObject.Future.PluralObjectType = sdk.PluralObjectType(pluralObjectType)
			if v, ok := future["in_database"]; ok && len(v.(string)) > 0 {
				logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting future in database")
				resourceID.InDatabase = true
				resourceID.DatabaseName = v.(string)
				on.SchemaObject.Future.InDatabase = sdk.Pointer(sdk.NewAccountObjectIdentifierFromFullyQualifiedName(v.(string)))
			}
			if v, ok := future["in_schema"]; ok && len(v.(string)) > 0 {
				logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: setting future in schema")
				resourceID.InSchema = true
				resourceID.SchemaName = v.(string)
				on.SchemaObject.Future.InSchema = sdk.Pointer(sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(v.(string)))
			}
		}

		privilegesToGrant = setRolePrivilegeOptions(privileges, allPrivileges, false, false, false, true)
		return privilegesToGrant, &on, nil
	}
	logging.DebugLogger.Printf("[DEBUG] Configuring account role grant privileges options: invalid")
	return nil, nil, fmt.Errorf("invalid grant options")
}

func setRolePrivilegeOptions(privileges []string, allPrivileges bool, onAccount bool, onAccountObject bool, onSchema bool, onSchemaObject bool) *sdk.AccountRoleGrantPrivileges {
	privilegesToGrant := &sdk.AccountRoleGrantPrivileges{}
	if allPrivileges {
		logging.DebugLogger.Printf("[DEBUG] Setting all privileges on privileges to grant")
		privilegesToGrant.AllPrivileges = sdk.Bool(true)
		return privilegesToGrant
	}
	if onAccount {
		logging.DebugLogger.Printf("[DEBUG] Setting global privileges on privileges to grant")
		privilegesToGrant.GlobalPrivileges = []sdk.GlobalPrivilege{}
		for _, privilege := range privileges {
			privilegesToGrant.GlobalPrivileges = append(privilegesToGrant.GlobalPrivileges, sdk.GlobalPrivilege(privilege))
		}
		return privilegesToGrant
	}
	if onAccountObject {
		logging.DebugLogger.Printf("[DEBUG] Setting account object privileges on privileges to grant")
		privilegesToGrant.AccountObjectPrivileges = []sdk.AccountObjectPrivilege{}
		for _, privilege := range privileges {
			privilegesToGrant.AccountObjectPrivileges = append(privilegesToGrant.AccountObjectPrivileges, sdk.AccountObjectPrivilege(privilege))
		}
		return privilegesToGrant
	}
	if onSchema {
		logging.DebugLogger.Printf("[DEBUG] Setting schema privileges on privileges to grant")
		privilegesToGrant.SchemaPrivileges = []sdk.SchemaPrivilege{}
		for _, privilege := range privileges {
			privilegesToGrant.SchemaPrivileges = append(privilegesToGrant.SchemaPrivileges, sdk.SchemaPrivilege(privilege))
		}
		return privilegesToGrant
	}
	if onSchemaObject {
		logging.DebugLogger.Printf("[DEBUG] Setting schema object privileges on privileges to grant")
		privilegesToGrant.SchemaObjectPrivileges = []sdk.SchemaObjectPrivilege{}
		for _, privilege := range privileges {
			privilegesToGrant.SchemaObjectPrivileges = append(privilegesToGrant.SchemaObjectPrivileges, sdk.SchemaObjectPrivilege(privilege))
		}
		return privilegesToGrant
	}
	logging.DebugLogger.Printf("[DEBUG] Not setting any privileges on privileges to grant")
	return nil
}

func readRoleGrantPrivileges(ctx context.Context, client *sdk.Client, grantedOn sdk.ObjectType, id GrantPrivilegesToRoleID, opts *sdk.ShowGrantOptions, d *schema.ResourceData) error {
	logging.DebugLogger.Printf("[DEBUG] About to show grants")
	grants, err := client.Grants.Show(ctx, opts)
	logging.DebugLogger.Printf("[DEBUG] After showing grants: err = %v", err)
	if err != nil {
		return fmt.Errorf("error retrieving grants for account role: %w", err)
	}

	withGrantOption := d.Get("with_grant_option").(bool)
	roleName := d.Get("role_name").(string)

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
		// Only consider privileges that are already present in the ID so we
		// don't delete privileges managed by other resources.
		if !slices.Contains(expectedPrivileges, grant.Privilege) {
			continue
		}
		if grant.GrantOption == withGrantOption && grant.GranteeName.Name() == sdk.NewAccountObjectIdentifier(roleName).Name() {
			// future grants do not have grantedBy, only current grants do. If grantedby
			// is an empty string it means the grant could not have been created by terraform
			if !id.Future && grant.GrantedBy.Name() == "" {
				continue
			}
			// grant_on is for future grants, granted_on is for current grants. They function the same way though in a test for matching the object type
			if grantedOn == grant.GrantedOn || grantedOn == grant.GrantOn {
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

	logging.DebugLogger.Printf("[DEBUG] Setting privileges on account")
	if err := d.Set("privileges", actualPrivileges); err != nil {
		logging.DebugLogger.Printf("[DEBUG] Error setting privileges for account role: err = %v", err)
		return fmt.Errorf("error setting privileges for account role: %w", err)
	}
	return nil
}
