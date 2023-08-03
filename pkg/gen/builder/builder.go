package builder

import (
	"reflect"
)

func QueryStruct(name string) *StructBuilder {
	return &StructBuilder{name: name}
}

type StructBuilder struct {
	name   string
	fields []FieldBuilder
	nodes  []Node
}
type Node struct {
	kind Kind
}
type Kind struct {
}

type EmptyFieldBuilder struct {
	structBuilder *StructBuilder
	name          string
}

type FieldBuilderWithSql struct {
	name string
	sql  string
}
type FieldBuilderWithRequired struct {
	structBuilder *StructBuilder
	name          string
	isRequired    bool
}
type FieldBuilderWithRequiredAndSql struct {
	name       string
	isRequired bool
	sql        string
}

func (sb *StructBuilder) Create() *StructBuilder {
	return sb.Field("create").Required().Sql("CREATE").NoValue().End()
}
func (sb *StructBuilder) OrReplace() *StructBuilder {
	return sb.Field("OrReplace").Required().Sql("OR REPLACE").Value(Var(T[bool]()).End()).End()
}
func (sb *StructBuilder) OneOf(fieldsFuncs ...func(builder *StructBuilder) FieldBuilder) *StructBuilder {
	var fields []FieldBuilder
	for _, ff := range fieldsFuncs {
		fields = append(fields, ff(sb))
	}
	sb.fields = append(sb.fields)
	return sb
}

func Field(name string) *EmptyFieldBuilder {
	return &EmptyFieldBuilder{name: name}

}

//	func (sb *StructBuilder) Field(name string) func(builder *StructBuilder) *EmptyFieldBuilder {
//		return func(builder *StructBuilder) *EmptyFieldBuilder {
//			return &EmptyFieldBuilder{name: name, structBuilder: sb}
//		}
//
// }
func (sb *StructBuilder) Field(name string) *EmptyFieldBuilder {
	return &EmptyFieldBuilder{name: name, structBuilder: sb}

}

func QueryAdt(name string) *AdtBuilder {
	return &AdtBuilder{name: name}
}

type EnumBuilder[T any] struct {
	name   string
	values []EnumValue[T]
}

type EnumValue[T any] struct {
	Name  string
	Value T
}

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

func (adtBuilder *AdtBuilder) With(name string, adtCase *AdtCase) *AdtBuilder {
	adtBuilder.cases = append(adtBuilder.cases, AdtCaseWithName{name: name, adtCase: adtCase})
	return adtBuilder
}
func (adtBuilder *AdtBuilder) Type() Type {
	return Type{name: adtBuilder.name}
}

type AdtBuilder struct {
	name  string
	cases []AdtCaseWithName
}
type AdtCaseWithName struct {
	name    string
	adtCase *AdtCase
}
type AdtCase struct {
	name   string
	values []AdtValue
}
type AdtValue struct {
	//one of
	stringValue string
	intValue    int
	variable    VarAdtValue
}
type VarAdtValue struct {
	_type Type
	name  string
}

func (adtCase *AdtCase) Var(_type Type, name string) *AdtCase {
	adtCase.values = append(adtCase.values,
		AdtValue{
			variable: VarAdtValue{
				_type: _type,
				name:  name,
			},
		},
	)
	return adtCase
}
func (adtCase *AdtCase) String(value string) *AdtCase {
	adtCase.values = append(adtCase.values,
		AdtValue{
			stringValue: value,
		},
	)
	return adtCase
}

func (adtCase *AdtCase) Int(value int) *AdtCase {
	adtCase.values = append(adtCase.values,
		AdtValue{
			intValue: value,
		},
	)
	return adtCase
}

// inna paczka, zeby bylo d.Var
func DVar(_type Type, name string) *AdtCase {
	adtCase := &AdtCase{}
	adtCase.values = append(adtCase.values,
		AdtValue{
			variable: VarAdtValue{
				_type: _type,
				name:  name,
			},
		},
	)
	return adtCase
}
func DString(value string) *AdtCase {
	adtCase := &AdtCase{}
	adtCase.values = append(adtCase.values,
		AdtValue{
			stringValue: value,
		},
	)
	return adtCase
}

func DInt(value int) *AdtCase {
	adtCase := &AdtCase{}
	adtCase.values = append(adtCase.values,
		AdtValue{
			intValue: value,
		},
	)
	return adtCase
}

