package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FunctionJava() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.FunctionJava, CreateContextFunctionJava),
		ReadContext:   TrackingReadWrapper(resources.FunctionJava, ReadContextFunctionJava),
		UpdateContext: TrackingUpdateWrapper(resources.FunctionJava, UpdateContextFunctionJava),
		DeleteContext: TrackingDeleteWrapper(resources.FunctionJava, DeleteFunction),
		Description:   "Resource used to manage java function objects. For more information, check [function documentation](https://docs.snowflake.com/en/sql-reference/sql/create-function).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionJava, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(javaFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javaFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// TODO[SNOW-1348103]: recreate when type changed externally
		)),

		Schema: collections.MergeMaps(javaFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	database := d.Get("database").(string)
	sc := d.Get("schema").(string)
	name := d.Get("name").(string)

	argumentRequests, diags := parseFunctionArgumentsCommon(d)
	if diags != nil {
		return diags
	}
	returns, diags := parseFunctionReturnsCommon(d)
	if diags != nil {
		return diags
	}
	handler := d.Get("handler").(string)

	argumentDataTypes := collections.Map(argumentRequests, func(r sdk.FunctionArgumentRequest) datatypes.DataType { return r.ArgDataType })
	id := sdk.NewSchemaObjectIdentifierWithArgumentsNormalized(database, sc, name, argumentDataTypes...)
	request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
		WithArguments(argumentRequests)

	if v, ok := d.GetOk("statement"); ok {
		request.WithFunctionDefinitionWrapped(v.(string))
	}

	if err := client.Functions.CreateForJava(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	// TODO [this PR]: handle parameters

	d.SetId(id.FullyQualifiedName())
	return ReadContextFunctionJava(ctx, d, meta)

	// Set optionals
	//if v, ok := d.GetOk("is_secure"); ok {
	//	request.WithSecure(v.(bool))
	//}
	//if v, ok := d.GetOk("null_input_behavior"); ok {
	//	request.WithNullInputBehavior(sdk.NullInputBehavior(v.(string)))
	//}
	//if v, ok := d.GetOk("return_behavior"); ok {
	//	request.WithReturnResultsBehavior(sdk.ReturnResultsBehavior(v.(string)))
	//}
	//if v, ok := d.GetOk("runtime_version"); ok {
	//	request.WithRuntimeVersion(v.(string))
	//}
	//if v, ok := d.GetOk("comment"); ok {
	//	request.WithComment(v.(string))
	//}
	//if _, ok := d.GetOk("imports"); ok {
	//	var imports []sdk.FunctionImportRequest
	//	for _, item := range d.Get("imports").([]interface{}) {
	//		imports = append(imports, *sdk.NewFunctionImportRequest().WithImport(item.(string)))
	//	}
	//	request.WithImports(imports)
	//}
	//if _, ok := d.GetOk("packages"); ok {
	//	var packages []sdk.FunctionPackageRequest
	//	for _, item := range d.Get("packages").([]interface{}) {
	//		packages = append(packages, *sdk.NewFunctionPackageRequest().WithPackage(item.(string)))
	//	}
	//	request.WithPackages(packages)
	//}
	//if v, ok := d.GetOk("target_path"); ok {
	//	request.WithTargetPath(v.(string))
	//}
}

func ReadContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func UpdateContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
