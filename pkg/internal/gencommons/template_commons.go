package gencommons

import (
	"strings"
	"text/template"
)

func firstLetterLowercase(in string) string {
	return strings.ToLower(in[:1]) + in[1:]
}

func runMapper(mapper Mapper, in ...string) string {
	return mapper(strings.Join(in, ""))
}

var FirstLetterLowercaseEntry = map[string]any{
	"firstLetterLowercase": firstLetterLowercase,
}

var RunMapperEntry = map[string]any{
	"runMapper": runMapper,
}

func MergeFuncsMap(funcs ...map[string]any) template.FuncMap {
	var allFuncs = make(map[string]any)
	for _, f := range funcs {
		for k, v := range f {
			allFuncs[k] = v
		}
	}
	return allFuncs
}
