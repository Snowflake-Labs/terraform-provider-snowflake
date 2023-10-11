package internal

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func FindOne[T any](collection []T, condition func(T) bool) (*T, error) {
	for _, o := range collection {
		if condition(o) {
			return &o, nil
		}
	}
	return nil, sdk.ErrObjectNotExistOrAuthorized
}
