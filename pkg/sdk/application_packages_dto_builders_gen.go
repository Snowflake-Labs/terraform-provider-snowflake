// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateApplicationPackageRequest(
	name AccountObjectIdentifier,
) *CreateApplicationPackageRequest {
	s := CreateApplicationPackageRequest{}
	s.name = name
	return &s
}

func (s *CreateApplicationPackageRequest) WithIfNotExists(IfNotExists *bool) *CreateApplicationPackageRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateApplicationPackageRequest) WithDataRetentionTimeInDays(DataRetentionTimeInDays *int) *CreateApplicationPackageRequest {
	s.DataRetentionTimeInDays = DataRetentionTimeInDays
	return s
}

func (s *CreateApplicationPackageRequest) WithMaxDataExtensionTimeInDays(MaxDataExtensionTimeInDays *int) *CreateApplicationPackageRequest {
	s.MaxDataExtensionTimeInDays = MaxDataExtensionTimeInDays
	return s
}

func (s *CreateApplicationPackageRequest) WithDefaultDdlCollation(DefaultDdlCollation *string) *CreateApplicationPackageRequest {
	s.DefaultDdlCollation = DefaultDdlCollation
	return s
}

func (s *CreateApplicationPackageRequest) WithComment(Comment *string) *CreateApplicationPackageRequest {
	s.Comment = Comment
	return s
}

func (s *CreateApplicationPackageRequest) WithDistribution(Distribution *Distribution) *CreateApplicationPackageRequest {
	s.Distribution = Distribution
	return s
}

func (s *CreateApplicationPackageRequest) WithTag(Tag []TagAssociation) *CreateApplicationPackageRequest {
	s.Tag = Tag
	return s
}

func NewAlterApplicationPackageRequest(
	name AccountObjectIdentifier,
) *AlterApplicationPackageRequest {
	s := AlterApplicationPackageRequest{}
	s.name = name
	return &s
}

func (s *AlterApplicationPackageRequest) WithIfExists(IfExists *bool) *AlterApplicationPackageRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterApplicationPackageRequest) WithSet(Set *ApplicationPackageSetRequest) *AlterApplicationPackageRequest {
	s.Set = Set
	return s
}

func (s *AlterApplicationPackageRequest) WithUnsetDataRetentionTimeInDays(UnsetDataRetentionTimeInDays *bool) *AlterApplicationPackageRequest {
	s.UnsetDataRetentionTimeInDays = UnsetDataRetentionTimeInDays
	return s
}

func (s *AlterApplicationPackageRequest) WithUnsetMaxDataExtensionTimeInDays(UnsetMaxDataExtensionTimeInDays *bool) *AlterApplicationPackageRequest {
	s.UnsetMaxDataExtensionTimeInDays = UnsetMaxDataExtensionTimeInDays
	return s
}

func (s *AlterApplicationPackageRequest) WithUnsetDefaultDdlCollation(UnsetDefaultDdlCollation *bool) *AlterApplicationPackageRequest {
	s.UnsetDefaultDdlCollation = UnsetDefaultDdlCollation
	return s
}

func (s *AlterApplicationPackageRequest) WithUnsetComment(UnsetComment *bool) *AlterApplicationPackageRequest {
	s.UnsetComment = UnsetComment
	return s
}

func (s *AlterApplicationPackageRequest) WithUnsetDistribution(UnsetDistribution *bool) *AlterApplicationPackageRequest {
	s.UnsetDistribution = UnsetDistribution
	return s
}

func (s *AlterApplicationPackageRequest) WithModifyReleaseDirective(ModifyReleaseDirective *ModifyReleaseDirectiveRequest) *AlterApplicationPackageRequest {
	s.ModifyReleaseDirective = ModifyReleaseDirective
	return s
}

func (s *AlterApplicationPackageRequest) WithSetDefaultReleaseDirective(SetDefaultReleaseDirective *SetDefaultReleaseDirectiveRequest) *AlterApplicationPackageRequest {
	s.SetDefaultReleaseDirective = SetDefaultReleaseDirective
	return s
}

func (s *AlterApplicationPackageRequest) WithSetReleaseDirective(SetReleaseDirective *SetReleaseDirectiveRequest) *AlterApplicationPackageRequest {
	s.SetReleaseDirective = SetReleaseDirective
	return s
}

func (s *AlterApplicationPackageRequest) WithUnsetReleaseDirective(UnsetReleaseDirective *UnsetReleaseDirectiveRequest) *AlterApplicationPackageRequest {
	s.UnsetReleaseDirective = UnsetReleaseDirective
	return s
}

func (s *AlterApplicationPackageRequest) WithAddVersion(AddVersion *AddVersionRequest) *AlterApplicationPackageRequest {
	s.AddVersion = AddVersion
	return s
}

func (s *AlterApplicationPackageRequest) WithDropVersion(DropVersion *DropVersionRequest) *AlterApplicationPackageRequest {
	s.DropVersion = DropVersion
	return s
}

