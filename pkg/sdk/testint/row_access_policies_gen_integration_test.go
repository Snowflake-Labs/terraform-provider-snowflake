package testint

import "testing"

func TestInt_RowAccessPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create row access policy: no optionals", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("create row access policy: full", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop row access policy: existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop row access policy: non-existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter row access policy: rename", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter row access policy: set and unset comment", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter row access policy: set and unset body", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter row access policy: set and unset tags", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show row access policy: default", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show row access policy: with options", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe row access policy: existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe row access policy: non-existing", func(t *testing.T) {
		// TODO: fill me
	})
}
