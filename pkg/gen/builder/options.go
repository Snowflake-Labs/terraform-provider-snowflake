package builder

// First - can repeat options

type keywordOptions struct {
	sqlPrefix string
	quotes    string
}

func KeywordOptions() *keywordOptions {
	return new(keywordOptions)
}

func (opts *keywordOptions) SQLPrefix(sql string) *keywordOptions {
	opts.sqlPrefix = sql
	return opts
}

func (opts *keywordOptions) SingleQuotes() *keywordOptions {
	opts.quotes = "single_quotes"
	return opts
}

func (opts *keywordOptions) DoubleQuotes() *keywordOptions {
	opts.quotes = "double_quotes"
	return opts
}

func (opts *keywordOptions) NoQuotes() *keywordOptions {
	opts.quotes = "no_quotes"
	return opts
}

func (opts keywordOptions) transform(fb *FieldBuilder) *FieldBuilder {
	if len(opts.quotes) != 0 {
		fb.Tags["ddl"] = append(fb.Tags["ddl"], opts.quotes)
	}
	if len(opts.sqlPrefix) != 0 {
		fb.Tags["sql"] = append(fb.Tags["sql"], opts.sqlPrefix)
	}
	return fb
}

type parameterOptions struct {
	equals string
	quotes string
}

func ParameterOptions() *parameterOptions {
	return new(parameterOptions)
}

func (opts *parameterOptions) Equals() *parameterOptions {
	opts.equals = "equals"
	return opts
}

func (opts *parameterOptions) NoEquals() *parameterOptions {
	opts.equals = "no_equals"
	return opts
}

func (opts *parameterOptions) SingleQuotes() *parameterOptions {
	opts.quotes = "single_quotes"
	return opts
}

func (opts *parameterOptions) DoubleQuotes() *parameterOptions {
	opts.quotes = "double_quotes"
	return opts
}

func (opts *parameterOptions) NoQuotes() *parameterOptions {
	opts.quotes = "no_quotes"
	return opts
}

func (opts *parameterOptions) transform(fb *FieldBuilder) *FieldBuilder {
	if len(opts.quotes) != 0 {
		fb.Tags["ddl"] = append(fb.Tags["ddl"], opts.quotes)
	}
	if len(opts.equals) != 0 {
		fb.Tags["ddl"] = append(fb.Tags["ddl"], opts.equals)
	}
	return fb
}

type listOptions struct {
	paren string
}

func ListOptions() *listOptions {
	return new(listOptions)
}

func (opts *listOptions) Parentheses() *listOptions {
	opts.paren = "parentheses"
	return opts
}

func (opts *listOptions) NoParentheses() *listOptions {
	opts.paren = "no_parentheses"
	return opts
}

func (opts *listOptions) transform(fb *FieldBuilder) *FieldBuilder {
	if len(opts.paren) != 0 {
		fb.Tags["ddl"] = append(fb.Tags["ddl"], opts.paren)
	}
	return fb
}

// TODO: second approach - limit the number of options after given option was set

