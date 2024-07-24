# Better tests poc
This package contains a quick implementation of helpers that should allow us a quicker, more pleasant, and more readable implementation of tests, mainly the acceptance ones.
It contains the following packages:
- `assert` - all the assertions reside here. Also, the utilities to build assertions for new objects. All the assertions will be ultimately generated; the ones presented in this folder were manually created. The currently supported assertions are:
  - resource assertions (currently, created manually)
  - show output assertions (currently, created manually)
  - resource parameters assertions (currently, created manually)
  - Snowflake object assertions (generated in subpackage `objectassert`)
  - Snowflake object parameters assertions (generated in subpackage `objectparametersassert`)

- `config` - the new ResourceModel abstraction resides here. It provides models for objects and the builder methods allowing better config preparation in the acceptance tests.
It aims to be more readable than using `Config:` with hardcoded string or `ConfigFile:` for file that is not directly reachable from the test body. Also, it should be easier to reuse the models and prepare convenience extension methods.
All the models will be ultimately generated; the ones presented for warehouse were manually created.

## Usage
You can check the current example usage in `TestAcc_Warehouse_BasicFlows` and the `create: complete` inside `TestInt_Warehouses`. To see the output after invalid assertions:
- add the following to the first step of `TestAcc_Warehouse_BasicFlows`
```go
    // bad checks below
    assert.WarehouseResource(t, "snowflake_warehouse.w").
        HasType(string(sdk.WarehouseTypeSnowparkOptimized)).
        HasSize(string(sdk.WarehouseSizeMedium)),
    assert.WarehouseShowOutput(t, "snowflake_warehouse.w").
        HasType(sdk.WarehouseTypeSnowparkOptimized),
    assert.WarehouseParameters(t, "snowflake_warehouse.w").
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
    assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", string(sdk.WarehouseTypeSnowparkOptimized))),
```
it will result in:
```
    warehouse_acceptance_test.go:46: Step 1/8 error: Check failed: check 6/10 error:
        snowflake_warehouse.w resource assertion [1/2]: failed with error: Attribute 'warehouse_type' not found
        snowflake_warehouse.w resource assertion [2/2]: failed with error: Attribute 'warehouse_size' not found
        check 7/10 error:
        snowflake_warehouse.w show_output assertion [2/2]: failed with error: Attribute 'show_output.0.type' expected "SNOWPARK-OPTIMIZED", got "STANDARD"
        check 8/10 error:
        snowflake_warehouse.w parameters assertion [2/3]: failed with error: Attribute 'parameters.0.max_concurrency_level.0.value' expected "16", got "8"
        snowflake_warehouse.w parameters assertion [3/3]: failed with error: Attribute 'parameters.0.max_concurrency_level.0.level' expected "WAREHOUSE", got ""
        check 9/10 error:
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [1/13]: failed with error: expected name: bad name; got: URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [2/13]: failed with error: expected state: SUSPENDED; got: STARTED
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [3/13]: failed with error: expected type: SNOWPARK-OPTIMIZED; got: STANDARD
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [4/13]: failed with error: expected size: MEDIUM; got: XSMALL
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [5/13]: failed with error: expected max cluster count: 12; got: 1
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [6/13]: failed with error: expected min cluster count: 13; got: 1
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [7/13]: failed with error: expected type: ECONOMY; got: STANDARD
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [8/13]: failed with error: expected auto suspend: 123; got: 600
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [9/13]: failed with error: expected auto resume: false; got: true
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [10/13]: failed with error: expected resource monitor: some-id; got: 
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [11/13]: failed with error: expected comment: bad comment; got: From furthermore rarely cast anything those you could also whoever.
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [12/13]: failed with error: expected enable query acceleration: true; got: false
        object WAREHOUSE["URVBDDAT_E7589B32_6534_1F93_DC1B_9E94FB8D27D7"] assertion [13/13]: failed with error: expected query acceleration max scale factor: 12; got: 8
        check 10/10 error:
        snowflake_warehouse.w: Attribute 'warehouse_type' not found
```

