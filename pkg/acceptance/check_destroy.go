package acceptance

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
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
			if id == nil {
				return fmt.Errorf("could not get the id of %s", resource)
			}
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
		return sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(rs.Primary.ID)
	case resources.Function:
		return sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(rs.Primary.ID)
	case resources.Procedure:
		return sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(rs.Primary.ID)
	default:
		return helpers.DecodeSnowflakeID(rs.Primary.ID)
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
		return runShowById(ctx, id, client.ApiIntegrations.ShowByID)
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
	resources.FileFormat: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.FileFormats.ShowByID)
	},
	resources.Function: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Functions.ShowByID)
	},
	resources.ManagedAccount: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ManagedAccounts.ShowByID)
	},
	resources.MaskingPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.MaskingPolicies.ShowByID)
	},
	resources.MaterializedView: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.MaterializedViews.ShowByID)
	},
	resources.NetworkPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.NetworkPolicies.ShowByID)
	},
	resources.NetworkRule: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.NetworkRules.ShowByID)
	},
	resources.NotificationIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.NotificationIntegrations.ShowByID)
	},
	resources.PasswordPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.PasswordPolicies.ShowByID)
	},
	resources.Pipe: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Pipes.ShowByID)
	},
	resources.Procedure: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Procedures.ShowByID)
	},
	resources.ResourceMonitor: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ResourceMonitors.ShowByID)
	},
	resources.Role: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Roles.ShowByID)
	},
	resources.RowAccessPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.RowAccessPolicies.ShowByID)
	},
	resources.Schema: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Schemas.ShowByID)
	},
	resources.ScimSecurityIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.SecondaryDatabase: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Databases.ShowByID)
	},
	resources.Sequence: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Sequences.ShowByID)
	},
	resources.Share: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Shares.ShowByID)
	},
	resources.SharedDatabase: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Databases.ShowByID)
	},
	resources.Stage: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Stages.ShowByID)
	},
	resources.StorageIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.StorageIntegrations.ShowByID)
	},
	resources.Stream: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Streams.ShowByID)
	},
	resources.Table: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Tables.ShowByID)
	},
	resources.Tag: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Tags.ShowByID)
	},
	resources.Task: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Tasks.ShowByID)
	},
	resources.User: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Users.ShowByID)
	},
	resources.View: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Views.ShowByID)
	},
	resources.Warehouse: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Warehouses.ShowByID)
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

// CheckGrantAccountRoleDestroy is a custom checks that should be later incorporated into generic CheckDestroy
func CheckGrantAccountRoleDestroy(t *testing.T) func(*terraform.State) error {
	t.Helper()
	client := Client(t)

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_grant_account_role" {
				continue
			}
			ctx := context.Background()
			parts := strings.Split(rs.Primary.ID, "|")
			roleName := parts[0]
			roleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(roleName)
			objectType := parts[1]
			targetIdentifier := parts[2]
			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				Of: &sdk.ShowGrantsOf{
					Role: roleIdentifier,
				},
			})
			if err != nil {
				return nil
			}

			var found bool
			for _, grant := range grants {
				if grant.GrantedTo == sdk.ObjectType(objectType) {
					if grant.GranteeName.FullyQualifiedName() == targetIdentifier {
						found = true
						break
					}
				}
			}
			if found {
				return fmt.Errorf("role grant %v still exists", rs.Primary.ID)
			}
		}
		return nil
	}
}

// CheckGrantDatabaseRoleDestroy is a custom checks that should be later incorporated into generic CheckDestroy
func CheckGrantDatabaseRoleDestroy(t *testing.T) func(*terraform.State) error {
	t.Helper()
	client := Client(t)

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_grant_database_role" {
				continue
			}
			ctx := context.Background()
			id := rs.Primary.ID
			ids := strings.Split(id, "|")
			databaseRoleName := ids[0]
			objectType := ids[1]
			parentRoleName := ids[2]
			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				Of: &sdk.ShowGrantsOf{
					DatabaseRole: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(databaseRoleName),
				},
			})
			if err != nil {
				continue
			}
			for _, grant := range grants {
				if grant.GrantedTo == sdk.ObjectType(objectType) {
					if grant.GranteeName.FullyQualifiedName() == parentRoleName {
						return fmt.Errorf("database role grant %v still exists", grant)
					}
				}
			}
		}
		return nil
	}
}

