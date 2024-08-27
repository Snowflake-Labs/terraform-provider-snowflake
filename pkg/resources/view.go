package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var viewSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the view; must be unique for the schema in which the view is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the view."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the view."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"or_replace": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Overwrites the View if it exists.",
	},
	// TODO [SNOW-1348118: this is used only during or_replace, we would like to change the behavior before v1
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Retains the access permissions from the original view when a new view is created using the OR REPLACE clause. OR REPLACE must be set when COPY GRANTS is set.",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return oldValue != "" && oldValue != newValue
		},
		RequiredWith: []string{"or_replace"},
	},
	"is_secure": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_secure"),
		Description:      booleanStringFieldDescription("Specifies that the view is secure. By design, the Snowflake's `SHOW VIEWS` command does not provide information about secure views (consult [view usage notes](https://docs.snowflake.com/en/sql-reference/sql/create-view#usage-notes)) which is essential to manage/import view with Terraform. Use the role owning the view while managing secure views."),
	},
	"is_temporary": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Specifies that the view persists only for the duration of the session that you created it in. A temporary view and all its contents are dropped at the end of the session. In context of this provider, it means that it's dropped after a Terraform operation. This results in a permanent plan with object creation."),
	},
	"is_recursive": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Specifies that the view can refer to itself using recursive syntax without necessarily using a CTE (common table expression)."),
	},
	"change_tracking": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("change_tracking", func(x any) any {
			return x.(string) == "ON"
		}),
		Description: booleanStringFieldDescription("Specifies to enable or disable change tracking on the table."),
	},
	// TODO(next pr): support remaining fields
	// "data_metric_functions": {
	// 	Type:     schema.TypeSet,
	// 	Optional: true,
	// 	Elem: &schema.Resource{
	// 		Schema: map[string]*schema.Schema{
	// 			"metric_name": {
	// 				Type:        schema.TypeString,
	// 				Optional:    true,
	// 				Description: "Identifier of the data metric function to add to the table or view or drop from the table or view.",
	// 			},
	// 			"column_name": {
	// 				Type:        schema.TypeString,
	// 				Optional:    true,
	// 				Description: "The table or view columns on which to associate the data metric function. The data types of the columns must match the data types of the columns specified in the data metric function definition.",
	// 			},
	// 		},
	// 	},
	// 	Description: "Data metric functions used for the view.",
	// },
	// "data_metric_schedule": {
	// 	Type:     schema.TypeList,
	// 	Optional: true,
	// 	MaxItems: 1,
	// 	Elem: &schema.Resource{
	// 		Schema: map[string]*schema.Schema{
	// 			"minutes": {
	// 				Type:        schema.TypeInt,
	// 				Optional:    true,
	// 				Description: "Specifies an interval (in minutes) of wait time inserted between runs of the data metric function. Conflicts with `using_cron` and `trigger_on_changes`.",
	// 				// TODO: move to sdk
	// 				ValidateFunc:  validation.IntInSlice([]int{5, 15, 30, 60, 720, 1440}),
	// 				ConflictsWith: []string{"data_metric_schedule.using_cron", "data_metric_schedule.trigger_on_changes"},
	// 			},
	// 			"using_cron": {
	// 				Type:        schema.TypeString,
	// 				Optional:    true,
	// 				Description: "Specifies a cron expression and time zone for periodically running the data metric function. Supports a subset of standard cron utility syntax. Conflicts with `minutes` and `trigger_on_changes`.",
	// 				// TODO: validate?
	// 				ConflictsWith: []string{"data_metric_schedule.minutes", "data_metric_schedule.trigger_on_changes"},
	// 			},
	// 			"trigger_on_changes": {
	// 				Type:          schema.TypeString,
	// 				Optional:      true,
	// 				Default:       BooleanDefault,
	// 				Description:   booleanStringFieldDescription("Specifies that the DMF runs when a DML operation modifies the table, such as inserting a new row or deleting a row. Conflicts with `minutes` and `using_cron`."),
	// 				ConflictsWith: []string{"data_metric_schedule.minutes", "data_metric_schedule.using_cron"},
	// 			},
	// 		},
	// 	},
	// 	Description: "Specifies the schedule to run the data metric function periodically.",
	// },
	// "columns": {
	// 	Type:     schema.TypeList,
	// 	Optional: true,
	// 	Elem: &schema.Resource{
	// 		Schema: map[string]*schema.Schema{
	// 			"column_name": {
	// 				Type:        schema.TypeString,
	// 				Required:    true,
	// 				Description: "Specifies affected column name.",
	// 			},
	// 			"masking_policy": {
	// 				Type:     schema.TypeList,
	// 				Optional: true,
	// 				Elem: &schema.Resource{
	// 					Schema: map[string]*schema.Schema{
	// 						// TODO: change to `name`? in other policies as well
	// 						"policy_name": {
	// 							Type:        schema.TypeString,
	// 							Required:    true,
	// 							Description: "Specifies the masking policy to set on a column.",
	// 						},
	// 						"using": {
	// 							Type:     schema.TypeList,
	// 							Optional: true,
	// 							Elem: &schema.Schema{
	// 								Type: schema.TypeString,
	// 							},
	// 							Description: "Specifies the arguments to pass into the conditional masking policy SQL expression. The first column in the list specifies the column for the policy conditions to mask or tokenize the data and must match the column to which the masking policy is set. The additional columns specify the columns to evaluate to determine whether to mask or tokenize the data in each row of the query result when a query is made on the first column. If the USING clause is omitted, Snowflake treats the conditional masking policy as a normal masking policy.",
	// 						},
	// 					},
	// 				},
	// 			},
	// 			"projection_policy": {
	// 				Type:             schema.TypeString,
	// 				Optional:         true,
	// 				DiffSuppressFunc: DiffSuppressStatement,
	// 				Description:      "Specifies the projection policy to set on a column.",
	// 			},
	// "comment": {
	// 	Type:        schema.TypeString,
	// 	Optional:    true,
	// 	Description: "Specifies a comment for the column.",
	// },
	// 		},
	// 	},
	// 	Description: "If you want to change the name of a column or add a comment to a column in the new view, include a column list that specifies the column names and (if needed) comments about the columns. (You do not need to specify the data types of the columns.)",
	// },
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the view.",
	},
	"row_access_policy": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy_name": {
					Type:             schema.TypeString,
					Required:         true,
					DiffSuppressFunc: suppressIdentifierQuoting,
					Description:      "Row access policy name.",
				},
				"on": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Defines which columns are affected by the policy.",
				},
			},
		},
		Description: "Specifies the row access policy to set on a view.",
	},
	"aggregation_policy": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy_name": {
					Type:             schema.TypeString,
					Required:         true,
					DiffSuppressFunc: suppressIdentifierQuoting,
					Description:      "Aggregation policy name.",
				},
				"entity_key": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Defines which columns uniquely identify an entity within the view.",
				},
			},
		},
		Description: "Specifies the aggregation policy to set on a view.",
	},
	"statement": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the query used to create the view.",
		DiffSuppressFunc: DiffSuppressStatement,
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW VIEW` for the given view.",
		Elem: &schema.Resource{
			Schema: schemas.ShowViewSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE VIEW` for the given view.",
		Elem: &schema.Resource{
			Schema: schemas.ViewDescribeSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// View returns a pointer to the resource representing a view.
func View() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateView(false),
		ReadContext:   ReadView(true),
		UpdateContext: UpdateView,
		DeleteContext: DeleteView,
		Description:   "Resource used to manage view objects. For more information, check [view documentation](https://docs.snowflake.com/en/sql-reference/sql/create-view).",

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "comment", "change_tracking", "is_secure", "is_temporary", "is_recursive", "statement"),
			ComputedIfAnyAttributeChanged(FullyQualifiedNameAttributeName, "name"),
		),

		Schema: viewSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportView,
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v0_94_1_ViewStateUpgrader,
			},
		},
	}
}

