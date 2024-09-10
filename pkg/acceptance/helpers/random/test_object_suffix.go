package random

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
)

// TODO [SNOW-1356199]: add generation tests
// TODO [SNOW-1356199]: use the same fallback suffix for acceptance and integration tests (now two different ones are generated if the env is missing)
var (
	AcceptanceTestsSuffix  = acceptanceTestsSuffix()
	IntegrationTestsSuffix = integrationTestsSuffix()
)

func acceptanceTestsSuffix() string {
	suffix := "AT_" + objectSuffix()
	log.Printf("[DEBUG] Suffix for the given test run is: %s", suffix)
	return suffix
}

func integrationTestsSuffix() string {
	suffix := "IT_" + objectSuffix()
	log.Printf("[DEBUG] Suffix for the given test run is: %s", suffix)
	return suffix
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
