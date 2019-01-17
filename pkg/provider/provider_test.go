package provider_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	a := assert.New(t)
	p := provider.Provider()
	err := p.InternalValidate()
	a.NoError(err)
}
