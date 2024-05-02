package docs

import (
	"fmt"
	"regexp"
)

// deprecationMessageRegex is the message that should be used in resource/datasource DeprecationMessage to get a nice link in the documentation to the replacing resource.
var deprecationMessageRegex = regexp.MustCompile(`Please use (snowflake_(\w+)) instead.`)

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
