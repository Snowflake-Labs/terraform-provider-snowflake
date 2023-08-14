package sdk

func NewCreateDatabaseRoleRequest(
	name DatabaseObjectIdentifier,
) *CreateDatabaseRoleRequest {
	s := CreateDatabaseRoleRequest{}
	s.name = name
	return &s
}

func (s *CreateDatabaseRoleRequest) WithOrReplace(orReplace bool) *CreateDatabaseRoleRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateDatabaseRoleRequest) WithIfNotExists(ifNotExists bool) *CreateDatabaseRoleRequest {
	s.ifNotExists = ifNotExists
	return s
}

func (s *CreateDatabaseRoleRequest) WithComment(comment *string) *CreateDatabaseRoleRequest {
	s.comment = comment
	return s
}

func NewAlterDatabaseRoleRequest(
	name DatabaseObjectIdentifier,
) *AlterDatabaseRoleRequest {
	s := AlterDatabaseRoleRequest{}
	s.name = name
	return &s
}

func (s *AlterDatabaseRoleRequest) WithIfExists(ifExists bool) *AlterDatabaseRoleRequest {
	s.ifExists = ifExists
	return s
}

func (s *AlterDatabaseRoleRequest) WithRename(rename *DatabaseRoleRenameRequest) *AlterDatabaseRoleRequest {
	s.clearForOneOf()
	s.rename = rename
	return s
}

func (s *AlterDatabaseRoleRequest) WithSet(set *DatabaseRoleSetRequest) *AlterDatabaseRoleRequest {
	s.clearForOneOf()
	s.set = set
	return s
}

func (s *AlterDatabaseRoleRequest) WithUnsetComment() *AlterDatabaseRoleRequest {
	s.clearForOneOf()
	s.unset = NewDatabaseRoleUnsetRequest()
	return s
}

func (s *AlterDatabaseRoleRequest) clearForOneOf() *AlterDatabaseRoleRequest {
	s.rename = nil
	s.set = nil
	s.unset = nil
	return s
}

func NewDatabaseRoleRenameRequest(
	name DatabaseObjectIdentifier,
) *DatabaseRoleRenameRequest {
	s := DatabaseRoleRenameRequest{}
	s.name = name
	return &s
}

func NewDatabaseRoleSetRequest(
	comment string,
) *DatabaseRoleSetRequest {
	s := DatabaseRoleSetRequest{}
	s.comment = comment
	return &s
}

func NewDatabaseRoleUnsetRequest() *DatabaseRoleUnsetRequest {
	return &DatabaseRoleUnsetRequest{}
}

func NewDropDatabaseRoleRequest(
	name DatabaseObjectIdentifier,
) *DropDatabaseRoleRequest {
	s := DropDatabaseRoleRequest{}
	s.name = name
	return &s
}

func (s *DropDatabaseRoleRequest) WithIfExists(ifExists bool) *DropDatabaseRoleRequest {
	s.ifExists = ifExists
	return s
}

func NewShowDatabaseRoleRequest(
	database AccountObjectIdentifier,
) *ShowDatabaseRoleRequest {
	s := ShowDatabaseRoleRequest{}
	s.database = database
	return &s
}

func (s *ShowDatabaseRoleRequest) WithLike(pattern string) *ShowDatabaseRoleRequest {
	s.like = &Like{
		Pattern: String(pattern),
	}
	return s
}
