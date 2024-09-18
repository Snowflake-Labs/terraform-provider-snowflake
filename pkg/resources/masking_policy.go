package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var maskingPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the masking policy; must be unique for the database and schema in which the masking policy is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the masking policy."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the masking policy."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"argument": {
		Type:     schema.TypeList,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The argument name",
					ForceNew:    true,
				},
				// TODO(SNOW-1596962): Fully support VECTOR data type sdk.ParseFunctionArgumentsFromString could be a base for another function that takes argument names into consideration.
				"type": {
					Type:             schema.TypeString,
					Required:         true,
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToDataType),
					ValidateDiagFunc: sdkValidation(sdk.ToDataType),
					Description:      dataTypeFieldDescription("The argument type. VECTOR data types are not yet supported."),
					ForceNew:         true,
				},
			},
		},
		Required:    true,
		Description: "List of the arguments for the masking policy. The first column and its data type always indicate the column data type values to mask or tokenize in the subsequent policy conditions. Note that you can not specify a virtual column as the first column argument in a conditional masking policy.",
		ForceNew:    true,
	},
	"body": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      diffSuppressStatementFieldDescription("Specifies the SQL expression that transforms the data."),
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"return_data_type": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      dataTypeFieldDescription("The return data type must match the input data type of the first column that is specified as an input column."),
		ForceNew:         true,
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToDataType),
		ValidateDiagFunc: sdkValidation(sdk.ToDataType),
	},
	"exempt_other_policies": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("exempt_other_policies"),
		Description:      booleanStringFieldDescription("Specifies whether the row access policy or conditional masking policy can reference a column that is already protected by a masking policy. Due to Snowflake limitations, when value is chenged, the resource is recreated."),
		ForceNew:         true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the masking policy.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW MASKING POLICY` for the given masking policy.",
		Elem: &schema.Resource{
			Schema: schemas.ShowMaskingPolicySchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE MASKING POLICY` for the given masking policy.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeMaskingPolicySchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// MaskingPolicy returns a pointer to the resource representing a masking policy.
func MaskingPolicy() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateMaskingPolicy,
		ReadContext:   ReadMaskingPolicy(true),
		UpdateContext: UpdateMaskingPolicy,
		DeleteContext: DeleteMaskingPolicy,
		Description:   "Resource used to manage masking policies. For more information, check [masking policies documentation](https://docs.snowflake.com/en/sql-reference/sql/create-masking-policy).",

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(maskingPolicySchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(maskingPolicySchema, DescribeOutputAttributeName, "name", "body"),
			ComputedIfAnyAttributeChanged(maskingPolicySchema, FullyQualifiedNameAttributeName, "name"),
		),

		Schema: maskingPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportMaskingPolicy,
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v0_95_0_MaskingPolicyStateUpgrader,
			},
		},
	}
}

func ImportMaskingPolicy(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Starting masking policy import")
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	policy, err := client.MaskingPolicies.ShowByID(ctx, id)
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
	if err := d.Set("exempt_other_policies", booleanStringFromBool(policy.ExemptOtherPolicies)); err != nil {
		return nil, err
	}
	policyDescription, err := client.MaskingPolicies.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("body", policyDescription.Body); err != nil {
		return nil, err
	}
	if err := d.Set("argument", schemas.MaskingPolicyArgumentsToSchema(policyDescription.Signature)); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateMaskingPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	expression := d.Get("body").(string)
	returnDataType := d.Get("return_data_type").(string)

	arguments := d.Get("argument").([]any)
	args := make([]sdk.TableColumnSignature, 0)
	for _, arg := range arguments {
		v := arg.(map[string]any)
		dataType, err := sdk.ToDataType(v["type"].(string))
		if err != nil {
			return diag.FromErr(err)
		}
		args = append(args, sdk.TableColumnSignature{
			Name: v["name"].(string),
			Type: dataType,
		})
	}

	returns, err := sdk.ToDataType(returnDataType)
	if err != nil {
		return diag.FromErr(err)
	}

	// set optionals
	opts := &sdk.CreateMaskingPolicyOptions{}
	if comment, ok := d.Get("comment").(string); ok {
		opts.Comment = sdk.String(comment)
	}
	if v := d.Get("exempt_other_policies").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}

		opts.ExemptOtherPolicies = sdk.Pointer(parsed)
	}

	err = client.MaskingPolicies.Create(ctx, id, args, returns, expression, opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadMaskingPolicy(false)(ctx, d, meta)
}

func ReadMaskingPolicy(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		maskingPolicy, err := client.MaskingPolicies.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query masking policy. Marking the resource as removed.",
						Detail:   fmt.Sprintf("masking policy name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}
		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("comment", maskingPolicy.Comment); err != nil {
			return diag.FromErr(err)
		}

		maskingPolicyDescription, err := client.MaskingPolicies.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("body", maskingPolicyDescription.Body); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("return_data_type", maskingPolicyDescription.ReturnType); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("argument", schemas.MaskingPolicyArgumentsToSchema(maskingPolicyDescription.Signature)); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"exempt_other_policies", "exempt_other_policies", maskingPolicy.ExemptOtherPolicies, booleanStringFromBool(maskingPolicy.ExemptOtherPolicies), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, maskingPolicySchema, []string{
			"exempt_other_policies",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.MaskingPolicyToSchema(maskingPolicy)}); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.MaskingPolicyDescriptionToSchema(*maskingPolicyDescription)}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func UpdateMaskingPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("name") {
		newID := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.MaskingPolicies.Alter(ctx, id, &sdk.AlterMaskingPolicyOptions{
			NewName: &newID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newID))
		id = newID
	}

	if d.HasChange("body") {
		alterOptions := &sdk.AlterMaskingPolicyOptions{
			Set: &sdk.MaskingPolicySet{
				Body: sdk.Pointer(d.Get("body").(string)),
			},
		}
		err := client.MaskingPolicies.Alter(ctx, id, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("comment") {
		alterOptions := &sdk.AlterMaskingPolicyOptions{}
		if v, ok := d.GetOk("comment"); ok {
			alterOptions.Set = &sdk.MaskingPolicySet{
				Comment: sdk.String(v.(string)),
			}
		} else {
			alterOptions.Unset = &sdk.MaskingPolicyUnset{
				Comment: sdk.Bool(true),
			}
		}
		err := client.MaskingPolicies.Alter(ctx, id, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// exempt_other_policies is handled by ForceNew

	return ReadMaskingPolicy(false)(ctx, d, meta)
}

func DeleteMaskingPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.MaskingPolicies.Drop(ctx, id, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Pointer(true)})
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting masking policy",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