- add the following to the second step of `TestAcc_Warehouse_BasicFlows`
```go
    // bad checks below
    assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(warehouseId.Name(), "bad name", name)),
    assert.ImportedWarehouseResource(t, warehouseId.Name()).
        HasName("bad name").
        HasType(string(sdk.WarehouseTypeSnowparkOptimized)).
        HasSize(string(sdk.WarehouseSizeMedium)).
        HasMaxClusterCount("2").
        HasMinClusterCount("3").
        HasScalingPolicy(string(sdk.ScalingPolicyEconomy)).
        HasAutoSuspend("123").
        HasAutoResume("false").
        HasResourceMonitor("abc").
        HasComment("bad comment").
        HasEnableQueryAcceleration("true").
        HasQueryAccelerationMaxScaleFactor("16"),
    assert.ImportedWarehouseParameters(t, warehouseId.Name()).
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
```
it will result in:
```
    warehouse_acceptance_test.go:46: check 6/9 error:
        attribute bad name not found in instance state
        check 7/9 error:
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [1/12]: failed with error: expected: bad name, got: RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [2/12]: failed with error: expected: SNOWPARK-OPTIMIZED, got: STANDARD
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [3/12]: failed with error: expected: MEDIUM, got: XSMALL
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [4/12]: failed with error: expected: 2, got: 1
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [5/12]: failed with error: expected: 3, got: 1
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [6/12]: failed with error: expected: ECONOMY, got: STANDARD
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [7/12]: failed with error: expected: 123, got: 600
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [8/12]: failed with error: expected: false, got: true
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [9/12]: failed with error: expected: abc, got: 
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [10/12]: failed with error: expected: bad comment, got: School huh one here entirely mustering where crew though wealth.
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [11/12]: failed with error: expected: true, got: false
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported resource assertion [12/12]: failed with error: expected: 16, got: 8
        check 8/9 error:
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported parameters assertion [2/7]: failed with error: expected: 1, got: 8
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported parameters assertion [3/7]: failed with error: expected: WAREHOUSE, got: 
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported parameters assertion [4/7]: failed with error: expected: 23, got: 0
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported parameters assertion [5/7]: failed with error: expected: WAREHOUSE, got: 
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported parameters assertion [6/7]: failed with error: expected: 1232, got: 172800
        RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844 imported parameters assertion [7/7]: failed with error: expected: WAREHOUSE, got: 
        check 9/9 error:
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [1/13]: failed with error: expected name: bad name; got: RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [2/13]: failed with error: expected state: SUSPENDED; got: STARTED
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [3/13]: failed with error: expected type: SNOWPARK-OPTIMIZED; got: STANDARD
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [4/13]: failed with error: expected size: MEDIUM; got: XSMALL
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [5/13]: failed with error: expected max cluster count: 12; got: 1
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [6/13]: failed with error: expected min cluster count: 13; got: 1
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [7/13]: failed with error: expected type: ECONOMY; got: STANDARD
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [8/13]: failed with error: expected auto suspend: 123; got: 600
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [9/13]: failed with error: expected auto resume: false; got: true
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [10/13]: failed with error: expected resource monitor: some-id; got: 
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [11/13]: failed with error: expected comment: bad comment; got: School huh one here entirely mustering where crew though wealth.
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [12/13]: failed with error: expected enable query acceleration: true; got: false
        object WAREHOUSE["RQYLJJAT_04646516_1F33_50E9_CC19_D6B14E374844"] assertion [13/13]: failed with error: expected query acceleration max scale factor: 12; got: 8
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

## Adding new resource assertions
For object `abc` create the following files with the described content in the `assert` package:
- `abc_resource.go`
```go
type AbcResourceAssert struct {
    *ResourceAssert
}

func AbcResource(t *testing.T, name string) *AbcResourceAssert {
    t.Helper()

    return &AbcResourceAssert{
        ResourceAssert: NewResourceAssert(name, "resource"),
    }
}

func ImportedAbcResource(t *testing.T, id string) *AbcResourceAssert {
    t.Helper()

    return &AbcResourceAssert{
        ResourceAssert: NewImportedResourceAssert(id, "imported resource"),
    }
}
```
Two methods for each parameter (let's say parameter name is xyz):
```go
func (w *AbcResourceAssert) HasXyz(expected string) *AbcResourceAssert {
    w.assertions = append(w.assertions, valueSet("xyz", expected))
    return w
}

