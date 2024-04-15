package acceptance

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckDestroy(t *testing.T, resourceType string) func(*terraform.State) error {
	t.Helper()
	client := Client(t)

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			ctx := context.Background()
			id := helpers.DecodeSnowflakeID(rs.Primary.Attributes["id"])
			if cleanups[resourceType](client, ctx, id) == nil {
				return fmt.Errorf("%s %v still exists", resourceType, id.FullyQualifiedName())
			}
		}
		return nil
	}
}

type cleanupFunc func(*sdk.Client, context.Context, sdk.ObjectIdentifier) error

var cleanups = map[string]cleanupFunc{
	"snowflake_view": func(client *sdk.Client, ctx context.Context, id sdk.ObjectIdentifier) error {
		return run[*sdk.View, sdk.SchemaObjectIdentifier](id, ctx, client.Views.ShowByID)
	},
}

func run[T any, U sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier](id sdk.ObjectIdentifier, ctx context.Context, show func(ctx context.Context, id U) (T, error)) error {
	idCast, err := asId[U](id)
	if err != nil {
		return err
	}
	_, err = show(ctx, *idCast)
	return err
}

func asId[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier](id sdk.ObjectIdentifier) (*T, error) {
	if idCast, ok := id.(T); !ok {
		return nil, fmt.Errorf("expected %s identifier type, but got: %T", reflect.TypeOf(new(T)).Elem().Name(), id)
	} else {
		return &idCast, nil
	}
}
