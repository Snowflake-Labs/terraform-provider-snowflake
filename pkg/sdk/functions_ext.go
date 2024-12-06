package sdk

import (
	"context"
	"strconv"
)

func (v *Function) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, v.ArgumentsOld...)
}

// FunctionDetails contains aggregated describe results for the given function.
// TODO [this PR]: do we keep *Property or types directly? -> types
type FunctionDetails struct {
	Signature                  *StringProperty
	Returns                    *StringProperty
	Language                   *StringProperty
	NullHandling               *StringProperty
	Volatility                 *StringProperty
	Body                       *StringProperty
	ExternalAccessIntegrations *StringProperty // list
	Secrets                    *StringProperty // map
	Imports                    *StringProperty // list
	Handler                    *StringProperty
	RuntimeVersion             *StringProperty
	Packages                   *StringProperty // list
	InstalledPackages          *StringProperty // list
	IsAggregate                *BoolProperty
	TargetPath                 *StringProperty
}

// TODO [this PR]: handle errors
func functionDetailsFromRows(rows []FunctionDetail) (*FunctionDetails, error) {
	v := &FunctionDetails{}
	for _, row := range rows {
		switch row.Property {
		case "signature":
			v.Signature = row.toStringProperty()
		case "returns":
			v.Returns = row.toStringProperty()
		case "language":
			v.Language = row.toStringProperty()
		case "null handling":
			v.NullHandling = row.toStringProperty()
		case "volatility":
			v.Volatility = row.toStringProperty()
		case "body":
			v.Body = row.toStringProperty()
		case "external_access_integrations":
			v.ExternalAccessIntegrations = row.toStringProperty()
		case "secrets":
			v.Secrets = row.toStringProperty()
		case "imports":
			v.Imports = row.toStringProperty()
		case "handler":
			v.Handler = row.toStringProperty()
		case "runtime_version":
			v.RuntimeVersion = row.toStringProperty()
		case "packages":
			v.Packages = row.toStringProperty()
		case "installed_packages":
			v.InstalledPackages = row.toStringProperty()
		case "is_aggregate":
			v.IsAggregate = row.toBoolProperty()
		case "targetPath":
			v.TargetPath = row.toStringProperty()
		}
	}
	return v, nil
}

func (v *functions) DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*FunctionDetails, error) {
	rows, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return functionDetailsFromRows(rows)
}

func (d *FunctionDetail) toStringProperty() *StringProperty {
	return &StringProperty{
		Value:       d.Value,
		Description: d.Property,
	}
}

func (d *FunctionDetail) toIntProperty() *IntProperty {
	var value *int
	v, err := strconv.Atoi(d.Value)
	if err == nil {
		value = &v
	} else {
		value = nil
	}
	return &IntProperty{
		Value:       value,
		Description: d.Property,
	}
}

func (d *FunctionDetail) toFloatProperty() *FloatProperty {
	var value *float64
	v, err := strconv.ParseFloat(d.Value, 64)
	if err == nil {
		value = &v
	} else {
		value = nil
	}
	return &FloatProperty{
		Value:       value,
		Description: d.Property,
	}
}

func (d *FunctionDetail) toBoolProperty() *BoolProperty {
	var value bool
	if d.Value != "" && d.Value != "null" {
		value = ToBool(d.Value)
	} else {
		value = false
	}
	return &BoolProperty{
		Value:       value,
		Description: d.Property,
	}
}

//python function describe:
//- signature
//- returns
//- language
//- null handling
//- volatility
//- [hidden for secure] body
//- external_access_integrations
//- secrets
//- [hidden for secure] imports
//- [hidden for secure] handler
//- [hidden for secure] runtime_version
//- [hidden for secure] packages
//- [hidden for secure] installed_packages
//- is_aggregate
//
//SQL function describe:
//- signature
//- returns
//- language
//- [hidden for secure] body
//
//scala function describe:
//- signature
//- returns
//- language
//- null handling
//- volatility
//- [hidden for secure] body
//- [hidden for secure] imports
//- [hidden for secure] handler
//- [hidden for secure] target_path
//- [hidden for secure] runtime_version
//- [hidden for secure] packages
//- external_access_integrations
//- secrets
//
//java:
//- signature
//- returns
//- language
//- null handling
//- volatility
//- [hidden for secure] body
//- [hidden for secure] imports
//- [hidden for secure] handler
//- [hidden for secure] target_path
//- [hidden for secure] runtime_version
//- [hidden for secure] packages
//- external_access_integrations
//- secrets
//
//javascript:
//- signature
//- returns
//- language
//- null handling
//- volatility
//- [hidden for secure] body

// CreateForJavaFunctionOptions
// TODO [SNOW-1348103 - this PR]: test setting the paths for all types (like imports, target paths)
// TODO [SNOW-1348103 - this PR]: test weird names for arg name - lower/upper if used with double quotes, to upper without quotes, dots, spaces, and both quotes not permitted
// TODO [SNOW-1348103 - next PRs]: check data type mappings https://docs.snowflake.com/en/sql-reference/sql/create-function#all-languages (signature + returns)
// TODO [SNOW-1348103 - this PR]: setting RUNTIME_VERSION (only 11.x, 17.x supported, 11.x being the default)
// TODO [SNOW-1348103 - this PR]: packages: package_name:version_number; do we validate? - check SELECT * FROM INFORMATION_SCHEMA.PACKAGES WHERE LANGUAGE = 'java';
// TODO [SNOW-1348103 - next PRs]: add to the resource docs https://docs.snowflake.com/en/sql-reference/sql/create-function#access-control-requirements
// TODO [SNOW-1348103 - this PR]: what delimiter do we use for <function_definition>: ' versus $$? - we use $$ as tasks
// TODO [SNOW-1348103 - this PR]: escaping single quotes test - don't have to do this with $$
// TODO [SNOW-1348103 - this PR]: validation of JAR (check https://docs.snowflake.com/en/sql-reference/sql/create-function#id6)
// TODO [SNOW-1348103 - next PRs]: active warehouse vs validations
// TODO [SNOW-1348103 - this PR]: check creation of all functions (using examples and more)

// CreateForPythonFunctionOptions
// TODO [SNOW-1348103 - this PR]: test aggregate func creation
// TODO [SNOW-1348103 - this PR]: what about [==<version>] - SDK level or resource level? check also: SELECT * FROM INFORMATION_SCHEMA.PACKAGES WHERE LANGUAGE = 'python';
// TODO [SNOW-1348103 - this PR]: what about preview feature >= ?
// TODO [SNOW-1348103 - this PR]: what about '<module_file_name>.<function_name>' for non-inline functions?
// TODO [SNOW-1348103 - this PR]: setting RUNTIME_VERSION (only 3.8, 3.9, 3.10, 3.11 supported, which one is a default?)

// CreateForScalaFunctionOptions
// TODO [SNOW-1348103 - this PR]: setting RUNTIME_VERSION (only 2.12 supported, which is the default)

// AlterFunctionOptions
// TODO [this PR]: can we run multiple sets/unsets? - yes, parameters + all besides SECURE
// TODO [this PR]: add setting EXTERNAL_ACCESS_INTEGRATIONS/SECRETS
// TODO [this PR]: unset EXTERNAL_ACCESS_INTEGRATIONS or SECRETS? - works for external access integrations, passes for secrets but does nothing SET to () works for secrets
// TODO [this PR]: EXTERNAL_ACCESS_INTEGRATIONS or SECRETS in Javascript or SQL - not working, working in SCALA though
