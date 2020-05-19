package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringFromTaskID(t *testing.T) {
	r := require.New(t)
	task := taskID{DatabaseName: "test_db", SchemaName: "test_schema", TaskName: "test_task"}
	id, err := task.String()
	r.NoError(err)
	r.Equal(id, "test_db|test_schema|test_task")
}

func TestTaskIDFromString(t *testing.T) {
	r := require.New(t)

	id := "test_db|test_schema|test_task"
	task, err := taskIDFromString(id)
	r.NoError(err)
	r.Equal("test_db", task.DatabaseName)
	r.Equal("test_schema", task.SchemaName)
	r.Equal("test_task", task.TaskName)

	id = "test_db"
	_, err = taskIDFromString(id)
	r.Equal(fmt.Errorf("3 fields allowed"), err)

	// Bad ID
	id = "|"
	_, err = taskIDFromString(id)
	r.Equal(fmt.Errorf("3 fields allowed"), err)

	// 0 lines
	id = ""
	_, err = taskIDFromString(id)
	r.Equal(fmt.Errorf("1 line per task"), err)

	// 2 lines
	id = `database|schema|task
		database|schema|task`
	_, err = taskIDFromString(id)
	r.Equal(fmt.Errorf("1 line per task"), err)

}
