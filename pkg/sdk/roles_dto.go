package sdk

type CreateRoleRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	Comment     *string
	Tag         []TagAssociation
}

func NewCreateRoleRequest(name AccountObjectIdentifier) *CreateRoleRequest {
	return &CreateRoleRequest{
		name: name,
	}
}

func (s *CreateRoleRequest) WithOrReplace(OrReplace bool) *CreateRoleRequest {
	s.OrReplace = Bool(OrReplace)
	return s
}

func (s *CreateRoleRequest) WithIfNotExists(IfNotExists bool) *CreateRoleRequest {
	s.IfNotExists = Bool(IfNotExists)
	return s
}

func (s *CreateRoleRequest) WithComment(Comment string) *CreateRoleRequest {
	s.Comment = String(Comment)
	return s
}

func (s *CreateRoleRequest) WithTag(Tag []TagAssociation) *CreateRoleRequest {
	s.Tag = Tag
	return s
}

type AlterRoleRequest struct {
	IfExists     *bool
	name         AccountObjectIdentifier // required
	RenameTo     *AccountObjectIdentifier
	SetComment   *string
	SetTags      []TagAssociation
	UnsetComment *bool
	UnsetTags    []ObjectIdentifier
}

func NewAlterRoleRequest(name AccountObjectIdentifier) *AlterRoleRequest {
	return &AlterRoleRequest{
		name: name,
	}
}

func (s *AlterRoleRequest) WithIfExists(IfExists bool) *AlterRoleRequest {
	s.IfExists = Bool(IfExists)
	return s
}

func (s *AlterRoleRequest) WithRenameTo(RenameTo AccountObjectIdentifier) *AlterRoleRequest {
	s.RenameTo = &RenameTo
	return s
}

func (s *AlterRoleRequest) WithSetComment(SetComment string) *AlterRoleRequest {
	s.SetComment = String(SetComment)
	return s
}

func (s *AlterRoleRequest) WithSetTags(SetTags []TagAssociation) *AlterRoleRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterRoleRequest) WithUnsetComment(UnsetComment bool) *AlterRoleRequest {
	s.UnsetComment = Bool(UnsetComment)
	return s
}

func (s *AlterRoleRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterRoleRequest {
	s.UnsetTags = UnsetTags
	return s
}

type DropRoleRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

func NewDropRoleRequest(name AccountObjectIdentifier) *DropRoleRequest {
	return &DropRoleRequest{
		name: name,
	}
}

func (s *DropRoleRequest) WithIfExists(IfExists bool) *DropRoleRequest {
	s.IfExists = Bool(IfExists)
	return s
}

type ShowRoleRequest struct {
	Like    *Like
	InClass *RolesInClass
}

func NewShowRoleRequest() *ShowRoleRequest {
	return &ShowRoleRequest{}
}

func (s *ShowRoleRequest) WithLike(like *LikeRequest) *ShowRoleRequest {
	s.Like = &Like{
		Pattern: String(like.pattern),
	}
	return s
}

func (s *ShowRoleRequest) WithInClass(inClass RolesInClass) *ShowRoleRequest {
	s.InClass = &inClass
	return s
}

type LikeRequest struct {
	pattern string // required
}

func NewLikeRequest(pattern string) *LikeRequest {
	return &LikeRequest{
		pattern: pattern,
	}
}

type ShowRoleByIdRequest struct {
	id AccountObjectIdentifier // required
}

func NewShowByIdRoleRequest(id AccountObjectIdentifier) *ShowRoleByIdRequest {
	return &ShowRoleByIdRequest{
		id: id,
	}
}

type GrantRoleRequest struct {
	name  AccountObjectIdentifier // required
	Grant GrantRole               // required
}

func NewGrantRoleRequest(name AccountObjectIdentifier, Grant GrantRole) *GrantRoleRequest {
	return &GrantRoleRequest{
		name:  name,
		Grant: Grant,
	}
}

type RevokeRoleRequest struct {
	name   AccountObjectIdentifier // required
	Revoke RevokeRole              // required
}

func NewRevokeRoleRequest(name AccountObjectIdentifier, Revoke RevokeRole) *RevokeRoleRequest {
	return &RevokeRoleRequest{
		name:   name,
		Revoke: Revoke,
	}
}

type UseRoleRequest struct {
	id AccountObjectIdentifier // required
}

func NewUseRoleRequest(id AccountObjectIdentifier) *UseRoleRequest {
	return &UseRoleRequest{
		id: id,
	}
}

type UseSecondaryRolesRequest struct {
	option SecondaryRoleOption // required
}

func NewUseSecondaryRolesRequest(option SecondaryRoleOption) *UseSecondaryRolesRequest {
	return &UseSecondaryRolesRequest{
		option: option,
	}
}
