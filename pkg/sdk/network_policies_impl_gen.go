package sdk

import "context"

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

func (v *networkPolicies) Describe(ctx context.Context, id AccountObjectIdentifier) ([]NetworkPolicyDescription, error) {
	opts := &DescribeNetworkPolicyOptions{
		// TODO enforce this convention in the DSL (field "name" is queryStruct identifier)
		name: id,
	}
	s, err := validateAndQuery[describeNetworkPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	result := make([]NetworkPolicyDescription, len(*s))
	for i, value := range *s {
		result[i] = *value.convert()
	}
	return result, nil
}

func (r *CreateNetworkPolicyRequest) toOpts() *CreateNetworkPolicyOptions {
	opts := &CreateNetworkPolicyOptions{
		OrReplace: r.OrReplace,
		name:      r.name,

		Comment: r.Comment,
	}
	if r.AllowedIpList != nil {
		s := make([]IP, len(r.AllowedIpList))
		for i, v := range r.AllowedIpList {
			s[i] = IP{
				IP: v.IP,
			}
		}
		opts.AllowedIpList = s
	}
	if r.BlockedIpList != nil {
		s := make([]IP, len(r.BlockedIpList))
		for i, v := range r.BlockedIpList {
			s[i] = IP{
				IP: v.IP,
			}
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

			Comment: r.Set.Comment,
		}
		if r.Set.AllowedIpList != nil {
			s := make([]IP, len(r.Set.AllowedIpList))
			for i, v := range r.Set.AllowedIpList {
				s[i] = IP{
					IP: v.IP,
				}
			}
			opts.Set.AllowedIpList = s
		}
		if r.Set.BlockedIpList != nil {
			s := make([]IP, len(r.Set.BlockedIpList))
			for i, v := range r.Set.BlockedIpList {
				s[i] = IP{
					IP: v.IP,
				}
			}
			opts.Set.BlockedIpList = s
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
		CreatedOn:              r.CreatedOn,
		Name:                   r.Name,
		Comment:                r.Comment,
		EntriesInAllowedIpList: r.EntriesInAllowedIpList,
		EntriesInBlockedIpList: r.EntriesInBlockedIpList,
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
