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
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Retains the access permissions from the original view when a new view is created using the OR REPLACE clause.",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return oldValue != "" && oldValue != newValue
		},
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
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Specifies that the view persists only for the duration of the session that you created it in. A temporary view and all its contents are dropped at the end of the session. In context of this provider, it means that it's dropped after a Terraform operation. This results in a permanent plan with object creation."),
	},
	"is_recursive": {
		Type:             schema.TypeString,
		Optional:         true,
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
	"data_metric_functions": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"function_name": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Identifier of the data metric function to add to the table or view or drop from the table or view. This function identifier must be provided without arguments in parenthesis.",
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
				"on": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "The table or view columns on which to associate the data metric function. The data types of the columns must match the data types of the columns specified in the data metric function definition.",
				},
				// TODO (next pr)
				// "schedule_status": {
				// 	Type:             schema.TypeString,
				// 	Optional:         true,
				// 	ValidateDiagFunc: sdkValidation(sdk.ToAllowedDataMetricScheduleStatusOption),
				// 	Description:      fmt.Sprintf("The status of the metrics association. Valid values are: %v. When status of a data metric function is changed, it is being reassigned with `DROP DATA METRIC FUNCTION` and `ADD DATA METRIC FUNCTION`, and then its status is changed by `MODIFY DATA METRIC FUNCTION` ", possibleValuesListed(sdk.AllAllowedDataMetricScheduleStatusOptions)),
				// 	DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToAllowedDataMetricScheduleStatusOption), func(_, oldValue, newValue string, _ *schema.ResourceData) bool {
				// 		if newValue == "" {
				// 			return true
				// 		}
				// 		return false
				// 	}),
				// },
			},
		},
		Description:  "Data metric functions used for the view.",
		RequiredWith: []string{"data_metric_schedule"},
	},
	"data_metric_schedule": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"minutes": {
					Type:             schema.TypeInt,
					Optional:         true,
					Description:      fmt.Sprintf("Specifies an interval (in minutes) of wait time inserted between runs of the data metric function. Conflicts with `using_cron`. Valid values are: %s. Due to Snowflake limitations, changes in this field is not managed by the provider. Please consider using [taint](https://developer.hashicorp.com/terraform/cli/commands/taint) command, `using_cron` field, or [replace_triggered_by](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#replace_triggered_by) metadata argument.", possibleValuesListedInt(sdk.AllViewDataMetricScheduleMinutes)),
					ValidateDiagFunc: IntInSlice(sdk.AllViewDataMetricScheduleMinutes),
					ConflictsWith:    []string{"data_metric_schedule.using_cron"},
				},
				"using_cron": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "Specifies a cron expression and time zone for periodically running the data metric function. Supports a subset of standard cron utility syntax. Conflicts with `minutes`.",
					ConflictsWith: []string{"data_metric_schedule.minutes"},
				},
			},
		},
		Description:  "Specifies the schedule to run the data metric functions periodically.",
		RequiredWith: []string{"data_metric_functions"},
	},
	// TODO (next pr): add columns
	// "column": {
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
	// 			"comment": {
	// 				Type:        schema.TypeString,
	// 				Optional:    true,
	// 				Description: "Specifies a comment for the column.",
	// 			},
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
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	v, err := client.Views.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("name", id.Name()); err != nil {
		return nil, err
	}
	if err := d.Set("database", id.DatabaseName()); err != nil {
		return nil, err
	}
	if err := d.Set("schema", id.SchemaName()); err != nil {
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

		if orReplace {
			req.WithOrReplace(true)
		}

		if v := d.Get("copy_grants"); v.(bool) {
			req.WithCopyGrants(true).WithOrReplace(true)
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
			id, columns, err := extractPolicyWithColumns(v, "on")
			if err != nil {
				return diag.FromErr(err)
			}
			req.WithRowAccessPolicy(*sdk.NewViewRowAccessPolicyRequest(id, columns))
		}

		if v := d.Get("aggregation_policy"); len(v.([]any)) > 0 {
			id, columns, err := extractPolicyWithColumns(v, "entity_key")
			if err != nil {
				return diag.FromErr(err)
			}
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

		d.SetId(helpers.EncodeResourceIdentifier(id))

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

		if v := d.Get("data_metric_schedule"); len(v.([]any)) > 0 {
			var req *sdk.ViewSetDataMetricScheduleRequest
			dmsConfig := v.([]any)[0].(map[string]any)
			if v, ok := dmsConfig["minutes"]; ok && v.(int) > 0 {
				req = sdk.NewViewSetDataMetricScheduleRequest(fmt.Sprintf("%d MINUTE", v.(int)))
			} else if v, ok := dmsConfig["using_cron"]; ok {
				req = sdk.NewViewSetDataMetricScheduleRequest(fmt.Sprintf("USING CRON %s", v.(string)))
			}
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetDataMetricSchedule(*req))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting data matric schedule in view %v err = %w", id.Name(), err))
			}
		}

		if v, ok := d.GetOk("data_metric_functions"); ok {
			addedRaw, err := extractDataMetricFunctions(v.(*schema.Set).List())
			if err != nil {
				return diag.FromErr(err)
			}
			added := make([]sdk.ViewDataMetricFunction, len(addedRaw))
			for i := range addedRaw {
				added[i] = sdk.ViewDataMetricFunction{
					DataMetricFunction: addedRaw[i].DataMetricFunction,
					On:                 addedRaw[i].On,
				}
			}
			err = client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithAddDataMetricFunction(*sdk.NewViewAddDataMetricFunctionRequest(added)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error adding data matric functions in view %v err = %w", id.Name(), err))
			}
			changeSchedule := make([]sdk.ViewModifyDataMetricFunction, 0, len(addedRaw))
			for i := range addedRaw {
				if addedRaw[i].ScheduleStatus != "" {
					expectedStatus, err := sdk.ToAllowedDataMetricScheduleStatusOption(addedRaw[i].ScheduleStatus)
					if err != nil {
						return diag.FromErr(err)
					}
					var statusCmd sdk.ViewDataMetricScheduleStatusOperationOption
					switch expectedStatus {
					case sdk.DataMetricScheduleStatusStarted:
						statusCmd = sdk.ViewDataMetricScheduleStatusOperationResume
					case sdk.DataMetricScheduleStatusSuspended:
						statusCmd = sdk.ViewDataMetricScheduleStatusOperationSuspend
					default:
						return diag.FromErr(fmt.Errorf("unexpected data metric function status: %v", expectedStatus))
					}
					changeSchedule = append(changeSchedule, sdk.ViewModifyDataMetricFunction{
						DataMetricFunction: addedRaw[i].DataMetricFunction,
						On:                 addedRaw[i].On,
						ViewDataMetricScheduleStatusOperationOption: statusCmd,
					})
				}
			}
			if len(changeSchedule) > 0 {
				err = client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithModifyDataMetricFunction(*sdk.NewViewModifyDataMetricFunctionsRequest(changeSchedule)))
				if err != nil {
					return diag.FromErr(fmt.Errorf("error adding data matric functions in view %v err = %w", id.Name(), err))
				}
			}
		}

		return ReadView(false)(ctx, d, meta)
	}
}

