package gen

import (
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)

	//go:embed templates/definition.tmpl
	definitionTemplateContent string
	DefinitionTemplate, _     = template.New("definitionTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
	)).Parse(definitionTemplateContent)

	// TODO: builders temple reused for now
	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate, gen.BuildersTemplate}
)
