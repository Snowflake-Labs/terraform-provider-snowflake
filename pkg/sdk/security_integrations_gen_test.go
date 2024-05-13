package sdk

import (
	"testing"
)

func TestSecurityIntegrations_CreateSCIM(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateSCIMSecurityIntegrationOptions
	defaultOpts := func() *CreateSCIMSecurityIntegrationOptions {
		return &CreateSCIMSecurityIntegrationOptions{
			name:       id,
			Enabled:    true,
			ScimClient: "GENERIC",
			RunAsRole:  "GENERIC_SCIM_PROVISIONER",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSCIMSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSCIMSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION %s TYPE = SCIM ENABLED = true SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = 'GENERIC_SCIM_PROVISIONER'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		networkPolicyID := randomAccountObjectIdentifier()
		opts.NetworkPolicy = Pointer(networkPolicyID)
		opts.SyncPassword = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION %s TYPE = SCIM ENABLED = true SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = 'GENERIC_SCIM_PROVISIONER'"+
			" NETWORK_POLICY = %s SYNC_PASSWORD = true", id.FullyQualifiedName(), networkPolicyID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterSCIMIntegration(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterSCIMIntegrationSecurityIntegrationOptions
	defaultOpts := func() *AlterSCIMIntegrationSecurityIntegrationOptions {
		return &AlterSCIMIntegrationSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSCIMIntegrationSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set.Enabled opts.Set.NetworkPolicy opts.Set.SyncPassword opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSCIMIntegrationSecurityIntegrationOptions.Set", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.NetworkPolicy opts.Unset.SyncPassword opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSCIMIntegrationSecurityIntegrationOptions.Unset", "NetworkPolicy", "SyncPassword", "Comment"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		networkPolicyID := randomAccountObjectIdentifier()
		opts.Set = &SCIMIntegrationSet{
			Enabled:       Pointer(true),
			NetworkPolicy: Pointer(networkPolicyID),
			SyncPassword:  Pointer(true),
			Comment:       Pointer("test"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true NETWORK_POLICY = %s SYNC_PASSWORD = true COMMENT = 'test'",
			id.FullyQualifiedName(), networkPolicyID.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &SCIMIntegrationUnset{
			NetworkPolicy: Pointer(true),
			SyncPassword:  Pointer(true),
			Comment:       Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET NETWORK_POLICY SYNC_PASSWORD COMMENT", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid DropSecurityIntegrationOptions
	defaultOpts := func() *DropSecurityIntegrationOptions {
		return &DropSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP SECURITY INTEGRATION IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid DescribeSecurityIntegrationOptions
	defaultOpts := func() *DescribeSecurityIntegrationOptions {
		return &DescribeSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE SECURITY INTEGRATION %s", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Show(t *testing.T) {
	// Minimal valid ShowSecurityIntegrationOptions
	defaultOpts := func() *ShowSecurityIntegrationOptions {
		return &ShowSecurityIntegrationOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECURITY INTEGRATIONS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("some pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECURITY INTEGRATIONS LIKE 'some pattern'")
	})
}
