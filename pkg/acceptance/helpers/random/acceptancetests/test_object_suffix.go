package acceptancetests

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

var ObjectsSuffix = acceptanceTestsSuffix()

func acceptanceTestsSuffix() string {
	suffix := "AT_" + random.ObjectSuffix()
	log.Printf("[DEBUG] Suffix for the given test run is: %s", suffix)
	return suffix
}
