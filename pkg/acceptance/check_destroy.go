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
			id := decodeSnowflakeId(rs, resource)
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

func decodeSnowflakeId(rs *terraform.ResourceState, resource resources.Resource) sdk.ObjectIdentifier {
	switch resource {
	case resources.ExternalFunction:
		return sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["id"])
	default:
		return helpers.DecodeSnowflakeID(rs.Primary.Attributes["id"])
	}
}

type showByIdFunc func(context.Context, *sdk.Client, sdk.ObjectIdentifier) error

var showByIdFunctions = map[resources.Resource]showByIdFunc{
	resources.Account: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Accounts.ShowByID)
	},
	resources.Alert: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Alerts.ShowByID)
	},
	resources.ApiIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Applications.ShowByID)
	},
	resources.Database: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Databases.ShowByID)
	},
	resources.DatabaseRole: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.DatabaseRoles.ShowByID)
	},
	resources.DynamicTable: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.DynamicTables.ShowByID)
	},
	resources.EmailNotificationIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.NotificationIntegrations.ShowByID)
	},
	resources.ExternalFunction: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ExternalFunctions.ShowByID)
	},
	resources.ExternalTable: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ExternalTables.ShowByID)
	},
	resources.FailoverGroup: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.FailoverGroups.ShowByID)
	},
	resources.Schema: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Schemas.ShowByID)
	},
	resources.Stage: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Stages.ShowByID)
	},
	resources.View: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Views.ShowByID)
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
