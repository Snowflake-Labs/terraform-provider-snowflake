package testint

import (
	"errors"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AlertsShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	alertTest, alertCleanup := createAlert(t, client, testDb(t), testSchema(t), testWarehouse(t))
	t.Cleanup(alertCleanup)

	alert2Test, alert2Cleanup := createAlert(t, client, testDb(t), testSchema(t), testWarehouse(t))
	t.Cleanup(alert2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		alerts, err := client.Alerts.Show(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, 2, len(alerts))
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowAlertOptions{
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, alerts, *alertTest)
		assert.Contains(t, alerts, *alert2Test)
		assert.Equal(t, 2, len(alerts))
	})

	t.Run("with show options and like", func(t *testing.T) {
		showOptions := &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alertTest.Name),
			},
			In: &sdk.In{
				Database: testDb(t).ID(),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, alerts, *alertTest)
		assert.Equal(t, 1, len(alerts))
	})

	t.Run("when searching a non-existent alert", func(t *testing.T) {
		showOptions := &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(alerts))
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowAlertOptions{
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
			Limit: sdk.Int(1),
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
	})
}

func TestInt_AlertCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test complete case", func(t *testing.T) {
		name := random.String()
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		comment := random.Comment()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.Alerts.Create(ctx, id, testWarehouse(t).ID(), schedule, condition, action, &sdk.CreateAlertOptions{
			OrReplace:   sdk.Bool(true),
			IfNotExists: sdk.Bool(false),
			Comment:     sdk.String(comment),
		})
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testWarehouse(t).Name, alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, comment, *alertDetails.Comment)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alert))
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, comment, *alert[0].Comment)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		name := random.String()
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		comment := random.Comment()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.Alerts.Create(ctx, id, testWarehouse(t).ID(), schedule, condition, action, &sdk.CreateAlertOptions{
			OrReplace:   sdk.Bool(false),
			IfNotExists: sdk.Bool(true),
			Comment:     sdk.String(comment),
		})
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testWarehouse(t).Name, alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, comment, *alertDetails.Comment)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alert))
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, comment, *alert[0].Comment)
	})

	t.Run("test no options", func(t *testing.T) {
		name := random.String()
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.Alerts.Create(ctx, id, testWarehouse(t).ID(), schedule, condition, action, nil)
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testWarehouse(t).Name, alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alert))
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, "", *alert[0].Comment)
	})

	t.Run("test multiline action", func(t *testing.T) {
		name := random.String()
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := `
			select
				case
					when true then
						1
					else
						2
				end
		`
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.Alerts.Create(ctx, id, testWarehouse(t).ID(), schedule, condition, action, nil)
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testWarehouse(t).Name, alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, strings.TrimSpace(action), alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alert))
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, "", *alert[0].Comment)
	})
}

func TestInt_AlertDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	alert, alertCleanup := createAlert(t, client, testDb(t), testSchema(t), testWarehouse(t))
	t.Cleanup(alertCleanup)

	t.Run("when alert exists", func(t *testing.T) {
		alertDetails, err := client.Alerts.Describe(ctx, alert.ID())
		require.NoError(t, err)
		assert.Equal(t, alert.Name, alertDetails.Name)
	})

	t.Run("when alert does not exist", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")
		_, err := client.Alerts.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_AlertAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when setting and unsetting a value", func(t *testing.T) {
		alert, alertCleanup := createAlert(t, client, testDb(t), testSchema(t), testWarehouse(t))
		t.Cleanup(alertCleanup)
		newSchedule := "USING CRON * * * * TUE,FRI GMT"

		alterOptions := &sdk.AlterAlertOptions{
			Set: &sdk.AlertSet{
				Schedule: &newSchedule,
			},
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.Equal(t, newSchedule, alerts[0].Schedule)
	})

	t.Run("when modifying condition and action", func(t *testing.T) {
		alert, alertCleanup := createAlert(t, client, testDb(t), testSchema(t), testWarehouse(t))
		t.Cleanup(alertCleanup)
		newCondition := "select * from DUAL where false"

		alterOptions := &sdk.AlterAlertOptions{
			ModifyCondition: &[]string{newCondition},
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.Equal(t, newCondition, alerts[0].Condition)

		newAction := "create table FOO(ID INT)"

		alterOptions = &sdk.AlterAlertOptions{
			ModifyAction: &newAction,
		}

		err = client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err = client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.Equal(t, newAction, alerts[0].Action)
	})

	t.Run("resume and then suspend", func(t *testing.T) {
		alert, alertCleanup := createAlert(t, client, testDb(t), testSchema(t), testWarehouse(t))
		t.Cleanup(alertCleanup)

		alterOptions := &sdk.AlterAlertOptions{
			Action: &sdk.AlertActionResume,
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.True(t, alerts[0].State == sdk.AlertStateStarted)

		alterOptions = &sdk.AlterAlertOptions{
			Action: &sdk.AlertActionSuspend,
		}

		err = client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err = client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.True(t, alerts[0].State == sdk.AlertStateSuspended)
	})
}

func TestInt_AlertDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when alert exists", func(t *testing.T) {
		alert, _ := createAlert(t, client, testDb(t), testSchema(t), testWarehouse(t))
		id := alert.ID()
		err := client.Alerts.Drop(ctx, id)
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when alert does not exist", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")
		err := client.Alerts.Drop(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_AlertsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest, warehouseTest := testDb(t), testSchema(t), testWarehouse(t)
	cleanupAlertHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Alerts.Drop(ctx, id)
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createAlertHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		schedule, condition, action := "USING CRON * * * * * UTC", "SELECT 1", "SELECT 1"
		err := client.Alerts.Create(ctx, id, warehouseTest.ID(), schedule, condition, action, &sdk.CreateAlertOptions{})
		require.NoError(t, err)
		t.Cleanup(cleanupAlertHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaWithIdentifier(t, databaseTest, random.AlphaN(8))
		t.Cleanup(schemaCleanup)

		name := random.AlphaN(4)
		id1 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		id2 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, name)

		createAlertHandle(t, id1)
		createAlertHandle(t, id2)

		e1, err := client.Alerts.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Alerts.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