func extractPolicyWithColumns(v any, columnsKey string) (sdk.SchemaObjectIdentifier, []sdk.Column, error) {
	policyConfig := v.([]any)[0].(map[string]any)
	id, err := sdk.ParseSchemaObjectIdentifier(policyConfig["policy_name"].(string))
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, err
	}
	columnsRaw := expandStringList(policyConfig[columnsKey].(*schema.Set).List())
	columns := make([]sdk.Column, len(columnsRaw))
	for i := range columnsRaw {
		columns[i] = sdk.Column{Value: columnsRaw[i]}
	}
	return id, columns, nil
}

func ReadView(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

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
		err = handleDataMetricFunctions(ctx, client, id, d)
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
		case sdk.PolicyKindAggregationPolicy:
			var entityKey []string
			if p.RefArgColumnNames != nil {
				entityKey = sdk.ParseCommaSeparatedStringArray(*p.RefArgColumnNames, true)
			}
			aggregationPolicies = append(aggregationPolicies, map[string]any{
				"policy_name": policyName.FullyQualifiedName(),
				"entity_key":  entityKey,
			})
		case sdk.PolicyKindRowAccessPolicy:
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

func handleDataMetricFunctions(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, d *schema.ResourceData) error {
	dataMetricFunctionReferences, err := client.DataMetricFunctionReferences.GetForEntity(ctx, sdk.NewGetForEntityDataMetricFunctionReferenceRequest(id, sdk.DataMetricFuncionRefEntityDomainView))
	if err != nil {
		return err
	}
	if len(dataMetricFunctionReferences) == 0 {
		return d.Set("data_metric_schedule", nil)
	}
	dataMetricFunctions := make([]map[string]any, len(dataMetricFunctionReferences))
	var schedule string
	for i, dmfRef := range dataMetricFunctionReferences {
		dmfName := sdk.NewSchemaObjectIdentifier(dmfRef.MetricDatabaseName, dmfRef.MetricSchemaName, dmfRef.MetricName)
		var columns []string
		for _, v := range dmfRef.RefArguments {
			columns = append(columns, v.Name)
		}
		// TODO (next pr)
		// var scheduleStatus sdk.DataMetricScheduleStatusOption
		// status, err := sdk.ToDataMetricScheduleStatusOption(dmfRef.ScheduleStatus)
		// if err != nil {
		// 	return err
		// }
		// if slices.Contains(sdk.AllDataMetricScheduleStatusStartedOptions, status) {
		// 	scheduleStatus = sdk.DataMetricScheduleStatusStarted
		// }
		// if slices.Contains(sdk.AllDataMetricScheduleStatusSuspendedOptions, status) {
		// 	scheduleStatus = sdk.DataMetricScheduleStatusSuspended
		// }
		dataMetricFunctions[i] = map[string]any{
			"function_name": dmfName.FullyQualifiedName(),
			"on":            columns,
			// "schedule_status": string(scheduleStatus),
		}
		schedule = dmfRef.Schedule
	}
	if err = d.Set("data_metric_functions", dataMetricFunctions); err != nil {
		return err
	}

	return d.Set("data_metric_schedule", []map[string]any{
		{
			"using_cron": schedule,
		},
	})
}

type ViewDataMetricFunctionDDL struct {
	DataMetricFunction sdk.SchemaObjectIdentifier
	On                 []sdk.Column
	ScheduleStatus     string
}

func extractDataMetricFunctions(v any) (dmfs []ViewDataMetricFunctionDDL, err error) {
	for _, v := range v.([]any) {
		config := v.(map[string]any)
		columnsRaw := expandStringList(config["on"].(*schema.Set).List())
		columns := make([]sdk.Column, len(columnsRaw))
		for i := range columnsRaw {
			columns[i] = sdk.Column{Value: columnsRaw[i]}
		}
		id, err := sdk.ParseSchemaObjectIdentifier(config["function_name"].(string))
		if err != nil {
			return nil, err
		}
		dmfs = append(dmfs, ViewDataMetricFunctionDDL{
			DataMetricFunction: id,
			On:                 columns,
			// TODO (next pr)
			// ScheduleStatus:     config["schedule_status"].(string),
		})
	}
	return
}

func changedKeys(d *schema.ResourceData, keys []string) []string {
	changed := make([]string, 0, len(keys))
	for _, key := range keys {
		if d.HasChange(key) {
			changed = append(changed, key)
		}
	}
	return changed
}

func UpdateView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// change on these fields can not be ForceNew because then view is dropped explicitly and copying grants does not have effect
	if d.HasChange("statement") || d.HasChange("is_temporary") || d.HasChange("is_recursive") || d.HasChange("copy_grant") {
		log.Printf("[DEBUG] Detected change on %q, recreating...", changedKeys(d, []string{"statement", "is_temporary", "is_recursive", "copy_grant"}))
		return CreateView(true)(ctx, d, meta)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming view %v err = %w", d.Id(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
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

	if d.HasChange("data_metric_schedule") {
		if v := d.Get("data_metric_schedule"); len(v.([]any)) > 0 {
			var req *sdk.ViewSetDataMetricScheduleRequest
			dmsConfig := v.([]any)[0].(map[string]any)
			if v := dmsConfig["minutes"]; v.(int) > 0 {
				req = sdk.NewViewSetDataMetricScheduleRequest(fmt.Sprintf("%d MINUTE", v.(int)))
			} else if v, ok := dmsConfig["using_cron"]; ok {
				req = sdk.NewViewSetDataMetricScheduleRequest(fmt.Sprintf("USING CRON %s", v.(string)))
			}
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetDataMetricSchedule(*req))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting data matric schedule in view %v err = %w", id.Name(), err))
			}
		} else {
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetDataMetricSchedule(*sdk.NewViewUnsetDataMetricScheduleRequest()))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting data matric schedule in view %v err = %w", id.Name(), err))
			}
		}
	}

	if d.HasChange("data_metric_functions") {
		old, new := d.GetChange("data_metric_functions")
		removedRaw, addedRaw := old.(*schema.Set).List(), new.(*schema.Set).List()
		added, err := extractDataMetricFunctions(addedRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		removed, err := extractDataMetricFunctions(removedRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(removed) > 0 {
			removed2 := make([]sdk.ViewDataMetricFunction, len(removed))
			for i := range removed {
				removed2[i] = sdk.ViewDataMetricFunction{
					DataMetricFunction: removed[i].DataMetricFunction,
					On:                 removed[i].On,
				}
			}
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithDropDataMetricFunction(*sdk.NewViewDropDataMetricFunctionRequest(removed2)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error adding data matric functions in view %v err = %w", id.Name(), err))
			}
		}
		if len(added) > 0 {
			added2 := make([]sdk.ViewDataMetricFunction, len(added))
			for i := range added {
				added2[i] = sdk.ViewDataMetricFunction{
					DataMetricFunction: added[i].DataMetricFunction,
					On:                 added[i].On,
				}
			}
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithAddDataMetricFunction(*sdk.NewViewAddDataMetricFunctionRequest(added2)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error adding data matric functions in view %v err = %w", id.Name(), err))
			}
		}
	}

	if d.HasChange("row_access_policy") {
		var addReq *sdk.ViewAddRowAccessPolicyRequest
		var dropReq *sdk.ViewDropRowAccessPolicyRequest

		oldRaw, newRaw := d.GetChange("row_access_policy")
		if len(oldRaw.([]any)) > 0 {
			oldId, _, err := extractPolicyWithColumns(oldRaw, "on")
			if err != nil {
				return diag.FromErr(err)
			}
			dropReq = sdk.NewViewDropRowAccessPolicyRequest(oldId)
		}
		if len(newRaw.([]any)) > 0 {
			newId, newColumns, err := extractPolicyWithColumns(newRaw, "on")
			if err != nil {
				return diag.FromErr(err)
			}
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
			newId, newColumns, err := extractPolicyWithColumns(v, "entity_key")
			if err != nil {
				return diag.FromErr(err)
			}
			aggregationPolicyReq := sdk.NewViewSetAggregationPolicyRequest(newId)
			if len(newColumns) > 0 {
				aggregationPolicyReq.WithEntityKey(newColumns)
			}
			err = client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetAggregationPolicy(*aggregationPolicyReq.WithForce(true)))
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
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*provider.Context).Client

	err = client.Views.Drop(ctx, sdk.NewDropViewRequest(id).WithIfExists(true))
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
