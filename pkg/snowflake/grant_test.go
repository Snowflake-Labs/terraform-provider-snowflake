package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestDatabaseGrant(t *testing.T) {
	r := require.New(t)
	dg := snowflake.DatabaseGrant("testDB")
	r.Equal(dg.Name(), "testDB")

	s := dg.Show()
	r.Equal(`SHOW GRANTS ON DATABASE "testDB"`, s)

	s = dg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON DATABASE "testDB" TO ROLE "bob"`, s)

	s = dg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON DATABASE "testDB" FROM ROLE "bob"`, s)

	s = dg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON DATABASE "testDB" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = dg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON DATABASE "testDB" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)

	s = dg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON DATABASE "testDB" TO SHARE "bob"`, s)

	s = dg.Share("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON DATABASE "testDB" FROM SHARE "bob"`, s)
}

func TestSchemaGrant(t *testing.T) {
	r := require.New(t)
	sg := snowflake.SchemaGrant("test_db", "testSchema")
	r.Equal(sg.Name(), "testSchema")

	s := sg.Show()
	r.Equal(`SHOW GRANTS ON SCHEMA "test_db"."testSchema"`, s)

	s = sg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON SCHEMA "test_db"."testSchema" TO ROLE "bob"`, s)

	s = sg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON SCHEMA "test_db"."testSchema" FROM ROLE "bob"`, s)

	s = sg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON SCHEMA "test_db"."testSchema" TO SHARE "bob"`, s)

	s = sg.Share("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON SCHEMA "test_db"."testSchema" FROM SHARE "bob"`, s)

	s = sg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON SCHEMA "test_db"."testSchema" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = sg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON SCHEMA "test_db"."testSchema" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestViewGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.ViewGrant("test_db", "PUBLIC", "testView")
	r.Equal(vg.Name(), "testView")

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON VIEW "test_db"."PUBLIC"."testView"`, s)

	s = vg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON VIEW "test_db"."PUBLIC"."testView" TO ROLE "bob"`, s)

	s = vg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON VIEW "test_db"."PUBLIC"."testView" FROM ROLE "bob"`, s)

	s = vg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON VIEW "test_db"."PUBLIC"."testView" TO SHARE "bob"`, s)

	s = vg.Share("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON VIEW "test_db"."PUBLIC"."testView" FROM SHARE "bob"`, s)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testView" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = vg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testView" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestMaterializedViewGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.MaterializedViewGrant("test_db", "PUBLIC", "testMaterializedView")
	r.Equal(vg.Name(), "testMaterializedView")

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON MATERIALIZED VIEW "test_db"."PUBLIC"."testMaterializedView"`, s)

	s = vg.Role("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO ROLE "bob"`, s)

	s = vg.Role("bob").Revoke("SELECT")
	r.Equal(`REVOKE SELECT ON VIEW "test_db"."PUBLIC"."testMaterializedView" FROM ROLE "bob"`, s)

	s = vg.Share("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO SHARE "bob"`, s)

	s = vg.Share("bob").Revoke("SELECT")
	r.Equal(`REVOKE SELECT ON VIEW "test_db"."PUBLIC"."testMaterializedView" FROM SHARE "bob"`, s)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = vg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestExternalTableGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.ExternalTableGrant("test_db", "PUBLIC", "testExternalTable")
	r.Equal(vg.Name(), "testExternalTable")

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable"`, s)

	s = vg.Role("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO ROLE "bob"`, s)

	s = vg.Role("bob").Revoke("SELECT")
	r.Equal(`REVOKE SELECT ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" FROM ROLE "bob"`, s)

	s = vg.Share("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO SHARE "bob"`, s)

	s = vg.Share("bob").Revoke("SELECT")
	r.Equal(`REVOKE SELECT ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" FROM SHARE "bob"`, s)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = vg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestFileFormatGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.FileFormatGrant("test_db", "PUBLIC", "testFileFormat")
	r.Equal(vg.Name(), "testFileFormat")

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat"`, s)

	s = vg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO ROLE "bob"`, s)

	s = vg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" FROM ROLE "bob"`, s)

	s = vg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO SHARE "bob"`, s)

	s = vg.Share("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" FROM SHARE "bob"`, s)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = vg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestFunctionGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.FunctionGrant("test_db", "PUBLIC", "testFunction", []string{"ARRAY", "STRING"})
	r.Equal(vg.Name(), "testFunction")

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING)`, s)

	s = vg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO ROLE "bob"`, s)

	s = vg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) FROM ROLE "bob"`, s)

	s = vg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO SHARE "bob"`, s)

	s = vg.Share("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) FROM SHARE "bob"`, s)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = vg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestProcedureGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.ProcedureGrant("test_db", "PUBLIC", "testProcedure", []string{"ARRAY", "STRING"})
	r.Equal(vg.Name(), "testProcedure")

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING)`, s)

	s = vg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO ROLE "bob"`, s)

	s = vg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) FROM ROLE "bob"`, s)

	s = vg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO SHARE "bob"`, s)

	s = vg.Share("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) FROM SHARE "bob"`, s)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = vg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestWarehouseGrant(t *testing.T) {
	r := require.New(t)
	wg := snowflake.WarehouseGrant("test_warehouse")
	r.Equal(wg.Name(), "test_warehouse")

	s := wg.Show()
	r.Equal(`SHOW GRANTS ON WAREHOUSE "test_warehouse"`, s)

	s = wg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON WAREHOUSE "test_warehouse" TO ROLE "bob"`, s)

	s = wg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON WAREHOUSE "test_warehouse" FROM ROLE "bob"`, s)

	s = wg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON WAREHOUSE "test_warehouse" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = wg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON WAREHOUSE "test_warehouse" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)

}

