package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestResourceMonitor(t *testing.T) {
	a := assert.New(t)
	rm := snowflake.ResourceMonitor("resource_monitor")
	a.NotNil(rm)

	q := rm.Show()
	a.Equal(`SHOW RESOURCE MONITORS LIKE 'resource_monitor'`, q)

	q = rm.Create().Statement()
	a.Equal(`CREATE RESOURCE MONITOR "resource_monitor"`, q)

	q = rm.Drop()
	a.Equal(`DROP RESOURCE MONITOR "resource_monitor"`, q)

	ab := rm.Alter()
	ab.SetFloat("credit_quota", 66.6)
	q = ab.Statement()
	a.Equal(`ALTER RESOURCE MONITOR "resource_monitor" SET CREDIT_QUOTA=66.6`, q)

	cb := snowflake.ResourceMonitor("resource_monitor").Create()
	cb.NotifyAt(80).NotifyAt(90).SuspendAt(95).SuspendImmediatelyAt(100)
	cb.SetString("frequency", "YEARLY")
	cb.SetFloat("credit_quota", 666.66)
	q = cb.Statement()
	a.Equal(`CREATE RESOURCE MONITOR "resource_monitor" FREQUENCY='YEARLY' CREDIT_QUOTA=666.66 TRIGGERS ON 80 PERCENT DO NOTIFY ON 90 PERCENT DO NOTIFY ON 95 PERCENT DO SUSPEND ON 100 PERCENT DO SUSPEND_IMMEDIATE`, q)
}
