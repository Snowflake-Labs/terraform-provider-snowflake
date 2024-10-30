package docs

import (
	"fmt"
	"regexp"
	"strings"
)

// deprecationMessageRegex is the message that should be used in resource/datasource DeprecationMessage to get a nice link in the documentation to the replacing resource.
var deprecationMessageRegex = regexp.MustCompile(`Please use (snowflake_(\w+)) instead.`)

// TODO(SNOW-1465227): Should detect more than one replacements
// GetDeprecatedResourceReplacement allows us to get resource replacement based on the regex deprecationMessageRegex
func GetDeprecatedResourceReplacement(deprecationMessage string) (replacement string, replacementPage string, ok bool) {
	resourceReplacement := deprecationMessageRegex.FindStringSubmatch(deprecationMessage)
	if len(resourceReplacement) == 3 {
		return resourceReplacement[1], resourceReplacement[2], true
	} else {
		return "", "", false
	}
}

// RelativeLink allows us to get relative link to the resource/datasource in the same subtree. Will have to change when we introduce subcategories.
func RelativeLink(title string, path string) string {
	return fmt.Sprintf(`[%s](./%s)`, title, path)
}

func PossibleValuesListed[T ~string | ~int](values []T) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%v`", value)
	}
	return strings.Join(valuesWrapped, " | ")
}
