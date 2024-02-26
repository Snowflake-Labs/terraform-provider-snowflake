package sdk

import (
	"context"
	"fmt"
)

type SystemFunctions interface {
	GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) (string, error)
	GetTags(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) ([]string, error)
}

var _ SystemFunctions = (*systemFunctions)(nil)

type systemFunctions struct {
	client *Client
}

func (c *systemFunctions) GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) (string, error) {
	s := &struct {
		Tag string `db:"TAG"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$GET_TAG('%s', '%s', '%v') AS "TAG"`, tagID.FullyQualifiedName(), objectID.FullyQualifiedName(), objectType)
	err := c.client.queryOne(ctx, s, sql)
	if err != nil {
		return "", err
	}
	return s.Tag, nil
}

func (c *systemFunctions) GetTags(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) ([]string, error) {
	res := []struct {
		Tag string `db:"TAG"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$GET_TAG('%s', '%s', '%s') TAG_VALUE WHERE TAG_VALUE IS NOT NULL`, tagID.FullyQualifiedName(), objectID.FullyQualifiedName(), objectType)
	if err := c.client.query(ctx, res, sql); err != nil {
		return nil, err
	}
	tags := []string{}
	for _, item := range res {
		tags = append(tags, item.Tag)
	}
	return tags, nil
}
