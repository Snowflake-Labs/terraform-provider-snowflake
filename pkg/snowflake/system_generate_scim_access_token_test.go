package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemGenerateSCIMAccessToken(t *testing.T) {
	r := require.New(t)
	sb := SystemGenerateSCIMAccessToken("AAD_PROVISIONING")

	r.Equal(sb.Select(), `SELECT SYSTEM$GENERATE_SCIM_ACCESS_TOKEN('AAD_PROVISIONING') AS "token"`)
}
