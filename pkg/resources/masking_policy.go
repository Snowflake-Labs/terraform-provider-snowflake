package resources

import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
		ForceNew:    true,
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
	"qualified_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the qualified identifier for the masking policy.",
	},
}

// DatabaseName|SchemaName|MaskingPolicyName.

// MaskingPolicy returns a pointer to the resource representing a masking policy.
func MaskingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: CreateMaskingPolicy,
		Read:   ReadMaskingPolicy,
		Update: UpdateMaskingPolicy,
		Delete: DeleteMaskingPolicy,

		Schema: maskingPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateMaskingPolicy implements schema.CreateFunc.
func CreateMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

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
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	ctx := context.Background()
	maskingPolicy, err := client.MaskingPolicies.ShowByID(ctx, objectIdentifier)
	if err != nil {
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

	maskingPolicyDetails, err := client.MaskingPolicies.Describe(ctx, objectIdentifier)
	if err != nil {
		return err
	}

	if err := d.Set("masking_expression", maskingPolicyDetails.Body); err != nil {
		return err
	}

	if err := d.Set("return_data_type", maskingPolicyDetails.ReturnType); err != nil {
		return err
	}

	signature := []map[string]interface{}{}
	for _, s := range maskingPolicyDetails.Signature {
		signature = append(signature, map[string]interface{}{
			"column": []map[string]interface{}{
				{
					"name": s.Name,
					"type": s.Type,
				},
			},
		})
	}
	if err := d.Set("signature", signature); err != nil {
		return err
	}
	if err := d.Set("qualified_name", objectIdentifier.FullyQualifiedName()); err != nil {
		return err
	}

	return err
}

// UpdateMaskingPolicy implements schema.UpdateFunc.
func UpdateMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	ctx := context.Background()

	if d.HasChange("masking_expression") {
		alterOptions := &sdk.AlterMaskingPolicyOptions{}
		_, n := d.GetChange("masking_expression")
		alterOptions.Set = &sdk.MaskingPolicySet{
			Body: sdk.String(n.(string)),
		}
		err := client.MaskingPolicies.Alter(ctx, objectIdentifier, alterOptions)
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
		err := client.MaskingPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return err
		}
	}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		newName := n.(string)
		newID := sdk.NewSchemaObjectIdentifier(objectIdentifier.DatabaseName(), objectIdentifier.SchemaName(), newName)
		alterOptions := &sdk.AlterMaskingPolicyOptions{
			NewName: newID,
		}
		err := client.MaskingPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return err
		}
		d.SetId(helpers.EncodeSnowflakeID(newID))
	}

	return ReadMaskingPolicy(d, meta)
}

// DeleteMaskingPolicy implements schema.DeleteFunc.
func DeleteMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.MaskingPolicies.Drop(ctx, objectIdentifier)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
