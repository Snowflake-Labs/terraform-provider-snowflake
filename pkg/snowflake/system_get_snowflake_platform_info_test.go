package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemGetSnowflakePlatformInfoQuery(t *testing.T) {
	r := require.New(t)
	sb := SystemGetSnowflakePlatformInfoQuery()

	r.Equal(sb, `SELECT SYSTEM$GET_SNOWFLAKE_PLATFORM_INFO() AS "INFO"`)
}

func TestSystemGetSnowflakePlatformInfoGetStructuredConfigAws(t *testing.T) {
	r := require.New(t)

	raw := &RawSnowflakePlatformInfo{
		Info: `{"snowflake-vpc-id": ["vpc-1", "vpc-2"]}`,
	}

	c, e := raw.GetStructuredConfig()
	r.Nil(e)

	r.Equal([]string{"vpc-1", "vpc-2"}, c.AwsVpcIds)
	r.Equal([]string(nil), c.AzureVnetSubnetIds)
}

func TestSystemGetSnowflakePlatformInfoGetStructuredConfigAzure(t *testing.T) {
	r := require.New(t)

	raw := &RawSnowflakePlatformInfo{
		Info: `{"snowflake-vnet-subnet-id": ["/subscription/1/1", "/subscription/1/2"]}`,
	}

	c, e := raw.GetStructuredConfig()
	r.Nil(e)

	r.Equal([]string{"/subscription/1/1", "/subscription/1/2"}, c.AzureVnetSubnetIds)
	r.Equal([]string(nil), c.AwsVpcIds)
}