func (fb *EmptyFieldBuilder) Sql(sql string) FieldBuilderWithSql {
	return FieldBuilderWithSql{name: fb.name, sql: sql}
}
func (fb *EmptyFieldBuilder) NoSql() *FieldBuilderWithSql {
	return &FieldBuilderWithSql{name: fb.name, sql: ""}
}

type Foo func(*StructBuilder) EmptyFieldBuilder

func (f Foo) Required() func(builder *StructBuilder) *FieldBuilderWithRequired {
return func(builder *StructBuilder) *FieldBuilderWithRequired {
	a := f(builder)
	return &FieldBuilderWithRequired{name: a.name, isRequired: true, structBuilder: builder}
}
}

func (fb *EmptyFieldBuilder) Required() *FieldBuilderWithRequired {
	return &FieldBuilderWithRequired{name: fb.name, isRequired: true}
}
func (fb *EmptyFieldBuilder) Optional() *FieldBuilderWithRequired {
	return &FieldBuilderWithRequired{name: fb.name, isRequired: false}
}
func (fb *FieldBuilderWithRequired) Sql(sql string) *FieldBuilderWithRequiredAndSql {
	return &FieldBuilderWithRequiredAndSql{name: fb.name, sql: sql, isRequired: fb.isRequired}
}
func (fb *FieldBuilderWithRequired) NoSql() *FieldBuilderWithRequiredAndSql {
	return &FieldBuilderWithRequiredAndSql{name: fb.name, sql: "", isRequired: fb.isRequired}
}

func (fb *FieldBuilderWithRequiredAndSql) Value(value *Value) *FieldBuilder {
	return &FieldBuilder{name: fb.name, isRequired: fb.isRequired, sql: fb.sql, value: value}
}

func (fb *FieldBuilderWithRequiredAndSql) NoValue() *FieldBuilder {
	return &FieldBuilder{name: fb.name, isRequired: fb.isRequired, sql: fb.sql, value: nil}
}

func (fb *FieldBuilder) End() *StructBuilder {
	builder.fields = append(builder.fields, *fb)
	return builder
}

func Equals() *ValueBuilder {
	return &ValueBuilder{hasEquals: true}
}

func Var(_type Type) *VarValueBuilder {
	return &VarValueBuilder{hasEquals: false, _type: _type}

}
func (vb *ValueBuilder) Var(_type Type) *VarValueBuilder {
	return &VarValueBuilder{hasEquals: vb.hasEquals, _type: _type}

}

func (vb *ValueBuilder) List(_type Type) *ListValueBuilder {
	return &ListValueBuilder{hasEquals: vb.hasEquals, _type: _type}

}

func (vb *ValueBuilder) SchemaObjectIdentifier() *Value {
	return nil
}

type VarValueBuilder struct {
	hasEquals bool
	_type     Type
}

type ListValueBuilder struct {
	hasEquals bool
	_type     Type
}

type ListValueBuilderWithCommas struct {
	hasEquals bool
	_type     Type
	hasCommas bool
}

type ListValueBuilderWithParens struct {
	hasEquals bool
	_type     Type
	hasParens bool
}

type ListValueBuilderWithCommasAndParens struct {
	hasEquals bool
	_type     Type
	hasParens bool
	hasCommas bool
}

func (lb *ListValueBuilder) Parens() *ListValueBuilderWithParens {
	return &ListValueBuilderWithParens{_type: lb._type, hasEquals: lb.hasEquals, hasParens: true}
}
func (lb *ListValueBuilder) NoParens() *ListValueBuilderWithParens {
	return &ListValueBuilderWithParens{_type: lb._type, hasEquals: lb.hasEquals, hasParens: false}
}
func (lb *ListValueBuilder) Commas() *ListValueBuilderWithCommas {
	return &ListValueBuilderWithCommas{_type: lb._type, hasEquals: lb.hasEquals, hasCommas: true}
}
func (lb *ListValueBuilder) NoCommas() *ListValueBuilderWithCommas {
	return &ListValueBuilderWithCommas{_type: lb._type, hasEquals: lb.hasEquals, hasCommas: false}
}
func (lb *ListValueBuilderWithCommas) Parens() *ListValueBuilderWithCommasAndParens {
	return &ListValueBuilderWithCommasAndParens{_type: lb._type, hasEquals: lb.hasEquals, hasParens: true, hasCommas: lb.hasCommas}
}
func (lb *ListValueBuilderWithParens) Commas() *ListValueBuilderWithCommasAndParens {
	return &ListValueBuilderWithCommasAndParens{_type: lb._type, hasEquals: lb.hasEquals, hasCommas: true, hasParens: lb.hasParens}
}
func (lb *ListValueBuilderWithCommas) NoParens() *ListValueBuilderWithCommasAndParens {
	return &ListValueBuilderWithCommasAndParens{_type: lb._type, hasEquals: lb.hasEquals, hasParens: false, hasCommas: lb.hasCommas}
}
func (lb *ListValueBuilderWithParens) NoCommas() *ListValueBuilderWithCommasAndParens {
	return &ListValueBuilderWithCommasAndParens{_type: lb._type, hasEquals: lb.hasEquals, hasCommas: false, hasParens: lb.hasParens}
}

