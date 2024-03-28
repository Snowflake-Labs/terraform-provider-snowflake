package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	var opts *sdk.ShowGrantOptions
	var err error
	if grantsOn, ok := d.GetOk("grants_on"); ok {
		opts, err = buildOptsForGrantsOn(grantsOn.([]interface{})[0].(map[string]interface{}))
	}
	if grantsTo, ok := d.GetOk("grants_to"); ok {
		opts, err = buildOptsForGrantsTo(grantsTo.([]interface{})[0].(map[string]interface{}))
	}
	if grantsOf, ok := d.GetOk("grants_of"); ok {
		opts, err = buildOptsForGrantsOf(grantsOf.([]interface{})[0].(map[string]interface{}))
	}
	if futureGrantsIn, ok := d.GetOk("future_grants_in"); ok {
		opts, err = buildOptsForFutureGrantsIn(ctx, client, futureGrantsIn.([]interface{})[0].(map[string]interface{}))
	}
	if futureGrantsTo, ok := d.GetOk("future_grants_to"); ok {
		opts, err = buildOptsForFutureGrantsTo(futureGrantsTo.([]interface{})[0].(map[string]interface{}))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	grants, err := client.Grants.Show(ctx, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("grants", flattenGrants(grants))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("grants")
	return nil
}

func buildOptsForGrantsOn(grantsOn map[string]interface{}) (*sdk.ShowGrantOptions, error) {
	opts := new(sdk.ShowGrantOptions)

	objectType := grantsOn["object_type"].(string)
	objectName := grantsOn["object_name"].(string)
	account := grantsOn["account"].(bool)

	if account {
		opts.On = &sdk.ShowGrantsOn{
			Account: sdk.Bool(true),
		}
	} else {
		if objectType == "" || objectName == "" {
			return nil, fmt.Errorf("object_type (%s) or object_name (%s) missing", objectType, objectName)
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
	return opts, nil
}

func buildOptsForGrantsTo(grantsTo map[string]interface{}) (*sdk.ShowGrantOptions, error) {
	opts := new(sdk.ShowGrantOptions)

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
	return opts, nil
}

func buildOptsForGrantsOf(grantsOf map[string]interface{}) (*sdk.ShowGrantOptions, error) {
	opts := new(sdk.ShowGrantOptions)

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
	return opts, nil
}

func buildOptsForFutureGrantsIn(ctx context.Context, client *sdk.Client, futureGrantsIn map[string]interface{}) (*sdk.ShowGrantOptions, error) {
	opts := new(sdk.ShowGrantOptions)
	opts.Future = sdk.Bool(true)

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
	return opts, nil
}

func buildOptsForFutureGrantsTo(futureGrantsTo map[string]interface{}) (*sdk.ShowGrantOptions, error) {
	opts := new(sdk.ShowGrantOptions)
	opts.Future = sdk.Bool(true)

	if role := futureGrantsTo["role"].(string); role != "" {
		opts.To = &sdk.ShowGrantsTo{
			Role: sdk.NewAccountObjectIdentifier(role),
		}
	}
	if databaseRole := futureGrantsTo["database_role"].(string); databaseRole != "" {
		databaseRoleId, err := helpers.DecodeSnowflakeParameterID(databaseRole)
		if err != nil {
			return nil, err
		}
		validDatabaseRoleId, ok := databaseRoleId.(sdk.DatabaseObjectIdentifier)
		if !ok {
			return nil, fmt.Errorf("incorrect database role identifier (%s)", databaseRole)
		}
		opts.To = &sdk.ShowGrantsTo{
			DatabaseRole: validDatabaseRoleId,
		}
	}
	return opts, nil
}

func flattenGrants(grants []sdk.Grant) []map[string]interface{} {
	grantDetails := make([]map[string]interface{}, len(grants))
	for i, grant := range grants {
		grantDetails[i] = map[string]interface{}{
			"created_on":   grant.CreatedOn.String(),
			"privilege":    grant.Privilege,
			"granted_on":   grant.GrantedOn.String(),
			"name":         grant.Name.FullyQualifiedName(),
			"granted_to":   grant.GrantedTo.String(),
			"grantee_name": grant.GranteeName.FullyQualifiedName(),
			"grant_option": grant.GrantOption,
			"granted_by":   grant.GrantedBy.FullyQualifiedName(),
		}
	}
	return grantDetails
}
