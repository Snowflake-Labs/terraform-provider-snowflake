package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSCIMSecurityIntegrationOptions]           = new(CreateSCIMSecurityIntegrationRequest)
	_ optionsProvider[AlterSCIMIntegrationSecurityIntegrationOptions] = new(AlterSCIMIntegrationSecurityIntegrationRequest)
	_ optionsProvider[DropSecurityIntegrationOptions]                 = new(DropSecurityIntegrationRequest)
	_ optionsProvider[DescribeSecurityIntegrationOptions]             = new(DescribeSecurityIntegrationRequest)
	_ optionsProvider[ShowSecurityIntegrationOptions]                 = new(ShowSecurityIntegrationRequest)
)

type CreateSCIMSecurityIntegrationRequest struct {
	OrReplace     *bool
	IfNotExists   *bool
	name          AccountObjectIdentifier // required
	Enabled       bool                    // required
	ScimClient    string                  // required
	RunAsRole     string                  // required
	NetworkPolicy *AccountObjectIdentifier
	SyncPassword  *bool
	Comment       *string
}

type AlterSCIMIntegrationSecurityIntegrationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
	Set      *SCIMIntegrationSetRequest
	Unset    *SCIMIntegrationUnsetRequest
	SetTag   []TagAssociation
	UnsetTag []ObjectIdentifier
}

type SCIMIntegrationSetRequest struct {
	Enabled       *bool
	NetworkPolicy *AccountObjectIdentifier
	SyncPassword  *bool
	Comment       *string
}

type SCIMIntegrationUnsetRequest struct {
	NetworkPolicy *bool
	SyncPassword  *bool
	Comment       *bool
}

type DropSecurityIntegrationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type DescribeSecurityIntegrationRequest struct {
	name AccountObjectIdentifier // required
}

type ShowSecurityIntegrationRequest struct {
	Like *Like
}
