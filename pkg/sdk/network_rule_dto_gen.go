package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateNetworkRuleOptions]   = new(CreateNetworkRuleRequest)
	_ optionsProvider[AlterNetworkRuleOptions]    = new(AlterNetworkRuleRequest)
	_ optionsProvider[DropNetworkRuleOptions]     = new(DropNetworkRuleRequest)
	_ optionsProvider[ShowNetworkRuleOptions]     = new(ShowNetworkRuleRequest)
	_ optionsProvider[DescribeNetworkRuleOptions] = new(DescribeNetworkRuleRequest)
)

type CreateNetworkRuleRequest struct {
	OrReplace *bool
	name      SchemaObjectIdentifier // required
	Type      NetworkRuleType        // required
	ValueList []NetworkRuleValue     // required
	Mode      NetworkRuleMode        // required
	Comment   *string
}

func (r *CreateNetworkRuleRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type AlterNetworkRuleRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Set      *NetworkRuleSetRequest
	Unset    *NetworkRuleUnsetRequest
}

type NetworkRuleSetRequest struct {
	ValueList []NetworkRuleValue // required
	Comment   *string
}

type NetworkRuleUnsetRequest struct {
	ValueList *bool
	Comment   *bool
}

type DropNetworkRuleRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowNetworkRuleRequest struct {
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeNetworkRuleRequest struct {
	name SchemaObjectIdentifier // required
}
