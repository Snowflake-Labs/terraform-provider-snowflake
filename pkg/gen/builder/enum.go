package builder

import (
	"fmt"
	"strings"
)

func EnumType[T any](name string) *EnumBuilder[T] {
	return &EnumBuilder[T]{
		name:   name,
		values: make([]EnumValue[T], 0),
	}
}

func (eb *EnumBuilder[T]) With(variableName string, value T) *EnumBuilder[T] {
	eb.values = append(eb.values, EnumValue[T]{
		Name:  variableName,
		Value: value,
	})
	return eb
}

func (eb *EnumBuilder[T]) IntoFieldBuilder() []FieldBuilder {
	return []FieldBuilder{
		{
			Name:  eb.name,
			Typer: TypeOfString(eb.name),
			Tags: map[string][]string{
				"ddl": {
					"keyword",
				},
			},
		},
	}
}

func (eb *EnumBuilder[T]) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("type %s string\n\n", eb.name))
	sb.WriteString("const (\n")
	for _, v := range eb.values {
		switch any(v.Value).(type) {
		case string:
			sb.WriteString(fmt.Sprintf("\t%-20s %s = %s\n", v.Name, eb.name, fmt.Sprintf(`"%s"`, v.Value)))
		default:
			sb.WriteString(fmt.Sprintf("\t%-20s %s = %v\n", v.Name, eb.name, v.Value))
		}
	}
	sb.WriteString(")\n")
	return sb.String()
}
