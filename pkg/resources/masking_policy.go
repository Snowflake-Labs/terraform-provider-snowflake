package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var maskingPolicySchema = map[string]*schema.Schema{
	"or_replace": {
		Type:                  schema.TypeBool,
		Optional:              true,
		Default:               false,
		Description:           "Whether to override a previous masking policy with the same name.",
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return old != new
		},
	},
	"if_not_exists": {
		Type:                  schema.TypeBool,
		Optional:              true,
		Default:               false,
		Description:           "Prevent overwriting a previous masking policy with the same name.",
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return old != new
		},
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the masking policy; must be unique for the database and schema in which the masking policy is created.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the masking policy.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the masking policy.",
		ForceNew:    true,
	},
	"signature": {
		Type:        schema.TypeList,
		Required:    true,
		Description: "The signature for the masking policy; specifies the input columns and data types to evaluate at query runtime.",
		MinItems:    1,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"column": {
					Type:     schema.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the column name to mask.",
							},
							"type": {
								Type:             schema.TypeString,
								Required:         true,
								Description:      "Specifies the column type to mask.",
								ForceNew:         true,
								ValidateFunc:     dataTypeValidateFunc,
								DiffSuppressFunc: dataTypeDiffSuppressFunc,
							},
						},
					},
				},
			},
		},
	},
	"masking_expression": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the SQL expression that transforms the data.",
		DiffSuppressFunc: ignoreTrimSpaceSuppressFunc,
	},
	"return_data_type": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the data type to return.",
		ForceNew:         true,
		ValidateFunc:     dataTypeValidateFunc,
		DiffSuppressFunc: dataTypeDiffSuppressFunc,
	},
	"exempt_other_policies": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the row access policy or conditional masking policy can reference a column that is already protected by a masking policy.",
		Default:     false,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the masking policy.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// DatabaseName|SchemaName|MaskingPolicyName.

// MaskingPolicy returns a pointer to the resource representing a masking policy.
func MaskingPolicy() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        CreateMaskingPolicy,
		Read:          ReadMaskingPolicy,
		Update:        UpdateMaskingPolicy,
		Delete:        DeleteMaskingPolicy,

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(FullyQualifiedNameAttributeName, "name"),
		),

		Schema: maskingPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v0_94_1_MaskingPolicyStateUpgrader,
			},
		},
	}
}

// CreateMaskingPolicy implements schema.CreateFunc.
func CreateMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	expression := d.Get("masking_expression").(string)
	returnDataType := d.Get("return_data_type").(string)

	ctx := context.Background()
	objectIdentifier := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	signatureList := d.Get("signature").([]interface{})
	signature := []sdk.TableColumnSignature{}
	for _, s := range signatureList {
		m := s.(map[string]interface{})
		columns := m["column"].([]interface{})
		for _, c := range columns {
			cm := c.(map[string]interface{})
			dt, err := sdk.ToDataType(cm["type"].(string))
			if err != nil {
				return err
			}
			signature = append(signature, sdk.TableColumnSignature{
				Name: cm["name"].(string),
				Type: dt,
			})
		}
	}

	returns, err := sdk.ToDataType(returnDataType)
	if err != nil {
		return err
	}
	opts := &sdk.CreateMaskingPolicyOptions{}
	if comment, ok := d.Get("comment").(string); ok {
		opts.Comment = sdk.String(comment)
	}
	if exemptOtherPolicies := d.Get("exempt_other_policies").(bool); exemptOtherPolicies {
		opts.ExemptOtherPolicies = sdk.Bool(exemptOtherPolicies)
	}

	err = client.MaskingPolicies.Create(ctx, objectIdentifier, signature, returns, expression, opts)
	if err != nil {
		return err
	}
	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))

	return ReadMaskingPolicy(d, meta)
}

// ReadMaskingPolicy implements schema.ReadFunc.
func ReadMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	ctx := context.Background()
	maskingPolicy, err := client.MaskingPolicies.ShowByID(ctx, id)
	if err != nil {
		return err
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return err
	}

	if err := d.Set("name", maskingPolicy.Name); err != nil {
		return err
	}

	if err := d.Set("database", maskingPolicy.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("schema", maskingPolicy.SchemaName); err != nil {
		return err
	}

	if err := d.Set("exempt_other_policies", maskingPolicy.ExemptOtherPolicies); err != nil {
		return err
	}

	if err := d.Set("comment", maskingPolicy.Comment); err != nil {
		return err
	}

	maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, id)
	if err != nil {
		return err
	}

	if err := d.Set("masking_expression", maskingPolicyDetails.Body); err != nil {
		return err
	}

	if err := d.Set("return_data_type", maskingPolicyDetails.ReturnType); err != nil {
		return err
	}

	columns := []map[string]interface{}{}
	for _, s := range maskingPolicyDetails.Signature {
		columns = append(columns, map[string]interface{}{
			"name": s.Name,
			"type": s.Type,
		})
	}
	signature := []map[string]interface{}{
		{"column": columns},
	}
	if err := d.Set("signature", signature); err != nil {
		return err
	}

	return err
}

// UpdateMaskingPolicy implements schema.UpdateFunc.
func UpdateMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	ctx := context.Background()

	if d.HasChange("name") {
		newID := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.MaskingPolicies.Alter(ctx, id, &sdk.AlterMaskingPolicyOptions{
			NewName: &newID,
		})
		if err != nil {
			return err
		}

		d.SetId(helpers.EncodeSnowflakeID(newID))
		id = newID
	}

	if d.HasChange("masking_expression") {
		alterOptions := &sdk.AlterMaskingPolicyOptions{}
		_, n := d.GetChange("masking_expression")
		alterOptions.Set = &sdk.MaskingPolicySet{
			Body: sdk.String(n.(string)),
		}
		err := client.MaskingPolicies.Alter(ctx, id, alterOptions)
		if err != nil {
			return err
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
			return err
		}
	}

	return ReadMaskingPolicy(d, meta)
}

// DeleteMaskingPolicy implements schema.DeleteFunc.
func DeleteMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.MaskingPolicies.Drop(ctx, objectIdentifier, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Bool(true)})
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
