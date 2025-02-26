package datasources_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tasks_Like_RootTask(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	// Created to show LIKE is working
	_, standaloneTaskCleanup := acc.TestClient().Task.Create(t)
	t.Cleanup(standaloneTaskCleanup)

	createRootReq := sdk.NewCreateTaskRequest(acc.TestClient().Ids.RandomSchemaObjectIdentifier(), "SELECT 1").
		WithSchedule("1 MINUTE").
		WithComment("some comment").
		WithAllowOverlappingExecution(true).
		WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(acc.TestClient().Ids.WarehouseId()))
	rootTask, rootTaskCleanup := acc.TestClient().Task.CreateWithRequest(t, createRootReq)
	t.Cleanup(rootTaskCleanup)

	childTask, childTaskCleanup := acc.TestClient().Task.CreateWithAfter(t, rootTask.ID())
	t.Cleanup(childTaskCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: taskDatasourceLikeRootOnly(rootTask.ID().Name(), true),
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_tasks.test", "tasks.#", "1")),
					resourceshowoutputassert.TaskDatasourceShowOutput(t, "snowflake_tasks.test").
						HasName(rootTask.Name).
						HasSchemaName(rootTask.SchemaName).
						HasDatabaseName(rootTask.DatabaseName).
						HasCreatedOnNotEmpty().
						HasIdNotEmpty().
						HasOwnerNotEmpty().
						HasComment("some comment").
						HasWarehouse(acc.TestClient().Ids.WarehouseId()).
						HasSchedule("1 MINUTE").
						HasPredecessors().
						HasDefinition("SELECT 1").
						HasCondition("").
						HasAllowOverlappingExecution(true).
						HasErrorIntegrationEmpty().
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}).
						HasLastSuspendedReason(""),
					resourceparametersassert.TaskDatasourceParameters(t, "snowflake_tasks.test").
						HasAllDefaults(),
				),
			},
			{
				Config: taskDatasourceLikeRootOnly(childTask.ID().Name(), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tasks.test", "tasks.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Tasks_In_StartsWith(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := acc.TestClient().Ids.AlphaN(4)

	_, standaloneTaskCleanup := acc.TestClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix), "SELECT 1"))
	t.Cleanup(standaloneTaskCleanup)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	standaloneTask2, standaloneTask2Cleanup := acc.TestClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID()), "SELECT 1"))
	t.Cleanup(standaloneTask2Cleanup)

	_, standaloneTask3Cleanup := acc.TestClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()), "SELECT 1"))
	t.Cleanup(standaloneTask3Cleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// On account with prefix
			{
				Config: taskDatasourceOnAccountStartsWith(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tasks.test", "tasks.#", "2"),
				),
			},
			// On database with prefix
			{
				Config: taskDatasourceInDatabaseStartsWith(acc.TestClient().Ids.DatabaseId(), prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tasks.test", "tasks.#", "2"),
				),
			},
			// On schema with prefix
			{
				Config: taskDatasourceInSchemaStartsWith(schema.ID(), prefix),
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_tasks.test", "tasks.#", "1")),
					resourceshowoutputassert.TaskDatasourceShowOutput(t, "snowflake_tasks.test").
						HasName(standaloneTask2.Name).
						HasSchemaName(standaloneTask2.SchemaName).
						HasDatabaseName(standaloneTask2.DatabaseName),
				),
			},
			// On schema
			{
				Config: taskDatasourceInSchema(schema.ID()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tasks.test", "tasks.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Tasks_Limit(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := acc.TestClient().Ids.AlphaN(4)

	_, standaloneTaskCleanup := acc.TestClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix), "SELECT 1"))
	t.Cleanup(standaloneTaskCleanup)

	_, standaloneTask2Cleanup := acc.TestClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix), "SELECT 1"))
	t.Cleanup(standaloneTask2Cleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Limit with prefix
			{
				Config: taskDatasourceLimitWithPrefix(2, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tasks.test", "tasks.#", "2"),
				),
			},
			// Only limit
			{
				Config: taskDatasourceLimit(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_tasks.test", "tasks.#", "1"),
				),
			},
		},
	})
}

