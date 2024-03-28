package datasources

import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantsSchema = map[string]*schema.Schema{
	"grants_on": {
		Type:         schema.TypeList,
		MaxItems:     1,
		Optional:     true,
		Description:  "Lists all privileges that have been granted on an object or on an account.",
		ExactlyOneOf: []string{"grants_on", "grants_to", "grants_of", "future_grants_in", "future_grants_to"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_name": {
					Type:         schema.TypeString,
					Optional:     true,
					RequiredWith: []string{"grants_on.0.object_type"},
					ExactlyOneOf: []string{"grants_on.0.object_name", "grants_on.0.account"},
					Description:  "Name of object to list privileges on.",
				},
				"object_type": {
					Type:          schema.TypeString,
					Optional:      true,
					RequiredWith:  []string{"grants_on.0.object_name"},
					Description:   "Type of object to list privileges on.",
					ConflictsWith: []string{"grants_on.0.account"},
				},
				"account": {
					Type:          schema.TypeBool,
					Optional:      true,
					Description:   "Object hierarchy to list privileges on. The only valid value is: ACCOUNT. Setting this attribute lists all the account-level (i.e. global) privileges that have been granted to roles.",
					ExactlyOneOf:  []string{"grants_on.0.object_name", "grants_on.0.account"},
					ConflictsWith: []string{"grants_on.0.object_type"},
				},
			},
		},
	},
	"grants_to": {
		Type:         schema.TypeList,
		MaxItems:     1,
		Optional:     true,
		ExactlyOneOf: []string{"grants_on", "grants_to", "grants_of", "future_grants_in", "future_grants_to"},
		Description:  "Lists all privileges granted to the object.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"application": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all the privileges and roles granted to the application.",
					ExactlyOneOf: []string{
						"grants_to.0.application",
						"grants_to.0.application_role",
						"grants_to.0.role",
						"grants_to.0.user",
						"grants_to.0.share",
					},
				},
				"application_role": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all the privileges and roles granted to the application role. Note: if fully qualified application role identifier is not specified, i.e. only the application role name is given, Snowflake uses the current application. If the application is not a database, this command does not return results. Consult with the proper section in the [docs](https://docs.snowflake.com/en/sql-reference/sql/show-grants#variants).",
					ExactlyOneOf: []string{
						"grants_to.0.application",
						"grants_to.0.application_role",
						"grants_to.0.role",
						"grants_to.0.user",
						"grants_to.0.share",
					},
				},
				"role": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all privileges and roles granted to the role.",
					ExactlyOneOf: []string{
						"grants_to.0.application",
						"grants_to.0.application_role",
						"grants_to.0.role",
						"grants_to.0.user",
						"grants_to.0.share",
					},
				},
				"user": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all the roles granted to the user. Note that the PUBLIC role, which is automatically available to every user, is not listed.",
					ExactlyOneOf: []string{
						"grants_to.0.application",
						"grants_to.0.application_role",
						"grants_to.0.role",
						"grants_to.0.user",
						"grants_to.0.share",
					},
				},
				"share": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Lists all the privileges granted to the share.",
					ExactlyOneOf: []string{
						"grants_to.0.application",
						"grants_to.0.application_role",
						"grants_to.0.role",
						"grants_to.0.user",
						"grants_to.0.share",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"share_name": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Lists all of the privileges and roles granted to the specified share.",
							},
							"in_application_package": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Lists all of the privileges and roles granted to a share in the specified application package.",
							},
						},
					},
				},
			},
		},
	},
	"grants_of": {
		Type:         schema.TypeList,
		MaxItems:     1,
		Optional:     true,
		ExactlyOneOf: []string{"grants_on", "grants_to", "grants_of", "future_grants_in", "future_grants_to"},
		Description:  "Lists all objects to which the given object has been granted.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all users and roles to which the role has been granted.",
					ExactlyOneOf: []string{
						"grants_of.0.role",
						"grants_of.0.application_role",
						"grants_of.0.share",
					},
				},
				"application_role": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all the users and roles to which the application role has been granted. Note: if fully qualified application role identifier is not specified, i.e. only the application role name is given, Snowflake uses the current application. If the application is not a database, this command does not return results. Consult with the proper section in the [docs](https://docs.snowflake.com/en/sql-reference/sql/show-grants#variants).",
					ExactlyOneOf: []string{
						"grants_of.0.role",
						"grants_of.0.application_role",
						"grants_of.0.share",
					},
				},
				"share": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all the accounts for the share and indicates the accounts that are using the share.",
					ExactlyOneOf: []string{
						"grants_of.0.role",
						"grants_of.0.application_role",
						"grants_of.0.share",
					},
				},
			},
		},
	},
	"future_grants_in": {
		Type:         schema.TypeList,
		MaxItems:     1,
		Optional:     true,
		ExactlyOneOf: []string{"grants_on", "grants_to", "grants_of", "future_grants_in", "future_grants_to"},
		Description:  "Lists all privileges on new (i.e. future) objects.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all privileges on new (i.e. future) objects of a specified type in the database granted to a role.",
					ExactlyOneOf: []string{
						"future_grants_in.0.database",
						"future_grants_in.0.schema",
					},
				},
				"schema": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Lists all privileges on new (i.e. future) objects of a specified type in the schema granted to a role.",
					ExactlyOneOf: []string{
						"future_grants_in.0.database",
						"future_grants_in.0.schema",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"schema_name": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The name of the schema to list all privileges of new (ie. future) objects granted to.",
							},
							"database_name": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The database in which the scehma resides. Optional when querying a schema in the current database.",
							},
						},
					},
				},
			},
		},
	},
	"future_grants_to": {
		Type:         schema.TypeList,
		MaxItems:     1,
		Optional:     true,
		ExactlyOneOf: []string{"grants_on", "grants_to", "grants_of", "future_grants_in", "future_grants_to"},
		Description:  "Lists all privileges granted to the object on new (i.e. future) objects.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all privileges on new (i.e. future) objects of a specified type in a database or schema granted to the role.",
					ExactlyOneOf: []string{
						"future_grants_to.0.role",
						"future_grants_to.0.database_role",
					},
				},
				"database_role": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all privileges on new (i.e. future) objects granted to the database role.",
					ExactlyOneOf: []string{
						"future_grants_to.0.role",
						"future_grants_to.0.database_role",
					},
				},
			},
		},
	},
	"grants": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The list of grants",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"created_on": {
					Type:        schema.TypeString,
					Description: "The date and time the grant was created.",
					Computed:    true,
				},
				"privilege": {
					Type:        schema.TypeString,
					Description: "The privilege granted.",
					Computed:    true,
				},
				"granted_on": {
					Type:        schema.TypeString,
					Description: "The object on which the privilege was granted.",
					Computed:    true,
				},
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the object on which the privilege was granted.",
					Computed:    true,
				},
				"granted_to": {
					Type:        schema.TypeString,
					Description: "The role to which the privilege was granted.",
					Computed:    true,
				},
				"grantee_name": {
					Type:        schema.TypeString,
					Description: "The name of the role to which the privilege was granted.",
					Computed:    true,
				},
				"grant_option": {
					Type:        schema.TypeBool,
					Description: "Whether the grantee can grant the privilege to others.",
					Computed:    true,
				},
				"granted_by": {
					Type:        schema.TypeString,
					Description: "The role that granted the privilege.",
					Computed:    true,
				},
			},
		},
	},
}

