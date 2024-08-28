package resources

import (
	"fmt"
	"strings"
)

func possibleValuesListed[T ~string | ~int](values []T) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%v`", value)
	}
	return strings.Join(valuesWrapped, " | ")
}

func characterList(values []rune) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%c`", value)
	}
	return strings.Join(valuesWrapped, ", ")
}

func booleanStringFieldDescription(description string) string {
	return fmt.Sprintf(`%s Available options are: "%s" or "%s". When the value is not set in the configuration the provider will put "%s" there which means to use the Snowflake default for this value.`, description, BooleanTrue, BooleanFalse, BooleanDefault)
}

func externalChangesNotDetectedFieldDescription(description string) string {
	return fmt.Sprintf(`%s External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".`, description)
}

func withPrivilegedRolesDescription(description, paramName string) string {
	return fmt.Sprintf(`%s By default, this list includes the ACCOUNTADMIN, ORGADMIN and SECURITYADMIN roles. To remove these privileged roles from the list, use the ALTER ACCOUNT command to set the %s account parameter to FALSE. `, description, paramName)
}

func blocklistedCharactersFieldDescription(description string) string {
	return fmt.Sprintf(`%s Due to technical limitations (read more [here](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/docs/technical-documentation/identifiers_rework_design_decisions.md#known-limitations-and-identifier-recommendations)), avoid using the following characters: %s`, description, characterList([]rune{'|', '.', '(', ')', '"'}))
}
