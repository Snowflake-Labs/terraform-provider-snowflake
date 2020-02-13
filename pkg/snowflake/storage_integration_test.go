package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestStorageIntegration(t *testing.T) {
	a := assert.New(t)
	builder := snowflake.StorageIntegration("aws")
	a.NotNil(builder)

	q := builder.Show()
	a.Equal("SHOW STORAGE INTEGRATIONS LIKE 'aws'", q)

	c := builder.Create()

	c.SetString(`type`, `EXTERNAL_STAGE`)
	c.SetBool(`enabled`, true)
	q = c.Statement()

	a.Equal(`CREATE STORAGE INTEGRATION "aws" TYPE='EXTERNAL_STAGE' ENABLED=true`, q)
}
