package acceptance

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckDestroy(t *testing.T, resource resources.Resource) func(*terraform.State) error {
	t.Helper()
	client := Client(t)
	t.Logf("running check destroy for resource %s", resource)

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resource.String() {
				continue
			}
			t.Logf("found resource %s in state", resource)
			ctx := context.Background()
			id := helpers.DecodeSnowflakeID(rs.Primary.Attributes["id"])
			showById, ok := showByIdFunctions[resource]
			if !ok {
				return fmt.Errorf("unsupported show by id in cleanup for %s, with id %v", resource, id.FullyQualifiedName())
			}
			if showById(ctx, client, id) == nil {
				return fmt.Errorf("%s %v still exists", resource, id.FullyQualifiedName())
			} else {
				t.Logf("resource %s (%v) was dropped successfully in Snowflake", resource, id.FullyQualifiedName())
			}
		}
		return nil
	}
}

type showByIdFunc func(context.Context, *sdk.Client, sdk.ObjectIdentifier) error

var showByIdFunctions = map[resources.Resource]showByIdFunc{
	resources.View: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Views.ShowByID)
	},
	resources.Schema: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Schemas.ShowByID)
	},
}

func runShowById[T any, U sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier](ctx context.Context, id sdk.ObjectIdentifier, show func(ctx context.Context, id U) (T, error)) error {
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
