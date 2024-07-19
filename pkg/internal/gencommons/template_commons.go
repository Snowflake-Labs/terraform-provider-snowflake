package gencommons

import (
	"reflect"
	"runtime"
	"strings"
	"text/template"
)

func FirstLetterLowercase(in string) string {
	return strings.ToLower(in[:1]) + in[1:]
}

func RunMapper(mapper Mapper, in ...string) string {
	return mapper(strings.Join(in, ""))
}

func BuildTemplateFuncMap(funcs ...any) template.FuncMap {
	var allFuncs = make(map[string]any)
	for _, f := range funcs {
		allFuncs[getFunctionName(f)] = f
	}
	return allFuncs
}

func getFunctionName(f any) string {
	fullFuncName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	parts := strings.Split(fullFuncName, ".")
	return parts[len(parts)-1]
}
