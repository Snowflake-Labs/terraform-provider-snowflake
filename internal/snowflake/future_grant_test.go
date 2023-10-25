// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/stretchr/testify/require"
)

func TestFutureSchemaGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureSchemaGrant("test_db")
	r.Equal("test_db", fvg.Name())

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUTURE SCHEMAS IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke := fvg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FUTURE SCHEMAS IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)
}

func TestFutureTableGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureTableGrant("test_db", "PUBLIC")
	r.Equal("PUBLIC", fvg.Name())

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUTURE TABLES IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	revoke := fvg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FUTURE TABLES IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`}, revoke)

	b := require.New(t)
	fvgd := snowflake.FutureTableGrant("test_db", "")
	b.Equal("test_db", fvgd.Name())

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("USAGE", false)
	b.Equal(`GRANT USAGE ON FUTURE TABLES IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke = fvgd.Role("bob").Revoke("USAGE")
	b.Equal([]string{`REVOKE USAGE ON FUTURE TABLES IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)
}

func TestFutureMaterializedViewGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureMaterializedViewGrant("test_db", "PUBLIC")
	r.Equal("PUBLIC", fvg.Name())

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON FUTURE MATERIALIZED VIEWS IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	revoke := fvg.Role("bob").Revoke("SELECT")
	r.Equal([]string{`REVOKE SELECT ON FUTURE MATERIALIZED VIEWS IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`}, revoke)

	b := require.New(t)
	fvgd := snowflake.FutureMaterializedViewGrant("test_db", "")
	b.Equal("test_db", fvgd.Name())

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("SELECT", false)
	b.Equal(`GRANT SELECT ON FUTURE MATERIALIZED VIEWS IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke = fvgd.Role("bob").Revoke("SELECT")
	b.Equal([]string{`REVOKE SELECT ON FUTURE MATERIALIZED VIEWS IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)
}

func TestFutureViewGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureViewGrant("test_db", "PUBLIC")
	r.Equal("PUBLIC", fvg.Name())

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUTURE VIEWS IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	revoke := fvg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FUTURE VIEWS IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`}, revoke)

	b := require.New(t)
	fvgd := snowflake.FutureViewGrant("test_db", "")
	b.Equal("test_db", fvgd.Name())

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("USAGE", false)
	b.Equal(`GRANT USAGE ON FUTURE VIEWS IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke = fvgd.Role("bob").Revoke("USAGE")
	b.Equal([]string{`REVOKE USAGE ON FUTURE VIEWS IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)
}

func TestFutureStageGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureStageGrant("test_db", "PUBLIC")
	r.Equal("PUBLIC", fvg.Name())

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUTURE STAGES IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	revoke := fvg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FUTURE STAGES IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`}, revoke)

	b := require.New(t)
	fvgd := snowflake.FutureStageGrant("test_db", "")
	b.Equal("test_db", fvgd.Name())

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("USAGE", false)
	b.Equal(`GRANT USAGE ON FUTURE STAGES IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke = fvgd.Role("bob").Revoke("USAGE")
	b.Equal([]string{`REVOKE USAGE ON FUTURE STAGES IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)
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

func TestFutureExternalTableGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureExternalTableGrant("test_db", "PUBLIC")
	r.Equal("PUBLIC", fvg.Name())

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("SELECT", false)
	r.Equal(`GRANT SELECT ON FUTURE EXTERNAL TABLES IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	revoke := fvg.Role("bob").Revoke("SELECT")
	r.Equal([]string{`REVOKE SELECT ON FUTURE EXTERNAL TABLES IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`}, revoke)

	b := require.New(t)
	fvgd := snowflake.FutureExternalTableGrant("test_db", "")
	b.Equal("test_db", fvgd.Name())

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("SELECT", false)
	b.Equal(`GRANT SELECT ON FUTURE EXTERNAL TABLES IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke = fvgd.Role("bob").Revoke("SELECT")
	b.Equal([]string{`REVOKE SELECT ON FUTURE EXTERNAL TABLES IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)
}

func TestFutureFileFormatGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureFileFormatGrant("test_db", "PUBLIC")
	r.Equal("PUBLIC", fvg.Name())

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUTURE FILE FORMATS IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	revoke := fvg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FUTURE FILE FORMATS IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`}, revoke)

	b := require.New(t)
	fvgd := snowflake.FutureFileFormatGrant("test_db", "")
	b.Equal("test_db", fvgd.Name())

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("USAGE", false)
	b.Equal(`GRANT USAGE ON FUTURE FILE FORMATS IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke = fvgd.Role("bob").Revoke("USAGE")
	b.Equal([]string{`REVOKE USAGE ON FUTURE FILE FORMATS IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)
}

func TestFutureTaskGrant(t *testing.T) {
	r := require.New(t)
	fvg := snowflake.FutureTaskGrant("test_db", "PUBLIC")
	r.Equal("PUBLIC", fvg.Name())

	s := fvg.Show()
	r.Equal(`SHOW FUTURE GRANTS IN SCHEMA "test_db"."PUBLIC"`, s)

	s = fvg.Role("bob").Grant("USAGE", false)
	r.Equal(`GRANT USAGE ON FUTURE TASKS IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	revoke := fvg.Role("bob").Revoke("USAGE")
	r.Equal([]string{`REVOKE USAGE ON FUTURE TASKS IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`}, revoke)

	s = fvg.Role("bob").Grant("OWNERSHIP", false)
	r.Equal(`GRANT OWNERSHIP ON FUTURE TASKS IN SCHEMA "test_db"."PUBLIC" TO ROLE "bob"`, s)

	revoke = fvg.Role("bob").RevokeOwnership("OWNERSHIP")
	r.Equal([]string{`REVOKE OWNERSHIP ON FUTURE TASKS IN SCHEMA "test_db"."PUBLIC" FROM ROLE "bob"`}, revoke)

	b := require.New(t)
	fvgd := snowflake.FutureTaskGrant("test_db", "")
	b.Equal("test_db", fvgd.Name())

	s = fvgd.Show()
	b.Equal(`SHOW FUTURE GRANTS IN DATABASE "test_db"`, s)

	s = fvgd.Role("bob").Grant("USAGE", false)
	b.Equal(`GRANT USAGE ON FUTURE TASKS IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke = fvgd.Role("bob").Revoke("USAGE")
	b.Equal([]string{`REVOKE USAGE ON FUTURE TASKS IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)

	s = fvgd.Role("bob").Grant("OWNERSHIP", false)
	b.Equal(`GRANT OWNERSHIP ON FUTURE TASKS IN DATABASE "test_db" TO ROLE "bob"`, s)

	revoke = fvgd.Role("bob").RevokeOwnership("OWNERSHIP")
	b.Equal([]string{`REVOKE OWNERSHIP ON FUTURE TASKS IN DATABASE "test_db" FROM ROLE "bob"`}, revoke)
}
