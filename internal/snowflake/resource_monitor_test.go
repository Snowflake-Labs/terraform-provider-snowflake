// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/stretchr/testify/require"
)

func TestResourceMonitor(t *testing.T) {
	r := require.New(t)
	rm := snowflake.NewResourceMonitorBuilder("resource_monitor")
	r.NotNil(rm)

	q := rm.Show()
	r.Equal(`SHOW RESOURCE MONITORS LIKE 'resource_monitor'`, q)

	q = rm.Create().Statement()
	r.Equal(`CREATE RESOURCE MONITOR "resource_monitor"`, q)

	q = rm.Drop()
	r.Equal(`DROP RESOURCE MONITOR "resource_monitor"`, q)

	ab := rm.Alter()
	ab.SetInt("credit_quota", 66)
	q = ab.Statement()
	r.Equal(`ALTER RESOURCE MONITOR "resource_monitor" SET CREDIT_QUOTA=66`, q)

	ab = rm.Alter()
	ab.SetStringList("notify_users", []string{"USERONE", "USERTWO"})
	q = ab.Statement()
	r.Equal(`ALTER RESOURCE MONITOR "resource_monitor" SET NOTIFY_USERS=('USERONE', 'USERTWO')`, q)

	cb := snowflake.NewResourceMonitorBuilder("resource_monitor").Create()
	cb.NotifyAt(80).NotifyAt(90).SuspendAt(95).SuspendImmediatelyAt(100)
	cb.SetString("frequency", "YEARLY")

	cb.SetInt("credit_quota", 666)
	cb.SetStringList("notify_users", []string{"USERONE", "USERTWO"})
	q = cb.Statement()
	r.Equal(`CREATE RESOURCE MONITOR "resource_monitor" FREQUENCY='YEARLY' CREDIT_QUOTA=666 NOTIFY_USERS=('USERONE', 'USERTWO') TRIGGERS ON 80 PERCENT DO NOTIFY ON 90 PERCENT DO NOTIFY ON 95 PERCENT DO SUSPEND ON 100 PERCENT DO SUSPEND_IMMEDIATE`, q)
}

func TestResourceMonitorSetOnAccount(t *testing.T) {
	r := require.New(t)
	s := snowflake.NewResourceMonitorBuilder("test_resource_monitor")
	r.NotNil(s)

	q := s.Create().SetOnAccount()
	r.Equal(`ALTER ACCOUNT SET RESOURCE_MONITOR = "test_resource_monitor"`, q)
}

func TestResourceMonitorSetOnWarehouse(t *testing.T) {
	r := require.New(t)
	s := snowflake.NewResourceMonitorBuilder("test_resource_monitor")
	r.NotNil(s)

	q := s.Create().SetOnWarehouse("test_warehouse")
	r.Equal(`ALTER WAREHOUSE "test_warehouse" SET RESOURCE_MONITOR = "test_resource_monitor"`, q)
}
