package sdk

import (
	"context"
	"strconv"
)

func (v *Function) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, v.ArgumentsOld...)
}

// FunctionDetails contains aggregated describe results for the given function.
// TODO [this PR]: fill out
type FunctionDetails struct {
	A *StringProperty
	B *IntProperty
	C *FloatProperty
	D *BoolProperty
}

func functionDetailsFromRows(rows []FunctionDetail) *FunctionDetails {
	v := &FunctionDetails{}
	for _, row := range rows {
		switch row.Property {
		case "A":
			v.A = row.toStringProperty()
		case "B":
			v.B = row.toIntProperty()
		case "C":
			v.C = row.toFloatProperty()
		case "D":
			v.D = row.toBoolProperty()
		}
	}
	return v
}

func (v *functions) DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*FunctionDetails, error) {
	rows, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return functionDetailsFromRows(rows), nil
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
