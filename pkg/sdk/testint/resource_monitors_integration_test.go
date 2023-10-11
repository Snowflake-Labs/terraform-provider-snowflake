package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ResourceMonitorsShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	resourceMonitorTest, resourceMonitorCleanup := createResourceMonitor(t, client)
	t.Cleanup(resourceMonitorCleanup)

	t.Run("with like", func(t *testing.T) {
		showOptions := &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(resourceMonitorTest.Name),
			},
		}
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, resourceMonitors, *resourceMonitorTest)
		assert.Equal(t, 1, len(resourceMonitors))
	})

	t.Run("when searching a non-existent resource monitor", func(t *testing.T) {
		showOptions := &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(resourceMonitors))
	})
}

func TestInt_ResourceMonitorCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test complete case", func(t *testing.T) {
		name := random.RandomString(t)
		id := sdk.NewAccountObjectIdentifier(name)
		frequency, err := sdk.FrequencyFromString("Monthly")
		require.NoError(t, err)
		startTimeStamp := "IMMEDIATELY"
		creditQuota := 100
		endTimeStamp := "2024-01-01 12:34"

		triggers := []sdk.TriggerDefinition{
			{
				Threshold:     30,
				TriggerAction: sdk.TriggerActionSuspend,
			},
			{
				Threshold:     50,
				TriggerAction: sdk.TriggerActionSuspendImmediate,
			},
			{
				Threshold:     100,
				TriggerAction: sdk.TriggerActionNotify,
			},
		}
		err = client.ResourceMonitors.Create(ctx, id, &sdk.CreateResourceMonitorOptions{
			OrReplace: sdk.Bool(true),
			With: &sdk.ResourceMonitorWith{
				Frequency:      frequency,
				CreditQuota:    &creditQuota,
				StartTimestamp: &startTimeStamp,
				EndTimestamp:   &endTimeStamp,
				// Users' emails need to be verified in order to use them for notification
				NotifyUsers: nil,
				Triggers:    triggers,
			},
		})

		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
		})

		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor := resourceMonitors[0]
		require.NoError(t, err)
		assert.Equal(t, name, resourceMonitor.Name)
		assert.Equal(t, *frequency, resourceMonitor.Frequency)
		assert.Equal(t, creditQuota, int(resourceMonitor.CreditQuota))
		assert.NotEmpty(t, resourceMonitor.StartTime)
		assert.NotEmpty(t, resourceMonitor.EndTime)
		assert.Equal(t, creditQuota, int(resourceMonitor.CreditQuota))
		var allThresholds []int
		allThresholds = append(allThresholds, *resourceMonitor.SuspendAt)
		allThresholds = append(allThresholds, *resourceMonitor.SuspendImmediateAt)
		allThresholds = append(allThresholds, resourceMonitor.NotifyTriggers...)
		var thresholds []int
		for _, trigger := range triggers {
			thresholds = append(thresholds, trigger.Threshold)
		}
		assert.Equal(t, thresholds, allThresholds)

		t.Cleanup(func() {
			err = client.ResourceMonitors.Drop(ctx, id)
			require.NoError(t, err)
		})
	})

	t.Run("test no options", func(t *testing.T) {
		name := random.RandomString(t)
		id := sdk.NewAccountObjectIdentifier(name)

		err := client.ResourceMonitors.Create(ctx, id, nil)

		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
		})

		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor := resourceMonitors[0]
		require.NoError(t, err)
		assert.Equal(t, name, resourceMonitor.Name)
		assert.NotEmpty(t, resourceMonitor.StartTime)
		assert.Empty(t, resourceMonitor.EndTime)
		assert.Empty(t, resourceMonitor.CreditQuota)
		assert.Equal(t, sdk.FrequencyMonthly, resourceMonitor.Frequency)
		assert.Empty(t, resourceMonitor.NotifyUsers)
		assert.Empty(t, resourceMonitor.NotifyTriggers)
		assert.Empty(t, resourceMonitor.SuspendAt)
		assert.Empty(t, resourceMonitor.SuspendImmediateAt)

		t.Cleanup(func() {
			err = client.ResourceMonitors.Drop(ctx, id)
			require.NoError(t, err)
		})
	})
}

