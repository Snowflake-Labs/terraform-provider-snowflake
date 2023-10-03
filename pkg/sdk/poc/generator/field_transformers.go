package generator

import "golang.org/x/exp/slices"

type FieldTransformer interface {
	Transform(f *Field) *Field
}

type KeywordTransformer struct {
	required  bool
	sqlPrefix string
	quotes    string
}

func KeywordOptions() *KeywordTransformer {
	return new(KeywordTransformer)
}

func (v *KeywordTransformer) Required() *KeywordTransformer {
	v.required = true
	return v
}

func (v *KeywordTransformer) SQL(sqlPrefix string) *KeywordTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *KeywordTransformer) SingleQuotes() *KeywordTransformer {
	v.quotes = "single_quotes"
	return v
}

func (v *KeywordTransformer) DoubleQuotes() *KeywordTransformer {
	v.quotes = "double_quotes"
	return v
}

func (v *KeywordTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "keyword")
	if v.required {
		f.Required = true
	}
	addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	addTagIfMissing(f.Tags, "ddl", v.quotes)
	return f
}

type ParameterTransformer struct {
	required    bool
	sqlPrefix   string
	quotes      string
	parentheses string
}

func ParameterOptions() *ParameterTransformer {
	return new(ParameterTransformer)
}

func (v *ParameterTransformer) Required() *ParameterTransformer {
	v.required = true
	return v
}

func (v *ParameterTransformer) SQL(sqlPrefix string) *ParameterTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *ParameterTransformer) NoQuotes() *ParameterTransformer {
	v.quotes = "no_quotes"
	return v
}

func (v *ParameterTransformer) SingleQuotes() *ParameterTransformer {
	v.quotes = "single_quotes"
	return v
}

func (v *ParameterTransformer) DoubleQuotes() *ParameterTransformer {
	v.quotes = "double_quotes"
	return v
}

func (v *ParameterTransformer) Parentheses() *ParameterTransformer {
	v.quotes = "parentheses"
	return v
}

func (v *ParameterTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "parameter")
	if v.required {
		f.Required = true
	}
	addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	addTagIfMissing(f.Tags, "ddl", v.quotes)
	addTagIfMissing(f.Tags, "ddl", v.parentheses)
	return f
}

type ListTransformer struct {
	required    bool
	sqlPrefix   string
	parentheses string
}

func ListOptions() *ListTransformer {
	return new(ListTransformer)
}

func (v *ListTransformer) Required() *ListTransformer {
	v.required = true
	return v
}

func (v *ListTransformer) NoParens() *ListTransformer {
	v.parentheses = "no_parentheses"
	return v
}

func (v *ListTransformer) SQL(sqlPrefix string) *ListTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *ListTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "list")
	if v.required {
		f.Required = true
	}
	addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	addTagIfMissing(f.Tags, "ddl", v.parentheses)
	return f
}

type IdentifierTransformer struct {
	required  bool
	sqlPrefix string
	quotes    string
}

func IdentifierOptions() *IdentifierTransformer {
	return new(IdentifierTransformer)
}

func (v *IdentifierTransformer) SQL(sqlPrefix string) *IdentifierTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *IdentifierTransformer) SingleQuotes() *IdentifierTransformer {
	v.quotes = "single_quotes"
	return v
}

func (v *IdentifierTransformer) DoubleQuotes() *IdentifierTransformer {
	v.quotes = "double_quotes"
	return v
}

func (v *IdentifierTransformer) Required() *IdentifierTransformer {
	v.required = true
	return v
}

func (v *IdentifierTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "identifier")
	if v.required {
		f.Required = true
	}
	addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	addTagIfMissing(f.Tags, "ddl", v.quotes)
	return f
}

func addTagIfMissing(m map[string][]string, key string, value string) {
	if len(value) > 0 {
		if val, ok := m[key]; ok {
			if !slices.Contains(val, value) {
				m[key] = append(val, value)
			}
		} else {
			m[key] = []string{value}
		}
	}
}
