package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {
	a := assert.New(t)

	a.Equal(`\'`, snowflake.EscapeString(`'`))
	a.Equal(`\\\'`, snowflake.EscapeString(`\'`))
}
