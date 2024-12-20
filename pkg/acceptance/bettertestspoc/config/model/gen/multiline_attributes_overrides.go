package gen

var multilineAttributesOverrides = map[string][]string{
	"User":                {"rsa_public_key", "rsa_public_key_2"},
	"ServiceUser":         {"rsa_public_key", "rsa_public_key_2"},
	"LegacyServiceUser":   {"rsa_public_key", "rsa_public_key_2"},
	"FunctionJava":        {"function_definition"},
	"FunctionJavascript":  {"function_definition"},
	"FunctionPython":      {"function_definition"},
	"FunctionScala":       {"function_definition"},
	"FunctionSql":         {"function_definition"},
	"ProcedureJava":       {"procedure_definition"},
	"ProcedureJavascript": {"procedure_definition"},
	"ProcedurePython":     {"procedure_definition"},
	"ProcedureScala":      {"procedure_definition"},
	"ProcedureSql":        {"procedure_definition"},
}
