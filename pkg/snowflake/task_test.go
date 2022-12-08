package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskCreate(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`"test_db"."test_schema"."test_task"`, st.QualifiedName())

	st.WithWarehouse("test_wh")
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh"`, st.Create())

	st.WithSchedule("USING CRON 0 9-17 * * SUN America/Los_Angeles")
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles'`, st.Create())

	st.WithSessionParameters(map[string]interface{}{"TIMESTAMP_INPUT_FORMAT": "YYYY-MM-DD HH24"})
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24"`, st.Create())

	st.WithComment("test comment")
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment'`, st.Create())

	st.WithTimeout(12)
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' USER_TASK_TIMEOUT_MS = 12`, st.Create())

	st.WithAfter([]string{"other_task"})
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' USER_TASK_TIMEOUT_MS = 12 AFTER "test_db"."test_schema"."other_task"`, st.Create())

	st.WithCondition("SYSTEM$STREAM_HAS_DATA('MYSTREAM')")
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' USER_TASK_TIMEOUT_MS = 12 AFTER "test_db"."test_schema"."other_task" WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM')`, st.Create())

	st.WithStatement("SELECT * FROM table WHERE column = 'name'")
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' USER_TASK_TIMEOUT_MS = 12 AFTER "test_db"."test_schema"."other_task" WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT * FROM table WHERE column = 'name'`, st.Create())

	st.WithAllowOverlappingExecution(true)
	r.Equal(`CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "test_wh" SCHEDULE = 'USING CRON 0 9-17 * * SUN America/Los_Angeles' TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24" COMMENT = 'test comment' ALLOW_OVERLAPPING_EXECUTION = TRUE USER_TASK_TIMEOUT_MS = 12 AFTER "test_db"."test_schema"."other_task" WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT * FROM table WHERE column = 'name'`, st.Create())
}

func TestChangeWarehouse(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SET WAREHOUSE = "much_wh"`, st.ChangeWarehouse("much_wh"))
}

func TestSwitchWarehouseToManaged(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SET WAREHOUSE = null`, st.SwitchWarehouseToManaged())
}

func TestSwitchManagedWithInitialSize(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SET USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = 'SMALL'`, st.SwitchManagedWithInitialSize("SMALL"))
}

func TestChangeSchedule(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SET SCHEDULE = 'USING CRON 0 9-17 * * SUN America/New_York'`, st.ChangeSchedule("USING CRON 0 9-17 * * SUN America/New_York"))
}

func TestRemoveSchedule(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" UNSET SCHEDULE`, st.RemoveSchedule())
}

func TestChangeTimeout(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SET USER_TASK_TIMEOUT_MS = 100`, st.ChangeTimeout(100))
}

func TestRemoveTimeout(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" UNSET USER_TASK_TIMEOUT_MS`, st.RemoveTimeout())
}

func TestChangeComment(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SET COMMENT = 'much comment wow'`, st.ChangeComment("much comment wow"))
}

func TestRemoveComment(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" UNSET COMMENT`, st.RemoveComment())
}

func TestAddAfter(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" ADD AFTER "test_db"."test_schema"."other_task"`, st.AddAfter([]string{"other_task"}))
}

func TestRemoveAfter(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" REMOVE AFTER "test_db"."test_schema"."first_me_task"`, st.RemoveAfter([]string{"first_me_task"}))
}

func TestAddSessionParameters(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	params := map[string]interface{}{"TIMESTAMP_INPUT_FORMAT": "YYYY-MM-DD HH24", "CLIENT_TIMESTAMP_TYPE_MAPPING": "TIMESTAMP_LTZ"}
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SET CLIENT_TIMESTAMP_TYPE_MAPPING = "TIMESTAMP_LTZ", TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24"`, st.AddSessionParameters(params))
}

func TestRemoveSessionParameters(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	params := map[string]interface{}{"TIMESTAMP_INPUT_FORMAT": "YYYY-MM-DD HH24", "CLIENT_TIMESTAMP_TYPE_MAPPING": "TIMESTAMP_LTZ"}
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" UNSET CLIENT_TIMESTAMP_TYPE_MAPPING, TIMESTAMP_INPUT_FORMAT`, st.RemoveSessionParameters(params))
}

func TestChangeCondition(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" MODIFY WHEN TRUE = TRUE`, st.ChangeCondition("TRUE = TRUE"))
}

func TestChangeSqlStatement(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" MODIFY AS SELECT * FROM table`, st.ChangeSQLStatement("SELECT * FROM table"))
}

func TestSuspend(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SUSPEND`, st.Suspend())
}

func TestResume(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" RESUME`, st.Resume())
}

func TestShowParameters(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`SHOW PARAMETERS IN TASK "test_db"."test_schema"."test_task"`, st.ShowParameters())
}

func TestDrop(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`DROP TASK "test_db"."test_schema"."test_task"`, st.Drop())
}

func TestDescribe(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`DESCRIBE TASK "test_db"."test_schema"."test_task"`, st.Describe())
}

func TestShow(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`SHOW TASKS LIKE 'test_task' IN SCHEMA "test_db"."test_schema"`, st.Show())
}

func TestSetAllowOverlappingExecution(t *testing.T) {
	r := require.New(t)
	st := NewTaskBuilder("test_task", "test_db", "test_schema")
	r.Equal(`ALTER TASK "test_db"."test_schema"."test_task" SET ALLOW_OVERLAPPING_EXECUTION = TRUE`, st.SetAllowOverlappingExecutionParameter())
}
