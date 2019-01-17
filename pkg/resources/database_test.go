package resources_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
)

func TestDatabase(t *testing.T) {
	resources.Database().InternalValidate(provider.Provider().Schema, false)
}