func taskDatasourceLikeRootOnly(like string, rootOnly bool) string {
	return taskDatasourceConfig(like, false, sdk.AccountObjectIdentifier{}, sdk.DatabaseObjectIdentifier{}, "", rootOnly, nil)
}

func taskDatasourceOnAccountStartsWith(startsWith string) string {
	return taskDatasourceConfig("", true, sdk.AccountObjectIdentifier{}, sdk.DatabaseObjectIdentifier{}, startsWith, false, nil)
}

func taskDatasourceInDatabaseStartsWith(databaseId sdk.AccountObjectIdentifier, startsWith string) string {
	return taskDatasourceConfig("", false, databaseId, sdk.DatabaseObjectIdentifier{}, startsWith, false, nil)
}

func taskDatasourceInSchemaStartsWith(schemaId sdk.DatabaseObjectIdentifier, startsWith string) string {
	return taskDatasourceConfig("", false, sdk.AccountObjectIdentifier{}, schemaId, startsWith, false, nil)
}

func taskDatasourceInSchema(schemaId sdk.DatabaseObjectIdentifier) string {
	return taskDatasourceConfig("", false, sdk.AccountObjectIdentifier{}, schemaId, "", false, nil)
}

func taskDatasourceLimit(limit int) string {
	return taskDatasourceConfig("", false, sdk.AccountObjectIdentifier{}, sdk.DatabaseObjectIdentifier{}, "", false, &sdk.LimitFrom{
		Rows: sdk.Int(limit),
	})
}

func taskDatasourceLimitWithPrefix(limit int, prefix string) string {
	return taskDatasourceConfig("", false, sdk.AccountObjectIdentifier{}, sdk.DatabaseObjectIdentifier{}, "", false, &sdk.LimitFrom{
		Rows: sdk.Int(limit),
		From: sdk.String(prefix),
	})
}

func taskDatasourceConfig(like string, onAccount bool, onDatabase sdk.AccountObjectIdentifier, onSchema sdk.DatabaseObjectIdentifier, startsWith string, rootOnly bool, limitFrom *sdk.LimitFrom) string {
	var likeString string
	if len(like) > 0 {
		likeString = fmt.Sprintf("like = \"%s\"", like)
	}

	var startsWithString string
	if len(startsWith) > 0 {
		startsWithString = fmt.Sprintf("starts_with = \"%s\"", startsWith)
	}

	var inString string
	if onAccount || (onDatabase != sdk.AccountObjectIdentifier{}) || (onSchema != sdk.DatabaseObjectIdentifier{}) {
		inStringBuffer := new(bytes.Buffer)
		inStringBuffer.WriteString("in {\n")
		switch {
		case onAccount:
			inStringBuffer.WriteString("account = true\n")
		case onDatabase != sdk.AccountObjectIdentifier{}:
			inStringBuffer.WriteString(fmt.Sprintf("database = %s\n", strconv.Quote(onDatabase.FullyQualifiedName())))
		case onSchema != sdk.DatabaseObjectIdentifier{}:
			inStringBuffer.WriteString(fmt.Sprintf("schema = %s\n", strconv.Quote(onSchema.FullyQualifiedName())))
		}
		inStringBuffer.WriteString("}\n")
		inString = inStringBuffer.String()
	}

	var rootOnlyString string
	if rootOnly {
		rootOnlyString = fmt.Sprintf("root_only = %t", rootOnly)
	}

	var limitFromString string
	if limitFrom != nil {
		inStringBuffer := new(bytes.Buffer)
		inStringBuffer.WriteString("limit {\n")
		inStringBuffer.WriteString(fmt.Sprintf("rows = %d\n", *limitFrom.Rows))
		if limitFrom.From != nil {
			inStringBuffer.WriteString(fmt.Sprintf("from = \"%s\"\n", *limitFrom.From))
		}
		inStringBuffer.WriteString("}\n")
		limitFromString = inStringBuffer.String()
	}

	return fmt.Sprintf(`
	data "snowflake_tasks" "test" {
		%[1]s
		%[2]s
		%[3]s
		%[4]s
		%[5]s
	}`, likeString, inString, startsWithString, rootOnlyString, limitFromString)
}
