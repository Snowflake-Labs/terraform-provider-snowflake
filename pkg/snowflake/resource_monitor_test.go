package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestResourceMonitor(t *testing.T) {
	r := require.New(t)
	rm := snowflake.ResourceMonitor("resource_monitor")
	r.NotNil(rm)

	q := rm.Show()
	r.Equal(`SHOW RESOURCE MONITORS LIKE 'resource_monitor'`, q)

	q = rm.Create().Statement()
	r.Equal(`CREATE RESOURCE MONITOR "resource_monitor"`, q)

	q = rm.Drop()
	r.Equal(`DROP RESOURCE MONITOR "resource_monitor"`, q)

	ab := rm.Alter()
	ab.SetFloat("credit_quota", 66.6)
	q = ab.Statement()
	r.Equal(`ALTER RESOURCE MONITOR "resource_monitor" SET CREDIT_QUOTA=66.60`, q)

	cb := snowflake.ResourceMonitor("resource_monitor").Create()
	cb.NotifyAt(80).NotifyAt(90).SuspendAt(95).SuspendImmediatelyAt(100)
	cb.SetString("frequency", "YEARLY")

	cb.SetFloat("credit_quota", 666.66666666)
	q = cb.Statement()
	r.Equal(`CREATE RESOURCE MONITOR "resource_monitor" FREQUENCY='YEARLY' CREDIT_QUOTA=666.67 TRIGGERS ON 80 PERCENT DO NOTIFY ON 90 PERCENT DO NOTIFY ON 95 PERCENT DO SUSPEND ON 100 PERCENT DO SUSPEND_IMMEDIATE`, q)

	// Check if credit quota can be parsed correctly to float if given an integer
	cb.SetFloat("credit_quota", 666)
	q = cb.Statement()
	r.Equal(`CREATE RESOURCE MONITOR "resource_monitor" FREQUENCY='YEARLY' CREDIT_QUOTA=666.00 TRIGGERS ON 80 PERCENT DO NOTIFY ON 90 PERCENT DO NOTIFY ON 95 PERCENT DO SUSPEND ON 100 PERCENT DO SUSPEND_IMMEDIATE`, q)
}
