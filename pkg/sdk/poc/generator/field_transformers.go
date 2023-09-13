package generator

import "golang.org/x/exp/slices"

type FieldTransformer interface {
	Transform(f *Field) *Field
}

type KeywordTransformer struct {
	sqlPrefix string
	quotes    string
}

func KeywordOptions() *KeywordTransformer {
	return new(KeywordTransformer)
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
	if len(v.sqlPrefix) != 0 {
		addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	}
	if len(v.quotes) != 0 {
		addTagIfMissing(f.Tags, "ddl", v.quotes)
	}
	return f
}

type ParameterTransformer struct {
	sqlPrefix   string
	quotes      string
	parentheses string
}

func ParameterOptions() *ParameterTransformer {
	return new(ParameterTransformer)
}

func (v *ParameterTransformer) SQL(sqlPrefix string) *ParameterTransformer {
	v.sqlPrefix = sqlPrefix
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
	if len(v.sqlPrefix) != 0 {
		addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	}
	if len(v.quotes) != 0 {
		addTagIfMissing(f.Tags, "ddl", v.quotes)
	}
	if len(v.parentheses) != 0 {
		addTagIfMissing(f.Tags, "ddl", v.parentheses)
	}
	return f
}

type ListTransformer struct {
	sqlPrefix   string
	parentheses string
}

func ListOptions() *ListTransformer {
	return new(ListTransformer)
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
	if len(v.sqlPrefix) != 0 {
		addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	}
	if len(v.parentheses) != 0 {
		addTagIfMissing(f.Tags, "ddl", v.parentheses)
	}
	return f
}

func addTagIfMissing(m map[string][]string, key string, value string) {
	if val, ok := m[key]; ok {
		if !slices.Contains(val, value) {
			m[key] = append(val, value)
		}
	} else {
		m[key] = []string{value}
	}
}
