package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
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
	managedAccounts, err := v.Show(ctx, NewShowManagedAccountRequest().WithLike(&Like{String(id.Name())}))
	if err != nil {
		return nil, err
	}
	return collections.FindOne(managedAccounts, func(r ManagedAccount) bool { return r.Name == id.Name() })
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
		name: r.name,
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
		Name:              r.Name,
		Cloud:             r.Cloud,
		Region:            r.Region,
		Locator:           r.Locator,
		CreatedOn:         r.CreatedOn,
		URL:               r.Url,
		AccountLocatorURL: r.AccountLocatorUrl,
		IsReader:          r.IsReader,
	}
	if r.Comment.Valid {
		managedAccount.Comment = r.Comment.String
	}
	return managedAccount
}
