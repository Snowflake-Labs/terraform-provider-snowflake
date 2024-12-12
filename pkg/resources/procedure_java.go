package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProcedureJava() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ProcedureJava, CreateContextProcedureJava),
		ReadContext:   TrackingReadWrapper(resources.ProcedureJava, ReadContextProcedureJava),
		UpdateContext: TrackingUpdateWrapper(resources.ProcedureJava, UpdateContextProcedureJava),
		DeleteContext: TrackingDeleteWrapper(resources.ProcedureJava, DeleteProcedure),
		Description:   "Resource used to manage java procedure objects. For more information, check [procedure documentation](https://docs.snowflake.com/en/sql-reference/sql/create-procedure).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ProcedureJava, customdiff.All(
			// TODO[SNOW-1348103]: ComputedIfAnyAttributeChanged(javaProcedureSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javaProcedureSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(procedureParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllProcedureParameters), strings.ToLower)...),
			procedureParametersCustomDiff,
			// The language check is more for the future.
			// Currently, almost all attributes are marked as forceNew.
			// When language changes, these attributes also change, causing the object to recreate either way.
			// The only option is java staged <-> scala staged (however scala need runtime_version which may interfere).
			RecreateWhenResourceStringFieldChangedExternally("procedure_language", "JAVA"),
		)),

		Schema: collections.MergeMaps(javaProcedureSchema, procedureParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextProcedureJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	database := d.Get("database").(string)
	sc := d.Get("schema").(string)
	name := d.Get("name").(string)

	argumentRequests, err := parseProcedureArgumentsCommon(d)
	if err != nil {
		return diag.FromErr(err)
	}
	returns, err := parseProcedureReturnsCommon(d)
	if err != nil {
		return diag.FromErr(err)
	}
	handler := d.Get("handler").(string)
	runtimeVersion := d.Get("runtime_version").(string)
	// TODO [this PR]: handle real packages
	packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0")}

	argumentDataTypes := collections.Map(argumentRequests, func(r sdk.ProcedureArgumentRequest) datatypes.DataType { return r.ArgDataType })
	id := sdk.NewSchemaObjectIdentifierWithArgumentsNormalized(database, sc, name, argumentDataTypes...)
	request := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, runtimeVersion, packages, handler).
		WithArguments(argumentRequests)

	errs := errors.Join(
		booleanStringAttributeCreateBuilder(d, "is_secure", request.WithSecure),
		attributeMappedValueCreateBuilder[string](d, "null_input_behavior", request.WithNullInputBehavior, sdk.ToNullInputBehavior),
		// TODO [SNOW-1348103]: handle the rest of the attributes
		// comment
		setProcedureImportsInBuilder(d, request.WithImports),
		// packages
		// external_access_integrations
		// secrets
		setProcedureTargetPathInBuilder(d, request.WithTargetPath),
		stringAttributeCreateBuilder(d, "procedure_definition", request.WithProcedureDefinitionWrapped),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.Procedures.CreateForJava(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// parameters do not work in create procedure (query does not fail but parameters stay unchanged)
	setRequest := sdk.NewProcedureSetRequest()
	if parametersCreateDiags := handleProcedureParametersCreate(d, setRequest); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}
	if !reflect.DeepEqual(*setRequest, *sdk.NewProcedureSetRequest()) {
		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithSet(*setRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextProcedureJava(ctx, d, meta)
}

func ReadContextProcedureJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	allProcedureDetails, diags := queryAllProcedureDetailsCommon(ctx, d, client, id)
	if diags != nil {
		return diags
	}

	// TODO [SNOW-1348103]: handle external changes marking
	// TODO [SNOW-1348103]: handle setting state to value from config

	errs := errors.Join(
		// TODO [SNOW-1348103]: set the rest of the fields
		// not reading is_secure on purpose (handled as external change to show output)
		readFunctionOrProcedureArguments(d, allProcedureDetails.procedureDetails.NormalizedArguments),
		d.Set("return_type", allProcedureDetails.procedureDetails.ReturnDataType.ToSql()),
		// not reading null_input_behavior on purpose (handled as external change to show output)
		setRequiredFromStringPtr(d, "runtime_version", allProcedureDetails.procedureDetails.RuntimeVersion),
		// comment
		readFunctionOrProcedureImports(d, allProcedureDetails.procedureDetails.NormalizedImports),
		// packages
		setRequiredFromStringPtr(d, "handler", allProcedureDetails.procedureDetails.Handler),
		// external_access_integrations
		// secrets
		readFunctionOrProcedureTargetPath(d, allProcedureDetails.procedureDetails.NormalizedTargetPath),
		setOptionalFromStringPtr(d, "procedure_definition", allProcedureDetails.procedureDetails.Body),
		d.Set("procedure_language", allProcedureDetails.procedureDetails.Language),

		handleProcedureParameterRead(d, allProcedureDetails.procedureParameters),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ProcedureToSchema(allProcedureDetails.procedure)}),
		d.Set(ParametersAttributeName, []map[string]any{schemas.ProcedureParametersToSchema(allProcedureDetails.procedureParameters)}),
	)
	if errs != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextProcedureJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(id.SchemaId(), d.Get("name").(string), id.ArgumentDataTypes()...)

		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithRenameTo(newId.SchemaObjectId()))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming procedure %v err = %w", d.Id(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// Batch SET operations and UNSET operations
	setRequest := sdk.NewProcedureSetRequest()
	unsetRequest := sdk.NewProcedureUnsetRequest()

	// TODO [SNOW-1348103]: handle all updates
	// secure
	// external access integration
	// secrets
	// comment

	if updateParamDiags := handleProcedureParametersUpdate(d, setRequest, unsetRequest); len(updateParamDiags) > 0 {
		return updateParamDiags
	}

	// Apply SET and UNSET changes
	if !reflect.DeepEqual(*setRequest, *sdk.NewProcedureSetRequest()) {
		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithSet(*setRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unsetRequest, *sdk.NewProcedureUnsetRequest()) {
		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithUnset(*unsetRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextProcedureJava(ctx, d, meta)
}
