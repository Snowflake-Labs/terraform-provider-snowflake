package sdk

import "testing"

func TestManagedAccounts_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid CreateManagedAccountOptions
	defaultOpts := func() *CreateManagedAccountOptions {
		return &CreateManagedAccountOptions{
			name: id,
			CreateManagedAccountParams: CreateManagedAccountParams{
				AdminName:     "admin",
				AdminPassword: "password",
			},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateManagedAccountOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: [opts.CreateManagedAccountParams.AdminName] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.CreateManagedAccountParams.AdminName = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateManagedAccountOptions.CreateManagedAccountParams", "AdminName"))
	})

	t.Run("validation: [opts.CreateManagedAccountParams.AdminPassword] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.CreateManagedAccountParams.AdminPassword = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateManagedAccountOptions.CreateManagedAccountParams", "AdminPassword"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE MANAGED ACCOUNT %s ADMIN_NAME = 'admin', ADMIN_PASSWORD = 'password', TYPE = READER", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.CreateManagedAccountParams.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE MANAGED ACCOUNT %s ADMIN_NAME = 'admin', ADMIN_PASSWORD = 'password', TYPE = READER, COMMENT = 'comment'", id.FullyQualifiedName())
	})
}

func TestManagedAccounts_Drop(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid DropManagedAccountOptions
	defaultOpts := func() *DropManagedAccountOptions {
		return &DropManagedAccountOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropManagedAccountOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP MANAGED ACCOUNT %s", id.FullyQualifiedName())
	})
}

func TestManagedAccounts_Show(t *testing.T) {
	// Minimal valid ShowManagedAccountOptions
	defaultOpts := func() *ShowManagedAccountOptions {
		return &ShowManagedAccountOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowManagedAccountOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW MANAGED ACCOUNTS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("myaccount"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MANAGED ACCOUNTS LIKE 'myaccount'")
	})
}
