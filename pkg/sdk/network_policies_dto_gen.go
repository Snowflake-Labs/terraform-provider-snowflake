package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateNetworkPolicyOptions]   = new(CreateNetworkPolicyRequest)
	_ optionsProvider[AlterNetworkPolicyOptions]    = new(AlterNetworkPolicyRequest)
	_ optionsProvider[DropNetworkPolicyOptions]     = new(DropNetworkPolicyRequest)
	_ optionsProvider[ShowNetworkPolicyOptions]     = new(ShowNetworkPolicyRequest)
	_ optionsProvider[DescribeNetworkPolicyOptions] = new(DescribeNetworkPolicyRequest)
)

type CreateNetworkPolicyRequest struct {
	OrReplace     *bool
	name          AccountObjectIdentifier // required
	AllowedIpList []string                // required
	BlockedIpList []string
	Comment       *string
}

type AlterNetworkPolicyRequest struct {
	IfExists     *bool
	name         AccountObjectIdentifier // required
	Set          *NetworkPolicySetRequest
	UnsetComment *bool
	RenameTo     AccountObjectIdentifier
}

type NetworkPolicySetRequest struct {
	AllowedIpList []string
	BlockedIpList []string
	Comment       *string
}

type DropNetworkPolicyRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowNetworkPolicyRequest struct {
}

type DescribeNetworkPolicyRequest struct {
	name AccountObjectIdentifier // required
}
