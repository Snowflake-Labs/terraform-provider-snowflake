package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestFutureSchemaGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureSchemaGrant("test_db")
	r.Equal(fvg.Name(), "test_db")

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvg.Role("bob").Grant("USAGE")
	r.Equal(`GRANT USAGE ON FUTURE SCHEMAS IN DATABASE "test_db" TO ROLE "bob"`, s)

	s = fvg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON FUTURE SCHEMAS IN DATABASE "test_db" FROM ROLE "bob"`, s)
}

func TestFutureTableGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureTableGrant("test_db", "PUBLIC")
	r.Equal(fvg.Name(), "PUBLIC")

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE")
	r.Equal(`GRANT USAGE ON FUTURE TABLES IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	s = fvg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON FUTURE TABLES IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`, s)

	b := require.New(t)
	fvgd := snowflake.FutureTableGrant("test_db", "")
	b.Equal(fvgd.Name(), "test_db")

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("USAGE")
	b.Equal(`GRANT USAGE ON FUTURE TABLES IN DATABASE "test_db" TO ROLE "bob"`, s)

	s = fvgd.Role("bob").Revoke("USAGE")
	b.Equal(`REVOKE USAGE ON FUTURE TABLES IN DATABASE "test_db" FROM ROLE "bob"`, s)
}

func TestFutureViewGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureViewGrant("test_db", "PUBLIC")
	r.Equal(fvg.Name(), "PUBLIC")

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE")
	r.Equal(`GRANT USAGE ON FUTURE VIEWS IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	s = fvg.Role("bob").Revoke("USAGE")
	r.Equal(`REVOKE USAGE ON FUTURE VIEWS IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`, s)

	b := require.New(t)
	fvgd := snowflake.FutureViewGrant("test_db", "")
	b.Equal(fvgd.Name(), "test_db")

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("USAGE")
	b.Equal(`GRANT USAGE ON FUTURE VIEWS IN DATABASE "test_db" TO ROLE "bob"`, s)

	s = fvgd.Role("bob").Revoke("USAGE")
	b.Equal(`REVOKE USAGE ON FUTURE VIEWS IN DATABASE "test_db" FROM ROLE "bob"`, s)
}

func TestShowFutureGrantsInSchema(t *testing.T) {
	r := require.New(t)
	s := snowflake.FutureTableGrant("test_db", "PUBLIC").Role("testRole").Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = snowflake.FutureTableGrant("test_db", "").Role("testRole").Show()
	r.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = snowflake.FutureViewGrant("test_db", "PUBLIC").Role("testRole").Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = snowflake.FutureViewGrant("test_db", "").Role("testRole").Show()
	r.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)
}
