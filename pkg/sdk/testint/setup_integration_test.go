package testint

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-uuid"
)

var itc integrationTestContext

func TestMain(m *testing.M) {
	exitVal := execute(m)
	os.Exit(exitVal)
}

func execute(m *testing.M) int {
	defer timer("tests")()
	setup()
	exitVal := m.Run()
	cleanup()
	return exitVal
}

func setup() {
	log.Println("Running integration tests setup")

	err := itc.initialize()
	if err != nil {
		log.Printf("Integration test context initialisation failed with %s\n", err)
		cleanup()
		os.Exit(1)
	}
}

func cleanup() {
	log.Println("Running integration tests cleanup")
	if itc.databaseCleanup != nil {
		defer itc.databaseCleanup()
	}
}

type integrationTestContext struct {
	client *sdk.Client
	ctx    context.Context

	database        *sdk.Database
	databaseCleanup func()
}

func (itc *integrationTestContext) initialize() error {
	log.Println("Initializing integration test context")
	var err error
	itc.client, err = sdk.NewDefaultClient()
	itc.ctx = context.Background()

	db, dbCleanup, err := createDb(itc.client, itc.ctx)
	itc.database = db
	itc.databaseCleanup = dbCleanup

	return err
}

func createDb(client *sdk.Client, ctx context.Context) (*sdk.Database, func(), error) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, err
	}
	id := sdk.NewAccountObjectIdentifier("int_test_db_" + u)
	err = client.Databases.Create(ctx, id, nil)
	if err != nil {
		return nil, nil, err
	}
	database, err := client.Databases.ShowByID(ctx, id)
	return database, func() {
		_ = client.Databases.Drop(ctx, id, nil)
	}, err
}

// timer measures time from invocation point to the end of method.
// It's supposed to be used like:
//
//	defer timer("something to measure name")()
func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("[DEBUG] %s took %v\n", name, time.Since(start))
	}
}

// TODO: Discuss after this initial change is merged.
// This is temporary way to move all integration tests to this package without doing revolution in a single PR.
func testClient(t *testing.T) *sdk.Client {
	t.Helper()
	return itc.client
}

// TODO: Discuss after this initial change is merged.
// This is temporary way to move all integration tests to this package without doing revolution in a single PR.
func testContext(t *testing.T) context.Context {
	t.Helper()
	return itc.ctx
}

// TODO: Discuss after this initial change is merged.
// This is temporary way to move all integration tests to this package without doing revolution in a single PR.
func testDb(t *testing.T) *sdk.Database {
	t.Helper()
	return itc.database
}
