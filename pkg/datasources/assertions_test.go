package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func assertThat(t *testing.T, fs ...assert.TestCheckFuncProvider) resource.TestCheckFunc {
	t.Helper()
	return assert.AssertThat(t, acc.TestClient(), fs...)
}
