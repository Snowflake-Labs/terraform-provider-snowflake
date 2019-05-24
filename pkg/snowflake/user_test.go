package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	a := assert.New(t)
	u := snowflake.User("user1")
	a.NotNil(u)

	q := u.Show()
	a.Equal("SHOW USERS LIKE 'user1'", q)

	q = u.Drop()
	a.Equal(`DROP USER "user1"`, q)

	q = u.Rename("user2")
	a.Equal(`ALTER USER "user1" RENAME TO "user2"`, q)

	ab := u.Alter()
	a.NotNil(ab)

	ab.SetString(`foo`, `bar`)
	q = ab.Statement()

	a.Equal(`ALTER USER "user1" SET FOO='bar'`, q)

	ab.SetBool(`bam`, false)
	q = ab.Statement()

	a.Equal(`ALTER USER "user1" SET FOO='bar' BAM=false`, q)

	c := u.Create()
	c.SetString("foo", "bar")
	c.SetBool("bam", false)
	q = c.Statement()
	a.Equal(`CREATE USER "user1" FOO='bar' BAM=false`, q)
}