// CheckAccountRolePrivilegesRevoked is a custom checks that should be later incorporated into generic CheckDestroy
func CheckAccountRolePrivilegesRevoked(t *testing.T) func(*terraform.State) error {
	t.Helper()
	client := Client(t)

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_grant_privileges_to_account_role" {
				continue
			}
			ctx := context.Background()

			id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["account_role_name"])
			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				To: &sdk.ShowGrantsTo{
					Role: id,
				},
			})
			if err != nil {
				if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
					continue
				}
				return err
			}
			var grantedPrivileges []string
			for _, grant := range grants {
				grantedPrivileges = append(grantedPrivileges, grant.Privilege)
			}
			if len(grantedPrivileges) > 0 {
				return fmt.Errorf("account role (%s) is still granted, granted privileges %v", id.FullyQualifiedName(), grantedPrivileges)
			}
		}
		return nil
	}
}

// CheckDatabaseRolePrivilegesRevoked is a custom checks that should be later incorporated into generic CheckDestroy
func CheckDatabaseRolePrivilegesRevoked(t *testing.T) func(*terraform.State) error {
	t.Helper()
	client := Client(t)

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_grant_privileges_to_database_role" {
				continue
			}
			ctx := context.Background()

			id := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["database_role_name"])
			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				To: &sdk.ShowGrantsTo{
					DatabaseRole: id,
				},
			})
			if err != nil {
				return err
			}
			var grantedPrivileges []string
			for _, grant := range grants {
				// usage is the default privilege available after creation (it won't be revoked)
				if grant.Privilege != "USAGE" {
					grantedPrivileges = append(grantedPrivileges, grant.Privilege)
				}
			}
			if len(grantedPrivileges) > 0 {
				return fmt.Errorf("database role (%s) is still granted, granted privileges %v", id.FullyQualifiedName(), grantedPrivileges)
			}
		}
		return nil
	}
}

// CheckSharePrivilegesRevoked is a custom checks that should be later incorporated into generic CheckDestroy
func CheckSharePrivilegesRevoked(t *testing.T) func(*terraform.State) error {
	t.Helper()
	client := Client(t)

	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "snowflake_grant_privileges_to_share" {
				continue
			}
			ctx := context.Background()

			id := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["to_share"])
			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				To: &sdk.ShowGrantsTo{
					Share: &sdk.ShowGrantsToShare{
						Name: sdk.NewAccountObjectIdentifier(id.Name()),
					},
				},
			})
			if err != nil {
				return err
			}
			var grantedPrivileges []string
			for _, grant := range grants {
				grantedPrivileges = append(grantedPrivileges, grant.Privilege)
			}
			if len(grantedPrivileges) > 0 {
				return fmt.Errorf("share (%s) is still granted with privileges: %v", id.FullyQualifiedName(), grantedPrivileges)
			}
		}
		return nil
	}
}

// CheckUserPasswordPolicyAttachmentDestroy is a custom checks that should be later incorporated into generic CheckDestroy
func CheckUserPasswordPolicyAttachmentDestroy(t *testing.T) func(*terraform.State) error {
	t.Helper()
	client := Client(t)

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_user_password_policy_attachment" {
				continue
			}
			ctx := context.Background()
			policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(
				sdk.NewAccountObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["user_name"]),
				sdk.PolicyEntityDomainUser,
			))
			if err != nil {
				if strings.Contains(err.Error(), "does not exist or not authorized") {
					// Note: this can happen if the Policy Reference or the User has been deleted as well; in this case, ignore the error
					continue
				}
				return err
			}
			if len(policyReferences) > 0 {
				return fmt.Errorf("user password policy attachment %v still exists", policyReferences[0].PolicyName)
			}
		}
		return nil
	}
}

func TestAccCheckGrantApplicationRoleDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_grant_application_role" {
			continue
		}
		ctx := context.Background()
		id := rs.Primary.ID
		ids := strings.Split(id, "|")
		applicationRoleName := ids[0]
		objectType := ids[1]
		parentRoleName := ids[2]
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{
				ApplicationRole: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(applicationRoleName),
			},
		})
		if err != nil {
			continue
		}
		for _, grant := range grants {
			if grant.GrantedTo == sdk.ObjectType(objectType) {
				if grant.GranteeName.FullyQualifiedName() == parentRoleName {
					return fmt.Errorf("application role grant %v still exists", grant)
				}
			}
		}
	}
	return nil
}
