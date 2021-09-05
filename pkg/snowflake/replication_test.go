package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestReplication(t *testing.T) {
	r := require.New(t)
	s := snowflake.Replication("replication1")
	r.NotNil(s)

	q := s.Show()
	r.Equal("SHOW REPLICATION DATABASES LIKE 'replication1'", q)
}
