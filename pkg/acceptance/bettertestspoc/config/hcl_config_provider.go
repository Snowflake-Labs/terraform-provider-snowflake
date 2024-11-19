package config

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	hclv1printer "github.com/hashicorp/hcl/hcl/printer"
	hclv1parser "github.com/hashicorp/hcl/json/parser"
)

var DefaultHclProvider = NewHclV1ConfigProvider(unquoteBlockType, fixBlockArguments, unquoteArguments, unquoteArguments, removeDoubleNewlines, unquoteDependsOnReferences)

// HclProvider defines methods to generate .tf config from .tf.json configs.
type HclProvider interface {
	HclFromJson(json []byte) (string, error)
}

type HclFormatter func(string) (string, error)

type hclV1ConfigProvider struct {
	formatters []HclFormatter
}

func NewHclV1ConfigProvider(formatters ...HclFormatter) HclProvider {
	return &hclV1ConfigProvider{
		formatters: formatters,
	}
}

func (h *hclV1ConfigProvider) HclFromJson(json []byte) (string, error) {
	hcl, err := convertJsonToHclStringV1(json)
	if err != nil {
		return "", err
	}

	for _, formatter := range h.formatters {
		hcl, err = formatter(hcl)
		if err != nil {
			return "", err
		}
	}

	return hcl, nil
}

func convertJsonToHclStringV1(jsonBytes []byte) (string, error) {
	parsed, err := hclv1parser.Parse(jsonBytes)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = hclv1printer.Fprint(&buffer, parsed)
	if err != nil {
		return "", err
	}

	formatted, err := hclv1printer.Format(buffer.Bytes())
	if err != nil {
		return "", err
	}

	return string(formatted[:]), nil
}

// Conversion to HCL using hcl v1 does not unquote block types (i.e. `"resource"` instead of expected `resource`).
// Check experiments subpackage for details.
func unquoteBlockType(s string) (string, error) {
	blockTypeRegex := regexp.MustCompile(`"(resource|data|provider)"(( "\w+"){1,2} {)`)
	return blockTypeRegex.ReplaceAllString(s, `$1$2`), nil
}

// Conversion to HCL using hcl v1 uses `=` sign for objects and lists of objects.
// Check experiments subpackage for details.
func fixBlockArguments(s string) (string, error) {
	argumentRegex := regexp.MustCompile(`( +)"(\w+)"( +)= ({\n)`)
	return argumentRegex.ReplaceAllString(s, `$1$2$3$4`), nil
}

// Conversion to HCL using hcl v1  does not unquote arguments.
// Check experiments subpackage for details.
func unquoteArguments(s string) (string, error) {
	argumentRegex := regexp.MustCompile(`( +)"(\w+)"( +=)`)
	return argumentRegex.ReplaceAllString(s, `$1$2$3`), nil
}

// Conversion to HCL using hcl v1 leaves double newlines between each attribute.
// Check experiments subpackage for details.
func removeDoubleNewlines(input string) (string, error) {
	return fmt.Sprintf("%s", strings.ReplaceAll(input, "\n\n", "\n")), nil
}

// Based on https://developer.hashicorp.com/terraform/language/syntax/json#depends_on should be processed in a special way, but it isn't.
// Check experiments subpackage for details.
// TODO [SNOW-1501905]: use FindAllStringSubmatch or fix this method in other way
func unquoteDependsOnReferences(s string) (string, error) {
	dependsOnRegex := regexp.MustCompile(`("?depends_on"? = )(\["\w+\.\w+(.\w+)?"(, "\w+\.\w+(.\w+)?")*])`)
	submatches := dependsOnRegex.FindStringSubmatch(s)
	if len(submatches) < 2 {
		return s, nil
	} else {
		withoutQuotes := strings.ReplaceAll(submatches[2], `"`, "")
		return dependsOnRegex.ReplaceAllString(s, fmt.Sprintf(`$1%s`, withoutQuotes)), nil
	}
}
