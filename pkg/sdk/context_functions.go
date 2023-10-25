package sdk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type ContextFunctions interface {
	// Session functions.
	CurrentAccount(ctx context.Context) (string, error)
	CurrentRole(ctx context.Context) (string, error)
	CurrentSecondaryRoles(ctx context.Context) (*CurrentSecondaryRoles, error)
	CurrentRegion(ctx context.Context) (string, error)
	CurrentSession(ctx context.Context) (string, error)
	CurrentUser(ctx context.Context) (string, error)

	// Session Object functions.
	CurrentDatabase(ctx context.Context) (string, error)
	CurrentSchema(ctx context.Context) (string, error)
	CurrentWarehouse(ctx context.Context) (string, error)
	IsRoleInSession(ctx context.Context, role AccountObjectIdentifier) (bool, error)
}

var _ ContextFunctions = (*contextFunctions)(nil)

type contextFunctions struct {
	client *Client
}

func (c *contextFunctions) CurrentAccount(ctx context.Context) (string, error) {
	s := &struct {
		CurrentAccount string `db:"CURRENT_ACCOUNT"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_ACCOUNT() as CURRENT_ACCOUNT")
	if err != nil {
		return "", err
	}
	return s.CurrentAccount, nil
}

func (c *contextFunctions) CurrentRole(ctx context.Context) (string, error) {
	s := &struct {
		CurrentRole string `db:"CURRENT_ROLE"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_ROLE() as CURRENT_ROLE")
	if err != nil {
		return "", err
	}
	return s.CurrentRole, nil
}

type CurrentSecondaryRoles struct {
	Roles []AccountObjectIdentifier
	Value SecondaryRoleOption
}

func (c *contextFunctions) CurrentSecondaryRoles(ctx context.Context) (*CurrentSecondaryRoles, error) {
	s := &struct {
		CurrentRole string `db:"CURRENT_ROLES"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_SECONDARY_ROLES() as CURRENT_ROLES")
	if err != nil {
		return nil, err
	}

	jsonRoles := &struct {
		Roles string `json:"roles"`
		Value string `json:"value"`
	}{}
	err = json.Unmarshal([]byte(s.CurrentRole), jsonRoles)
	if err != nil {
		return nil, err
	}

	var roles []AccountObjectIdentifier
	if len(jsonRoles.Roles) > 0 {
		stringRoles := strings.Split(jsonRoles.Roles, ",")
		for _, role := range stringRoles {
			roles = append(roles, NewAccountObjectIdentifier(role))
		}
	}

	var value SecondaryRoleOption
	switch jsonRoles.Value {
	case "ALL":
		value = SecondaryRolesAll
	default:
		value = SecondaryRolesNone
	}
	secondaryRoles := &CurrentSecondaryRoles{
		Roles: roles,
		Value: value,
	}
	return secondaryRoles, nil
}

func (c *contextFunctions) CurrentRegion(ctx context.Context) (string, error) {
	s := &struct {
		CurrentRegion string `db:"CURRENT_REGION"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_REGION() AS CURRENT_REGION")
	if err != nil {
		return "", err
	}
	return s.CurrentRegion, nil
}

func (c *contextFunctions) CurrentSession(ctx context.Context) (string, error) {
	s := &struct {
		CurrentSession string `db:"CURRENT_SESSION"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_SESSION() as CURRENT_SESSION")
	if err != nil {
		return "", err
	}
	return s.CurrentSession, nil
}

func (c *contextFunctions) CurrentUser(ctx context.Context) (string, error) {
	s := &struct {
		CurrentUser string `db:"CURRENT_USER"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_USER() as CURRENT_USER")
	if err != nil {
		return "", err
	}
	return s.CurrentUser, nil
}

func (c *contextFunctions) CurrentDatabase(ctx context.Context) (string, error) {
	s := &struct {
		CurrentDatabase sql.NullString `db:"CURRENT_DATABASE"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_DATABASE() as CURRENT_DATABASE")
	if err != nil {
		return "", err
	}
	if !s.CurrentDatabase.Valid {
		return "", nil
	}
	return s.CurrentDatabase.String, nil
}

func (c *contextFunctions) CurrentSchema(ctx context.Context) (string, error) {
	s := &struct {
		CurrentSchema sql.NullString `db:"CURRENT_SCHEMA"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_SCHEMA() as CURRENT_SCHEMA")
	if err != nil {
		return "", err
	}
	if !s.CurrentSchema.Valid {
		return "", nil
	}
	return s.CurrentSchema.String, nil
}

func (c *contextFunctions) CurrentWarehouse(ctx context.Context) (string, error) {
	s := &struct {
		CurrentWarehouse sql.NullString `db:"CURRENT_WAREHOUSE"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_WAREHOUSE() as CURRENT_WAREHOUSE")
	if err != nil {
		return "", err
	}
	if !s.CurrentWarehouse.Valid {
		return "", nil
	}
	return s.CurrentWarehouse.String, nil
}

func (c *contextFunctions) IsRoleInSession(ctx context.Context, role AccountObjectIdentifier) (bool, error) {
	s := &struct {
		IsRoleInSession bool `db:"IS_ROLE_IN_SESSION"`
	}{}
	sql := fmt.Sprintf("SELECT IS_ROLE_IN_SESSION('%s') AS IS_ROLE_IN_SESSION", role.FullyQualifiedName())
	err := c.client.queryOne(ctx, s, sql)
	if err != nil {
		return false, err
	}
	return s.IsRoleInSession, nil
}
