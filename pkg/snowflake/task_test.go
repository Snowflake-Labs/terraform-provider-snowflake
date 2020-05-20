package snowflake

import (
	"database/sql"
	"testing"
)

func Test_buildFullyQualifiedTaskName(t *testing.T) {
	type args struct {
		name     string
		schema   string
		database string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty-name",
			args: args{
				name:     "",
				schema:   "schema",
				database: "database",
			},
			want: "",
		},
		{
			name: "empty-schema",
			args: args{
				name:     "name",
				schema:   "",
				database: "database",
			},
			want: "",
		},
		{
			name: "empty-database",
			args: args{
				name:     "name",
				schema:   "schema",
				database: "",
			},
			want: "",
		},
		{
			name: "all-lower",
			args: args{
				name:     "name",
				schema:   "schema",
				database: "database",
			},
			want: "\"database\".\"schema\".\"name\"",
		},
		{
			name: "all-cap-database",
			args: args{
				name:     "name",
				schema:   "schema",
				database: "DATABASE",
			},
			want: "DATABASE.\"schema\".\"name\"",
		},
		{
			name: "all-cap-schema",
			args: args{
				name:     "name",
				schema:   "SCHEMA",
				database: "database",
			},
			want: "\"database\".SCHEMA.\"name\"",
		},
		{
			name: "all-cap-name",
			args: args{
				name:     "NAME",
				schema:   "schema",
				database: "database",
			},
			want: "\"database\".\"schema\".NAME",
		},
		{
			name: "all-cap-name-schema",
			args: args{
				name:     "NAME",
				schema:   "SCHEMA",
				database: "database",
			},
			want: "\"database\".SCHEMA.NAME",
		},
		{
			name: "all-cap-schema-database",
			args: args{
				name:     "name",
				schema:   "SCHEMA",
				database: "DATABASE",
			},
			want: "DATABASE.SCHEMA.\"name\"",
		},
		{
			name: "mixed-case-database",
			args: args{
				name:     "name",
				schema:   "SCHEMA",
				database: "Database",
			},
			want: "\"Database\".SCHEMA.\"name\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildFullyQualifiedTaskName(tt.args.name, tt.args.schema, tt.args.database); got != tt.want {
				t.Errorf("buildFullyQualifiedTaskName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRow_QualifiedPredecessorName(t *testing.T) {
	type fields struct {
		CreatedOn    string
		TaskName     string
		TaskID       string
		DatabaseName string
		SchemaName   string
		Owner        string
		Comment      sql.NullString
		Warehouse    string
		Schedule     sql.NullString
		Predecessor  sql.NullString
		State        string
		Definition   string
		Condition    sql.NullString
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple-lower",
			fields: fields{
				TaskName:     "task2",
				SchemaName:   "schema",
				DatabaseName: "db",
				Predecessor: sql.NullString{
					Valid:  true,
					String: "task1",
				},
			},
			want: "\"db\".\"schema\".\"task1\"",
		},
		{
			name: "simple-name-upper-schema-mixed",
			fields: fields{
				TaskName:     "Task2",
				SchemaName:   "Schema",
				DatabaseName: "db",
				Predecessor: sql.NullString{
					Valid:  true,
					String: "TASK1",
				},
			},
			want: "\"db\".\"Schema\".TASK1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TaskRow{
				CreatedOn:    tt.fields.CreatedOn,
				TaskName:     tt.fields.TaskName,
				TaskID:       tt.fields.TaskID,
				DatabaseName: tt.fields.DatabaseName,
				SchemaName:   tt.fields.SchemaName,
				Owner:        tt.fields.Owner,
				Comment:      tt.fields.Comment,
				Warehouse:    tt.fields.Warehouse,
				Schedule:     tt.fields.Schedule,
				Predecessor:  tt.fields.Predecessor,
				State:        tt.fields.State,
				Definition:   tt.fields.Definition,
				Condition:    tt.fields.Condition,
			}
			if got := tr.QualifiedPredecessorName(); got != tt.want {
				t.Errorf("TaskRow.QualifiedPredecessorName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRow_QualifiedName(t *testing.T) {
	type fields struct {
		CreatedOn    string
		TaskName     string
		TaskID       string
		DatabaseName string
		SchemaName   string
		Owner        string
		Comment      sql.NullString
		Warehouse    string
		Schedule     sql.NullString
		Predecessor  sql.NullString
		State        string
		Definition   string
		Condition    sql.NullString
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple-lower",
			fields: fields{
				TaskName:     "task2",
				SchemaName:   "schema",
				DatabaseName: "db",
			},
			want: "\"db\".\"schema\".\"task2\"",
		},
		{
			name: "simple-name-schema-mixed",
			fields: fields{
				TaskName:     "Task2",
				SchemaName:   "Schema",
				DatabaseName: "db",
			},
			want: "\"db\".\"Schema\".\"Task2\"",
		},
		{
			name: "simple-name-schema-mixed-db-upper",
			fields: fields{
				TaskName:     "Task2",
				SchemaName:   "Schema",
				DatabaseName: "DB",
			},
			want: "DB.\"Schema\".\"Task2\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TaskRow{
				CreatedOn:    tt.fields.CreatedOn,
				TaskName:     tt.fields.TaskName,
				TaskID:       tt.fields.TaskID,
				DatabaseName: tt.fields.DatabaseName,
				SchemaName:   tt.fields.SchemaName,
				Owner:        tt.fields.Owner,
				Comment:      tt.fields.Comment,
				Warehouse:    tt.fields.Warehouse,
				Schedule:     tt.fields.Schedule,
				Predecessor:  tt.fields.Predecessor,
				State:        tt.fields.State,
				Definition:   tt.fields.Definition,
				Condition:    tt.fields.Condition,
			}
			if got := tr.QualifiedName(); got != tt.want {
				t.Errorf("TaskRow.QualifiedName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRow_IsEnabled(t *testing.T) {
	type fields struct {
		CreatedOn    string
		TaskName     string
		TaskID       string
		DatabaseName string
		SchemaName   string
		Owner        string
		Comment      sql.NullString
		Warehouse    string
		Schedule     sql.NullString
		Predecessor  sql.NullString
		State        string
		Definition   string
		Condition    sql.NullString
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "suspended",
			fields: fields{
				State: "Suspended",
			},
			want: false,
		},
		{
			name: "started",
			fields: fields{
				State: "Started",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TaskRow{
				CreatedOn:    tt.fields.CreatedOn,
				TaskName:     tt.fields.TaskName,
				TaskID:       tt.fields.TaskID,
				DatabaseName: tt.fields.DatabaseName,
				SchemaName:   tt.fields.SchemaName,
				Owner:        tt.fields.Owner,
				Comment:      tt.fields.Comment,
				Warehouse:    tt.fields.Warehouse,
				Schedule:     tt.fields.Schedule,
				Predecessor:  tt.fields.Predecessor,
				State:        tt.fields.State,
				Definition:   tt.fields.Definition,
				Condition:    tt.fields.Condition,
			}
			if got := tr.IsEnabled(); got != tt.want {
				t.Errorf("TaskRow.IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskBuilder_Create(t *testing.T) {
	type fields struct {
		name           string
		schema         string
		database       string
		warehouse      string
		schedule       string
		scheduleSet    bool
		timeout        int
		timeoutSet     bool
		comment        string
		commentSet     bool
		predecessor    string
		predecessorSet bool
		conditional    string
		conditionalSet bool
		definition     string
		enabled        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "minimum-fields",
			fields: fields{
				name:       "task",
				schema:     "sch",
				database:   "db",
				warehouse:  "wh",
				definition: "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' AS select * from table",
		},
		{
			name: "include-schedule",
			fields: fields{
				name:       "task",
				schema:     "sch",
				database:   "db",
				warehouse:  "wh",
				schedule:   "5 MINUTE",
				definition: "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' SCHEDULE = '5 MINUTE' AS select * from table",
		},
		{
			name: "include-comment",
			fields: fields{
				name:       "task",
				schema:     "sch",
				database:   "db",
				warehouse:  "wh",
				comment:    "simple task (test)",
				definition: "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' COMMENT = 'simple task (test)' AS select * from table",
		},
		{
			name: "include-predecessor",
			fields: fields{
				name:        "task",
				schema:      "sch",
				database:    "db",
				warehouse:   "wh",
				predecessor: "ROOT_TASK",
				definition:  "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' AFTER \"db\".\"sch\".ROOT_TASK AS select * from table",
		},
		{
			name: "include-task-timeout",
			fields: fields{
				name:       "task",
				schema:     "sch",
				database:   "db",
				warehouse:  "wh",
				timeout:    600,
				timeoutSet: true,
				definition: "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' USER_TASK_TIMEOUT_MS = 600 AS select * from table",
		},
		{
			name: "include-conditional",
			fields: fields{
				name:           "task",
				schema:         "sch",
				database:       "db",
				warehouse:      "wh",
				conditional:    "SYSTEM$STREAM_HAS_DATA('MYSTREAM')",
				conditionalSet: true,
				definition:     "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS select * from table",
		},
		{
			name: "include-conditional-predecessor",
			fields: fields{
				name:           "task",
				schema:         "sch",
				database:       "db",
				warehouse:      "wh",
				predecessor:    "root_task",
				conditional:    "SYSTEM$STREAM_HAS_DATA('MYSTREAM')",
				conditionalSet: true,
				definition:     "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' AFTER \"db\".\"sch\".\"root_task\" WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS select * from table",
		},
		{
			name: "include-conditional-schedule",
			fields: fields{
				name:           "task",
				schema:         "sch",
				database:       "db",
				warehouse:      "wh",
				schedule:       "5 MINUTE",
				conditional:    "SYSTEM$STREAM_HAS_DATA('MYSTREAM')",
				conditionalSet: true,
				definition:     "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' SCHEDULE = '5 MINUTE' WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS select * from table",
		},
		{
			name: "include-conditional-schedule-timeout",
			fields: fields{
				name:           "task",
				schema:         "sch",
				database:       "db",
				warehouse:      "wh",
				timeout:        600,
				timeoutSet:     true,
				schedule:       "5 MINUTE",
				conditional:    "SYSTEM$STREAM_HAS_DATA('MYSTREAM')",
				conditionalSet: true,
				definition:     "select * from table",
			},
			want: "CREATE TASK \"db\".\"sch\".\"task\" WAREHOUSE = 'wh' SCHEDULE = '5 MINUTE' USER_TASK_TIMEOUT_MS = 600 WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS select * from table",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TaskBuilder{
				name:           tt.fields.name,
				schema:         tt.fields.schema,
				database:       tt.fields.database,
				warehouse:      tt.fields.warehouse,
				schedule:       tt.fields.schedule,
				scheduleSet:    tt.fields.scheduleSet,
				timeout:        tt.fields.timeout,
				timeoutSet:     tt.fields.timeoutSet,
				comment:        tt.fields.comment,
				commentSet:     tt.fields.commentSet,
				predecessor:    tt.fields.predecessor,
				predecessorSet: tt.fields.predecessorSet,
				conditional:    tt.fields.conditional,
				conditionalSet: tt.fields.conditionalSet,
				definition:     tt.fields.definition,
				enabled:        tt.fields.enabled,
			}
			if got := tb.Create(); got != tt.want {
				t.Errorf("TaskBuilder.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskBuilder_ChangeState(t *testing.T) {
	type fields struct {
		name           string
		schema         string
		database       string
		warehouse      string
		schedule       string
		scheduleSet    bool
		timeout        int
		timeoutSet     bool
		comment        string
		commentSet     bool
		predecessor    string
		predecessorSet bool
		conditional    string
		conditionalSet bool
		definition     string
		enabled        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "suspend",
			fields: fields{
				name:     "task1",
				schema:   "sch",
				database: "db",
				enabled:  false,
			},
			want: "ALTER TASK \"db\".\"sch\".\"task1\" SUSPEND",
		},
		{
			name: "resume",
			fields: fields{
				name:     "task1",
				schema:   "sch",
				database: "db",
				enabled:  true,
			},
			want: "ALTER TASK \"db\".\"sch\".\"task1\" RESUME",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TaskBuilder{
				name:           tt.fields.name,
				schema:         tt.fields.schema,
				database:       tt.fields.database,
				warehouse:      tt.fields.warehouse,
				schedule:       tt.fields.schedule,
				scheduleSet:    tt.fields.scheduleSet,
				timeout:        tt.fields.timeout,
				timeoutSet:     tt.fields.timeoutSet,
				comment:        tt.fields.comment,
				commentSet:     tt.fields.commentSet,
				predecessor:    tt.fields.predecessor,
				predecessorSet: tt.fields.predecessorSet,
				conditional:    tt.fields.conditional,
				conditionalSet: tt.fields.conditionalSet,
				definition:     tt.fields.definition,
				enabled:        tt.fields.enabled,
			}
			if got := tb.ChangeState(); got != tt.want {
				t.Errorf("TaskBuilder.ChangeState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskBuilder_ChangeWarehouseAndSchedule(t *testing.T) {
	type fields struct {
		name           string
		schema         string
		database       string
		warehouse      string
		schedule       string
		scheduleSet    bool
		timeout        int
		timeoutSet     bool
		comment        string
		commentSet     bool
		predecessor    string
		predecessorSet bool
		conditional    string
		conditionalSet bool
		definition     string
		enabled        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "change-warehouse",
			fields: fields{
				name:      "TASK",
				schema:    "sch",
				database:  "db",
				warehouse: "wh2",
			},
			want: "ALTER TASK \"db\".\"sch\".TASK SET WAREHOUSE = 'wh2' ",
		},
		{
			name: "change-schedule",
			fields: fields{
				name:        "TASK",
				schema:      "sch",
				database:    "db",
				schedule:    "10 MINUTE",
				scheduleSet: true,
			},
			want: "ALTER TASK \"db\".\"sch\".TASK SET SCHEDULE = '10 MINUTE' ",
		},
		{
			name: "remove-schedule",
			fields: fields{
				name:        "TASK",
				schema:      "sch",
				database:    "db",
				scheduleSet: true,
			},
			want: "ALTER TASK \"db\".\"sch\".TASK SET SCHEDULE = NULL ",
		},
		{
			name: "change-schedule-warehouse",
			fields: fields{
				name:        "TASK",
				schema:      "sch",
				database:    "db",
				warehouse:   "wh3",
				schedule:    "10 MINUTE",
				scheduleSet: true,
			},
			want: "ALTER TASK \"db\".\"sch\".TASK SET WAREHOUSE = 'wh3' SCHEDULE = '10 MINUTE' ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TaskBuilder{
				name:           tt.fields.name,
				schema:         tt.fields.schema,
				database:       tt.fields.database,
				warehouse:      tt.fields.warehouse,
				schedule:       tt.fields.schedule,
				scheduleSet:    tt.fields.scheduleSet,
				timeout:        tt.fields.timeout,
				timeoutSet:     tt.fields.timeoutSet,
				comment:        tt.fields.comment,
				commentSet:     tt.fields.commentSet,
				predecessor:    tt.fields.predecessor,
				predecessorSet: tt.fields.predecessorSet,
				conditional:    tt.fields.conditional,
				conditionalSet: tt.fields.conditionalSet,
				definition:     tt.fields.definition,
				enabled:        tt.fields.enabled,
			}
			if got := tb.ChangeWarehouseAndSchedule(); got != tt.want {
				t.Errorf("TaskBuilder.ChangeWarehouseAndSchedule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskBuilder_UpdateConditional(t *testing.T) {
	type fields struct {
		name           string
		schema         string
		database       string
		warehouse      string
		schedule       string
		scheduleSet    bool
		timeout        int
		timeoutSet     bool
		comment        string
		commentSet     bool
		predecessor    string
		predecessorSet bool
		conditional    string
		conditionalSet bool
		definition     string
		enabled        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "remove",
			fields: fields{
				name:           "task1",
				schema:         "sch",
				database:       "db",
				conditionalSet: true,
			},
			want: "ALTER TASK \"db\".\"sch\".\"task1\" MODIFY WHEN NULL",
		},
		{
			name: "update",
			fields: fields{
				name:           "task1",
				schema:         "sch",
				database:       "db",
				conditionalSet: true,
				conditional:    "SYSTEM$STREAM_HAS_DATA('MYSTREAM')",
			},
			want: "ALTER TASK \"db\".\"sch\".\"task1\" MODIFY WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM')",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TaskBuilder{
				name:           tt.fields.name,
				schema:         tt.fields.schema,
				database:       tt.fields.database,
				warehouse:      tt.fields.warehouse,
				schedule:       tt.fields.schedule,
				scheduleSet:    tt.fields.scheduleSet,
				timeout:        tt.fields.timeout,
				timeoutSet:     tt.fields.timeoutSet,
				comment:        tt.fields.comment,
				commentSet:     tt.fields.commentSet,
				predecessor:    tt.fields.predecessor,
				predecessorSet: tt.fields.predecessorSet,
				conditional:    tt.fields.conditional,
				conditionalSet: tt.fields.conditionalSet,
				definition:     tt.fields.definition,
				enabled:        tt.fields.enabled,
			}
			if got := tb.UpdateConditional(); got != tt.want {
				t.Errorf("TaskBuilder.UpdateConditional() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskBuilder_UpdateSQL(t *testing.T) {
	type fields struct {
		name           string
		schema         string
		database       string
		warehouse      string
		schedule       string
		scheduleSet    bool
		timeout        int
		timeoutSet     bool
		comment        string
		commentSet     bool
		predecessor    string
		predecessorSet bool
		conditional    string
		conditionalSet bool
		definition     string
		enabled        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "update",
			fields: fields{
				name:       "task1",
				schema:     "sch",
				database:   "db",
				definition: "select * from table",
			},
			want: "ALTER TASK \"db\".\"sch\".\"task1\" MODIFY AS select * from table",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TaskBuilder{
				name:           tt.fields.name,
				schema:         tt.fields.schema,
				database:       tt.fields.database,
				warehouse:      tt.fields.warehouse,
				schedule:       tt.fields.schedule,
				scheduleSet:    tt.fields.scheduleSet,
				timeout:        tt.fields.timeout,
				timeoutSet:     tt.fields.timeoutSet,
				comment:        tt.fields.comment,
				commentSet:     tt.fields.commentSet,
				predecessor:    tt.fields.predecessor,
				predecessorSet: tt.fields.predecessorSet,
				conditional:    tt.fields.conditional,
				conditionalSet: tt.fields.conditionalSet,
				definition:     tt.fields.definition,
				enabled:        tt.fields.enabled,
			}
			if got := tb.UpdateSQL(); got != tt.want {
				t.Errorf("TaskBuilder.UpdateSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskBuilder_RemovePredecessor(t *testing.T) {
	type fields struct {
		name           string
		schema         string
		database       string
		warehouse      string
		schedule       string
		scheduleSet    bool
		timeout        int
		timeoutSet     bool
		comment        string
		commentSet     bool
		predecessor    string
		predecessorSet bool
		conditional    string
		conditionalSet bool
		definition     string
		enabled        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "remove",
			fields: fields{
				name:        "task1",
				schema:      "sch",
				database:    "db",
				predecessor: "task2",
			},
			want: "ALTER TASK \"db\".\"sch\".\"task1\" REMOVE AFTER \"db\".\"sch\".\"task2\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TaskBuilder{
				name:           tt.fields.name,
				schema:         tt.fields.schema,
				database:       tt.fields.database,
				warehouse:      tt.fields.warehouse,
				schedule:       tt.fields.schedule,
				scheduleSet:    tt.fields.scheduleSet,
				timeout:        tt.fields.timeout,
				timeoutSet:     tt.fields.timeoutSet,
				comment:        tt.fields.comment,
				commentSet:     tt.fields.commentSet,
				predecessor:    tt.fields.predecessor,
				predecessorSet: tt.fields.predecessorSet,
				conditional:    tt.fields.conditional,
				conditionalSet: tt.fields.conditionalSet,
				definition:     tt.fields.definition,
				enabled:        tt.fields.enabled,
			}
			if got := tb.RemovePredecessor(); got != tt.want {
				t.Errorf("TaskBuilder.RemovePredecessor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskBuilder_Drop(t *testing.T) {
	type fields struct {
		name           string
		schema         string
		database       string
		warehouse      string
		schedule       string
		scheduleSet    bool
		timeout        int
		timeoutSet     bool
		comment        string
		commentSet     bool
		predecessor    string
		predecessorSet bool
		conditional    string
		conditionalSet bool
		definition     string
		enabled        bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "delete",
			fields: fields{
				name:     "task1",
				schema:   "sch",
				database: "db",
			},
			want: "DROP TASK \"db\".\"sch\".\"task1\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TaskBuilder{
				name:           tt.fields.name,
				schema:         tt.fields.schema,
				database:       tt.fields.database,
				warehouse:      tt.fields.warehouse,
				schedule:       tt.fields.schedule,
				scheduleSet:    tt.fields.scheduleSet,
				timeout:        tt.fields.timeout,
				timeoutSet:     tt.fields.timeoutSet,
				comment:        tt.fields.comment,
				commentSet:     tt.fields.commentSet,
				predecessor:    tt.fields.predecessor,
				predecessorSet: tt.fields.predecessorSet,
				conditional:    tt.fields.conditional,
				conditionalSet: tt.fields.conditionalSet,
				definition:     tt.fields.definition,
				enabled:        tt.fields.enabled,
			}
			if got := tb.Drop(); got != tt.want {
				t.Errorf("TaskBuilder.Drop() = %v, want %v", got, tt.want)
			}
		})
	}
}