func Grants() *schema.Resource {
	return &schema.Resource{
		ReadContext: ReadGrants,
		Schema:      grantsSchema,
	}
}

func ReadGrants(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB

	var grantDetails []snowflake.GrantDetail
	var grants []sdk.Grant
	_ = grants
	var err error
	if v, ok := d.GetOk("grants_on"); ok {
		grants, err = handleGrantsOn(ctx, client, v)
	}

	if v, ok := d.GetOk("grants_to"); ok {
		grants, err = handleGrantsTo(ctx, client, v)
	}

	if v, ok := d.GetOk("grants_of"); ok {
		grants, err = handleGrantsOf(ctx, client, v)
	}

	if v, ok := d.GetOk("future_grants_in"); ok {
		grants, err = handleFutureGrantsIn(ctx, client, v)
	}

	if v, ok := d.GetOk("future_grants_to"); ok {
		grantDetails, err = handleFutureGrantsTo(ctx, client, v, db)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("grants", flattenGrants(grantDetails))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("grants")
	return nil
}

func handleGrantsOn(ctx context.Context, client *sdk.Client, v any) ([]sdk.Grant, error) {
	var err error
	opts := new(sdk.ShowGrantOptions)

	grantsOn := v.([]interface{})[0].(map[string]interface{})
	objectType := grantsOn["object_type"].(string)
	objectName := grantsOn["object_name"].(string)
	account := grantsOn["account"].(bool)

	if account {
		opts.On = &sdk.ShowGrantsOn{
			Account: sdk.Bool(true),
		}
	} else {
		if objectType == "" || objectName == "" {
			return nil, err
		}
		objectId, err := helpers.DecodeSnowflakeParameterID(objectName)
		if err != nil {
			return nil, err
		}
		opts.On = &sdk.ShowGrantsOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectType(objectType),
				Name:       objectId,
			},
		}
	}
	return client.Grants.Show(ctx, opts)
}