func TestInt_ResourceMonitorAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when adding a new trigger", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := createResourceMonitor(t, client)
		t.Cleanup(resourceMonitorCleanup)

		var oldNotifyTriggers []sdk.TriggerDefinition
		for _, threshold := range resourceMonitor.NotifyTriggers {
			oldNotifyTriggers = append(oldNotifyTriggers, sdk.TriggerDefinition{Threshold: threshold, TriggerAction: sdk.TriggerActionNotify})
		}

		var oldTriggers []sdk.TriggerDefinition
		oldTriggers = append(oldTriggers, oldNotifyTriggers...)
		oldTriggers = append(oldTriggers, sdk.TriggerDefinition{Threshold: *resourceMonitor.SuspendAt, TriggerAction: sdk.TriggerActionSuspend})
		oldTriggers = append(oldTriggers, sdk.TriggerDefinition{Threshold: *resourceMonitor.SuspendImmediateAt, TriggerAction: sdk.TriggerActionSuspendImmediate})
		newTriggers := oldTriggers
		newTriggers = append(newTriggers, sdk.TriggerDefinition{Threshold: 30, TriggerAction: sdk.TriggerActionNotify})
		alterOptions := &sdk.AlterResourceMonitorOptions{
			Triggers: newTriggers,
		}
		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), alterOptions)
		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(resourceMonitor.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor = &resourceMonitors[0]
		var newNotifyTriggers []sdk.TriggerDefinition
		for _, threshold := range resourceMonitor.NotifyTriggers {
			newNotifyTriggers = append(newNotifyTriggers, sdk.TriggerDefinition{Threshold: threshold, TriggerAction: sdk.TriggerActionNotify})
		}
		var allTriggers []sdk.TriggerDefinition
		allTriggers = append(allTriggers, newNotifyTriggers...)
		allTriggers = append(allTriggers, sdk.TriggerDefinition{Threshold: *resourceMonitor.SuspendAt, TriggerAction: sdk.TriggerActionSuspend})
		allTriggers = append(allTriggers, sdk.TriggerDefinition{Threshold: *resourceMonitor.SuspendImmediateAt, TriggerAction: sdk.TriggerActionSuspendImmediate})
		assert.ElementsMatch(t, newTriggers, allTriggers)
	})

	t.Run("when setting credit quota", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := createResourceMonitor(t, client)
		t.Cleanup(resourceMonitorCleanup)
		creditQuota := 100
		alterOptions := &sdk.AlterResourceMonitorOptions{
			Set: &sdk.ResourceMonitorSet{
				CreditQuota: &creditQuota,
			},
		}
		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), alterOptions)
		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(resourceMonitor.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor = &resourceMonitors[0]
		assert.Equal(t, creditQuota, int(resourceMonitor.CreditQuota))
	})
	t.Run("when changing scheduling info", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := createResourceMonitor(t, client)
		t.Cleanup(resourceMonitorCleanup)
		frequency, err := sdk.FrequencyFromString("NEVER")
		require.NoError(t, err)
		startTimeStamp := "2025-01-01 12:34"
		endTimeStamp := "2026-01-01 12:34"

		alterOptions := &sdk.AlterResourceMonitorOptions{
			Set: &sdk.ResourceMonitorSet{
				Frequency:      frequency,
				StartTimestamp: &startTimeStamp,
				EndTimestamp:   &endTimeStamp,
			},
		}
		err = client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), alterOptions)
		require.NoError(t, err)
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(resourceMonitor.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(resourceMonitors))
		resourceMonitor = &resourceMonitors[0]
		assert.Equal(t, *frequency, resourceMonitor.Frequency)
		startTime, err := ParseTimestampWithOffset(resourceMonitor.StartTime)
		require.NoError(t, err)
		endTime, err := ParseTimestampWithOffset(resourceMonitor.EndTime)
		require.NoError(t, err)
		assert.Equal(t, startTimeStamp, startTime.Format("2006-01-01 15:04"))
		assert.Equal(t, endTimeStamp, endTime.Format("2006-01-01 15:04"))
	})
}

func TestInt_ResourceMonitorDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when resource monitor exists", func(t *testing.T) {
		resourceMonitor, _ := createResourceMonitor(t, client)
		id := resourceMonitor.ID()
		err := client.ResourceMonitors.Drop(ctx, id)
		require.NoError(t, err)
		_, err = client.ResourceMonitors.ShowByID(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when resource monitor does not exist", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")
		err := client.ResourceMonitors.Drop(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
