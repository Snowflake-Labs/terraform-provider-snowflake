package gen

import (
	"text/template"

	_ "embed"

	resourcemodel "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model/gen"

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

	//go:embed templates/marshal_json.tmpl
	marshalJsonTemplateContent string
	MarshalJsonTemplate, _     = template.New("marshalJsonTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
	)).Parse(marshalJsonTemplateContent)

	// TODO [SNOW-1501905]: consider duplicating the builders template from resource (currently same template used for datasources and provider which limits the customization possibilities for just one block type)
	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate, MarshalJsonTemplate, resourcemodel.BuildersTemplate}
)
