package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var rowAccessPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the row access policy; must be unique for the database and schema in which the row access policy is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the row access policy.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the row access policy.",
		ForceNew:    true,
	},
	// TODO [SNOW-1020074]: Implement DiffSuppressFunc and test after https://github.com/hashicorp/terraform-plugin-sdk/issues/477 is solved.
	"signature": {
		Type:        schema.TypeMap,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		ForceNew:    true,
		Description: "Specifies signature (arguments) for the row access policy (uppercase and sorted to avoid recreation of resource). A signature specifies a set of attributes that must be considered to determine whether the row is accessible. The attribute values come from the database object (e.g. table or view) to be protected by the row access policy.",
	},
	"row_access_expression": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the SQL expression. The expression can be any boolean-valued SQL expression.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the row access policy.",
	},
}

// RowAccessPolicy returns a pointer to the resource representing a row access policy.
func RowAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: CreateRowAccessPolicy,
		Read:   ReadRowAccessPolicy,
		Update: UpdateRowAccessPolicy,
		Delete: DeleteRowAccessPolicy,

		Schema: rowAccessPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateRowAccessPolicy implements schema.CreateFunc.
func CreateRowAccessPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	signature := d.Get("signature").(map[string]any)
	rowAccessExpression := d.Get("row_access_expression").(string)

	args := make([]sdk.CreateRowAccessPolicyArgsRequest, 0)
	for k, v := range signature {
		dataType := sdk.DataType(v.(string))
		args = append(args, *sdk.NewCreateRowAccessPolicyArgsRequest(k, dataType))
	}

	createRequest := sdk.NewCreateRowAccessPolicyRequest(id, args, rowAccessExpression)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	err := client.RowAccessPolicies.Create(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("error creating row access policy %v err = %w", name, err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadRowAccessPolicy(d, meta)
}

// ReadRowAccessPolicy implements schema.ReadFunc.
func ReadRowAccessPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	rowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] row access policy (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("name", rowAccessPolicy.Name); err != nil {
		return err
	}

	if err := d.Set("database", rowAccessPolicy.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("schema", rowAccessPolicy.SchemaName); err != nil {
		return err
	}

	if err := d.Set("comment", rowAccessPolicy.Comment); err != nil {
		return err
	}

	rowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, id)
	if err != nil {
		return err
	}

	if err := d.Set("row_access_expression", rowAccessPolicyDescription.Body); err != nil {
		return err
	}

	if err := d.Set("signature", parseSignature(rowAccessPolicyDescription.Signature)); err != nil {
		return err
	}

	return err
}

// UpdateRowAccessPolicy implements schema.UpdateFunc.
func UpdateRowAccessPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("comment") {
		comment := d.Get("comment")
		if c := comment.(string); c == "" {
			err := client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithUnsetComment(sdk.Bool(true)))
			if err != nil {
				return fmt.Errorf("error unsetting comment for row access policy on %v err = %w", d.Id(), err)
			}
		} else {
			err := client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithSetComment(sdk.String(c)))
			if err != nil {
				return fmt.Errorf("error updating comment for row access policy on %v err = %w", d.Id(), err)
			}
		}
	}

	if d.HasChange("row_access_expression") {
		rowAccessExpression := d.Get("row_access_expression").(string)
		err := client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.String(rowAccessExpression)))
		if err != nil {
			return fmt.Errorf("error updating row access policy expression on %v err = %w", d.Id(), err)
		}
	}

	return ReadRowAccessPolicy(d, meta)
}

// DeleteRowAccessPolicy implements schema.DeleteFunc.
func DeleteRowAccessPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

// TODO [SNOW-1020074]: should we put signature parsing to the SDK?
func parseSignature(signature string) map[string]interface{} {
	// Format in database is `(column <data_type>)`
	plainSignature := strings.ReplaceAll(signature, "(", "")
	plainSignature = strings.ReplaceAll(plainSignature, ")", "")
	signatureParts := strings.Split(plainSignature, ", ")
	signatureMap := map[string]interface{}{}

	for _, e := range signatureParts {
		parts := strings.Split(e, " ")
		signatureMap[parts[0]] = parts[1]
	}

	return signatureMap
}