// lintignore:AT003
func TestAccountGrant(t *testing.T) {
	r := require.New(t)
	wg := snowflake.AccountGrant()
	r.Equal(wg.Name(), "")

	// There's an extra space after "ACCOUNT"
	//  because accounts don't have names

	s := wg.Show()
	r.Equal("SHOW GRANTS ON ACCOUNT ", s)

	s = wg.Role("bob").Grant("MANAGE GRANTS", false)
	r.Equal(`GRANT MANAGE GRANTS ON ACCOUNT  TO ROLE "bob"`, s)

	s = wg.Role("bob").Revoke("MONITOR USAGE")
	r.Equal(`REVOKE MONITOR USAGE ON ACCOUNT  FROM ROLE "bob"`, s)
}

func TestIntegrationGrant(t *testing.T) {
	r := require.New(t)
	wg := snowflake.IntegrationGrant("test_integration")
	r.Equal(wg.Name(), "test_integration")

	s := wg.Show()
	r.Equal(`SHOW GRANTS ON INTEGRATION "test_integration"`, s)

	s = wg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON INTEGRATION "test_integration" TO ROLE "bob"`, s)

	s = wg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON INTEGRATION "test_integration" FROM ROLE "bob"`, s)

	s = wg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON INTEGRATION "test_integration" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = wg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON INTEGRATION "test_integration" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestResourceMonitorGrant(t *testing.T) {
	r := require.New(t)
	wg := snowflake.ResourceMonitorGrant("test_monitor")
	r.Equal(wg.Name(), "test_monitor")

	s := wg.Show()
	r.Equal(`SHOW GRANTS ON RESOURCE MONITOR "test_monitor"`, s)

	s = wg.Role("bob").Grant("MONITOR", false)
	r.Equal(`GRANT MONITOR ON RESOURCE MONITOR "test_monitor" TO ROLE "bob"`, s)

	s = wg.Role("bob").Revoke("MODIFY")
	r.Equal(`REVOKE MODIFY ON RESOURCE MONITOR "test_monitor" FROM ROLE "bob"`, s)

	s = wg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON RESOURCE MONITOR "test_monitor" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	s = wg.Role("bob").Revoke("OWNERSHIP")
	r.Equal(`SET currentRole=CURRENT_ROLE(); GRANT OWNERSHIP ON RESOURCE MONITOR "test_monitor" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`, s)
}

func TestShowGrantsOf(t *testing.T) {
	r := require.New(t)
	s := snowflake.ViewGrant("test_db", "PUBLIC", "testView").Role("testRole").Show()
	r.Equal(`SHOW GRANTS OF ROLE "testRole"`, s)

	s = snowflake.ViewGrant("test_db", "PUBLIC", "testView").Share("testShare").Show()
	r.Equal(`SHOW GRANTS OF SHARE "testShare"`, s)
}
