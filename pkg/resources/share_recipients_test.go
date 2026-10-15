package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_reconcileShareRecipients(t *testing.T) {
	t.Run("account usage is current: full list from it, unmanaged recipient included", func(t *testing.T) {
		visible := []string{"SOME_ORG.ACCOUNT_A", "SOME_ORG.ACCOUNT_B", "SOME_ORG.ACCOUNT_C"}
		usage := []string{"SOME_ORG.ACCOUNT_A", "SOME_ORG.ACCOUNT_B", "SOME_ORG.ACCOUNT_C", "SOME_ORG.ACCOUNT_D"}
		expected := []string{"SOME_ORG.ACCOUNT_B", "SOME_ORG.ACCOUNT_A", "SOME_ORG.ACCOUNT_D"}
		require.Equal(t, usage, reconcileShareRecipients(expected, visible, usage, 4))
	})

	t.Run("hidden recipient that is not expected is reported, not hidden", func(t *testing.T) {
		visible := []string{"ORG.A", "ORG.B", "ORG.C"}
		usage := []string{"ORG.A", "ORG.B", "ORG.C", "ORG.X"}
		require.Equal(t, []string{"ORG.A", "ORG.B", "ORG.C", "ORG.X"}, reconcileShareRecipients([]string{"ORG.A", "ORG.B", "ORG.C", "ORG.D"}, visible, usage, 4))
	})

	t.Run("recently added recipient not in account usage yet is taken from expected", func(t *testing.T) {
		visible := []string{"ORG.A", "ORG.B", "ORG.C"}
		usage := []string{"ORG.A", "ORG.B", "ORG.C"}
		require.Equal(t, []string{"ORG.A", "ORG.B", "ORG.C", "ORG.D"}, reconcileShareRecipients([]string{"ORG.A", "ORG.B", "ORG.C", "ORG.D"}, visible, usage, 4))
	})

	t.Run("recently removed recipient still in account usage is dropped", func(t *testing.T) {
		visible := []string{"ORG.A", "ORG.B", "ORG.C"}
		usage := []string{"ORG.A", "ORG.B", "ORG.C", "ORG.D", "ORG.OLD"}
		require.Equal(t, []string{"ORG.A", "ORG.B", "ORG.C", "ORG.D"}, reconcileShareRecipients([]string{"ORG.A", "ORG.B", "ORG.C", "ORG.D"}, visible, usage, 4))
	})

	t.Run("visible recipients are never dropped and matching ignores case", func(t *testing.T) {
		visible := []string{"org.a", "ORG.B", "ORG.C"}
		usage := []string{"ORG.A", "ORG.B", "ORG.C", "ORG.D"}
		require.Equal(t, []string{"org.a", "ORG.B", "ORG.C", "ORG.D"}, reconcileShareRecipients(nil, visible, usage, 4))
	})
}