//type keywordOptions struct{}
//type keywordOptionsWithQuotes struct{}
//type keywordOptionsWithoutQuotes struct{}
//type keywordOptionsWithoutSQLPrefix struct {
//	prefix string
//}
//
//func KeywordOptions() keywordOptions {
//	return keywordOptions{}
//}
//
//func (opts keywordOptions) Quotes() keywordOptionsWithQuotes {
//	return keywordOptionsWithQuotes{}
//}
//
//func (opts keywordOptions) NoQuotes() keywordOptionsWithoutQuotes {
//	return keywordOptionsWithoutQuotes{}
//}
//
//func (opts keywordOptions) SQLPrefix(sql string) keywordOptionsWithoutSQLPrefix {
//	return keywordOptionsWithoutSQLPrefix{
//		prefix: sql,
//	}
//}
//
//func (opts keywordOptionsWithQuotes) Transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "quotes")
//	return fb
//}
//
//func (opts keywordOptionsWithoutQuotes) Transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "no_quotes")
//	return fb
//}
//
//func (opts keywordOptionsWithoutSQLPrefix) Transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["sql"] = append(fb.Tags["sql"], opts.prefix)
//	return fb
//}
//
//type parameterOptions struct{}
//
//func ParameterOptions() parameterOptions {
//	return parameterOptions{}
//}
//
//func (opts parameterOptions) Quotes() parameterOptionsWithQuotes {
//	return parameterOptionsWithQuotes{}
//}
//
//func (opts parameterOptions) NoQuotes() parameterOptionsWithNoQuotes {
//	return parameterOptionsWithNoQuotes{}
//}
//
//func (opts parameterOptions) Equals() parameterOptionsWithEquals {
//	return parameterOptionsWithEquals{}
//}
//
//func (opts parameterOptions) NoEquals() parameterOptionsWithNoEquals {
//	return parameterOptionsWithNoEquals{}
//}
//
//type parameterOptionsWithEquals struct{}
//type parameterOptionsWithNoEquals struct{}
//type parameterOptionsWithQuotes struct{}
//type parameterOptionsWithQuotesAndWithEquals struct{}
//type parameterOptionsWithQuotesAndNoEquals struct{}
//type parameterOptionsWithNoQuotes struct{}
//type parameterOptionsWithNoQuotesAndWithEquals struct{}
//type parameterOptionsWithNoQuotesAndNoEquals struct{}
//
//func (opts parameterOptionsWithQuotes) Equals() parameterOptionsWithQuotesAndWithEquals {
//	return parameterOptionsWithQuotesAndWithEquals{}
//}
//
//func (opts parameterOptionsWithQuotes) NoEquals() parameterOptionsWithQuotesAndNoEquals {
//	return parameterOptionsWithQuotesAndNoEquals{}
//}
//
//func (opts parameterOptionsWithNoQuotes) Equals() parameterOptionsWithNoQuotesAndWithEquals {
//	return parameterOptionsWithNoQuotesAndWithEquals{}
//}
//
//func (opts parameterOptionsWithNoQuotes) NoEquals() parameterOptionsWithNoQuotesAndNoEquals {
//	return parameterOptionsWithNoQuotesAndNoEquals{}
//}
//
//func (opts parameterOptionsWithEquals) Quotes() parameterOptionsWithNoQuotesAndWithEquals {
//	return parameterOptionsWithNoQuotesAndWithEquals{}
//}
//
//func (opts parameterOptionsWithEquals) NoQuotes() parameterOptionsWithNoQuotesAndWithEquals {
//	return parameterOptionsWithNoQuotesAndWithEquals{}
//}
//
//func (opts parameterOptionsWithNoEquals) Quotes() parameterOptionsWithQuotesAndNoEquals {
//	return parameterOptionsWithQuotesAndNoEquals{}
//}
//
//func (opts parameterOptionsWithNoEquals) NoQuotes() parameterOptionsWithNoQuotesAndNoEquals {
//	return parameterOptionsWithNoQuotesAndNoEquals{}
//}
//
//func (opts parameterOptionsWithQuotes) transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "quotes")
//	return fb
//}
//
//func (opts parameterOptionsWithNoQuotes) transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "no_quotes")
//	return fb
//}
//
//func (opts parameterOptionsWithEquals) transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "equals")
//	return fb
//}
//
//func (opts parameterOptionsWithNoEquals) transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "no_equals")
//	return fb
//}
//
//func (opts parameterOptionsWithQuotesAndWithEquals) transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "quotes")
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "equals")
//	return fb
//}
//
//func (opts parameterOptionsWithQuotesAndNoEquals) transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "quotes")
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "no_equals")
//	return fb
//}
//
//func (opts parameterOptionsWithNoQuotesAndWithEquals) transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "no_quotes")
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "equals")
//	return fb
//}
//
//func (opts parameterOptionsWithNoQuotesAndNoEquals) transform(fb *FieldBuilder) *FieldBuilder {
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "no_quotes")
//	fb.Tags["ddl"] = append(fb.Tags["ddl"], "no_equals")
//	return fb
//}
