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

func (v *networkPolicies) Drop(ctx context.Context, request *DropNetworkPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkPolicies) Show(ctx context.Context, request *ShowNetworkPolicyRequest) (*NetworkPolicy, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showNetworkPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showNetworkPolicyDBRow, NetworkPolicy](dbRows)
	return resultList, nil
}

// TODO Task interface identifier kind
func (v *networkPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*NetworkPolicy, error) {
	opts := &DescribeNetworkPolicyOptions{
		// TODO enforce this convention in the DSL (field "name" is queryStruct identifier)
		name: id,
	}
	result, err := validateAndQueryOne[describeNetworkPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
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
	// TODO: Mapping
	return &NetworkPolicy{}
}

func (r *DescribeNetworkPolicyRequest) toOpts() *DescribeNetworkPolicyOptions {
	opts := &DescribeNetworkPolicyOptions{
		name: r.name,
	}
	return opts
}
