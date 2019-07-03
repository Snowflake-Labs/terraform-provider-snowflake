package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseGrant(t *testing.T) {
	a := assert.New(t)
	dg := snowflake.DatabaseGrant("testDB")
	a.Equal(dg.Name(), "testDB")

	s := dg.Show()
	a.Equal(`SHOW GRANTS ON DATABASE "testDB"`, s)

	s = dg.Role("bob").Grant("USAGE")
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
	sg := snowflake.SchemaGrant("testSchema")
	a.Equal(sg.Name(), "testSchema")

	s := sg.Show()
	a.Equal(`SHOW GRANTS ON SCHEMA "testSchema"`, s)

	s = sg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON SCHEMA "testSchema" TO ROLE "bob"`, s)

	s = sg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON SCHEMA "testSchema" FROM ROLE "bob"`, s)

	s = sg.Share("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON SCHEMA "testSchema" TO SHARE "bob"`, s)

	s = sg.Share("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON SCHEMA "testSchema" FROM SHARE "bob"`, s)
}

func TestViewGrant(t *testing.T) {
	a := assert.New(t)
	vg := snowflake.ViewGrant("test_db", "PUBLIC", "testView")
	a.Equal(vg.Name(), "testView")

	s := vg.Show()
	a.Equal(`SHOW GRANTS ON VIEW "test_db"."PUBLIC"."testView"`, s)

	s = vg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON VIEW "test_db"."PUBLIC"."testView" TO ROLE "bob"`, s)

	s = vg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON VIEW "test_db"."PUBLIC"."testView" FROM ROLE "bob"`, s)

	s = vg.Share("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON VIEW "test_db"."PUBLIC"."testView" TO SHARE "bob"`, s)

	s = vg.Share("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON VIEW "test_db"."PUBLIC"."testView" FROM SHARE "bob"`, s)
}

func TestWarehouseGrant(t *testing.T) {
	a := assert.New(t)
	wg := snowflake.WarehouseGrant("test_warehouse")
	a.Equal(wg.Name(), "test_warehouse")

	s := wg.Show()
	a.Equal(`SHOW GRANTS ON WAREHOUSE "test_warehouse"`, s)

	s = wg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON WAREHOUSE "test_warehouse" TO ROLE "bob"`, s)

	s = wg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON WAREHOUSE "test_warehouse" FROM ROLE "bob"`, s)
}

func TestShowGrantsOf(t *testing.T) {
	a := assert.New(t)
	s := snowflake.ViewGrant("test_db", "PUBLIC", "testView").Role("testRole").Show()
	a.Equal(`SHOW GRANTS OF ROLE "testRole"`, s)

	s = snowflake.ViewGrant("test_db", "PUBLIC", "testView").Share("testShare").Show()
	a.Equal(`SHOW GRANTS OF SHARE "testShare"`, s)
}
