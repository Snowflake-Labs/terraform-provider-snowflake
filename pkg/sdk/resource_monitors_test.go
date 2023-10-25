package sdk

import (
	"testing"
	"time"
)

func TestResourceMonitorCreate(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &CreateResourceMonitorOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("with complete options", func(t *testing.T) {
		creditQuota := Int(100)
		frequency := FrequencyMonthly
		startTimeStamp := "IMMIEDIATELY"
		endTimeStamp := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).String()
		notifiedUsers := []NotifiedUser{{Name: "FIRST_USER"}, {Name: "SECOND_USER"}}
		triggers := []TriggerDefinition{
			{
				Threshold:     50,
				TriggerAction: TriggerActionSuspendImmediate,
			},
			{
				Threshold:     100,
				TriggerAction: TriggerActionNotify,
			},
		}

		opts := &CreateResourceMonitorOptions{
			OrReplace: Bool(true),
			With: &ResourceMonitorWith{
				CreditQuota:    creditQuota,
				Frequency:      &frequency,
				StartTimestamp: &startTimeStamp,
				EndTimestamp:   &endTimeStamp,
				NotifyUsers:    &NotifyUsers{notifiedUsers},
				Triggers:       triggers,
			},
			name: id,
		}

		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE RESOURCE MONITOR %s WITH CREDIT_QUOTA = 100 FREQUENCY = MONTHLY START_TIMESTAMP = 'IMMIEDIATELY' END_TIMESTAMP = '%s' NOTIFY_USERS = ("FIRST_USER", "SECOND_USER") TRIGGERS ON 50 PERCENT DO SUSPEND_IMMEDIATE ON 100 PERCENT DO NOTIFY`,
			id.FullyQualifiedName(),
			endTimeStamp,
		)
	})
}

func TestResourceMonitorAlter(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{
			name: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterResourceMonitorOptions", "Set", "NotifyUsers", "Triggers"))
	})

	t.Run("with a single set", func(t *testing.T) {
		newCreditQuota := Int(50)
		opts := &AlterResourceMonitorOptions{
			name: id,
			Set: &ResourceMonitorSet{
				CreditQuota: newCreditQuota,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = %d", id.FullyQualifiedName(), *newCreditQuota)
	})

	t.Run("with a multitple set", func(t *testing.T) {
		newCreditQuota := Int(50)
		newFrequency := FrequencyYearly
		newStartTimeStamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).String()
		opts := &AlterResourceMonitorOptions{
			name: id,
			Set: &ResourceMonitorSet{
				CreditQuota:    newCreditQuota,
				Frequency:      &newFrequency,
				StartTimestamp: &newStartTimeStamp,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = %d FREQUENCY = %s START_TIMESTAMP = '%s'", id.FullyQualifiedName(), *newCreditQuota, newFrequency, newStartTimeStamp)
	})
}

func TestResourceMonitorDrop(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &dropResourceMonitorOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &dropResourceMonitorOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP RESOURCE MONITOR %s", id.FullyQualifiedName())
	})
}

func TestResourceMonitorShow(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowResourceMonitorOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW RESOURCE MONITORS")
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW RESOURCE MONITORS LIKE '%s'", id.Name())
	})
}
