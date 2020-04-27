package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestEscapeString(t *testing.T) {
	r := require.New(t)

	r.Equal(`\'`, snowflake.EscapeString(`'`))
	r.Equal(`\\\'`, snowflake.EscapeString(`\'`))
}
