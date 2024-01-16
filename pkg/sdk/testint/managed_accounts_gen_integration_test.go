package testint

import "testing"

func TestInt_ManagedAccounts(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create managed account: no optionals", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("create managed account: full", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop managed account: existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop managed account: non-existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show managed account: default", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show managed account: with like", func(t *testing.T) {
		// TODO: fill me
	})
}
