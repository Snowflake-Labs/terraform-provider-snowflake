package sdk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/assert"
)

// TODO [SNOW-867247]: move the sweepers outside of the sdk package
// TODO [SNOW-867247]: use test client helpers in sweepers?
func TestSweepAll(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)
	testenvs.AssertEnvSet(t, string(testenvs.TestObjectsSuffix))

	t.Run("sweep after tests", func(t *testing.T) {
		client := defaultTestClient(t)
		secondaryClient := secondaryTestClient(t)

		err := SweepAfterIntegrationTests(client, random.IntegrationTestsSuffix)
		assert.NoError(t, err)

		err = SweepAfterIntegrationTests(secondaryClient, random.IntegrationTestsSuffix)
		assert.NoError(t, err)

		err = SweepAfterAcceptanceTests(client, random.AcceptanceTestsSuffix)
		assert.NoError(t, err)

		err = SweepAfterAcceptanceTests(secondaryClient, random.AcceptanceTestsSuffix)
		assert.NoError(t, err)
	})
}

func Test_Sweeper_NukeStaleObjects(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)

	client := defaultTestClient(t)
	secondaryClient := secondaryTestClient(t)
	thirdClient := thirdTestClient(t)
	fourthClient := fourthTestClient(t)

	allClients := []*Client{client, secondaryClient, thirdClient, fourthClient}

	// can't use extracted IntegrationTestPrefix and AcceptanceTestPrefix until sweepers reside in the SDK package (cyclic)
	const integrationTestPrefix = "int_test_"
	const acceptanceTestPrefix = "acc_test_"

	t.Run("sweep integration test precreated objects", func(t *testing.T) {
		integrationTestWarehousesPrefix := fmt.Sprintf("%swh_%%", integrationTestPrefix)
		integrationTestDatabasesPrefix := fmt.Sprintf("%sdb_%%", integrationTestPrefix)

		for _, c := range allClients {
			err := nukeWarehouses(c, integrationTestWarehousesPrefix)()
			assert.NoError(t, err)

			err = nukeDatabases(c, integrationTestDatabasesPrefix)()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep acceptance tests precreated objects", func(t *testing.T) {
		acceptanceTestWarehousesPrefix := fmt.Sprintf("%swh_%%", acceptanceTestPrefix)
		acceptanceTestDatabasesPrefix := fmt.Sprintf("%sdb_%%", acceptanceTestPrefix)

		for _, c := range allClients {
			err := nukeWarehouses(c, acceptanceTestWarehousesPrefix)()
			assert.NoError(t, err)

			err = nukeDatabases(c, acceptanceTestDatabasesPrefix)()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep users", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeUsers(c)()
			assert.NoError(t, err)
		}
	})

	// TODO [SNOW-867247]: unskip
	t.Run("sweep databases", func(t *testing.T) {
		t.Skipf("Used for manual sweeping; will be addressed during SNOW-867247")
		for _, c := range allClients {
			err := nukeDatabases(c, "")()
			assert.NoError(t, err)
		}
	})

	// TODO [SNOW-867247]: unskip
	t.Run("sweep warehouses", func(t *testing.T) {
		t.Skipf("Used for manual sweeping; will be addressed during SNOW-867247")
		for _, c := range allClients {
			err := nukeWarehouses(c, "")()
			assert.NoError(t, err)
		}
	})

	// TODO [SNOW-867247]: nuke stale objects (e.g. created more than 2 weeks ago)

	// TODO [SNOW-867247]: nuke external oauth integrations because of errors like
	// Error: 003524 (22023): SQL execution error: An integration with the given issuer already exists for this account
}

