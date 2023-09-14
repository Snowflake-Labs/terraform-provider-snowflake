package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateNetworkPolicyOptions]   = new(CreateNetworkPolicyRequest)
	_ optionsProvider[DropNetworkPolicyOptions]     = new(DropNetworkPolicyRequest)
	_ optionsProvider[ShowNetworkPolicyOptions]     = new(ShowNetworkPolicyRequest)
	_ optionsProvider[DescribeNetworkPolicyOptions] = new(DescribeNetworkPolicyRequest)
)

type CreateNetworkPolicyRequest struct {
	OrReplace     *bool
	name          AccountObjectIdentifier
	AllowedIpList []string
	Comment       *string
}

type DropNetworkPolicyRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier
}

type ShowNetworkPolicyRequest struct {
}

type DescribeNetworkPolicyRequest struct {
	name AccountObjectIdentifier
}