func ImportView(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Starting view import")
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	v, err := client.Views.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("name", v.Name); err != nil {
		return nil, err
	}

	if err := d.Set("change_tracking", booleanStringFromBool(v.IsChangeTracking())); err != nil {
		return nil, err
	}
	if err := d.Set("is_recursive", booleanStringFromBool(v.IsRecursive())); err != nil {
		return nil, err
	}
	if err := d.Set("is_secure", booleanStringFromBool(v.IsSecure)); err != nil {
		return nil, err
	}
	if err := d.Set("is_temporary", booleanStringFromBool(v.IsTemporary())); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateView(orReplace bool) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		statement := d.Get("statement").(string)
		req := sdk.NewCreateViewRequest(id, statement)

		// TODO(next pr): remove or_replace field
		if v := d.Get("or_replace"); v.(bool) || orReplace {
			req.WithOrReplace(true)
		}

		if v := d.Get("copy_grants"); v.(bool) {
			req.WithCopyGrants(true)
		}

		if v := d.Get("is_secure").(string); v != BooleanDefault {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			req.WithSecure(parsed)
		}

		if v := d.Get("is_temporary").(string); v != BooleanDefault {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			req.WithTemporary(parsed)
		}

		if v := d.Get("is_recursive").(string); v != BooleanDefault {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			req.WithRecursive(parsed)
		}

		if v := d.Get("comment").(string); len(v) > 0 {
			req.WithComment(v)
		}

		if v := d.Get("row_access_policy"); len(v.([]any)) > 0 {
			req.WithRowAccessPolicy(*sdk.NewViewRowAccessPolicyRequest(extractPolicyWithColumns(v, "on")))
		}

		if v := d.Get("aggregation_policy"); len(v.([]any)) > 0 {
			id, columns := extractPolicyWithColumns(v, "entity_key")
			aggregationPolicyReq := sdk.NewViewAggregationPolicyRequest(id)
			if len(columns) > 0 {
				aggregationPolicyReq.WithEntityKey(columns)
			}
			req.WithAggregationPolicy(*aggregationPolicyReq)
		}

		err := client.Views.Create(ctx, req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating view %v err = %w", id.Name(), err))
		}

		d.SetId(helpers.EncodeSnowflakeID(id))

		if v := d.Get("change_tracking").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}

			err = client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetChangeTracking(parsed))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting change_tracking in view %v err = %w", id.Name(), err))
			}
		}

		return ReadView(false)(ctx, d, meta)
	}
}

