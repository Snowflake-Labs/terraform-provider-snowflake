package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateApplicationPackageOptions] = new(CreateApplicationPackageRequest)
	_ optionsProvider[AlterApplicationPackageOptions]  = new(AlterApplicationPackageRequest)
	_ optionsProvider[DropApplicationPackageOptions]   = new(DropApplicationPackageRequest)
	_ optionsProvider[ShowApplicationPackageOptions]   = new(ShowApplicationPackageRequest)
)

type CreateApplicationPackageRequest struct {
	IfNotExists                *bool
	name                       AccountObjectIdentifier // required
	DataRetentionTimeInDays    *int
	MaxDataExtensionTimeInDays *int
	DefaultDdlCollation        *string
	Comment                    *string
	Distribution               *Distribution
	Tag                        []TagAssociation
}

type AlterApplicationPackageRequest struct {
	IfExists                        *bool
	name                            AccountObjectIdentifier // required
	Set                             *ApplicationPackageSetRequest
	UnsetDataRetentionTimeInDays    *bool
	UnsetMaxDataExtensionTimeInDays *bool
	UnsetDefaultDdlCollation        *bool
	UnsetComment                    *bool
	UnsetDistribution               *bool
	ModifyReleaseDirective          *ModifyReleaseDirectiveRequest
	SetDefaultReleaseDirective      *SetDefaultReleaseDirectiveRequest
	SetReleaseDirective             *SetReleaseDirectiveRequest
	UnsetReleaseDirective           *UnsetReleaseDirectiveRequest
	AddVersion                      *AddVersionRequest
	DropVersion                     *DropVersionRequest
	AddPatchForVersion              *AddPatchForVersionRequest
	SetTags                         []TagAssociation
	UnsetTags                       []ObjectIdentifier
}

type ApplicationPackageSetRequest struct {
	DataRetentionTimeInDays    *int
	MaxDataExtensionTimeInDays *int
	DefaultDdlCollation        *string
	Comment                    *string
	Distribution               *Distribution
}

type ModifyReleaseDirectiveRequest struct {
	ReleaseDirective string // required
	Version          string // required
	Patch            int    // required
}

type SetDefaultReleaseDirectiveRequest struct {
	Version string // required
	Patch   int    // required
}

type SetReleaseDirectiveRequest struct {
	ReleaseDirective string   // required
	Accounts         []string // required
	Version          string   // required
	Patch            int      // required
}

type UnsetReleaseDirectiveRequest struct {
	ReleaseDirective string // required
}

type AddVersionRequest struct {
	VersionIdentifier *string
	Using             string // required
	Label             *string
}

type DropVersionRequest struct {
	VersionIdentifier string // required
}

type AddPatchForVersionRequest struct {
	VersionIdentifier *string // required
	Using             string  // required
	Label             *string
}

type DropApplicationPackageRequest struct {
	name AccountObjectIdentifier // required
}

type ShowApplicationPackageRequest struct {
	Like       *Like
	StartsWith *string
	Limit      *LimitFrom
}
