package testint

import "testing"

func TestInt_Conntections(t *testing.T) {
	client := testClient(t)
	_ = client
	ctx := testContext(t)
	_ = ctx

	t.Run("CreateConnection", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		_ = id
	})

	t.Run("CreateReplicatedConnection", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("AlterConnectionFailover", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("AlterConnection", func(t *testing.T) {
		// TODO: fill me
	})
}
