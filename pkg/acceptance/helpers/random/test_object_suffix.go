package random

import (
	"fmt"
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
)

// TODO: test generation
var (
	AcceptanceTestsSuffix  = acceptanceTestsSuffix()
	IntegrationTestsSuffix = integrationTestsSuffix()
)

func acceptanceTestsSuffix() string {
	return "AT_" + objectSuffix()
}

func integrationTestsSuffix() string {
	return "IT_" + objectSuffix()
}

func objectSuffix() string {
	var suffix string
	testObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.TestObjectsSuffix))
	if testObjectSuffix != "" {
		suffix = strings.ToUpper(testObjectSuffix)
	} else {
		uuid := UUID()
		formattedUuidSuffix := strings.ToUpper(strings.ReplaceAll(uuid, "-", "_"))
		suffix = formattedUuidSuffix
	}
	return suffix
}
