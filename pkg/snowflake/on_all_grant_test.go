package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestOnAllTableGrant(t *testing.T) {
	a := assert.New(t)
	fvg := snowflake.OnAllTableGrant("test_db", "PUBLIC")
	a.Equal(fvg.Name(), "PUBLIC")

	s := fvg.Show()
	//a.Equal(`SHOW GRANTS ON SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON ALL TABLES IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	s = fvg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON ALL TABLES IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`, s)
}

func TestOnAllViewGrant(t *testing.T) {
	a := assert.New(t)
	fvg := snowflake.OnAllViewGrant("test_db", "PUBLIC")
	a.Equal(fvg.Name(), "PUBLIC")

	s := fvg.Show()
	//a.Equal(`SHOW GRANTS ON SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE")
	a.Equal(`GRANT USAGE ON ALL VIEWS IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	s = fvg.Role("bob").Revoke("USAGE")
	a.Equal(`REVOKE USAGE ON ALL VIEWS IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`, s)
}

func TestShowOnAllGrantsInSchema(t *testing.T) {
	a := assert.New(t)
	s := snowflake.OnAllTableGrant("test_db", "PUBLIC").Role("testRole").Show()
	a.Equal(`SHOW GRANTS TO ROLE testRole`, s)

	s = snowflake.OnAllViewGrant("test_db", "PUBLIC").Role("testRole").Show()
	a.Equal(`SHOW GRANTS TO ROLE testRole`, s)
}
