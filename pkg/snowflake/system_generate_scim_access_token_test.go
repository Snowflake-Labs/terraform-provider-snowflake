package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemGenerateSCIMAccessToken(t *testing.T) {
	r := require.New(t)
	sb := NewSystemGenerateSCIMAccessTokenBuilder("AAD_PROVISIONING")

	r.Equal(`SELECT SYSTEM$GENERATE_SCIM_ACCESS_TOKEN('AAD_PROVISIONING') AS "TOKEN"`, sb.Select())
}
