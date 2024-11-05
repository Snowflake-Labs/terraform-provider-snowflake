# Better tests poc

<!-- TOC -->
* [Better tests poc](#better-tests-poc)
  * [How it works](#how-it-works)
    * [Adding new resource assertions](#adding-new-resource-assertions)
    * [Adding new resource show output assertions](#adding-new-resource-show-output-assertions)
    * [Adding new resource parameters assertions](#adding-new-resource-parameters-assertions)
    * [Adding new Snowflake object assertions](#adding-new-snowflake-object-assertions)
    * [Adding new Snowflake object parameters assertions](#adding-new-snowflake-object-parameters-assertions)
    * [Adding new resource config model builders](#adding-new-resource-config-model-builders)
    * [Running the generators](#running-the-generators)
    * [Example usage in practice](#example-usage-in-practice)
  * [Known limitations/planned improvements](#known-limitationsplanned-improvements)
<!-- TOC -->

This package contains a quick implementation of helpers that should allow us a quicker, more pleasant, and more readable implementation of tests, mainly the acceptance ones.
It contains the following packages:
- `assert` - all the assertions reside here. Also, the utilities to build assertions for new objects. All the current assertions are generated. The currently supported assertions are:
  - Snowflake object assertions (generated in subpackage `objectassert`)
  - Snowflake object parameters assertions (generated in subpackage `objectparametersassert`)
  - resource assertions (generated in subpackage `resourceassert`)
  - resource parameters assertions (generated in subpackage `resourceparametersassert`)
  - show output assertions (generated in subpackage `resourceshowoutputassert`)

- `config` - the new `ResourceModel` abstraction resides here. It provides models for objects and the builder methods allowing better config preparation in the acceptance tests.
It aims to be more readable than using `Config:` with hardcoded string or `ConfigFile:` for file that is not directly reachable from the test body. Also, it should be easier to reuse the models and prepare convenience extension methods. The models are already generated.

## How it works
### Adding new resource assertions
Resource assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allStructs` slice in the `assert/resourceassert/gen/resource_schema_def.go`
- to add custom (not generated assertions) create file `abc_resource_ext.go` in the `assert/resourceassert` package. Example would be:
```go
func (w *WarehouseResourceAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceAssert {
    w.AddAssertion(assert.ValueSet("max_concurrency_level", "8"))
    return w
}
```

### Adding new resource show output assertions
Resource show output assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allResourceSchemaDefs` slice in the `assert/objectassert/gen/sdk_object_def.go`
- to add custom (not generated assertions) create file `abc_show_output_ext.go` in the `assert/resourceshowoutputassert` package. Example would be:
```go
func (u *UserShowOutputAssert) HasNameAndLoginName(expected string) *UserShowOutputAssert {
	return u.
		HasName(expected).
		HasLoginName(expected)
}
```

### Adding new resource parameters assertions
Resource parameters assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allObjectsParameters` slice in the `assert/objectparametersassert/gen/object_parameters_def.go`
- to add custom (not generated assertions) create file `warehouse_resource_parameters_ext.go` in the `assert/resourceparametersassert` package. Example would be:
```go
func (w *WarehouseResourceParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceParametersAssert {
	return w.
		HasMaxConcurrencyLevel(8).
		HasMaxConcurrencyLevelLevel("")
}
```

### Adding new Snowflake object assertions
Snowflake object assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allStructs` slice in the `assert/objectassert/gen/sdk_object_def.go`
- to add custom (not generated assertions) create file `abc_snowflake_ext.go` in the `objectassert` package. Example would be:
```go
func (w *WarehouseAssert) HasStateOneOf(expected ...sdk.WarehouseState) *WarehouseAssert {
    w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
        t.Helper()
        if !slices.Contains(expected, o.State) {
            return fmt.Errorf("expected state one of: %v; got: %v", expected, string(o.State))
        }
        return nil
    })
    return w
}
```

### Adding new Snowflake object parameters assertions
Snowflake object parameters assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allObjectsParameters` slice in the `assert/objectparametersassert/gen/main/main.go`
- make sure that test helper method `acc.TestClient().Parameter.ShowAbcParameters` exists in `/pkg/acceptance/helpers/parameter_client.go`
- to add custom (not generated) assertions create file `abc_parameters_snowflake_ext.go` in the `objectparametersassert` package. Example would be:
```go
func (w *WarehouseParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseParametersAssert {
    return w.
        HasMaxConcurrencyLevel(8).
        HasMaxConcurrencyLevelLevel("")
}
```

### Adding new resource config model builders
Resource config model builders can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allResourceSchemaDefs` slice in the `assert/resourceassert/gen/resource_schema_def.go`
- to add custom (not generated) config builder methods create file `warehouse_model_ext` in the `config/model` package. Example would be:
```go
func BasicWarehouseModel(
	name string,
	comment string,
) *WarehouseModel {
	return WarehouseWithDefaultMeta(name).WithComment(comment)
}

func (w *WarehouseModel) WithWarehouseSizeEnum(warehouseSize sdk.WarehouseSize) *WarehouseModel {
	return w.WithWarehouseSize(string(warehouseSize))
}
```

### Running the generators
Each of the above assertion types/config models has its own generator and cleanup entry in our Makefile.
You can generate config models with:
```shell
  make clean-resource-model-builders generate-resource-model-builders
```

You can use cli flags:
```shell
  make clean-resource-model-builders generate-resource-model-builders SF_TF_GENERATOR_ARGS='--dry-run --verbose'
```

To clean/generate all from this package run
```shell
  make clean-all-assertions-and-config-models generate-all-assertions-and-config-models
```

### Example usage in practice
You can check the current example usage in `TestAcc_Warehouse_BasicFlows` and the `create: complete` inside `TestInt_Warehouses`. To see the output after invalid assertions:
- add the following to the first step of `TestAcc_Warehouse_BasicFlows`
```go
    // bad checks below
    resourceassert.WarehouseResource(t, "snowflake_warehouse.w").
        HasWarehouseTypeString(string(sdk.WarehouseTypeSnowparkOptimized)).
        HasWarehouseSizeString(string(sdk.WarehouseSizeMedium)),
    resourceshowoutputassert.WarehouseShowOutput(t, "snowflake_warehouse.w").
        HasType(sdk.WarehouseTypeSnowparkOptimized),
    resourceparametersassert.WarehouseResourceParameters(t, "snowflake_warehouse.w").
        HasMaxConcurrencyLevel(16).
        HasMaxConcurrencyLevelLevel(sdk.ParameterTypeWarehouse),
    objectassert.Warehouse(t, warehouseId).
        HasName("bad name").
        HasState(sdk.WarehouseStateSuspended).
        HasType(sdk.WarehouseTypeSnowparkOptimized).
        HasSize(sdk.WarehouseSizeMedium).
        HasMaxClusterCount(12).
        HasMinClusterCount(13).
        HasScalingPolicy(sdk.ScalingPolicyEconomy).
        HasAutoSuspend(123).
        HasAutoResume(false).
        HasResourceMonitor(sdk.NewAccountObjectIdentifier("some-id")).
        HasComment("bad comment").
        HasEnableQueryAcceleration(true).
        HasQueryAccelerationMaxScaleFactor(12),
    objectparametersassert.WarehouseParameters(t, warehouseId).
        HasMaxConcurrencyLevel(16),
    assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", string(sdk.WarehouseTypeSnowparkOptimized))),
```
it will result in:
```
    warehouse_acceptance_test.go:51: Step 1/8 error: Check failed: check 7/12 error:
        snowflake_warehouse.w resource assertion [1/2]: failed with error: Attribute 'warehouse_type' not found
        snowflake_warehouse.w resource assertion [2/2]: failed with error: Attribute 'warehouse_size' not found
        check 8/12 error:
        snowflake_warehouse.w show_output assertion [2/2]: failed with error: Attribute 'show_output.0.type' expected "SNOWPARK-OPTIMIZED", got "STANDARD"
        check 9/12 error:
        snowflake_warehouse.w parameters assertion [2/3]: failed with error: Attribute 'parameters.0.max_concurrency_level.0.value' expected "16", got "8"
        snowflake_warehouse.w parameters assertion [3/3]: failed with error: Attribute 'parameters.0.max_concurrency_level.0.level' expected "WAREHOUSE", got ""
        check 10/12 error:
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [1/13]: failed with error: expected name: bad name; got: XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [2/13]: failed with error: expected state: SUSPENDED; got: STARTED
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [3/13]: failed with error: expected type: SNOWPARK-OPTIMIZED; got: STANDARD
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [4/13]: failed with error: expected size: MEDIUM; got: XSMALL
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [5/13]: failed with error: expected max cluster count: 12; got: 1
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [6/13]: failed with error: expected min cluster count: 13; got: 1
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [7/13]: failed with error: expected scaling policy: ECONOMY; got: STANDARD
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [8/13]: failed with error: expected auto suspend: 123; got: 600
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [9/13]: failed with error: expected auto resume: false; got: true
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [10/13]: failed with error: expected resource monitor: some-id; got: 
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [11/13]: failed with error: expected comment: bad comment; got: Who does encouraging eagerly annoying dream several their scold straightaway.
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [12/13]: failed with error: expected enable query acceleration: true; got: false
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [13/13]: failed with error: expected query acceleration max scale factor: 12; got: 8
        check 11/12 error:
        parameter assertion for WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"][MAX_CONCURRENCY_LEVEL][1/1] failed: expected value 16, got 8
        check 12/12 error:
        snowflake_warehouse.w: Attribute 'warehouse_type' not found
```

- add the following to the second step of `TestAcc_Warehouse_BasicFlows`
```go
    // bad checks below
    assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(warehouseId.Name(), "bad name", name)),
    resourceassert.ImportedWarehouseResource(t, warehouseId.Name()).
        HasNameString("bad name").
        HasWarehouseTypeString(string(sdk.WarehouseTypeSnowparkOptimized)).
        HasWarehouseSizeString(string(sdk.WarehouseSizeMedium)).
        HasMaxClusterCountString("2").
        HasMinClusterCountString("3").
        HasScalingPolicyString(string(sdk.ScalingPolicyEconomy)).
        HasAutoSuspendString("123").
        HasAutoResumeString("false").
        HasResourceMonitorString("abc").
        HasCommentString("bad comment").
        HasEnableQueryAccelerationString("true").
        HasQueryAccelerationMaxScaleFactorString("16"),
    resourceparametersassert.ImportedWarehouseResourceParameters(t, warehouseId.Name()).
        HasMaxConcurrencyLevel(1).
        HasMaxConcurrencyLevelLevel(sdk.ParameterTypeWarehouse).
        HasStatementQueuedTimeoutInSeconds(23).
        HasStatementQueuedTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse).
        HasStatementTimeoutInSeconds(1232).
        HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
    objectassert.Warehouse(t, warehouseId).
        HasName("bad name").
        HasState(sdk.WarehouseStateSuspended).
        HasType(sdk.WarehouseTypeSnowparkOptimized).
        HasSize(sdk.WarehouseSizeMedium).
        HasMaxClusterCount(12).
        HasMinClusterCount(13).
        HasScalingPolicy(sdk.ScalingPolicyEconomy).
        HasAutoSuspend(123).
        HasAutoResume(false).
        HasResourceMonitor(sdk.NewAccountObjectIdentifier("some-id")).
        HasComment("bad comment").
        HasEnableQueryAcceleration(true).
        HasQueryAccelerationMaxScaleFactor(12),
    objectparametersassert.WarehouseParameters(t, warehouseId).
        HasMaxConcurrencyLevel(1),
```
it will result in:
```
    warehouse_acceptance_test.go:51: check 7/11 error:
        attribute bad name not found in instance state
        check 8/11 error:
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [1/12]: failed with error: expected: bad name, got: WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [2/12]: failed with error: expected: SNOWPARK-OPTIMIZED, got: STANDARD
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [3/12]: failed with error: expected: MEDIUM, got: XSMALL
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [4/12]: failed with error: expected: 2, got: 1
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [5/12]: failed with error: expected: 3, got: 1
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [6/12]: failed with error: expected: ECONOMY, got: STANDARD
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [7/12]: failed with error: expected: 123, got: 600
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [8/12]: failed with error: expected: false, got: true
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [9/12]: failed with error: expected: abc, got: 
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [10/12]: failed with error: expected: bad comment, got: Promise my huh off certain you bravery dynasty with Roman.
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [11/12]: failed with error: expected: true, got: false
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [12/12]: failed with error: expected: 16, got: 8
        check 9/11 error:
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [2/7]: failed with error: expected: 1, got: 8
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [3/7]: failed with error: expected: WAREHOUSE, got: 
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [4/7]: failed with error: expected: 23, got: 0
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [5/7]: failed with error: expected: WAREHOUSE, got: 
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [6/7]: failed with error: expected: 1232, got: 172800
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [7/7]: failed with error: expected: WAREHOUSE, got: 
        check 10/11 error:
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [1/13]: failed with error: expected name: bad name; got: WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [2/13]: failed with error: expected state: SUSPENDED; got: STARTED
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [3/13]: failed with error: expected type: SNOWPARK-OPTIMIZED; got: STANDARD
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [4/13]: failed with error: expected size: MEDIUM; got: XSMALL
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [5/13]: failed with error: expected max cluster count: 12; got: 1
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [6/13]: failed with error: expected min cluster count: 13; got: 1
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [7/13]: failed with error: expected scaling policy: ECONOMY; got: STANDARD
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [8/13]: failed with error: expected auto suspend: 123; got: 600
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [9/13]: failed with error: expected auto resume: false; got: true
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [10/13]: failed with error: expected resource monitor: some-id; got: 
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [11/13]: failed with error: expected comment: bad comment; got: Promise my huh off certain you bravery dynasty with Roman.
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [12/13]: failed with error: expected enable query acceleration: true; got: false
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [13/13]: failed with error: expected query acceleration max scale factor: 12; got: 8
        check 11/11 error:
        parameter assertion for WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"][MAX_CONCURRENCY_LEVEL][1/1] failed: expected value 1, got 8
```

- add the following to the `create: complete` in `TestInt_Warehouses`:
```go
    // to show errors
    warehouseAssertionsBad := objectassert.Warehouse(t, id).
        HasName("bad name").
        HasState(sdk.WarehouseStateSuspended).
        HasType(sdk.WarehouseTypeSnowparkOptimized).
        HasSize(sdk.WarehouseSizeMedium).
        HasMaxClusterCount(12).
        HasMinClusterCount(13).
        HasScalingPolicy(sdk.ScalingPolicyStandard).
        HasAutoSuspend(123).
        HasAutoResume(false).
        HasResourceMonitor(sdk.NewAccountObjectIdentifier("some-id")).
        HasComment("bad comment").
        HasEnableQueryAcceleration(false).
        HasQueryAccelerationMaxScaleFactor(12)
    assertions.AssertThatObject(t, warehouseAssertionsBad)
```
it will result in:
```
    commons.go:101: 
        	Error Trace:	/Users/asawicki/Projects/terraform-provider-snowflake/pkg/sdk/testint/warehouses_integration_test.go:149
        	Error:      	Received unexpected error:
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [1/13]: failed with error: expected name: bad name; got: VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [2/13]: failed with error: expected state: SUSPENDED; got: STARTED
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [3/13]: failed with error: expected type: SNOWPARK-OPTIMIZED; got: STANDARD
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [4/13]: failed with error: expected size: MEDIUM; got: SMALL
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [5/13]: failed with error: expected max cluster count: 12; got: 8
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [6/13]: failed with error: expected min cluster count: 13; got: 2
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [7/13]: failed with error: expected type: STANDARD; got: ECONOMY
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [8/13]: failed with error: expected auto suspend: 123; got: 1000
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [9/13]: failed with error: expected auto resume: false; got: true
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [10/13]: failed with error: expected resource monitor: some-id; got: OOUJMDIT_535F314F_6549_348F_370E_AB430EE4BC7B
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [11/13]: failed with error: expected comment: bad comment; got: comment
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [12/13]: failed with error: expected enable query acceleration: false; got: true
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [13/13]: failed with error: expected query acceleration max scale factor: 12; got: 90
        	Test:       	TestInt_Warehouses/create:_complete
```

## Known limitations/planned improvements
- Test all the utilities for assertion/model construction (public interfaces, methods, functions).
- Verify if all the config types are supported.
- Consider a better implementation for the model conversion to config (TODO left in `config/config.go`).
- Support additional methods for references in models (TODO left in `config/config.go`).
- Support depends_on in models so that it can be chained like other resource fields (TODO left in `config/config.go`).
- Add a convenience function to concatenate multiple models (TODO left in `config/config.go`).
- Add function to support using `ConfigFile:` in the acceptance tests (TODO left in `config/config.go`).
- Replace `acceptance/snowflakechecks` with the new proposed Snowflake objects assertions.
- Support `showOutputValueUnset` and add a second function for each `show_output` attribute.
- Support `resourceAssertionTypeValueNotSet` for import checks (`panic` left currently).
- Add assertions for the `describe_output` and handle describe objects too.
- Add support for datasource tests (assertions and config builders).
- Consider overriding the assertions when invoking same check multiple times with different params (e.g. `Warehouse(...).HasType(X).HasType(Y)`; it could use the last-check-wins approach, to more easily reuse complex checks between the test steps).
- Consider not adding the check for `show_output` presence on creation (same with `parameters`). The majority of the use cases need it to be present but there are a few others (like conditional presence in the datasources). Currently, it seems that they should be always present in the resources, so no change is made. Later, with adding the support for the datasource tests, consider simple destructive implementation like:
- Add support for `set` so that assertions like e.g. `oauth_scopes.*` could be done.
```go
func (w *WarehouseDatasourceShowOutputAssert) IsEmpty() {
    w.assertions = make([]resourceAssertion, 0)
    w.assertions = append(w.assertions, valueSet("show_output.#", "0"))
}
```
- support other mappings if needed (TODO left in `assert/objectassert/gen/model.go`)
- consider extracting preamble model to commons (TODOs left in `assert/objectassert/gen/model.go` and in `assert/objectparametersassert/gen/model.go`)
- get a runtime name for the assertion creator (TODOs left in `assert/objectparametersassert/gen/model.go`)
- use a better definition for each objet's snowflake parameters (TODO left in `assert/objectparametersassert/gen/main/main.go`)
- add possibility to have enums generated in config builders (TODO left in `config/model/warehouse_model_ext.go`)
- handle situations where snowflake default behaves inconsistently (TODO left in `assert/objectparametersassert/gen/object_parameters_def.go`)
- handle attribute types in resource assertions (currently strings only; TODO left in `assert/resourceassert/gen/model.go`)
- distinguish between different enum types (TODO left in `assert/resourceshowoutputassert/gen/templates.go`)
- support the rest of attribute types in config model builders (TODO left in `config/model/gen/model.go`)
- parametrize test client helper used - integration versus acceptance tests - this has to be changed in the generator too (TODO left in `assert/objectassert/user_snowflake_ext.go`)
- Omit computed fields in the model (like FullyQualifiedName), because it doesn't make sense to set them
- There's an error when generating models, steps to reproduce:
  - Go to view resource code and change `data_metric_function` field to `testing` and make it required
  - During the generation, the following error appears: mixed named and unnamed parameters.
    It's a golang error indicating that the parameter has both unnamed and named parameters in function (e.g. `func(abc string, int)`).
    The error is a result of both things:
    1. Lists of objects are partially generated, and only parameter name is generated in some functions (the type has to be added manually).
    2. `testing` is a package name that makes Go think that we want to have unnamed parameter there, but we just didn't generate the type for that field in the function argument.
- generate assertions checking that time is not empty - we often do not compare time fields by value, but check if they are set
- utilize `ContainsExactlyInAnyOrder` function in `pkg/acceptance/bettertestspoc/assert/commons.go` to create asserts on collections that are order independent
- support generating provider config and use generated configs in `pkg/provider/provider_acceptance_test.go`