func handleGrantsTo(ctx context.Context, client *sdk.Client, v any) ([]sdk.Grant, error) {
	opts := new(sdk.ShowGrantOptions)
	grantsTo := v.([]interface{})[0].(map[string]interface{})

	// TODO: add database role?
	if application := grantsTo["application"].(string); application != "" {
		// TODO: unsupported SHOW GRANTS TO APPLICATION
	}
	if applicationRole := grantsTo["application_role"].(string); applicationRole != "" {
		// TODO: unsupported SHOW GRANTS TO APPLICATION ROLE
	}
	if role := grantsTo["role"].(string); role != "" {
		opts.To = &sdk.ShowGrantsTo{
			Role: sdk.NewAccountObjectIdentifier(role),
		}
	}
	if user := grantsTo["user"].(string); user != "" {
		opts.To = &sdk.ShowGrantsTo{
			User: sdk.NewAccountObjectIdentifier(user),
		}
	}
	if share := grantsTo["share"]; share != nil {
		shareMap := share.([]interface{})[0].(map[string]interface{})
		opts.To = &sdk.ShowGrantsTo{
			Share: sdk.NewAccountObjectIdentifier(shareMap["share_name"].(string)),
		}
		// TODO: unsupported IN APPLICATION PACKAGE
	}
	return client.Grants.Show(ctx, opts)
}

func handleGrantsOf(ctx context.Context, client *sdk.Client, v any) ([]sdk.Grant, error) {
	opts := new(sdk.ShowGrantOptions)
	grantsOf := v.([]interface{})[0].(map[string]interface{})

	// TODO: add database role?
	if role := grantsOf["role"].(string); role != "" {
		opts.Of = &sdk.ShowGrantsOf{
			Role: sdk.NewAccountObjectIdentifier(role),
		}
	}
	if applicationRole := grantsOf["application_role"].(string); applicationRole != "" {
		// TODO: unsupported SHOW GRANTS OF APPLICATION ROLE
	}
	if share := grantsOf["share"].(string); share != "" {
		opts.Of = &sdk.ShowGrantsOf{
			Share: sdk.NewAccountObjectIdentifier(share),
		}
	}
	return client.Grants.Show(ctx, opts)
}

func handleFutureGrantsIn(ctx context.Context, client *sdk.Client, v any) ([]sdk.Grant, error) {
	opts := new(sdk.ShowGrantOptions)
	opts.Future = sdk.Bool(true)
	futureGrantsIn := v.([]interface{})[0].(map[string]interface{})

	if db := futureGrantsIn["database"].(string); db != "" {
		opts.In = &sdk.ShowGrantsIn{
			Database: sdk.Pointer(sdk.NewAccountObjectIdentifier(db)),
		}
	}
	if sc := futureGrantsIn["schema"].([]interface{}); len(sc) > 0 {
		schemaMap := sc[0].(map[string]interface{})
		schemaName := schemaMap["schema_name"].(string)
		databaseName := schemaMap["database_name"].(string)
		if databaseName == "" {
			current, err := client.ContextFunctions.CurrentDatabase(ctx)
			if err != nil {
				return nil, err
			}
			databaseName = current
		}
		opts.In = &sdk.ShowGrantsIn{
			Schema: sdk.Pointer(sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)),
		}
	}
	return client.Grants.Show(ctx, opts)
}

