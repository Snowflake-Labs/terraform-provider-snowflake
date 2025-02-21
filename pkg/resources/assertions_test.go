package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func assertThatObject(t *testing.T, objectAssert assert.InPlaceAssertionVerifier) {
	t.Helper()
	assert.AssertThatObject(t, objectAssert, acc.TestClient())
}