func (w *AbcResourceAssert) HasNoXyz() *AbcResourceAssert {
    w.assertions = append(w.assertions, valueNotSet("xyz"))
    return w
}
```

- `abc_show_output.go`
```go
type AbcShowOutputAssert struct {
    *ResourceAssert
}

func AbcShowOutput(t *testing.T, name string) *AbcShowOutputAssert {
    t.Helper()
    w := AbcShowOutputAssert{
        NewResourceAssert(name, "show_output"),
    }
    w.assertions = append(w.assertions, valueSet("show_output.#", "1"))
    return &w
}

func ImportedAbcShowOutput(t *testing.T, id string) *AbcShowOutputAssert {
    t.Helper()
    w := AbcShowOutputAssert{
        NewImportedResourceAssert(id, "show_output"),
    }
    w.assertions = append(w.assertions, valueSet("show_output.#", "1"))
    return &w
}
```

A method for each parameter (let's say parameter name is xyz):
```go
func (w *AbcShowOutputAssert) HasXyz(expected string) *AbcShowOutputAssert {
    w.assertions = append(w.assertions, showOutputValueSet("xyz", string(expected)))
    return w
}
```

- `abc_parameters.go`
```go
type AbcParametersAssert struct {
    *ResourceAssert
}

func AbcParameters(t *testing.T, name string) *AbcParametersAssert {
    t.Helper()
    w := AbcParametersAssert{
        NewResourceAssert(name, "parameters"),
    }
    w.assertions = append(w.assertions, valueSet("parameters.#", "1"))
    return &w
}

func ImportedAbcParameters(t *testing.T, id string) *AbcParametersAssert {
    t.Helper()
    w := AbcParametersAssert{
        NewImportedResourceAssert(id, "imported parameters"),
    }
    w.assertions = append(w.assertions, valueSet("parameters.#", "1"))
    return &w
}
```
Two methods for each parameter (let's say parameter name is xyz):
```go
func (w *AbcParametersAssert) HasXyz(expected int) *AbcParametersAssert {
    w.assertions = append(w.assertions, parameterValueSet("xyz", strconv.Itoa(expected)))
    return w
}

