// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

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

func (s *CreateDatabaseRoleRequest) WithComment(comment string) *CreateDatabaseRoleRequest {
	s.comment = &comment
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

func (s *AlterDatabaseRoleRequest) WithRename(rename DatabaseObjectIdentifier) *AlterDatabaseRoleRequest {
	s.rename = &rename
	return s
}

func (s *AlterDatabaseRoleRequest) WithSet(set DatabaseRoleSetRequest) *AlterDatabaseRoleRequest {
	s.set = &set
	return s
}

func (s *AlterDatabaseRoleRequest) WithUnset(unset DatabaseRoleUnsetRequest) *AlterDatabaseRoleRequest {
	s.unset = &unset
	return s
}

func (s *AlterDatabaseRoleRequest) WithSetTags(setTags []TagAssociation) *AlterDatabaseRoleRequest {
	s.setTags = setTags
	return s
}

func (s *AlterDatabaseRoleRequest) WithUnsetTags(unsetTags []ObjectIdentifier) *AlterDatabaseRoleRequest {
	s.unsetTags = unsetTags
	return s
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

func (s *ShowDatabaseRoleRequest) WithLike(like Like) *ShowDatabaseRoleRequest {
	s.like = &like
	return s
}

func (s *ShowDatabaseRoleRequest) WithLimit(limit LimitFrom) *ShowDatabaseRoleRequest {
	s.limit = &limit
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
	s.databaseRole = &databaseRole
	return s
}

func (s *GrantDatabaseRoleRequest) WithAccountRole(accountRole AccountObjectIdentifier) *GrantDatabaseRoleRequest {
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
	s.databaseRole = &databaseRole
	return s
}

func (s *RevokeDatabaseRoleRequest) WithAccountRole(accountRole AccountObjectIdentifier) *RevokeDatabaseRoleRequest {
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