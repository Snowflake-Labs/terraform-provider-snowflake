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
	OrReplace              *bool
	name                   AccountObjectIdentifier // required
	AllowedNetworkRuleList []SchemaObjectIdentifier
	BlockedNetworkRuleList []SchemaObjectIdentifier
	AllowedIpList          []IPRequest
	BlockedIpList          []IPRequest
	Comment                *string
}

func (r *CreateNetworkPolicyRequest) GetName() AccountObjectIdentifier {
	return r.name
}

type IPRequest struct {
	IP string // required
}

type AlterNetworkPolicyRequest struct {
	IfExists     *bool
	name         AccountObjectIdentifier // required
	Set          *NetworkPolicySetRequest
	Add          *AddNetworkRuleRequest
	Remove       *RemoveNetworkRuleRequest
	UnsetComment *bool
	RenameTo     *AccountObjectIdentifier
}

type NetworkPolicySetRequest struct {
	AllowedNetworkRuleList []SchemaObjectIdentifier
	BlockedNetworkRuleList []SchemaObjectIdentifier
	AllowedIpList          []IPRequest
	BlockedIpList          []IPRequest
	Comment                *string
}

type AddNetworkRuleRequest struct {
	AllowedNetworkRuleList []SchemaObjectIdentifier
	BlockedNetworkRuleList []SchemaObjectIdentifier
}

type RemoveNetworkRuleRequest struct {
	AllowedNetworkRuleList []SchemaObjectIdentifier
	BlockedNetworkRuleList []SchemaObjectIdentifier
}

type DropNetworkPolicyRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowNetworkPolicyRequest struct{}

type DescribeNetworkPolicyRequest struct {
	name AccountObjectIdentifier // required
}