func extractPolicyWithColumns(v any, columnsKey string) (sdk.SchemaObjectIdentifier, []sdk.Column) {
	policyConfig := v.([]any)[0].(map[string]any)
	columnsRaw := expandStringList(policyConfig[columnsKey].(*schema.Set).List())
	columns := make([]sdk.Column, len(columnsRaw))
	for i := range columnsRaw {
		columns[i] = sdk.Column{Value: columnsRaw[i]}
	}
	return sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(policyConfig["policy_name"].(string)), columns
}

func ReadView(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

		view, err := client.Views.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query view. Marking the resource as removed.",
						Detail:   fmt.Sprintf("View: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		if err = d.Set("name", view.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("database", view.DatabaseName); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("schema", view.SchemaName); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("copy_grants", view.HasCopyGrants()); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("comment", view.Comment); err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"is_secure", "is_secure", view.IsSecure, booleanStringFromBool(view.IsSecure), nil},
				showMapping{"text", "is_recursive", view.IsRecursive(), booleanStringFromBool(view.IsRecursive()), func(x any) any {
					return strings.Contains(x.(string), "RECURSIVE")
				}},
				showMapping{"text", "is_temporary", view.IsTemporary(), booleanStringFromBool(view.IsTemporary()), func(x any) any {
					return strings.Contains(x.(string), "TEMPORARY")
				}},
				showMapping{"change_tracking", "change_tracking", view.IsChangeTracking(), booleanStringFromBool(view.IsChangeTracking()), func(x any) any {
					return x.(string) == "ON"
				}},
			); err != nil {
				return diag.FromErr(err)
			}
		}
		if err = setStateToValuesFromConfig(d, viewSchema, []string{
			"change_tracking",
			"is_recursive",
			"is_secure",
			"is_temporary",
		}); err != nil {
			return diag.FromErr(err)
		}

		err = handlePolicyReferences(ctx, client, id, d)
		if err != nil {
			return diag.FromErr(err)
		}
		if view.Text != "" {
			// Want to only capture the SELECT part of the query because before that is the CREATE part of the view.
			extractor := snowflake.NewViewSelectStatementExtractor(view.Text)
			statement, err := extractor.Extract()
			if err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("statement", statement); err != nil {
				return diag.FromErr(err)
			}
		} else {
			return diag.FromErr(fmt.Errorf("error reading view %v, err = %w, `text` is missing; if the view is secure then the role used by the provider must own the view (consult https://docs.snowflake.com/en/sql-reference/sql/create-view#usage-notes)", d.Id(), err))
		}

		describeResult, err := client.Views.Describe(ctx, view.ID())
		if err != nil {
			log.Printf("[DEBUG] describing view: %s, err: %s", id.FullyQualifiedName(), err)
		} else {
			if err = d.Set(DescribeOutputAttributeName, schemas.ViewDescriptionToSchema(describeResult)); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.ViewToSchema(view)}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func handlePolicyReferences(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, d *schema.ResourceData) error {
	policyRefs, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(id, sdk.PolicyEntityDomainView))
	if err != nil {
		return fmt.Errorf("getting policy references for view: %w", err)
	}
	var aggregationPolicies []map[string]any
	var rowAccessPolicies []map[string]any
	for _, p := range policyRefs {
		policyName := sdk.NewSchemaObjectIdentifier(*p.PolicyDb, *p.PolicySchema, p.PolicyName)
		switch p.PolicyKind {
		case string(sdk.PolicyKindAggregationPolicy):
			var entityKey []string
			if p.RefArgColumnNames != nil {
				entityKey = sdk.ParseCommaSeparatedStringArray(*p.RefArgColumnNames, true)
			}
			aggregationPolicies = append(aggregationPolicies, map[string]any{
				"policy_name": policyName.FullyQualifiedName(),
				"entity_key":  entityKey,
			})
		case string(sdk.PolicyKindRowAccessPolicy):
			var on []string
			if p.RefArgColumnNames != nil {
				on = sdk.ParseCommaSeparatedStringArray(*p.RefArgColumnNames, true)
			}
			rowAccessPolicies = append(rowAccessPolicies, map[string]any{
				"policy_name": policyName.FullyQualifiedName(),
				"on":          on,
			})
		default:
			log.Printf("[WARN] unexpected policy kind %v in policy references returned from Snowflake", p.PolicyKind)
		}
	}
	if err = d.Set("aggregation_policy", aggregationPolicies); err != nil {
		return err
	}
	if err = d.Set("row_access_policy", rowAccessPolicies); err != nil {
		return err
	}
	return err
}

