package builder

type oneof struct {
	fields []IntoFieldBuilder
}

func OneOf(fieldName string, fields ...IntoFieldBuilder) oneof {
	return oneof{
		fields: fields,
	}
}

func (o oneof) IntoFieldBuilder() []FieldBuilder {
	fbs := make([]FieldBuilder, 0)
	for _, f := range o.fields {
		for _, fb := range f.IntoFieldBuilder() {
			fbs = append(fbs, fb)
		}
	}
	return fbs
}
