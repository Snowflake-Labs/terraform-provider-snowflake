package datasources

import (
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantsSchema = map[string]*schema.Schema{
	"grants_on": {
		Type:          schema.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: []string{"grants_of", "grants_to", "future_grants_in", "future_grants_to"},
		Description:   "Lists all privileges that have been granted on an object or account",
		ExactlyOneOf:  []string{"grants_on", "grants_of", "grants_to", "future_grants_in", "future_grants_to"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_name": {
					Type:          schema.TypeString,
					Optional:      true,
					RequiredWith:  []string{"grants_on.0.object_type"},
					Description:   "Name of object to list privileges on",
					ConflictsWith: []string{"grants_on.0.account"},
					AtLeastOneOf:  []string{"grants_on.0.object_name", "grants_on.0.account"},
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
					ConflictsWith: []string{"grants_on.0.object_name", "grants_on.0.object_type"},
					AtLeastOneOf:  []string{"grants_on.0.object_name", "grants_on.0.account"},
				},
			},
		},
	},
	"grants_to": {
		Type:          schema.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: []string{"grants_on", "grants_of", "future_grants_in", "future_grants_to"},
		Description:   "Lists all privileges granted to the object",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all privileges and roles granted to the role",
					ConflictsWith: []string{
						"grants_to.0.user",
						"grants_to.0.share",
					},
					ExactlyOneOf: []string{
						"grants_to.0.role",
						"grants_to.0.user",
						"grants_to.0.share",
					},
				},
				"user": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all the roles granted to the user. Note that the PUBLIC role, which is automatically available to every user, is not listed",
					ConflictsWith: []string{
						"grants_to.0.role",
						"grants_to.0.share",
					},
					ExactlyOneOf: []string{
						"grants_to.0.role",
						"grants_to.0.user",
						"grants_to.0.share",
					},
				},
				"share": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all the privileges granted to the share",
					ConflictsWith: []string{
						"grants_to.0.role",
						"grants_to.0.user",
					},
					ExactlyOneOf: []string{
						"grants_to.0.role",
						"grants_to.0.user",
						"grants_to.0.share",
					},
				},
			},
		},
	},
	"grants_of": {
		Type:          schema.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: []string{"grants_on", "grants_to", "future_grants_in", "future_grants_to"},
		Description:   "Lists all objects to which the given object has been granted",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all users and roles to which the role has been granted",
					ConflictsWith: []string{
						"grants_of.0.share",
					},
					ExactlyOneOf: []string{
						"grants_of.0.role",
						"grants_of.0.share",
					},
				},
				"share": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all the accounts for the share and indicates the accounts that are using the share.",
					ConflictsWith: []string{
						"grants_of.0.role",
					},
					ExactlyOneOf: []string{
						"grants_of.0.role",
						"grants_of.0.share",
					},
				},
			},
		},
	},
	"future_grants_in": {
		Type:          schema.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: []string{"grants_on", "grants_of", "grants_to", "future_grants_to"},
		Description:   "Lists all privileges on new (i.e. future) objects",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"database": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Lists all privileges on new (i.e. future) objects of a specified type in the database granted to a role.",
					ConflictsWith: []string{
						"future_grants_in.0.schema",
					},
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
					ConflictsWith: []string{
						"future_grants_in.0.database",
					},
					ExactlyOneOf: []string{
						"future_grants_in.0.database",
						"future_grants_in.0.schema",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"schema_name": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The name of the schema to list all privileges of new (ie. future) objects granted to",
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
		Type:          schema.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: []string{"grants_on", "grants_of", "grants_to", "future_grants_in"},
		Description:   "Lists all privileges granted to the object on new (i.e. future) objects",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Lists all privileges on new (i.e. future) objects of a specified type in a database or schema granted to the role.",
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
					Description: "The date and time the grant was created",
					Computed:    true,
				},
				"privilege": {
					Type:        schema.TypeString,
					Description: "The privilege granted",
					Computed:    true,
				},
				"granted_on": {
					Type:        schema.TypeString,
					Description: "The object on which the privilege was granted",
					Computed:    true,
				},
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the object on which the privilege was granted",
					Computed:    true,
				},
				"granted_to": {
					Type:        schema.TypeString,
					Description: "The role to which the privilege was granted",
					Computed:    true,
				},
				"grantee_name": {
					Type:        schema.TypeString,
					Description: "The name of the role to which the privilege was granted",
					Computed:    true,
				},
				"grant_option": {
					Type:        schema.TypeBool,
					Description: "Whether the grantee can grant the privilege to others",
					Computed:    true,
				},
				"granted_by": {
					Type:        schema.TypeString,
					Description: "The role that granted the privilege",
					Computed:    true,
				},
			},
		},
	},
}

