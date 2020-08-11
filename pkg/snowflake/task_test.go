package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskCreate(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.QualifiedName(), `"test_db"."test_schema"."test_task"`)

	st.WithWarehouse("test_wh")
	r.Equal(st.Create(), `CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh"`)

	st.WithSchedule("USING CRON 0 9-17 * * SUN America/Los_Angeles")
	r.Equal(st.Create(), `CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles'`)

	st.WithSessionParameters(map[string]interface{}{"TIMESTAMP_INPUT_FORMAT": "YYYY-MM-DD HH24"})
	r.Equal(st.Create(), `CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24"`)

	st.WithComment("test comment")
	r.Equal(st.Create(), `CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment'`)

	st.WithTimeout(12)
	r.Equal(st.Create(), `CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' USER_TASK_TIMEOUT_MS = 12`)

	st.WithDependency("other_task")
	r.Equal(st.Create(), `CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' USER_TASK_TIMEOUT_MS = 12 AFTER "test_db"."test_schema"."other_task"`)

	st.WithCondition("SYSTEM$STREAM_HAS_DATA('MYSTREAM')")
	r.Equal(st.Create(), `CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' USER_TASK_TIMEOUT_MS = 12 AFTER "test_db"."test_schema"."other_task" WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM')`)

	st.WithStatement("SELECT * FROM table WHERE column = 'name'")
	r.Equal(st.Create(), `CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' USER_TASK_TIMEOUT_MS = 12 AFTER "test_db"."test_schema"."other_task" WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT * FROM table WHERE column = 'name'`)
}

func TestChangeWarehouse(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.ChangeWarehouse("much_wh"), `ALTER TASK "test_db"."test_schema"."test_task" SET WAREHOUSE = "much_wh"`)
}

func TestChangeSchedule(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.ChangeSchedule("USING CRON 0 9-17 * * SUN America/New_York"), `ALTER TASK "test_db"."test_schema"."test_task" SET SCHEDULE = 'USING CRON 0 9-17 * * SUN America/New_York'`)
}

func TestRemoveSchedule(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.RemoveSchedule(), `ALTER TASK "test_db"."test_schema"."test_task" UNSET SCHEDULE`)
}

func TestChangeTimeout(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.ChangeTimeout(100), `ALTER TASK "test_db"."test_schema"."test_task" SET USER_TASK_TIMEOUT_MS = 100`)
}

func TestRemoveTimeout(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.RemoveTimeout(), `ALTER TASK "test_db"."test_schema"."test_task" UNSET USER_TASK_TIMEOUT_MS`)
}

func TestChangeComment(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.ChangeComment("much comment wow"), `ALTER TASK "test_db"."test_schema"."test_task" SET COMMENT = 'much comment wow'`)
}

func TestRemoveComment(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.RemoveComment(), `ALTER TASK "test_db"."test_schema"."test_task" UNSET COMMENT`)
}

func TestAddDependency(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.AddDependency("other_task"), `ALTER TASK "test_db"."test_schema"."test_task" ADD AFTER "test_db"."test_schema"."other_task"`)
}

func TestRemoveDependency(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.RemoveDependency("first_me_task"), `ALTER TASK "test_db"."test_schema"."test_task" REMOVE AFTER "test_db"."test_schema"."first_me_task"`)
}

func TestAddSessionParameters(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	params := map[string]interface{}{"TIMESTAMP_INPUT_FORMAT": "YYYY-MM-DD HH24", "CLIENT_TIMESTAMP_TYPE_MAPPING": "TIMESTAMP_LTZ"}
	r.Equal(st.AddSessionParameters(params), `ALTER TASK "test_db"."test_schema"."test_task" SET CLIENT_TIMESTAMP_TYPE_MAPPING = "TIMESTAMP_LTZ", TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24"`)
}

func TestRemoveSessionParameters(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	params := map[string]interface{}{"TIMESTAMP_INPUT_FORMAT": "YYYY-MM-DD HH24", "CLIENT_TIMESTAMP_TYPE_MAPPING": "TIMESTAMP_LTZ"}
	r.Equal(st.RemoveSessionParameters(params), `ALTER TASK "test_db"."test_schema"."test_task" UNSET CLIENT_TIMESTAMP_TYPE_MAPPING, TIMESTAMP_INPUT_FORMAT`)
}

func TestChangeCondition(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.ChangeCondition("TRUE = TRUE"), `ALTER TASK "test_db"."test_schema"."test_task" MODIFY WHEN TRUE = TRUE`)
}

func TestChangeSqlStatement(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.ChangeSqlStatement("SELECT * FROM table"), `ALTER TASK "test_db"."test_schema"."test_task" MODIFY AS SELECT * FROM table`)
}

func TestSuspend(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.Suspend(), `ALTER TASK "test_db"."test_schema"."test_task" SUSPEND`)
}

func TestResume(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.Resume(), `ALTER TASK "test_db"."test_schema"."test_task" RESUME`)
}

func TestShowParameters(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.ShowParameters(), `SHOW PARAMETERS IN TASK "test_db"."test_schema"."test_task"`)
}

func TestDrop(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.Drop(), `DROP TASK "test_db"."test_schema"."test_task"`)
}

func TestDescribe(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.Describe(), `DESCRIBE TASK "test_db"."test_schema"."test_task"`)
}

func TestShow(t *testing.T) {
	r := require.New(t)
	st := Task("test_task", "test_db", "test_schema")
	r.Equal(st.Show(), `SHOW TASKS LIKE 'test_task' IN DATABASE "test_db"`)
}
