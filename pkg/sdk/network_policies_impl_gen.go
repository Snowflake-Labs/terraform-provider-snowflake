package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ NetworkPolicies = (*networkPolicies)(nil)

type networkPolicies struct {
	client *Client
}

func (v *networkPolicies) Create(ctx context.Context, request *CreateNetworkPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkPolicies) Alter(ctx context.Context, request *AlterNetworkPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkPolicies) Drop(ctx context.Context, request *DropNetworkPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkPolicies) Show(ctx context.Context, request *ShowNetworkPolicyRequest) ([]NetworkPolicy, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showNetworkPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showNetworkPolicyDBRow, NetworkPolicy](dbRows)
	return resultList, nil
}

func (v *networkPolicies) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*NetworkPolicy, error) {
	networkPolicies, err := v.Show(ctx, NewShowNetworkPolicyRequest())
	if err != nil {
		return nil, err
	}

	return collections.FindOne(networkPolicies, func(r NetworkPolicy) bool { return r.Name == id.Name() })
}

func (v *networkPolicies) Describe(ctx context.Context, id AccountObjectIdentifier) ([]NetworkPolicyDescription, error) {
	opts := &DescribeNetworkPolicyOptions{
		name: id,
	}
	rows, err := validateAndQuery[describeNetworkPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[describeNetworkPolicyDBRow, NetworkPolicyDescription](rows), nil
}

func (r *CreateNetworkPolicyRequest) toOpts() *CreateNetworkPolicyOptions {
	opts := &CreateNetworkPolicyOptions{
		OrReplace:              r.OrReplace,
		name:                   r.name,
		AllowedNetworkRuleList: r.AllowedNetworkRuleList,
		BlockedNetworkRuleList: r.BlockedNetworkRuleList,

		Comment: r.Comment,
	}
	if r.AllowedIpList != nil {
		s := make([]IP, len(r.AllowedIpList))
		for i, v := range r.AllowedIpList {
			s[i] = IP(v)
		}
		opts.AllowedIpList = s
	}
	if r.BlockedIpList != nil {
		s := make([]IP, len(r.BlockedIpList))
		for i, v := range r.BlockedIpList {
			s[i] = IP(v)
		}
		opts.BlockedIpList = s
	}
	return opts
}

func (r *AlterNetworkPolicyRequest) toOpts() *AlterNetworkPolicyOptions {
	opts := &AlterNetworkPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,

		UnsetComment: r.UnsetComment,
		RenameTo:     r.RenameTo,
	}
	if r.Set != nil {
		opts.Set = &NetworkPolicySet{
			AllowedNetworkRuleList: r.Set.AllowedNetworkRuleList,
			BlockedNetworkRuleList: r.Set.BlockedNetworkRuleList,

			Comment: r.Set.Comment,
		}
		if r.Set.AllowedIpList != nil {
			s := make([]IP, len(r.Set.AllowedIpList))
			for i, v := range r.Set.AllowedIpList {
				s[i] = IP(v)
			}
			opts.Set.AllowedIpList = s
		}
		if r.Set.BlockedIpList != nil {
			s := make([]IP, len(r.Set.BlockedIpList))
			for i, v := range r.Set.BlockedIpList {
				s[i] = IP(v)
			}
			opts.Set.BlockedIpList = s
		}
	}
	if r.Add != nil {
		opts.Add = &AddNetworkRule{
			AllowedNetworkRuleList: r.Add.AllowedNetworkRuleList,
			BlockedNetworkRuleList: r.Add.BlockedNetworkRuleList,
		}
	}
	if r.Remove != nil {
		opts.Remove = &RemoveNetworkRule{
			AllowedNetworkRuleList: r.Remove.AllowedNetworkRuleList,
			BlockedNetworkRuleList: r.Remove.BlockedNetworkRuleList,
		}
	}
	return opts
}

func (r *DropNetworkPolicyRequest) toOpts() *DropNetworkPolicyOptions {
	opts := &DropNetworkPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowNetworkPolicyRequest) toOpts() *ShowNetworkPolicyOptions {
	opts := &ShowNetworkPolicyOptions{}
	return opts
}

func (r showNetworkPolicyDBRow) convert() *NetworkPolicy {
	return &NetworkPolicy{
		CreatedOn:                    r.CreatedOn,
		Name:                         r.Name,
		Comment:                      r.Comment,
		EntriesInAllowedIpList:       r.EntriesInAllowedIpList,
		EntriesInBlockedIpList:       r.EntriesInBlockedIpList,
		EntriesInAllowedNetworkRules: r.EntriesInAllowedNetworkRules,
		EntriesInBlockedNetworkRules: r.EntriesInBlockedNetworkRules,
	}
}

func (r *DescribeNetworkPolicyRequest) toOpts() *DescribeNetworkPolicyOptions {
	opts := &DescribeNetworkPolicyOptions{
		name: r.name,
	}
	return opts
}

func (r describeNetworkPolicyDBRow) convert() *NetworkPolicyDescription {
	return &NetworkPolicyDescription{
		Name:  r.Name,
		Value: r.Value,
	}
}
