package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	r := require.New(t)
	u := snowflake.User("user1")
	r.NotNil(u)

	q := u.Show()
	r.Equal("SHOW USERS LIKE 'user1'", q)

	q = u.Drop()
	r.Equal(`DROP USER "user1"`, q)

	q = u.Rename("user2")
	r.Equal(`ALTER USER "user1" RENAME TO "user2"`, q)

	ab := u.Alter()
	r.NotNil(ab)

	ab.SetString(`foo`, `bar`)
	q = ab.Statement()

	r.Equal(`ALTER USER "user1" SET FOO='bar'`, q)

	ab.SetBool(`bam`, false)
	q = ab.Statement()

	r.Equal(`ALTER USER "user1" SET FOO='bar' BAM=false`, q)

	c := u.Create()
	c.SetString("foo", "bar")
	c.SetBool("bam", false)
	q = c.Statement()
	r.Equal(`CREATE USER "user1" FOO='bar' BAM=false`, q)
}
