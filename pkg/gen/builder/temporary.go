package builder

//// Fields
//
//func (s Struct) IntoFieldBuilder() []FieldBuilder {
//	return []FieldBuilder{
//		{
//			Name:  s.name,
//			Typer: s.name,
//			Tags:  make(map[string][]string),
//		},
//	}
//}
//
////// Enums
//
////// Options
//
//type option func(sb *StructBuilder, f *FieldBuilder)
//
////type staticOption func(sb *StructBuilder, f *FieldBuilder)
////type keywordOption func(sb *StructBuilder, f *FieldBuilder)
////type parameterOption func(sb *StructBuilder, f *FieldBuilder)
////type keywordOptions
////type staticOptions
////type parameterOptions
//
//func WithTag(tagName string, tagValue string) func(sb *StructBuilder, f *FieldBuilder) {
//	return func(sb *StructBuilder, f *FieldBuilder) {
//		if _, ok := f.Tags[tagName]; !ok {
//			f.Tags[tagName] = make([]string, 0)
//		}
//		f.Tags[tagName] = append(f.Tags[tagName], tagValue)
//	}
//}
//
//func WithSQL(sql string) option {
//	return WithTag("sql", sql)
//}
//
//func WithSQLPrefix(sql string) option {
//	return WithTag("sql", sql)
//}
//
//func WithDDL(ddl string) option {
//	return WithTag("ddl", ddl)
//}
//
//func WithParen() option {
//	return WithDDL("parentheses")
//}
//
//func WithNoQuotes() option {
//	return WithDDL("no_quotes")
//}
//
//func NoEquals() option {
//	return WithDDL("no_equals")
//}
//
////func (sb StructBuilder) Type() string {
////	return sb.Name
////}
//
//// func Optional(kinder Typer) func(sb *StructBuilder, f *FieldBuilder) {
//// 	return withType("*" + kinder.Typer())
//// }
//
//func Number() func(sb *StructBuilder, f *FieldBuilder) {
//	return func(sb *StructBuilder, f *FieldBuilder) {
//		f.Typer = TypeInt.Type()
//	}
//}
//
//func Text() func(sb *StructBuilder, f *FieldBuilder) {
//	return func(sb *StructBuilder, f *FieldBuilder) {
//		f.Typer = TypeInt.Type()
//	}
//}
//
//func SQLPrefix(s string) option {
//	return func(sb *StructBuilder, f *FieldBuilder) {
//		f.Typer = TypeInt.Type()
//	}
//}
//
//func SingleQuotedText() func(sb *StructBuilder, f *FieldBuilder) {
//	return func(sb *StructBuilder, f *FieldBuilder) {
//		f.Typer = TypeInt.Type()
//	}
//}
//
////func WithType(kinder Typer) func(sb *StructBuilder, f *FieldBuilder) {
////	return func(sb *StructBuilder, f *FieldBuilder) {
////		f.Typer = kinder.Typer()
////	}
////}
////
////func withType(kind string) func(sb *StructBuilder, f *FieldBuilder) {
////	return func(sb *StructBuilder, f *FieldBuilder) {
////		f.Typer = kind
////	}
////}
////
////func TypeT[T any]() func(sb *StructBuilder, f *FieldBuilder) {
////	return func(sb *StructBuilder, f *FieldBuilder) {
////		tType := reflect.TypeOfString(new(T))
////		fmt.Printf("LUL: %v", new(T))
////		f.Typer = tType.Name()
////	}
////}
////
////func WithPointerType(kinder Typer) func(sb *StructBuilder, f *FieldBuilder) {
////	return withType("*" + kinder.Typer())
////}
////
////func ListOf(kinder Typer) func(sb *StructBuilder, f *FieldBuilder) {
////	return withType("[]" + kinder.Typer())
////}
////
////func WithTypeBoolPtr() option {
////	return WithType(TypeBoolPtr)
////}
//
////// Building functions
//
////// SQL
////
////func (sb *StructBuilder) SQL(fieldName string, sql string, options ...option) *StructBuilder {
////	return sb.Field("static", fieldName, options...)
////}
////
////func (sb *StructBuilder) SQL(sql string, options ...intOptions) *StructBuilder {
////	return sb.Field("static", fieldName, options...)
////}
////
////func (sb *StructBuilder) SQL(sql string, options ...option) *StructBuilder {
////	return sb.Field("static", fieldName, options...)
////}
////
////func (sb *StructBuilder) Create(opts ...createOptions) *StructBuilder {
////	return sb.SQL("create", WithType(TypeBool), WithSQL("CREATE"))
////}
////
////// Keyword
////func (sb *StructBuilder) Tag() *StructBuilder {
////	return sb
////}
////
////func (sb *StructBuilder) Keyword(fieldName string, options ...option) *StructBuilder {
////	// TODO add sql tag ?
////	return sb.Field("keyword", fieldName, options...)
////}
////
////func (sb *StructBuilder) OptionalAssignment(fieldName string, options ...option) *StructBuilder {
////	// TODO add sql tag ?
////	return sb.Field("keyword", fieldName, options...)
////}
////
////func (sb *StructBuilder) Assignment(fieldName string, options ...option) *StructBuilder {
////	// TODO add sql tag ?
////	return sb.Field("keyword", fieldName, options...)
////}
////
////func (sb *StructBuilder) OptionalValue(fieldName string, value any, options ...option) *StructBuilder {
////	// TODO add sql tag ?
////	return sb.Field("keyword", fieldName, options...)
////}
////
////func (sb *StructBuilder) Value(fieldName string, value any, options ...option) *StructBuilder {
////	// TODO add sql tag ?
////	return sb.Field("keyword", fieldName, options...)
////}
////
////func (sb *StructBuilder) OptionalSQL(sql string, options ...option) *StructBuilder {
////	// TODO add sql tag ?
////	return sb.Field("keyword", sql, options...)
////}
////
////func Keyword(fieldName string, options ...option) FieldBuilder {
////	fb := FieldBuilder{Name: fieldName, Typer: TypeBoolPtr.Typer(), Tags: map[string][]string{"ddl": {"keyword"}}}
////	for _, opt := range options {
////		opt(nil, &fb)
////	}
////	return fb
////}
////
////func Value(fieldName string, options ...option) FieldBuilder {
////	fb := FieldBuilder{Name: fieldName, Typer: TypeBoolPtr.Typer(), Tags: map[string][]string{"ddl": {"keyword"}}}
////	for _, opt := range options {
////		opt(nil, &fb)
////	}
////	return fb
////}
////
////func OptionalSQL(sql string, options ...option) FieldBuilder {
////	fb := FieldBuilder{Name: sql, Typer: TypeBoolPtr.Typer(), Tags: map[string][]string{"ddl": {"keyword"}}}
////	for _, opt := range options {
////		opt(nil, &fb)
////	}
////	return fb
////}
////
////func OptionalText(sql string, options ...option) FieldBuilder {
////	fb := FieldBuilder{Name: sql, Typer: TypeBoolPtr.Typer(), Tags: map[string][]string{"ddl": {"keyword"}}}
////	for _, opt := range options {
////		opt(nil, &fb)
////	}
////	return fb
////}
////
////func OptionalValue(sql string, typeObj any, options ...option) FieldBuilder {
////	fb := FieldBuilder{Name: sql, Typer: TypeBoolPtr.Typer(), Tags: map[string][]string{"ddl": {"keyword"}}}
////	for _, opt := range options {
////		opt(nil, &fb)
////	}
////	return fb
////}
////
////func (sb *StructBuilder) OrReplace() *StructBuilder {
////	return sb.Keyword("OrReplace", WithTypeBoolPtr())
////}
////
////func (sb *StructBuilder) IfExists() *StructBuilder {
////	return sb.Keyword("IfExists", WithTypeBoolPtr())
////}
////
////func (sb *StructBuilder) IfNotExists() *StructBuilder {
////	return sb.Keyword("IfNotExists", WithTypeBoolPtr())
////}
////
////func (sb *StructBuilder) Transient() *StructBuilder {
////	return sb.Keyword("Transient", WithTypeBoolPtr())
////}
////
////func (sb *StructBuilder) Number(fieldName string, options ...option) *StructBuilder {
////	return sb.Field("identifier", fieldName, options...)
////}
////
////func (sb *StructBuilder) Text(fieldName string, options ...option) *StructBuilder {
////	return sb.Field("identifier", fieldName, options...)
////}
////
////func (sb *StructBuilder) OptionalText(fieldName string, options ...option) *StructBuilder {
////	return sb.Field("identifier", fieldName, options...)
////}
////
////// Identifier
////func (sb *StructBuilder) Identifier(fieldName string, typeObj any, options ...option) *StructBuilder {
////	return sb.Field("identifier", fieldName, options...)
////}
////
////// Parameter
////func (sb *StructBuilder) Parameter(fieldName string, options ...option) *StructBuilder {
////	// TODO add sql tag ?
////	return sb.Field("parameter", fieldName, options...)
////}
////
////// Parameter
////func (sb *StructBuilder) List(fieldName string, typeObj any, options ...option) *StructBuilder {
////	// TODO add sql tag ?
////	return sb.Field("list", fieldName, options...)
////}
////
////type Group struct {
////	IntoFieldBuilders []IntoFieldBuilder
////}
////
////func (g Group) IntoFieldBuilder() []FieldBuilder {
////	fbs := make([]FieldBuilder, 0)
////	for _, f := range g.IntoFieldBuilders {
////		fbs = append(fbs, f.IntoFieldBuilder()...)
////	}
////	return fbs
////}
////
////func GroupFields(fields ...IntoFieldBuilder) Group {
////	return Group{
////		IntoFieldBuilders: fields,
////	}
////}
////
////func (sb *StructBuilder) Build() Struct {
////	return Struct{}
////}
////
//////// API
////
////type API struct{}
////
////func BuildAPI() API {
////	return API{}
////}
////
////type staticOption struct {
////	ddl []string
////}
////type staticOptionBuilder struct {
////	ddl []string
////}
////
////func StaticOpts() *staticOptionBuilder {
////	return &staticOptionBuilder{
////		ddl: make([]string, 0),
////	}
////}
////
////func (v *staticOptionBuilder) Quotes() *staticOptionBuilder {
////	v.ddl = append(v.ddl, "quotes")
////	return v
////}
////
////func (sb *StructBuilder) Static2(fieldName string) *StructBuilder {
////	return sb.Field("static", fieldName, nil)
////}
////
////func (sb *StructBuilder) Static2(opts staticOptionBuilder) *StructBuilder {
////	return sb.Static2("create")
////}
////
////type keywordOption struct {
////	sql []string
////}
////type keywordOptionBuilder struct {
////	sql []string
////}
////
////func KeywordOpts() *keywordOptionBuilder {
////	return &keywordOptionBuilder{
////		sql: make([]string, 0),
////	}
////}
////
////func (v *keywordOptionBuilder) SQL(sql string) *keywordOptionBuilder {
////	v.sql = append(v.sql, sql)
////	return v
////}
////
////type parameterOption struct {
////	ddl []string
////}
////type parameterOptionBuilder struct {
////	ddl []string
////}
////
////func ParameterOpts() *parameterOptionBuilder {
////	return &parameterOptionBuilder{
////		ddl: make([]string, 0),
////	}
////}
////
////func (v *parameterOptionBuilder) Paren(paren bool) *parameterOptionBuilder {
////	if paren {
////		v.ddl = append(v.ddl, "parentheses")
////	} else {
////		v.ddl = append(v.ddl, "no_parentheses")
////	}
////	return v
////}
////
////func (v *parameterOptionBuilder) Equals(equals bool) *parameterOptionBuilder {
////	if equals {
////		v.ddl = append(v.ddl, "equals")
////	} else {
////		v.ddl = append(v.ddl, "no_equals")
////	}
////	return v
////}
////
////func (sb *StructBuilder) Keyword2(fieldName string, opts *keywordOptionBuilder) *StructBuilder {
////	return sb.Field("keyword", fieldName, nil)
////}
////
////func (sb *StructBuilder) Parameter2(fieldName string, opts *parameterOptionBuilder) *StructBuilder {
////	return sb.Field("parameter", fieldName, nil)
////}
