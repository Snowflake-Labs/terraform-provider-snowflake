package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestStorageIntegration(t *testing.T) {
	r := require.New(t)
	builder := snowflake.StorageIntegration("aws")
	r.NotNil(builder)

	q := builder.Show()
	r.Equal("SHOW STORAGE INTEGRATIONS LIKE 'aws'", q)

	c := builder.Create()

	c.SetString(`type`, `EXTERNAL_STAGE`)
	c.SetStringList(`storage_allowed_locations`, []string{"s3://my-bucket/my-path/", "s3://another-bucket/"})
	c.SetBool(`enabled`, true)
	q = c.Statement()

	r.Equal(`CREATE STORAGE INTEGRATION "aws" TYPE='EXTERNAL_STAGE' STORAGE_ALLOWED_LOCATIONS=('s3://my-bucket/my-path/', 's3://another-bucket/') ENABLED=true`, q)
}
