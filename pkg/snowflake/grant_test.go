package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestDatabaseGrant(t *testing.T) {
	r := require.New(t)
	dg := snowflake.DatabaseGrant("testDB")
	r.Equal("testDB", dg.Name())

	s := dg.Show()
	r.Equal(`SHOW GRANTS ON DATABASE "testDB"`, s)

	s = dg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON DATABASE "testDB" TO ROLE "bob"`, s)

	revoke := dg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON DATABASE "testDB" FROM ROLE "bob"`}, revoke)

	s = dg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON DATABASE "testDB" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = dg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON DATABASE "testDB" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = dg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON DATABASE "testDB" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)

	s = dg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON DATABASE "testDB" TO SHARE "bob"`, s)

	revoke = dg.Share("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON DATABASE "testDB" FROM SHARE "bob"`}, revoke)
}

func TestSchemaGrant(t *testing.T) {
	r := require.New(t)
	sg := snowflake.SchemaGrant("test_db", "testSchema")
	r.Equal("testSchema", sg.Name())

	s := sg.Show()
	r.Equal(`SHOW GRANTS ON SCHEMA "test_db"."testSchema"`, s)

	s = sg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON SCHEMA "test_db"."testSchema" TO ROLE "bob"`, s)

	revoke := sg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON SCHEMA "test_db"."testSchema" FROM ROLE "bob"`}, revoke)

	s = sg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON SCHEMA "test_db"."testSchema" TO SHARE "bob"`, s)

	revoke = sg.Share("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON SCHEMA "test_db"."testSchema" FROM SHARE "bob"`}, revoke)

	s = sg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON SCHEMA "test_db"."testSchema" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = sg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON SCHEMA "test_db"."testSchema" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = sg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON SCHEMA "test_db"."testSchema" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestViewGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.ViewGrant("test_db", "PUBLIC", "testView")
	r.Equal("testView", vg.Name())

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON VIEW "test_db"."PUBLIC"."testView"`, s)

	s = vg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON VIEW "test_db"."PUBLIC"."testView" TO ROLE "bob"`, s)

	revoke := vg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON VIEW "test_db"."PUBLIC"."testView" FROM ROLE "bob"`}, revoke)

	s = vg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON VIEW "test_db"."PUBLIC"."testView" TO SHARE "bob"`, s)

	revoke = vg.Share("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON VIEW "test_db"."PUBLIC"."testView" FROM SHARE "bob"`}, revoke)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testView" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = vg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testView" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = vg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testView" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestMaterializedViewGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.MaterializedViewGrant("test_db", "PUBLIC", "testMaterializedView")
	r.Equal("testMaterializedView", vg.Name())

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON MATERIALIZED VIEW "test_db"."PUBLIC"."testMaterializedView"`, s)

	s = vg.Role("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO ROLE "bob"`, s)

	revoke := vg.Role("bob").Revoke("SELECT")
	r.Equal([]string{`REVOKE SELECT ON VIEW "test_db"."PUBLIC"."testMaterializedView" FROM ROLE "bob"`}, revoke)

	s = vg.Share("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO SHARE "bob"`, s)

	revoke = vg.Share("bob").Revoke("SELECT")
	r.Equal([]string{`REVOKE SELECT ON VIEW "test_db"."PUBLIC"."testMaterializedView" FROM SHARE "bob"`}, revoke)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = vg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = vg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON VIEW "test_db"."PUBLIC"."testMaterializedView" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestExternalTableGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.ExternalTableGrant("test_db", "PUBLIC", "testExternalTable")
	r.Equal("testExternalTable", vg.Name())

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable"`, s)

	s = vg.Role("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO ROLE "bob"`, s)

	revoke := vg.Role("bob").Revoke("SELECT")
	r.Equal([]string{`REVOKE SELECT ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" FROM ROLE "bob"`}, revoke)

	s = vg.Share("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO SHARE "bob"`, s)

	revoke = vg.Share("bob").Revoke("SELECT")
	r.Equal([]string{`REVOKE SELECT ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" FROM SHARE "bob"`}, revoke)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = vg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = vg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON EXTERNAL TABLE "test_db"."PUBLIC"."testExternalTable" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestFileFormatGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.FileFormatGrant("test_db", "PUBLIC", "testFileFormat")
	r.Equal("testFileFormat", vg.Name())

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat"`, s)

	s = vg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO ROLE "bob"`, s)

	revoke := vg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" FROM ROLE "bob"`}, revoke)

	s = vg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO SHARE "bob"`, s)

	revoke = vg.Share("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" FROM SHARE "bob"`}, revoke)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = vg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = vg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON FILE FORMAT "test_db"."PUBLIC"."testFileFormat" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestFunctionGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.FunctionGrant("test_db", "PUBLIC", "testFunction", []string{"ARRAY", "STRING"})
	r.Equal("testFunction", vg.Name())

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING)`, s)

	s = vg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO ROLE "bob"`, s)

	revoke := vg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) FROM ROLE "bob"`}, revoke)

	s = vg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO SHARE "bob"`, s)

	revoke = vg.Share("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) FROM SHARE "bob"`}, revoke)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = vg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = vg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON FUNCTION "test_db"."PUBLIC"."testFunction"(ARRAY, STRING) TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestProcedureGrant(t *testing.T) {
	r := require.New(t)
	vg := snowflake.ProcedureGrant("test_db", "PUBLIC", "testProcedure", []string{"ARRAY", "STRING"})
	r.Equal("testProcedure", vg.Name())

	s := vg.Show()
	r.Equal(`SHOW GRANTS ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING)`, s)

	s = vg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO ROLE "bob"`, s)

	revoke := vg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) FROM ROLE "bob"`}, revoke)

	s = vg.Share("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO SHARE "bob"`, s)

	revoke = vg.Share("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) FROM SHARE "bob"`}, revoke)

	s = vg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = vg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = vg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON PROCEDURE "test_db"."PUBLIC"."testProcedure"(ARRAY, STRING) TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestWarehouseGrant(t *testing.T) {
	r := require.New(t)
	wg := snowflake.WarehouseGrant("test_warehouse")
	r.Equal("test_warehouse", wg.Name())

	s := wg.Show()
	r.Equal(`SHOW GRANTS ON WAREHOUSE "test_warehouse"`, s)

	s = wg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON WAREHOUSE "test_warehouse" TO ROLE "bob"`, s)

	revoke := wg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON WAREHOUSE "test_warehouse" FROM ROLE "bob"`}, revoke)

	s = wg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON WAREHOUSE "test_warehouse" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = wg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON WAREHOUSE "test_warehouse" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = wg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON WAREHOUSE "test_warehouse" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

// lintignore:AT003
func TestAccountGrant(t *testing.T) {
	r := require.New(t)
	wg := snowflake.AccountGrant()
	r.Equal("", wg.Name())

	// There's an extra space after "ACCOUNT"
	//  because accounts don't have names

	s := wg.Show()
	r.Equal("SHOW GRANTS ON ACCOUNT ", s)

	s = wg.Role("bob").Grant("MANAGE GRANTS", false)
	r.Equal(`GRANT MANAGE GRANTS ON ACCOUNT  TO ROLE "bob"`, s)

	revoke := wg.Role("bob").Revoke("MONITOR USAGE")
	r.Equal([]string{`REVOKE MONITOR USAGE ON ACCOUNT  FROM ROLE "bob"`}, revoke)

	s = wg.Role("bob").Grant("APPLY MASKING POLICY", false)
	r.Equal(`GRANT APPLY MASKING POLICY ON ACCOUNT  TO ROLE "bob"`, s)
}

func TestIntegrationGrant(t *testing.T) {
	r := require.New(t)
	wg := snowflake.IntegrationGrant("test_integration")
	r.Equal("test_integration", wg.Name())

	s := wg.Show()
	r.Equal(`SHOW GRANTS ON INTEGRATION "test_integration"`, s)

	s = wg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON INTEGRATION "test_integration" TO ROLE "bob"`, s)

	revoke := wg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON INTEGRATION "test_integration" FROM ROLE "bob"`}, revoke)

	s = wg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON INTEGRATION "test_integration" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = wg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON INTEGRATION "test_integration" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = wg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON INTEGRATION "test_integration" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestResourceMonitorGrant(t *testing.T) {
	r := require.New(t)
	wg := snowflake.ResourceMonitorGrant("test_monitor")
	r.Equal("test_monitor", wg.Name())

	s := wg.Show()
	r.Equal(`SHOW GRANTS ON RESOURCE MONITOR "test_monitor"`, s)

	s = wg.Role("bob").Grant("MONITOR", false)
	r.Equal(`GRANT MONITOR ON RESOURCE MONITOR "test_monitor" TO ROLE "bob"`, s)

	revoke := wg.Role("bob").Revoke("MODIFY")
	r.Equal([]string{`REVOKE MODIFY ON RESOURCE MONITOR "test_monitor" FROM ROLE "bob"`}, revoke)

	s = wg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON RESOURCE MONITOR "test_monitor" TO ROLE "bob" COPY CURRENT GRANTS`, s)

	revoke = wg.Role("bob").RevokeOwnership("revertBob")
	r.Equal([]string{`GRANT OWNERSHIP ON RESOURCE MONITOR "test_monitor" TO ROLE "revertBob" COPY CURRENT GRANTS`}, revoke)

	revoke = wg.Role("bob").RevokeOwnership("")
	r.Equal([]string{`SET currentRole=CURRENT_ROLE()`, `GRANT OWNERSHIP ON RESOURCE MONITOR "test_monitor" TO ROLE IDENTIFIER($currentRole) COPY CURRENT GRANTS`}, revoke)
}

func TestMaskingPolicyGrant(t *testing.T) {
	r := require.New(t)
	mg := snowflake.MaskingPolicyGrant("test_db", "PUBLIC", "testMaskingPolicy")
	r.Equal("testMaskingPolicy", mg.Name())

	s := mg.Show()
	r.Equal(`SHOW GRANTS ON MASKING POLICY "test_db"."PUBLIC"."testMaskingPolicy"`, s)

	s = mg.Role("bob").Grant("APPLY", false)
	r.Equal(`GRANT APPLY ON MASKING POLICY "test_db"."PUBLIC"."testMaskingPolicy" TO ROLE "bob"`, s)

	revoke := mg.Role("bob").Revoke("APPLY")
	r.Equal([]string{`REVOKE APPLY ON MASKING POLICY "test_db"."PUBLIC"."testMaskingPolicy" FROM ROLE "bob"`}, revoke)
}

func TestShowGrantsOf(t *testing.T) {
	r := require.New(t)
	s := snowflake.ViewGrant("test_db", "PUBLIC", "testView").Role("testRole").Show()
	r.Equal(`SHOW GRANTS OF ROLE "testRole"`, s)

	s = snowflake.ViewGrant("test_db", "PUBLIC", "testView").Share("testShare").Show()
	r.Equal(`SHOW GRANTS OF SHARE "testShare"`, s)
}
