package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestManagedAccount(t *testing.T) {
	a := assert.New(t)
	u := snowflake.ManagedAccount("managedaccount1")
	a.NotNil(u)

	q := u.Show()
	a.Equal("SHOW MANAGED ACCOUNTS LIKE 'managedaccount1'", q)

	q = u.Drop()
	a.Equal(`DROP MANAGED ACCOUNT "managedaccount1"`, q)

	c := u.Create()
	c.SetString("foo", "bar")
	c.SetBool("bam", false)
	q = c.Statement()
	a.Equal(`CREATE MANAGED ACCOUNT "managedaccount1" FOO='bar' BAM=false`, q)
}