type VarValueBuilderWithQuotes struct {
	hasEquals      bool
	_type          Type
	isSingleQuotes bool
	isNoQuotes     bool
	isDoubleQuotes bool
}

func (vb *VarValueBuilder) SingleQuotes() *VarValueBuilderWithQuotes {
	return &VarValueBuilderWithQuotes{
		_type:          vb._type,
		hasEquals:      vb.hasEquals,
		isSingleQuotes: true,
		isDoubleQuotes: false,
		isNoQuotes:     false,
	}
}

func (vb *VarValueBuilder) DoubleQuotes() *VarValueBuilderWithQuotes {
	return &VarValueBuilderWithQuotes{
		_type:          vb._type,
		hasEquals:      vb.hasEquals,
		isSingleQuotes: false,
		isDoubleQuotes: true,
		isNoQuotes:     false,
	}
}

func (vb *VarValueBuilder) NoQuotes() *VarValueBuilderWithQuotes {
	return &VarValueBuilderWithQuotes{
		_type:          vb._type,
		hasEquals:      vb.hasEquals,
		isSingleQuotes: false,
		isDoubleQuotes: false,
		isNoQuotes:     true,
	}
}
func (vb *VarValueBuilder) End() *Value {
	return nil
}

// albo to, albo nie ma typu jak *VarValueBuilderWithQuotes i dopuszczamy Var(..).SingleQuotes().DoubleQuotes()
// bez tego end tez bedzie problem z roznymi typami, jak lista i identifiers
func (vb *VarValueBuilderWithQuotes) End() *Value {
	return nil
}

func (vb *ListValueBuilderWithCommas) End() *Value {
	return nil
}
func (vb *ListValueBuilderWithParens) End() *Value {
	return nil
}
func (vb *ListValueBuilderWithCommasAndParens) End() *Value {
	return nil
}
func (vb *ListValueBuilder) End() *Value {
	return nil
}

func T[T any]() Type {
	tType := reflect.TypeOf(new(T))
	return Type{
		name: tType.Name(),
	}
}

type Type struct {
	name string
}

type Value struct {
	name string
}
type ValueBuilder struct {
	hasEquals bool
}

type FieldBuilder struct {
	name       string
	isRequired bool
	sql        string
	value      *Value
}

// func (fb *FieldBuilderWithSql) Required() *FieldBuilderWithRequiredAndSql {
// 	return &FieldBuilderWithRequiredAndSql{name: fb.name, sql: fb.sql, isRequired: true}
//
// }
// func (fb *FieldBuilderWithSql) Optional() *FieldBuilderWithRequiredAndSql {
// 	return &FieldBuilderWithRequiredAndSql{name: fb.name, sql: fb.sql, isRequired: false}
// }

func foo() {
	//enumy??

	b := QueryAdt("AAA").
		With("abc", DString("FOO").Var(T[string](), "a").Int(10)).
		//return fmt.Sprintf("FOO%v%d",a,10)
		With("abc", DString("FOO"))

	a := QueryStruct("Name").
		// Field("foo").Sql("FOO")
		Field("name").Optional().Sql("AAA").Value(Var(T[string]()).DoubleQuotes().End()).End()
		Field("row access policy").Required().Sql("CREATE").Value(Var(T[string]()).NoQuotes().End()).End().
		Field("foo").Optional().Sql("foo").Value(Equals().Var(T[string]()).NoQuotes().End()).End().
		Field("foo").Optional().Sql("foo").Value(Equals().List(T[string]()).Commas().NoParens().End()).End().
		Field("foo").Optional().Sql("foo").Value(Equals().List(b.Type()).Commas().NoParens().End()).End()
	//tutaj opcja z wieloma fieldami na jednej linijce

}
