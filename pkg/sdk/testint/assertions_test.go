package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func assertThatObject(t *testing.T, objectAssert assert.InPlaceAssertionVerifier) {
	t.Helper()
	assert.AssertThatObjectWithTestClient(t, objectAssert, testClientHelper())
}
