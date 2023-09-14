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

func (v *networkPolicies) Show(ctx context.Context, request *ShowNetworkPolicyRequest) (any, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[databaseNetworkPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[databaseNetworkPolicyDBRow, NetworkPolicy](dbRows)
	return resultList, nil
}

func (v databaseNetworkPolicyDBRow) convert() *NetworkPolicy {
	// TODO Generate (template at least)
	return nil
}

func (r *CreateNetworkPolicyRequest) toOpts() *CreateNetworkPolicyOptions {
	opts := &CreateNetworkPolicyOptions{
		OrReplace:     r.OrReplace,
		name:          r.name,
		AllowedIpList: r.AllowedIpList,
		Comment:       r.Comment,
	}
	return opts
}

func (r *ShowNetworkPolicyRequest) toOpts() *ShowNetworkPolicyOptions {
	opts := &ShowNetworkPolicyOptions{}
	return opts
}