func (w *AbcParametersAssert) HasXyzLevel(expected sdk.ParameterType) *AbcParametersAssert {
    w.assertions = append(w.assertions, parameterLevelSet("xyz", string(expected)))
    return w
}
```
- extensions should be put in `abc_resource_ext.go`, `abc_show_output_ext.go`, or `abc_parameters_ext.go`. We can put here the named aggregations of other assertions. It allows us extendability. Later, we may choose to generate some of these methods too. Currently, the split will help when we start the generation of aforementioned methods. Current examples for extension could be:
```go
func (w *WarehouseResourceAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceAssert {
    w.assertions = append(w.assertions, valueSet("max_concurrency_level", "8"))
    return w
}

func (w *WarehouseResourceAssert) HasDefaultStatementQueuedTimeoutInSeconds() *WarehouseResourceAssert {
    w.assertions = append(w.assertions, valueSet("statement_queued_timeout_in_seconds", "0"))
    return w
}

func (w *WarehouseResourceAssert) HasDefaultStatementTimeoutInSeconds() *WarehouseResourceAssert {
    w.assertions = append(w.assertions, valueSet("statement_timeout_in_seconds", "172800"))
    return w
}

func (w *WarehouseResourceAssert) HasAllDefault() *WarehouseResourceAssert {
    return w.HasDefaultMaxConcurrencyLevel().
        HasNoType().
        HasNoSize().
        HasNoMaxClusterCount().
        HasNoMinClusterCount().
        HasNoScalingPolicy().
        HasAutoSuspend(r.IntDefaultString).
        HasAutoResume(r.BooleanDefault).
        HasNoInitiallySuspended().
        HasNoResourceMonitor().
        HasEnableQueryAcceleration(r.BooleanDefault).
        HasQueryAccelerationMaxScaleFactor(r.IntDefaultString).
        HasDefaultMaxConcurrencyLevel().
        HasDefaultStatementQueuedTimeoutInSeconds().
        HasDefaultStatementTimeoutInSeconds()
}

func (w *WarehouseParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseParametersAssert {
    return w.
        HasMaxConcurrencyLevel(8).
        HasMaxConcurrencyLevelLevel("")
}

func (w *WarehouseParametersAssert) HasDefaultStatementQueuedTimeoutInSeconds() *WarehouseParametersAssert {
    return w.
        HasStatementQueuedTimeoutInSeconds(0).
        HasStatementQueuedTimeoutInSecondsLevel("")
}

func (w *WarehouseParametersAssert) HasDefaultStatementTimeoutInSeconds() *WarehouseParametersAssert {
    return w.
        HasStatementTimeoutInSeconds(172800).
        HasStatementTimeoutInSecondsLevel("")
}
```

## Adding new Snowflake object assertions
Snowflake object assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allStructs` slice in the `assert/objectassert/gen/main/main.go`
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

## Adding new Snowflake object parameters assertions
Snowflake object parameters assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allObjectsParameters` slice in the `assert/objectparametersassert/gen/main/main.go`
- make sure that test helper method `acc.TestClient().Parameter.ShowAbcParameters` exists in `/pkg/acceptance/helpers/parameter_client.go`
- to add custom (not generated assertions) create file `abc_parameters_snowflake_ext.go` in the `objectparametersassert` package. Example would be:
```go
func (w *WarehouseParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseParametersAssert {
    return w.
        HasMaxConcurrencyLevel(8).
        HasMaxConcurrencyLevelLevel("")
}
```

## Adding new models
For object `abc` create the following files with the described content in the `config` package:
- `abc_model.go`
```go
type AbcModel struct {
    Xyz config.Variable `json:"xyz,omitempty"`
    *resourceModelMeta
}
```
two builders with required params only:
```go
func NewAbcModel(
    resourceName string,
    xyz string,
) *AbcModel {
    m := &AbcModel{resourceModelMeta: meta(resourceName, resources.Abc)}
    m.WithXyz(xyz)
    return m
}

func NewDefaultAbcModel(
    xyz string,
) *AbcModel {
    m := &AbcModel{resourceModelMeta: defaultMeta(resources.Abc)}
    m.WithXyz(xyz)
    return m
}
```
Two methods for each param (with good value type and with any value type):
```go
func (m *AbcModel) WithXyz(xyz string) *AbcModel {
    m.Xyz = config.StringVariable(xyz)
    return m
}

func (m *AbcModel) WithXyzValue(value config.Variable) *AbcModel {
    m.Xyz = value
    return m
}
```

- `abc_model_ext.go` - for the easier separation later (when we start generating the models for each object). Example would be:
```go
func BasicWarehouseModel(
    name string,
    comment string,
) *WarehouseModel {
    return NewDefaultWarehouseModel(name).WithComment(comment)
}
```

## Known limitations/planned improvements
- Generate all missing assertions and models.
- Test all the utilities for assertion/model construction (public interfaces, methods, functions).
- Verify if all the config types are supported.
- Consider a better implementation for the model conversion to config (TODO left).
- Support additional methods for references in models (TODO left).
- Support depends_on in models (TODO left).
- Add a convenience function to concatenate multiple models (TODO left).
- Add function to support using `ConfigFile:` in the acceptance tests.
- Replace `acceptance/snowflakechecks` with the new proposed Snowflake objects assertions.
- Support `showOutputValueUnset` and add a second function for each `show_output` attribute.
- Support `resourceAssertionTypeValueNotSet` for import checks (`panic` left currently).
- Add assertions for the `describe_output`.
- Add support for datasource tests (assertions and config builders).
- Consider overriding the assertions when invoking same check multiple times with different params (e.g. `Warehouse(...).HasType(X).HasType(Y)`; it could use the last-check-wins approach, to more easily reuse complex checks between the test steps).
- Consider not adding the check for `show_output` presence on creation (same with `parameters`). The majority of the use cases need it to be present but there are a few others (like conditional presence in the datasources). Currently, it seems that they should be always present in the resources, so no change is made. Later, with adding the support for the datasource tests, consider simple destructive implementation like:
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