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

func TestSweepAll(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)
	testenvs.AssertEnvSet(t, string(testenvs.TestObjectsSuffix))

	t.Run("sweep after tests", func(t *testing.T) {
		client := defaultTestClient(t)
		secondaryClient := testSecondaryClient(t)

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
	secondaryClient := testSecondaryClient(t)
	thirdClient := testThirdClient(t)
	fourthClient := testFourthClient(t)

	allClients := []*Client{client, secondaryClient, thirdClient, fourthClient}

	t.Run("sweep integration test precreated objects", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeWarehouses(c, "int_test_wh_%")()
			assert.NoError(t, err)

			err = nukeDatabases(c, "int_test_db_%")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep acceptance tests precreated objects", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeWarehouses(c, "acc_test_wh_%")()
			assert.NoError(t, err)

			err = nukeDatabases(c, "acc_test_db_%")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep users", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeUsers(c)()
			assert.NoError(t, err)
		}
	})

	// TODO [SNOW-955520]: nuke stale objects (e.g. created more than 2 weeks ago)
}

// TODO [SNOW-955520]: generalize nuke methods (sweepers too)
func nukeWarehouses(client *Client, prefix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Nuking warehouses with prefix %s\n", prefix)
		ctx := context.Background()

		whs, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{Like: &Like{Pattern: String(prefix)}})
		if err != nil {
			return fmt.Errorf("sweeping warehouses ended with error, err = %w", err)
		}
		var errs []error
		log.Printf("[DEBUG] Found %d warehouses matching search criteria\n", len(whs))
		for idx, wh := range whs {
			log.Printf("[DEBUG] Processing warehouse [%d/%d]: %s...\n", idx+1, len(whs), wh.ID().FullyQualifiedName())
			if wh.Name != "SNOWFLAKE" && wh.CreatedOn.Before(time.Now().Add(-4*time.Hour)) {
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
	return func() error {
		log.Printf("[DEBUG] Nuking databases with prefix %s\n", prefix)
		ctx := context.Background()

		dbs, err := client.Databases.Show(ctx, &ShowDatabasesOptions{Like: &Like{Pattern: String(prefix)}})
		if err != nil {
			return fmt.Errorf("sweeping databases ended with error, err = %w", err)
		}
		var errs []error
		log.Printf("[DEBUG] Found %d databases matching search criteria\n", len(dbs))
		for idx, db := range dbs {
			log.Printf("[DEBUG] Processing database [%d/%d]: %s...\n", idx+1, len(dbs), db.ID().FullyQualifiedName())
			if db.Name != "SNOWFLAKE" && db.CreatedOn.Before(time.Now().Add(-4*time.Hour)) {
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
			if !slices.Contains(protectedUsers, user.Name) && user.CreatedOn.Before(time.Now().Add(-2*time.Hour)) {
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