func handleFutureGrantsTo(ctx context.Context, client *sdk.Client, v any, db *sql.DB) ([]snowflake.GrantDetail, error) {
	var grantDetails []snowflake.GrantDetail
	var err error

	futureGrantsTo := v.([]interface{})[0].(map[string]interface{})
	role := futureGrantsTo["role"].(string)
	if role != "" {
		grantDetails, err = snowflake.ShowFutureGrantsTo(db, "ROLE", role)
		if err != nil {
			return grantDetails, err
		}
	}
	return grantDetails, nil
}

func aaa() (*sdk.ShowGrantOptions, sdk.ObjectType) {
	opts := new(sdk.ShowGrantOptions)
	var grantedOn sdk.ObjectType

	//switch id.Kind {
	//case OnAccountObjectAccountRoleGrantKind:
	//	data := id.Data.(*OnAccountObjectGrantData)
	//	grantedOn = data.ObjectType
	//	opts.On = &sdk.ShowGrantsOn{
	//		Object: &sdk.Object{
	//			ObjectType: data.ObjectType,
	//			Name:       data.ObjectName,
	//		},
	//	}
	//case OnSchemaAccountRoleGrantKind:
	//	grantedOn = sdk.ObjectTypeSchema
	//	data := id.Data.(*OnSchemaGrantData)
	//
	//	switch data.Kind {
	//	case OnSchemaSchemaGrantKind:
	//		opts.On = &sdk.ShowGrantsOn{
	//			Object: &sdk.Object{
	//				ObjectType: sdk.ObjectTypeSchema,
	//				Name:       data.SchemaName,
	//			},
	//		}
	//	case OnAllSchemasInDatabaseSchemaGrantKind:
	//		log.Printf("[INFO] Show with on_schema.all_schemas_in_database option is skipped. No changes in privileges in Snowflake will be detected.")
	//		return nil, ""
	//	case OnFutureSchemasInDatabaseSchemaGrantKind:
	//		opts.Future = sdk.Bool(true)
	//		opts.In = &sdk.ShowGrantsIn{
	//			Database: data.DatabaseName,
	//		}
	//	}
	//case OnSchemaObjectAccountRoleGrantKind:
	//	data := id.Data.(*OnSchemaObjectGrantData)
	//
	//	switch data.Kind {
	//	case OnObjectSchemaObjectGrantKind:
	//		grantedOn = data.Object.ObjectType
	//		opts.On = &sdk.ShowGrantsOn{
	//			Object: data.Object,
	//		}
	//	case OnAllSchemaObjectGrantKind:
	//		log.Printf("[INFO] Show with on_schema_object.on_all option is skipped. No changes in privileges in Snowflake will be detected.")
	//		return nil, ""
	//	case OnFutureSchemaObjectGrantKind:
	//		grantedOn = data.OnAllOrFuture.ObjectNamePlural.Singular()
	//		opts.Future = sdk.Bool(true)
	//
	//		switch data.OnAllOrFuture.Kind {
	//		case InDatabaseBulkOperationGrantKind:
	//			opts.In = &sdk.ShowGrantsIn{
	//				Database: data.OnAllOrFuture.Database,
	//			}
	//		case InSchemaBulkOperationGrantKind:
	//			opts.In = &sdk.ShowGrantsIn{
	//				Schema: data.OnAllOrFuture.Schema,
	//			}
	//		}
	//	}
	//}

	return opts, grantedOn
}

func flattenGrants(grants []snowflake.GrantDetail) []map[string]interface{} {
	grantDetails := make([]map[string]interface{}, len(grants))
	for i, grant := range grants {
		grantDetails[i] = map[string]interface{}{
			"created_on":   grant.CreatedOn.String,
			"privilege":    grant.Privilege.String,
			"granted_on":   grant.GrantedOn.String,
			"name":         grant.Name.String,
			"granted_to":   grant.GrantedTo.String,
			"grantee_name": grant.GranteeName.String,
			"grant_option": grant.GrantOption.String == "true",
			"granted_by":   grant.GrantedBy.String,
		}
	}
	return grantDetails
}
