package resources_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
)

func TestWarehouse(t *testing.T) {
	resources.Warehouse().InternalValidate(provider.Provider().Schema, false)
}
