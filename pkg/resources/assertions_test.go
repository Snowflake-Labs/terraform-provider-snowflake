package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func assertThatObject(t *testing.T, objectAssert assert.InPlaceAssertionVerifier) {
	t.Helper()
	assert.AssertThatObject(t, objectAssert, acc.TestClient())
}

func assertThatImport(t *testing.T, fs ...assert.ImportStateCheckFuncProvider) resource.ImportStateCheckFunc {
	t.Helper()
	return assert.AssertThatImport(t, acc.TestClient(), fs...)
}
