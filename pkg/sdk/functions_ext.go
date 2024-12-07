package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

func (v *Function) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, v.ArgumentsOld...)
}

// FunctionDetails contains aggregated describe results for the given function.
type FunctionDetails struct {
	Signature                  string  // present for all function types
	Returns                    string  // present for all function types
	Language                   string  // present for all function types
	Body                       *string // present for all function types (hidden when SECURE)
	NullHandling               *string // present for all function types but SQL
	Volatility                 *string // present for all function types but SQL
	ExternalAccessIntegrations *string // list present for python, java, and scala
	Secrets                    *string // map present for python, java, and scala
	Imports                    *string // list present for python, java, and scala (hidden when SECURE)
	Handler                    *string // present for python, java, and scala (hidden when SECURE)
	RuntimeVersion             *string // present for python, java, and scala (hidden when SECURE)
	Packages                   *string // list // present for python, java, and scala
	TargetPath                 *string // list present for scala and java (hidden when SECURE)
	InstalledPackages          *string // list present for python (hidden when SECURE)
	IsAggregate                *bool   // present for python
}

func functionDetailsFromRows(rows []FunctionDetail) (*FunctionDetails, error) {
	v := &FunctionDetails{}
	var errs []error
	for _, row := range rows {
		switch row.Property {
		case "signature":
			errs = append(errs, row.setStringValueOrError("signature", &v.Signature))
		case "returns":
			errs = append(errs, row.setStringValueOrError("returns", &v.Returns))
		case "language":
			errs = append(errs, row.setStringValueOrError("language", &v.Language))
		case "null handling":
			v.NullHandling = row.Value
		case "volatility":
			v.Volatility = row.Value
		case "body":
			v.Body = row.Value
		case "external_access_integrations":
			v.ExternalAccessIntegrations = row.Value
		case "secrets":
			v.Secrets = row.Value
		case "imports":
			v.Imports = row.Value
		case "handler":
			v.Handler = row.Value
		case "runtime_version":
			v.RuntimeVersion = row.Value
		case "packages":
			v.Packages = row.Value
		case "installed_packages":
			v.InstalledPackages = row.Value
		case "is_aggregate":
			errs = append(errs, row.setOptionalBoolValueOrError("is_aggregate", &v.IsAggregate))
		case "targetPath":
			v.TargetPath = row.Value
		}
	}
	return v, errors.Join(errs...)
}

func (v *functions) DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*FunctionDetails, error) {
	rows, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return functionDetailsFromRows(rows)
}

func (d *FunctionDetail) setStringValueOrError(property string, field *string) error {
	if d.Value == nil {
		return fmt.Errorf("value expected for field %s", property)
	} else {
		*field = *d.Value
	}
	return nil
}

func (d *FunctionDetail) setOptionalBoolValueOrError(property string, field **bool) error {
	if d.Value != nil && *d.Value != "" {
		v, err := strconv.ParseBool(*d.Value)
		if err != nil {
			return fmt.Errorf("invalid value for field %s, err: %w", property, err)
		} else {
			*field = Bool(v)
		}
	}
	return nil
}

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
