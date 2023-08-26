package sdk

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AlertsShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	testWarehouse, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)

	alertTest, alertCleanup := createAlert(t, client, databaseTest, schemaTest, testWarehouse)
	t.Cleanup(alertCleanup)

	alert2Test, alert2Cleanup := createAlert(t, client, databaseTest, schemaTest, testWarehouse)
	t.Cleanup(alert2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		alerts, err := client.Alerts.Show(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, 2, len(alerts))
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &ShowAlertOptions{
			In: &In{
				Schema: schemaTest.ID(),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, alerts, alertTest)
		assert.Contains(t, alerts, alert2Test)
		assert.Equal(t, 2, len(alerts))
	})

	t.Run("with show options and like", func(t *testing.T) {
		showOptions := &ShowAlertOptions{
			Like: &Like{
				Pattern: String(alertTest.Name),
			},
			In: &In{
				Database: databaseTest.ID(),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, alerts, alertTest)
		assert.Equal(t, 1, len(alerts))
	})

	t.Run("when searching a non-existent alert", func(t *testing.T) {
		showOptions := &ShowAlertOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(alerts))
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &ShowAlertOptions{
			In: &In{
				Schema: schemaTest.ID(),
			},
			Limit: Int(1),
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
	})
}

func TestInt_AlertCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	testWarehouse, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)

	t.Run("test complete case", func(t *testing.T) {
		name := randomString(t)
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		comment := randomComment(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Alerts.Create(ctx, id, testWarehouse.ID(), schedule, condition, action, &CreateAlertOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(false),
			Comment:     String(comment),
		})
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testWarehouse.Name, alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, comment, *alertDetails.Comment)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(name),
			},
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alert))
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, comment, *alert[0].Comment)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		name := randomString(t)
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		comment := randomComment(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Alerts.Create(ctx, id, testWarehouse.ID(), schedule, condition, action, &CreateAlertOptions{
			OrReplace:   Bool(false),
			IfNotExists: Bool(true),
			Comment:     String(comment),
		})
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testWarehouse.Name, alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, comment, *alertDetails.Comment)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(name),
			},
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alert))
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, comment, *alert[0].Comment)
	})

	t.Run("test no options", func(t *testing.T) {
		name := randomString(t)
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Alerts.Create(ctx, id, testWarehouse.ID(), schedule, condition, action, nil)
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testWarehouse.Name, alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(name),
			},
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alert))
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, "", *alert[0].Comment)
	})

	t.Run("test multiline action", func(t *testing.T) {
		name := randomString(t)
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
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Alerts.Create(ctx, id, testWarehouse.ID(), schedule, condition, action, nil)
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testWarehouse.Name, alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, strings.TrimSpace(action), alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(name),
			},
			In: &In{
				Schema: schemaTest.ID(),
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
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)

	alert, alertCleanup := createAlert(t, client, databaseTest, schemaTest, warehouseTest)
	t.Cleanup(alertCleanup)

	t.Run("when alert exists", func(t *testing.T) {
		alertDetails, err := client.Alerts.Describe(ctx, alert.ID())
		require.NoError(t, err)
		assert.Equal(t, alert.Name, alertDetails.Name)
	})

	t.Run("when alert does not exist", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "does_not_exist")
		_, err := client.Alerts.Describe(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_AlertAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)

	t.Run("when setting and unsetting a value", func(t *testing.T) {
		alert, alertCleanup := createAlert(t, client, databaseTest, schemaTest, warehouseTest)
		t.Cleanup(alertCleanup)
		newSchedule := "USING CRON * * * * TUE,FRI GMT"

		alterOptions := &AlterAlertOptions{
			Set: &AlertSet{
				Schedule: &newSchedule,
			},
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(alert.Name),
			},
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.Equal(t, newSchedule, alerts[0].Schedule)
	})

	t.Run("when modifying condition and action", func(t *testing.T) {
		alert, alertCleanup := createAlert(t, client, databaseTest, schemaTest, warehouseTest)
		t.Cleanup(alertCleanup)
		newCondition := "select * from DUAL where false"

		alterOptions := &AlterAlertOptions{
			ModifyCondition: &[]string{newCondition},
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(alert.Name),
			},
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.Equal(t, newCondition, alerts[0].Condition)

		newAction := "create table FOO(ID INT)"

		alterOptions = &AlterAlertOptions{
			ModifyAction: &newAction,
		}

		err = client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err = client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(alert.Name),
			},
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.Equal(t, newAction, alerts[0].Action)
	})

	t.Run("resume and then suspend", func(t *testing.T) {
		alert, alertCleanup := createAlert(t, client, databaseTest, schemaTest, warehouseTest)
		t.Cleanup(alertCleanup)

		alterOptions := &AlterAlertOptions{
			Action: &AlertActionResume,
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(alert.Name),
			},
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.True(t, alerts[0].State == AlertStateStarted)

		alterOptions = &AlterAlertOptions{
			Action: &AlertActionSuspend,
		}

		err = client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err = client.Alerts.Show(ctx, &ShowAlertOptions{
			Like: &Like{
				Pattern: String(alert.Name),
			},
			In: &In{
				Schema: schemaTest.ID(),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(alerts))
		assert.True(t, alerts[0].State == AlertStateSuspended)
	})
}

func TestInt_AlertDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)

	t.Run("when alert exists", func(t *testing.T) {
		alert, _ := createAlert(t, client, databaseTest, schemaTest, warehouseTest)
		id := alert.ID()
		err := client.Alerts.Drop(ctx, id)
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("when alert does not exist", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "does_not_exist")
		err := client.Alerts.Drop(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}
