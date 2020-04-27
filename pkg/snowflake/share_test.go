package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestShare(t *testing.T) {
	a := require.New(t)
	s := snowflake.Share("share1")
	a.NotNil(s)

	q := s.Show()
	a.Equal("SHOW SHARES LIKE 'share1'", q)

	q = s.Drop()
	a.Equal(`DROP SHARE "share1"`, q)

	q = s.Rename("share2")
	a.Equal(`ALTER SHARE "share1" RENAME TO "share2"`, q)

	ab := s.Alter()
	a.NotNil(ab)

	ab.SetString(`foo`, `bar`)
	q = ab.Statement()

	a.Equal(`ALTER SHARE "share1" SET FOO='bar'`, q)

	ab.SetBool(`bam`, false)
	q = ab.Statement()

	a.Equal(`ALTER SHARE "share1" SET FOO='bar' BAM=false`, q)

	c := s.Create()
	c.SetString("foo", "bar")
	c.SetBool("bam", false)
	q = c.Statement()
	a.Equal(`CREATE SHARE "share1" FOO='bar' BAM=false`, q)
}
