package sdk

import "strconv"

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

func functionDetailsFromRows(rows []functionDetailRow) *FunctionDetails {
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

func (r functionDetailRow) toStringProperty() *StringProperty {
	prop := &StringProperty{}
	if r.Value.Valid {
		prop.Value = r.Value.String
	}
	return prop
}

func (r functionDetailRow) toIntProperty() *IntProperty {
	prop := &IntProperty{}
	if r.Value.Valid {
		var value *int
		v, err := strconv.Atoi(r.Value.String)
		if err == nil {
			value = &v
		} else {
			value = nil
		}
		prop.Value = value
	}
	return prop
}

func (r functionDetailRow) toFloatProperty() *FloatProperty {
	prop := &FloatProperty{}
	if r.Value.Valid {
		var value *float64
		v, err := strconv.ParseFloat(r.Value.String, 64)
		if err == nil {
			value = &v
		} else {
			value = nil
		}
		prop.Value = value
	}
	return prop
}

func (r functionDetailRow) toBoolProperty() *BoolProperty {
	prop := &BoolProperty{}
	if r.Value.Valid {
		var value bool
		if r.Value.String != "" && r.Value.String != "null" {
			value = ToBool(r.Value.String)
		} else {
			value = false
		}
		prop.Value = value
	}
	return prop
}
