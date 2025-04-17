package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ ManagedAccounts = (*managedAccounts)(nil)

type managedAccounts struct {
	client *Client
}

func (v *managedAccounts) Create(ctx context.Context, request *CreateManagedAccountRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *managedAccounts) Drop(ctx context.Context, request *DropManagedAccountRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *managedAccounts) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropManagedAccountRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *managedAccounts) Show(ctx context.Context, request *ShowManagedAccountRequest) ([]ManagedAccount, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[managedAccountDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[managedAccountDBRow, ManagedAccount](dbRows)
	return resultList, nil
}

func (v *managedAccounts) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ManagedAccount, error) {
	request := NewShowManagedAccountRequest().
		WithLike(Like{Pattern: String(id.Name())})
	managedAccounts, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(managedAccounts, func(r ManagedAccount) bool { return r.Name == id.Name() })
}

func (v *managedAccounts) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ManagedAccount, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateManagedAccountRequest) toOpts() *CreateManagedAccountOptions {
	opts := &CreateManagedAccountOptions{
		name: r.name,
	}
	opts.CreateManagedAccountParams = CreateManagedAccountParams{
		AdminName:     r.CreateManagedAccountParams.AdminName,
		AdminPassword: r.CreateManagedAccountParams.AdminPassword,
		Comment:       r.CreateManagedAccountParams.Comment,
	}
	return opts
}

func (r *DropManagedAccountRequest) toOpts() *DropManagedAccountOptions {
	opts := &DropManagedAccountOptions{
		name:     r.name,
		IfExists: r.IfExists,
	}
	return opts
}

func (r *ShowManagedAccountRequest) toOpts() *ShowManagedAccountOptions {
	opts := &ShowManagedAccountOptions{
		Like: r.Like,
	}
	return opts
}

func (r managedAccountDBRow) convert() *ManagedAccount {
	managedAccount := &ManagedAccount{
		Cloud:             r.Cloud,
		Region:            r.Region,
		CreatedOn:         r.CreatedOn,
		AccountLocatorURL: r.AccountLocatorUrl,
		IsReader:          r.IsReader,
	}

	if r.AccountName.Valid {
		managedAccount.Name = r.AccountName.String
	} else if r.Name.Valid {
		managedAccount.Name = r.Name.String
	}

	if r.AccountLocator.Valid {
		managedAccount.Locator = r.AccountLocator.String
	} else if r.Locator.Valid {
		managedAccount.Locator = r.Locator.String
	}

	if r.AccountUrl.Valid {
		managedAccount.URL = r.AccountUrl.String
	} else if r.Url.Valid {
		managedAccount.URL = r.Url.String
	}

	if r.Comment.Valid {
		managedAccount.Comment = &r.Comment.String
	}

	return managedAccount
}
