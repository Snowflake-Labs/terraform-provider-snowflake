package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ResourceMonitorsShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	resourceMonitorTest, resourceMonitorCleanup := createResourceMonitor(t, client)
	t.Cleanup(resourceMonitorCleanup)

	_, resourceMonitorCleanup2 := createResourceMonitor(t, client)
	t.Cleanup(resourceMonitorCleanup2)

	t.Run("without show options", func(t *testing.T) {
		useDatabaseCleanup := useDatabase(t, client, databaseTest.ID())
		t.Cleanup(useDatabaseCleanup)
		useSchemaCleanup := useSchema(t, client, schemaTest.ID())
		t.Cleanup(useSchemaCleanup)

		resourceMonitors, err := client.ResourceMonitors.Show(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, 2, len(resourceMonitors))
	})

	t.Run("with like", func(t *testing.T) {
		showOptions := &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String(resourceMonitorTest.Name),
			},
		}
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, resourceMonitors, resourceMonitorTest)
		assert.Equal(t, 1, len(resourceMonitors))
	})

	t.Run("when searching a non-existent resource monitor", func(t *testing.T) {
		showOptions := &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(resourceMonitors))
	})
}

func TestInt_ResourceMonitorCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("test complete case", func(t *testing.T) {
		name := randomString(t)
		id := NewAccountObjectIdentifier(name)
		frequency, err := FrequencyFromString("Monthly")
		require.NoError(t, err)
		startTimeStamp := "IMMEDIATELY"
		creditQuota := 100
		endTimeStamp := "2024-01-01 12:34"

		triggers := []TriggerDefinition{
			{
				Threshold:     50,
				TriggerAction: SuspendImmediate,
			},
			{
				Threshold:     100,
				TriggerAction: Notify,
			},
		}
		err = client.ResourceMonitors.Create(ctx, id, &CreateResourceMonitorOptions{
			OrReplace:      Bool(true),
			Frequency:      frequency,
			CreditQuota:    &creditQuota,
			StartTimeStamp: &startTimeStamp,
			EndTimeStamp:   &endTimeStamp,
			// Users' emails need to be verified in order to use them for notification
			NotifyUsers: nil,
			Triggers:    &triggers,
		})

		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String(name),
			},
		})

		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor := resourceMonitors[0]
		require.NoError(t, err)
		assert.Equal(t, name, resourceMonitor.Name)
		assert.Equal(t, frequency, resourceMonitor.Frequency)
		assert.Equal(t, creditQuota, int(*resourceMonitor.CreditQuota))
		assert.NotEmpty(t, resourceMonitor.StartTime)
		assert.NotEmpty(t, resourceMonitor.EndTime)
		allTriggers := resourceMonitor.SuspendTriggers
		allTriggers = append(allTriggers, resourceMonitor.SuspendImmediateTriggers...)
		assert.Equal(t, creditQuota, int(*resourceMonitor.CreditQuota))
		allTriggers = append(allTriggers, resourceMonitor.NotifyTriggers...)
		assert.Equal(t, triggers, allTriggers)

		t.Cleanup(func() {
			err = client.ResourceMonitors.Drop(ctx, id)
			require.NoError(t, err)
		})
	})

	t.Run("test no options", func(t *testing.T) {
		name := randomString(t)
		id := NewAccountObjectIdentifier(name)

		err := client.ResourceMonitors.Create(ctx, id, nil)

		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String(name),
			},
		})

		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor := resourceMonitors[0]
		require.NoError(t, err)
		assert.Equal(t, name, resourceMonitor.Name)
		assert.NotEmpty(t, resourceMonitor.StartTime)
		assert.Nil(t, resourceMonitor.EndTime)
		assert.Nil(t, resourceMonitor.CreditQuota)
		assert.Equal(t, Monthly, *resourceMonitor.Frequency)
		assert.Empty(t, resourceMonitor.NotifyUsers)
		assert.Empty(t, resourceMonitor.NotifyTriggers)
		assert.Empty(t, resourceMonitor.SuspendImmediateTriggers)
		assert.Empty(t, resourceMonitor.SuspendTriggers)

		t.Cleanup(func() {
			err = client.ResourceMonitors.Drop(ctx, id)
			require.NoError(t, err)
		})
	})
}

func TestInt_ResourceMonitorAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("when adding a new trigger", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := createResourceMonitor(t, client)
		t.Cleanup(resourceMonitorCleanup)

		oldTriggers := []TriggerDefinition{}
		oldTriggers = append(oldTriggers, resourceMonitor.NotifyTriggers...)
		oldTriggers = append(oldTriggers, resourceMonitor.SuspendTriggers...)
		oldTriggers = append(oldTriggers, resourceMonitor.SuspendImmediateTriggers...)
		newTriggers := oldTriggers
		newTriggers = append(newTriggers, TriggerDefinition{Threshold: 30, TriggerAction: Notify})
		alterOptions := &AlterResourceMonitorOptions{
			Triggers: &newTriggers,
		}
		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), alterOptions)
		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String(resourceMonitor.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor = resourceMonitors[0]
		allTriggers := resourceMonitor.SuspendImmediateTriggers
		allTriggers = append(allTriggers, resourceMonitor.NotifyTriggers...)
		allTriggers = append(allTriggers, resourceMonitor.SuspendTriggers...)
		assert.ElementsMatch(t, newTriggers, allTriggers)
	})

	t.Run("when setting credit quota", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := createResourceMonitor(t, client)
		t.Cleanup(resourceMonitorCleanup)
		creditQuota := 100
		alterOptions := &AlterResourceMonitorOptions{
			Set: &ResourceMonitorSet{
				CreditQuota: &creditQuota,
			},
		}
		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), alterOptions)
		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String(resourceMonitor.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor = resourceMonitors[0]
		assert.Equal(t, creditQuota, int(*resourceMonitor.CreditQuota))
	})
}

func TestInt_ResourceMonitorDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("when resource monitor exists", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := createResourceMonitor(t, client)
		t.Cleanup(resourceMonitorCleanup)
		id := resourceMonitor.ID()
		err := client.ResourceMonitors.Drop(ctx, id)
		require.NoError(t, err)
		_, err = client.ResourceMonitors.ShowByID(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("when resource monitor does not exist", func(t *testing.T) {
		id := NewAccountObjectIdentifier("does_not_exist")
		err := client.ResourceMonitors.Drop(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}
