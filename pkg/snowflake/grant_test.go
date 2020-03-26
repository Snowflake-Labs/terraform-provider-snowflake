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

	s = dg.Role("bob").Grant("OWNERSHIP")
	a.Equal(`GRANT OWNERSHIP ON DATABASE "testDB" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = dg.Share("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON DATABASE "testDB" TO SHARE "bob"`, s)

	s = dg.Share("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON DATABASE "testDB" FROM SHARE "bob"`, s)
}

func TestSchemaGrant(t *testing.T) {
	a := assert.New(t)
	sg := snowflake.SchemaGrant("test_db", "testSchema")
	a.Equal(sg.Name(), "testSchema")

	s := sg.Show()
	a.Equal(`SHOW GRANTS ON SCHEMA "test_db"."testSchema"`, s)

	s = sg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON SCHEMA "test_db"."testSchema" TO ROLE "bob"`, s)

	s = sg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON SCHEMA "test_db"."testSchema" FROM ROLE "bob"`, s)

	s = sg.Share("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON SCHEMA "test_db"."testSchema" TO SHARE "bob"`, s)

	s = sg.Share("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON SCHEMA "test_db"."testSchema" FROM SHARE "bob"`, s)

	s = sg.Role("bob").Grant("OWNERSHIP")
	a.Equal(`GRANT OWNERSHIP ON SCHEMA "test_db"."testSchema" TO ROLE "bob" COPY CURRENT GRANTS`, s)
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

	s = vg.Role("bob").Grant("OWNERSHIP")
	a.Equal(`GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testView" TO ROLE "bob" COPY CURRENT GRANTS`, s)
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

	s = wg.Role("bob").Grant("OWNERSHIP")
	a.Equal(`GRANT OWNERSHIP ON WAREHOUSE "test_warehouse" TO ROLE "bob" COPY CURRENT GRANTS`, s)

}

func TestAccountGrant(t *testing.T) {
	a := assert.New(t)
	wg := snowflake.AccountGrant()
	a.Equal(wg.Name(), "")

	// There's an extra space after "ACCOUNT"
	//  because accounts don't have names

	s := wg.Show()
	a.Equal("SHOW GRANTS ON ACCOUNT ", s)

	s = wg.Role("bob").Grant("MANAGE GRANTS")
	a.Equal(`GRANT MANAGE GRANTS ON ACCOUNT  TO ROLE "bob"`, s)

	s = wg.Role("bob").Revoke("MONITOR USAGE")
	a.Equal(`REVOKE MONITOR USAGE ON ACCOUNT  FROM ROLE "bob"`, s)
}

func TestIntegrationGrant(t *testing.T) {
	a := assert.New(t)
	wg := snowflake.IntegrationGrant("test_integration")
	a.Equal(wg.Name(), "test_integration")

	s := wg.Show()
	a.Equal(`SHOW GRANTS ON INTEGRATION "test_integration"`, s)

	s = wg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON INTEGRATION "test_integration" TO ROLE "bob"`, s)

	s = wg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON INTEGRATION "test_integration" FROM ROLE "bob"`, s)

	s = wg.Role("bob").Grant("OWNERSHIP")
	a.Equal(`GRANT OWNERSHIP ON INTEGRATION "test_integration" TO ROLE "bob" COPY CURRENT GRANTS`, s)
}

func TestResourceMonitorGrant(t *testing.T) {
	a := assert.New(t)
	wg := snowflake.ResourceMonitorGrant("test_monitor")
	a.Equal(wg.Name(), "test_monitor")

	s := wg.Show()
	a.Equal(`SHOW GRANTS ON RESOURCE MONITOR "test_monitor"`, s)

	s = wg.Role("bob").Grant("MONITOR")
	a.Equal(`GRANT MONITOR ON RESOURCE MONITOR "test_monitor" TO ROLE "bob"`, s)

	s = wg.Role("bob").Revoke("MODIFY")
	a.Equal(`REVOKE MODIFY ON RESOURCE MONITOR "test_monitor" FROM ROLE "bob"`, s)
}

func TestShowGrantsOf(t *testing.T) {
	a := assert.New(t)
	s := snowflake.ViewGrant("test_db", "PUBLIC", "testView").Role("testRole").Show()
	a.Equal(`SHOW GRANTS OF ROLE "testRole"`, s)

	s = snowflake.ViewGrant("test_db", "PUBLIC", "testView").Share("testShare").Show()
	a.Equal(`SHOW GRANTS OF SHARE "testShare"`, s)
}
