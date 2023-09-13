package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateNetworkPolicyOptions] = new(CreateNetworkPolicyRequest)
	_ optionsProvider[ShowNetworkPolicyOptions]   = new(ShowNetworkPolicyRequest)
)

type CreateNetworkPolicyRequest struct {
	OrReplace     *bool
	name          AccountObjectIdentifier
	AllowedIpList []string
	Comment       *string
}

type ShowNetworkPolicyRequest struct {
}
