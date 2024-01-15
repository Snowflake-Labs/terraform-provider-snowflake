package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateApplicationOptions]   = new(CreateApplicationRequest)
	_ optionsProvider[DropApplicationOptions]     = new(DropApplicationRequest)
	_ optionsProvider[AlterApplicationOptions]    = new(AlterApplicationRequest)
	_ optionsProvider[ShowApplicationOptions]     = new(ShowApplicationRequest)
	_ optionsProvider[DescribeApplicationOptions] = new(DescribeApplicationRequest)
)

type CreateApplicationRequest struct {
	name        AccountObjectIdentifier // required
	PackageName AccountObjectIdentifier // required
	Version     *ApplicationVersionRequest
	DebugMode   *bool
	Comment     *string
	Tag         []TagAssociation
}

type ApplicationVersionRequest struct {
	VersionDirectory *string
	VersionAndPatch  *VersionAndPatchRequest
}

type VersionAndPatchRequest struct {
	Version string // required
	Patch   *int   // required
}

type DropApplicationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
	Cascade  *bool
}

type AlterApplicationRequest struct {
	IfExists                     *bool
	name                         AccountObjectIdentifier // required
	Set                          *ApplicationSetRequest
	UnsetComment                 *bool
	UnsetShareEventsWithProvider *bool
	UnsetDebugMode               *bool
	Upgrade                      *bool
	UpgradeVersion               *ApplicationVersionRequest
	UnsetReferences              *ApplicationReferencesRequest
	SetTags                      []TagAssociation
	UnsetTags                    []ObjectIdentifier
}

type ApplicationSetRequest struct {
	Comment                 *string
	ShareEventsWithProvider *bool
	DebugMode               *bool
}

type ApplicationReferencesRequest struct {
	References []ApplicationReferenceRequest
}

type ApplicationReferenceRequest struct {
	Reference string
}

type ShowApplicationRequest struct {
	Like       *Like
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeApplicationRequest struct {
	name AccountObjectIdentifier // required
}
