package testint

import (
	"fmt"
	"testing"
)

func TestInt_MaterializedViews(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := createTable(t, client, testDb(t), testSchema(t))
	t.Cleanup(tableCleanup)

	sql := fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName())

	t.Run("create materialized view: no optionals", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("create materialized view: complete case", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop materialized view: existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop view: non-existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: rename", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: set cluster by", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: recluster suspend and resume", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: suspend and resume", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: set and unset values", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show materialized view: default", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show materialized view: with options", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe materialized view", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe materialized view: non-existing", func(t *testing.T) {
		// TODO: fill me
	})
}