// TODO [SNOW-867247]: generalize nuke methods (sweepers too)
// TODO [SNOW-1658402]: handle the ownership problem while handling the better role setup for tests
func nukeWarehouses(client *Client, prefix string) func() error {
	protectedWarehouses := []string{
		"SNOWFLAKE",
		"SYSTEM$STREAMLIT_NOTEBOOK_WH",
	}

	return func() error {
		log.Printf("[DEBUG] Nuking warehouses with prefix %s\n", prefix)
		ctx := context.Background()

		var like *Like = nil
		if prefix != "" {
			like = &Like{Pattern: String(prefix)}
		}

		whs, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{Like: like})
		if err != nil {
			return fmt.Errorf("sweeping warehouses ended with error, err = %w", err)
		}
		var errs []error
		log.Printf("[DEBUG] Found %d warehouses matching search criteria\n", len(whs))
		for idx, wh := range whs {
			log.Printf("[DEBUG] Processing warehouse [%d/%d]: %s...\n", idx+1, len(whs), wh.ID().FullyQualifiedName())
			if !slices.Contains(protectedWarehouses, wh.Name) && wh.CreatedOn.Before(time.Now().Add(-2*time.Hour)) {
				if wh.Owner != "ACCOUNTADMIN" {
					log.Printf("[DEBUG] Granting ownership on warehouse %s, to ACCOUNTADMIN", wh.ID().FullyQualifiedName())
					err := client.Grants.GrantOwnership(
						ctx,
						OwnershipGrantOn{Object: &Object{
							ObjectType: ObjectTypeWarehouse,
							Name:       wh.ID(),
						}},
						OwnershipGrantTo{
							AccountRoleName: Pointer(NewAccountObjectIdentifier("ACCOUNTADMIN")),
						},
						nil,
					)
					if err != nil {
						errs = append(errs, fmt.Errorf("granting ownership on warehouse %s ended with error, err = %w", wh.ID().FullyQualifiedName(), err))
						continue
					}
				}

				log.Printf("[DEBUG] Dropping warehouse %s, created at: %s\n", wh.ID().FullyQualifiedName(), wh.CreatedOn.String())
				if err := client.Warehouses.Drop(ctx, wh.ID(), &DropWarehouseOptions{IfExists: Bool(true)}); err != nil {
					log.Printf("[DEBUG] Dropping warehouse %s, resulted in error %v\n", wh.ID().FullyQualifiedName(), err)
					errs = append(errs, fmt.Errorf("sweeping warehouse %s ended with error, err = %w", wh.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping warehouse %s, created at: %s\n", wh.ID().FullyQualifiedName(), wh.CreatedOn.String())
			}
		}
		return errors.Join(errs...)
	}
}

func nukeDatabases(client *Client, prefix string) func() error {
	protectedDatabases := []string{
		"SNOWFLAKE",
		"MFA_ENFORCEMENT_POLICY",
	}

	return func() error {
		log.Printf("[DEBUG] Nuking databases with prefix %s\n", prefix)
		ctx := context.Background()

		var like *Like = nil
		if prefix != "" {
			like = &Like{Pattern: String(prefix)}
		}
		dbs, err := client.Databases.Show(ctx, &ShowDatabasesOptions{Like: like})
		if err != nil {
			return fmt.Errorf("sweeping databases ended with error, err = %w", err)
		}
		var errs []error
		log.Printf("[DEBUG] Found %d databases matching search criteria\n", len(dbs))
		for idx, db := range dbs {
			if db.Owner != "ACCOUNTADMIN" {
				log.Printf("[DEBUG] Granting ownership on database %s, to ACCOUNTADMIN", db.ID().FullyQualifiedName())
				err := client.Grants.GrantOwnership(
					ctx,
					OwnershipGrantOn{Object: &Object{
						ObjectType: ObjectTypeDatabase,
						Name:       db.ID(),
					}},
					OwnershipGrantTo{
						AccountRoleName: Pointer(NewAccountObjectIdentifier("ACCOUNTADMIN")),
					},
					nil,
				)
				if err != nil {
					errs = append(errs, fmt.Errorf("granting ownership on database %s ended with error, err = %w", db.ID().FullyQualifiedName(), err))
					continue
				}
			}

			log.Printf("[DEBUG] Processing database [%d/%d]: %s...\n", idx+1, len(dbs), db.ID().FullyQualifiedName())
			if !slices.Contains(protectedDatabases, db.Name) && db.CreatedOn.Before(time.Now().Add(-2*time.Hour)) {
				log.Printf("[DEBUG] Dropping database %s, created at: %s\n", db.ID().FullyQualifiedName(), db.CreatedOn.String())
				if err := client.Databases.Drop(ctx, db.ID(), &DropDatabaseOptions{IfExists: Bool(true)}); err != nil {
					log.Printf("[DEBUG] Dropping database %s, resulted in error %v\n", db.ID().FullyQualifiedName(), err)
					errs = append(errs, fmt.Errorf("sweeping database %s ended with error, err = %w", db.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping database %s, created at: %s\n", db.ID().FullyQualifiedName(), db.CreatedOn.String())
			}
		}
		return errors.Join(errs...)
	}
}

func nukeUsers(client *Client) func() error {
	protectedUsers := []string{
		"SNOWFLAKE",
		"ARTUR_SAWICKI",
		"ARTUR_SAWICKI_LEGACY",
		"JAKUB_MICHALAK",
		"JAKUB_MICHALAK_LEGACY",
		"JAN_CIESLAK",
		"JAN_CIESLAK_LEGACY",
		"TERRAFORM_SVC_ACCOUNT",
		"TEST_CI_SERVICE_USER",
		"FILIP_BUDZYNSKI",
		"FILIP_BUDZYNSKI_LEGACY",
	}

	return func() error {
		log.Println("[DEBUG] Nuking users")
		ctx := context.Background()

		users, err := client.Users.Show(ctx, &ShowUserOptions{})
		if err != nil {
			return fmt.Errorf("sweeping users ended with error, err = %w", err)
		}
		var errs []error
		log.Printf("[DEBUG] Found %d users\n", len(users))
		for idx, user := range users {
			log.Printf("[DEBUG] Processing user [%d/%d]: %s...\n", idx+1, len(users), user.ID().FullyQualifiedName())
			if !slices.Contains(protectedUsers, user.Name) && user.CreatedOn.Before(time.Now().Add(-15*time.Minute)) {
				log.Printf("[DEBUG] Dropping user %s\n", user.ID().FullyQualifiedName())
				if err := client.Users.Drop(ctx, user.ID(), &DropUserOptions{IfExists: Bool(true)}); err != nil {
					log.Printf("[DEBUG] Dropping user %s, resulted in error %v\n", user.ID().FullyQualifiedName(), err)
					errs = append(errs, fmt.Errorf("sweeping user %s ended with error, err = %w", user.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping user %s\n", user.ID().FullyQualifiedName())
			}
		}
		return errors.Join(errs...)
	}
}
