package sdk

import "context"

type ContextFunctions interface {
	CurrentSession(ctx context.Context) (string, error)
}

var _ ContextFunctions = (*contextFunctions)(nil)

type contextFunctions struct {
	client  *Client
	builder *sqlBuilder
}

func (c *contextFunctions) CurrentSession(ctx context.Context) (string, error) {
	s := &struct {
		CurrentSession string `db:"CURRENT_SESSION()"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_SESSION()")
	if err != nil {
		return "", err
	}
	return s.CurrentSession, nil
}
