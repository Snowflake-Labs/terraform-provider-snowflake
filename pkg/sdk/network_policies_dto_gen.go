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

type IPRequest struct {
	IP string // required
}

type AlterNetworkPolicyRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
	Set      *NetworkPolicySetRequest
	Unset    *NetworkPolicyUnsetRequest
	Add      *AddNetworkRuleRequest
	Remove   *RemoveNetworkRuleRequest
	RenameTo *AccountObjectIdentifier
}

type NetworkPolicySetRequest struct {
	AllowedNetworkRuleList *AllowedNetworkRuleListRequest
	BlockedNetworkRuleList *BlockedNetworkRuleListRequest
	AllowedIpList          *AllowedIPListRequest
	BlockedIpList          *BlockedIPListRequest
	Comment                *string
}

type AllowedNetworkRuleListRequest struct {
	AllowedNetworkRuleList []SchemaObjectIdentifier
}

type BlockedNetworkRuleListRequest struct {
	BlockedNetworkRuleList []SchemaObjectIdentifier
}

type AllowedIPListRequest struct {
	AllowedIPList []IPRequest
}

type BlockedIPListRequest struct {
	BlockedIPList []IPRequest
}

type NetworkPolicyUnsetRequest struct {
	AllowedNetworkRuleList *bool
	BlockedNetworkRuleList *bool
	AllowedIpList          *bool
	BlockedIpList          *bool
	Comment                *bool
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

type ShowNetworkPolicyRequest struct {
	Like *Like
}

type DescribeNetworkPolicyRequest struct {
	name AccountObjectIdentifier // required
}
