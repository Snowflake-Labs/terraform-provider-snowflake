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

func (s *AlterDatabaseRoleRequest) WithRename(name DatabaseObjectIdentifier) *AlterDatabaseRoleRequest {
	s.rename = NewDatabaseRoleRenameRequest(name)
	return s
}

func (s *AlterDatabaseRoleRequest) WithSetComment(comment string) *AlterDatabaseRoleRequest {
	s.set = NewDatabaseRoleSetRequest(comment)
	return s
}

func (s *AlterDatabaseRoleRequest) WithUnsetComment() *AlterDatabaseRoleRequest {
	s.unset = NewDatabaseRoleUnsetRequest()
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

func NewGrantDatabaseRoleRequest(
	name DatabaseObjectIdentifier,
) *GrantDatabaseRoleRequest {
	s := GrantDatabaseRoleRequest{}
	s.name = name
	return &s
}

func (s *GrantDatabaseRoleRequest) WithDatabaseRole(databaseRole DatabaseObjectIdentifier) *GrantDatabaseRoleRequest {
	s.accountRole = nil
	s.databaseRole = &databaseRole
	return s
}

func (s *GrantDatabaseRoleRequest) WithAccountRole(accountRole AccountObjectIdentifier) *GrantDatabaseRoleRequest {
	s.databaseRole = nil
	s.accountRole = &accountRole
	return s
}

func NewRevokeDatabaseRoleRequest(
	name DatabaseObjectIdentifier,
) *RevokeDatabaseRoleRequest {
	s := RevokeDatabaseRoleRequest{}
	s.name = name
	return &s
}

func (s *RevokeDatabaseRoleRequest) WithDatabaseRole(databaseRole DatabaseObjectIdentifier) *RevokeDatabaseRoleRequest {
	s.accountRole = nil
	s.databaseRole = &databaseRole
	return s
}

func (s *RevokeDatabaseRoleRequest) WithAccountRole(accountRole AccountObjectIdentifier) *RevokeDatabaseRoleRequest {
	s.databaseRole = nil
	s.accountRole = &accountRole
	return s
}

func NewGrantDatabaseRoleToShareRequest(
	name DatabaseObjectIdentifier,
	share AccountObjectIdentifier,
) *GrantDatabaseRoleToShareRequest {
	s := GrantDatabaseRoleToShareRequest{}
	s.name = name
	s.share = share
	return &s
}

func NewRevokeDatabaseRoleFromShareRequest(
	name DatabaseObjectIdentifier,
	share AccountObjectIdentifier,
) *RevokeDatabaseRoleFromShareRequest {
	s := RevokeDatabaseRoleFromShareRequest{}
	s.name = name
	s.share = share
	return &s
}
