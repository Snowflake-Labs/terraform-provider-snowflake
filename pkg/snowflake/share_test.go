package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestShare(t *testing.T) {
	r := require.New(t)
	s := snowflake.Share("share1")
	r.NotNil(s)

	q := s.Show()
	r.Equal("SHOW SHARES LIKE 'share1'", q)

	q = s.Drop()
	r.Equal(`DROP SHARE "share1"`, q)

	q = s.Rename("share2")
	r.Equal(`ALTER SHARE "share1" RENAME TO "share2"`, q)

	ab := s.Alter()
	r.NotNil(ab)

	ab.SetString(`foo`, `bar`)
	q = ab.Statement()

	r.Equal(`ALTER SHARE "share1" SET FOO='bar'`, q)

	ab.SetBool(`bam`, false)
	q = ab.Statement()

	r.Equal(`ALTER SHARE "share1" SET FOO='bar' BAM=false`, q)

	c := s.Create()
	c.SetString("foo", "bar")
	c.SetBool("bam", false)
	q = c.Statement()
	r.Equal(`CREATE SHARE "share1" FOO='bar' BAM=false`, q)
}