func Grants() *schema.Resource {
	return &schema.Resource{
		Read:   ReadGrants,
		Schema: grantsSchema,
	}
}

func ReadGrants(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	var grantDetails []snowflake.GrantDetail
	var err error
	if v, ok := d.GetOk("grants_on"); ok {
		grantsOn := v.([]interface{})[0].(map[string]interface{})
		objectType := grantsOn["object_type"].(string)
		objectName := grantsOn["object_name"].(string)
		account := grantsOn["account"].(bool)

		if account {
			grantDetails, err = snowflake.ShowGrantsOnAccount(db)
			if err != nil {
				return err
			}
		} else if objectType != "" && objectName != "" {
			grantDetails, err = snowflake.ShowGrantsOn(db, objectType, objectName)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("grants_to"); ok {
		grantsTo := v.([]interface{})[0].(map[string]interface{})
		role := grantsTo["role"].(string)
		if role != "" {
			grantDetails, err = snowflake.ShowGrantsTo(db, "ROLE", role)
			if err != nil {
				return err
			}
		}
		user := grantsTo["user"].(string)
		if user != "" {
			grantDetails, err = snowflake.ShowGrantsTo(db, "USER", user)
			if err != nil {
				return err
			}
		}
		share := grantsTo["share"].(string)
		if share != "" {
			grantDetails, err = snowflake.ShowGrantsTo(db, "SHARE", share)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("grants_of"); ok {
		grantsOf := v.([]interface{})[0].(map[string]interface{})
		role := grantsOf["role"].(string)
		if role != "" {
			grantDetails, err = snowflake.ShowGrantsOf(db, "ROLE", role)
			if err != nil {
				return err
			}
		}
		share := grantsOf["share"].(string)
		if share != "" {
			grantDetails, err = snowflake.ShowGrantsOf(db, "SHARE", share)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("future_grants_in"); ok {
		futureGrantsIn := v.([]interface{})[0].(map[string]interface{})
		database := futureGrantsIn["database"].(string)
		if database != "" {
			grantDetails, err = snowflake.ShowFutureGrantsIn(db, "DATABASE", database)
			if err != nil {
				return err
			}
		}
		schema := futureGrantsIn["schema"].([]interface{})
		if len(schema) > 0 {
			schemaMap := schema[0].(map[string]interface{})
			schemaName := schemaMap["schema_name"].(string)
			databaseName := schemaMap["database_name"].(string)
			if databaseName != "" {
				schemaName = databaseName + "." + schemaName
			}

			grantDetails, err = snowflake.ShowFutureGrantsIn(db, "SCHEMA", schemaName)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("future_grants_to"); ok {
		futureGrantsTo := v.([]interface{})[0].(map[string]interface{})
		role := futureGrantsTo["role"].(string)
		if role != "" {
			grantDetails, err = snowflake.ShowFutureGrantsTo(db, "ROLE", role)
			if err != nil {
				return err
			}
		}
	}

	err = d.Set("grants", flattenGrants(grantDetails))
	if err != nil {
		return err
	}
	d.SetId("grants")
	return nil
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
