package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (e *ExternalVolumeResourceAssert) HasStorageLocationLength(len int) *ExternalVolumeResourceAssert {
	e.AddAssertion(assert.ValueSet("storage_location.#", strconv.FormatInt(int64(len), 10)))
	return e
}