func UpdateView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	// change on these fields can not be ForceNew because then view is dropped explicitly and copying grants does not have effect
	if d.HasChange("statement") || d.HasChange("is_temporary") || d.HasChange("is_recursive") || d.HasChange("copy_grant") {
		return CreateView(true)(ctx, d, meta)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming view %v err = %w", d.Id(), err))
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	if d.HasChange("comment") {
		if comment := d.Get("comment").(string); comment == "" {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetComment(true))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting comment for view %v", d.Id()))
			}
		} else {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetComment(comment))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting comment for view %v", d.Id()))
			}
		}
	}

	if d.HasChange("is_secure") {
		if v := d.Get("is_secure").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			err = client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetSecure(parsed))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting is_secure for view %v: %w", d.Id(), err))
			}
		} else {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetSecure(true))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting is_secure for view %v: %w", d.Id(), err))
			}
		}
	}
	if d.HasChange("change_tracking") {
		if v := d.Get("change_tracking").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			err = client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetChangeTracking(parsed))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting change_tracking for view %v: %w", d.Id(), err))
			}
		} else {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetChangeTracking(false))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting change_tracking for view %v: %w", d.Id(), err))
			}
		}
	}

	if d.HasChange("row_access_policy") {
		var addReq *sdk.ViewAddRowAccessPolicyRequest
		var dropReq *sdk.ViewDropRowAccessPolicyRequest

		oldRaw, newRaw := d.GetChange("row_access_policy")
		if len(oldRaw.([]any)) > 0 {
			oldId, _ := extractPolicyWithColumns(oldRaw, "on")
			dropReq = sdk.NewViewDropRowAccessPolicyRequest(oldId)
		}
		if len(newRaw.([]any)) > 0 {
			newId, newColumns := extractPolicyWithColumns(newRaw, "on")
			addReq = sdk.NewViewAddRowAccessPolicyRequest(newId, newColumns)
		}
		req := sdk.NewAlterViewRequest(id)
		if addReq != nil && dropReq != nil { // nolint
			req.WithDropAndAddRowAccessPolicy(*sdk.NewViewDropAndAddRowAccessPolicyRequest(*dropReq, *addReq))
		} else if addReq != nil {
			req.WithAddRowAccessPolicy(*addReq)
		} else if dropReq != nil {
			req.WithDropRowAccessPolicy(*dropReq)
		}
		err := client.Views.Alter(ctx, req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error altering row_access_policy for view %v: %w", d.Id(), err))
		}
	}
	if d.HasChange("aggregation_policy") {
		if v, ok := d.GetOk("aggregation_policy"); ok {
			newId, newColumns := extractPolicyWithColumns(v, "entity_key")
			aggregationPolicyReq := sdk.NewViewSetAggregationPolicyRequest(newId)
			if len(newColumns) > 0 {
				aggregationPolicyReq.WithEntityKey(newColumns)
			}
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetAggregationPolicy(*aggregationPolicyReq.WithForce(true)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting aggregation policy for view %v: %w", d.Id(), err))
			}
		} else {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetAggregationPolicy(*sdk.NewViewUnsetAggregationPolicyRequest()))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting aggregation policy for view %v", d.Id()))
			}
		}
	}

	return ReadView(false)(ctx, d, meta)
}

func DeleteView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	client := meta.(*provider.Context).Client

	err := client.Views.Drop(ctx, sdk.NewDropViewRequest(id).WithIfExists(true))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting view",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
