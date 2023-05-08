package sdk

import (
	"context"
	"fmt"
)

type SystemFunctions interface {
	GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier) (string, error)
}

var _ SystemFunctions = (*systemFunctions)(nil)

type systemFunctions struct {
	client  *Client
	builder *sqlBuilder
}

func (c *systemFunctions) GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier) (string, error) {
	s := &struct {
		Tag string `db:"TAG"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$GET_TAG('%s', '%s', '%v') AS "TAG"`, tagID.FullyQualifiedName(), objectID.FullyQualifiedName(), objectID.ObjectType)
	err := c.client.queryOne(ctx, s, sql)
	if err != nil {
		return "", err
	}
	return s.Tag, nil
}