func (s *AlterApplicationPackageRequest) WithAddPatchForVersion(AddPatchForVersion *AddPatchForVersionRequest) *AlterApplicationPackageRequest {
	s.AddPatchForVersion = AddPatchForVersion
	return s
}

func (s *AlterApplicationPackageRequest) WithSetTags(SetTags []TagAssociation) *AlterApplicationPackageRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterApplicationPackageRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterApplicationPackageRequest {
	s.UnsetTags = UnsetTags
	return s
}

func NewApplicationPackageSetRequest() *ApplicationPackageSetRequest {
	return &ApplicationPackageSetRequest{}
}

func (s *ApplicationPackageSetRequest) WithDataRetentionTimeInDays(DataRetentionTimeInDays *int) *ApplicationPackageSetRequest {
	s.DataRetentionTimeInDays = DataRetentionTimeInDays
	return s
}

func (s *ApplicationPackageSetRequest) WithMaxDataExtensionTimeInDays(MaxDataExtensionTimeInDays *int) *ApplicationPackageSetRequest {
	s.MaxDataExtensionTimeInDays = MaxDataExtensionTimeInDays
	return s
}

func (s *ApplicationPackageSetRequest) WithDefaultDdlCollation(DefaultDdlCollation *string) *ApplicationPackageSetRequest {
	s.DefaultDdlCollation = DefaultDdlCollation
	return s
}

func (s *ApplicationPackageSetRequest) WithComment(Comment *string) *ApplicationPackageSetRequest {
	s.Comment = Comment
	return s
}

func (s *ApplicationPackageSetRequest) WithDistribution(Distribution *Distribution) *ApplicationPackageSetRequest {
	s.Distribution = Distribution
	return s
}

func NewModifyReleaseDirectiveRequest(
	ReleaseDirective string,
	Version string,
	Patch int,
) *ModifyReleaseDirectiveRequest {
	s := ModifyReleaseDirectiveRequest{}
	s.ReleaseDirective = ReleaseDirective
	s.Version = Version
	s.Patch = Patch
	return &s
}

func NewSetDefaultReleaseDirectiveRequest(
	Version string,
	Patch int,
) *SetDefaultReleaseDirectiveRequest {
	s := SetDefaultReleaseDirectiveRequest{}
	s.Version = Version
	s.Patch = Patch
	return &s
}

func NewSetReleaseDirectiveRequest(
	ReleaseDirective string,
	Accounts []string,
	Version string,
	Patch int,
) *SetReleaseDirectiveRequest {
	s := SetReleaseDirectiveRequest{}
	s.ReleaseDirective = ReleaseDirective
	s.Accounts = Accounts
	s.Version = Version
	s.Patch = Patch
	return &s
}

func NewUnsetReleaseDirectiveRequest(
	ReleaseDirective string,
) *UnsetReleaseDirectiveRequest {
	s := UnsetReleaseDirectiveRequest{}
	s.ReleaseDirective = ReleaseDirective
	return &s
}

func NewAddVersionRequest(
	Using string,
) *AddVersionRequest {
	s := AddVersionRequest{}
	s.Using = Using
	return &s
}

func (s *AddVersionRequest) WithVersionIdentifier(VersionIdentifier *string) *AddVersionRequest {
	s.VersionIdentifier = VersionIdentifier
	return s
}

func (s *AddVersionRequest) WithLabel(Label *string) *AddVersionRequest {
	s.Label = Label
	return s
}

func NewDropVersionRequest(
	VersionIdentifier string,
) *DropVersionRequest {
	s := DropVersionRequest{}
	s.VersionIdentifier = VersionIdentifier
	return &s
}

func NewAddPatchForVersionRequest(
	VersionIdentifier *string,
	Using string,
) *AddPatchForVersionRequest {
	s := AddPatchForVersionRequest{}
	s.VersionIdentifier = VersionIdentifier
	s.Using = Using
	return &s
}

func (s *AddPatchForVersionRequest) WithLabel(Label *string) *AddPatchForVersionRequest {
	s.Label = Label
	return s
}

func NewDropApplicationPackageRequest(
	name AccountObjectIdentifier,
) *DropApplicationPackageRequest {
	s := DropApplicationPackageRequest{}
	s.name = name
	return &s
}

func NewShowApplicationPackageRequest() *ShowApplicationPackageRequest {
	return &ShowApplicationPackageRequest{}
}

func (s *ShowApplicationPackageRequest) WithLike(Like *Like) *ShowApplicationPackageRequest {
	s.Like = Like
	return s
}

func (s *ShowApplicationPackageRequest) WithStartsWith(StartsWith *string) *ShowApplicationPackageRequest {
	s.StartsWith = StartsWith
	return s
}

func (s *ShowApplicationPackageRequest) WithLimit(Limit *LimitFrom) *ShowApplicationPackageRequest {
	s.Limit = Limit
	return s
}
