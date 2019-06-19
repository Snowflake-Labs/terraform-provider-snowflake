package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseGrant(t *testing.T) {
	a := assert.New(t)
	dg := snowflake.DatabaseGrant("testDB")

	s := dg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON DATABASE "testDB" TO ROLE "bob"`, s)

	s = dg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON DATABASE "testDB" FROM ROLE "bob"`, s)

	s = dg.Share("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON DATABASE "testDB" TO SHARE "bob"`, s)

	s = dg.Share("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON DATABASE "testDB" FROM SHARE "bob"`, s)
}

func TestSchemaGrant(t *testing.T) {
	a := assert.New(t)
	sg := snowflake.SchemaGrant("testDB")

	s := sg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON SCHEMA "testDB" TO ROLE "bob"`, s)

	s = sg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON SCHEMA "testDB" FROM ROLE "bob"`, s)

	s = sg.Share("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON SCHEMA "testDB" TO SHARE "bob"`, s)

	s = sg.Share("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON SCHEMA "testDB" FROM SHARE "bob"`, s)
}

func TestViewGrant(t *testing.T) {
	a := assert.New(t)
	vg := snowflake.ViewGrant("testDB")

	s := vg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON VIEW "testDB" TO ROLE "bob"`, s)

	s = vg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON VIEW "testDB" FROM ROLE "bob"`, s)

	s = vg.Share("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON VIEW "testDB" TO SHARE "bob"`, s)

	s = vg.Share("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON VIEW "testDB" FROM SHARE "bob"`, s)
}
