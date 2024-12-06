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
