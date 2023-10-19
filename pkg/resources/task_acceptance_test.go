package resources_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type (
	AccTaskTestSettings struct {
		WarehouseName string
		RootTask      *TaskSettings
		ChildTask     *TaskSettings
		SoloTask      *TaskSettings
	}

	TaskSettings struct {
		Name              string
		Enabled           bool
		Schema            string
		SQL               string
		Schedule          string
		Comment           string
		When              string
		SessionParams     map[string]string
		UserTaskTimeoutMs int64
	}
)

var (
	rootname      = "root_task"
	childname     = "child_task"
	soloname      = "standalone_task"
	warehousename = strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	initialState = &AccTaskTestSettings{ //nolint
		WarehouseName: warehousename,

		RootTask: &TaskSettings{
			Name:              rootname,
			Schema:            acc.TestSchemaName,
			SQL:               "SHOW FUNCTIONS",
			Enabled:           true,
			Schedule:          "5 MINUTE",
			UserTaskTimeoutMs: 1800000,
		},

		ChildTask: &TaskSettings{
			Name:    childname,
			SQL:     "SELECT 1",
			Enabled: false,
			Comment: "initial state",
		},

		SoloTask: &TaskSettings{
			Name:    soloname,
			Schema:  acc.TestSchemaName,
			SQL:     "SELECT 1",
			When:    "TRUE",
			Enabled: false,
		},
	}

	// Enables the Child and changes the SQL.
	stepOne = &AccTaskTestSettings{ //nolint
		WarehouseName: warehousename,

		RootTask: &TaskSettings{
			Name:              rootname,
			Schema:            acc.TestSchemaName,
			SQL:               "SHOW FUNCTIONS",
			Enabled:           true,
			Schedule:          "5 MINUTE",
			UserTaskTimeoutMs: 1800000,
		},

		ChildTask: &TaskSettings{
			Name:    childname,
			SQL:     "SELECT *",
			Enabled: true,
			Comment: "secondary state",
		},

		SoloTask: &TaskSettings{
			Name:    soloname,
			Schema:  acc.TestSchemaName,
			SQL:     "SELECT *",
			When:    "TRUE",
			Enabled: true,
			SessionParams: map[string]string{
				"TIMESTAMP_INPUT_FORMAT": "YYYY-MM-DD HH24",
			},
			Schedule:          "5 MINUTE",
			UserTaskTimeoutMs: 1800000,
		},
	}

	// Changes Root Schedule and SQL.
	stepTwo = &AccTaskTestSettings{ //nolint
		WarehouseName: warehousename,

		RootTask: &TaskSettings{
			Name:              rootname,
			Schema:            acc.TestSchemaName,
			SQL:               "SHOW TABLES",
			Enabled:           true,
			Schedule:          "15 MINUTE",
			UserTaskTimeoutMs: 1800000,
		},

		ChildTask: &TaskSettings{
			Name:    childname,
			SQL:     "SELECT 1",
			Enabled: true,
			Comment: "third state",
		},

		SoloTask: &TaskSettings{
			Name:              soloname,
			Schema:            acc.TestSchemaName,
			SQL:               "SELECT *",
			When:              "FALSE",
			Enabled:           true,
			Schedule:          "15 MINUTE",
			UserTaskTimeoutMs: 900000,
		},
	}

	stepThree = &AccTaskTestSettings{ //nolint
		WarehouseName: warehousename,

		RootTask: &TaskSettings{
			Name:              rootname,
			Schema:            acc.TestSchemaName,
			SQL:               "SHOW FUNCTIONS",
			Enabled:           false,
			Schedule:          "5 MINUTE",
			UserTaskTimeoutMs: 1800000,
		},

		ChildTask: &TaskSettings{
			Name:    childname,
			SQL:     "SELECT 1",
			Enabled: false,
			Comment: "reset",
		},

		SoloTask: &TaskSettings{
			Name:    soloname,
			Schema:  acc.TestSchemaName,
			SQL:     "SELECT 1",
			When:    "TRUE",
			Enabled: true,
			SessionParams: map[string]string{
				"TIMESTAMP_INPUT_FORMAT": "YYYY-MM-DD HH24",
			},
			Schedule:          "5 MINUTE",
			UserTaskTimeoutMs: 0,
		},
	}
)

func TestAcc_Task(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskConfig(initialState),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", true),
					checkBool("snowflake_task.child_task", "enabled", false),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", initialState.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", initialState.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "after.0", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", initialState.ChildTask.Comment),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", initialState.RootTask.Schedule),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", initialState.ChildTask.Schedule),
					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", initialState.RootTask.UserTaskTimeoutMs),
					resource.TestCheckNoResourceAttr("snowflake_task.solo_task", "user_task_timeout_ms"),
				),
			},
			{
				Config: taskConfig(stepOne),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", true),
					checkBool("snowflake_task.child_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", stepOne.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", stepOne.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepOne.ChildTask.Comment),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", stepOne.RootTask.Schedule),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", stepOne.ChildTask.Schedule),
					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", stepOne.RootTask.UserTaskTimeoutMs),
					checkInt64("snowflake_task.solo_task", "user_task_timeout_ms", stepOne.SoloTask.UserTaskTimeoutMs),
				),
			},
			{
				Config: taskConfig(stepTwo),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", true),
					checkBool("snowflake_task.child_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", stepTwo.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", stepTwo.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepTwo.ChildTask.Comment),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", stepTwo.RootTask.Schedule),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", stepTwo.ChildTask.Schedule),
					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", stepTwo.RootTask.UserTaskTimeoutMs),
					checkInt64("snowflake_task.solo_task", "user_task_timeout_ms", stepTwo.SoloTask.UserTaskTimeoutMs),
				),
			},
			{
				Config: taskConfig(stepThree),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", false),
					checkBool("snowflake_task.child_task", "enabled", false),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", stepThree.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", stepThree.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepThree.ChildTask.Comment),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", stepThree.RootTask.Schedule),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", stepThree.ChildTask.Schedule),
					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", stepThree.RootTask.UserTaskTimeoutMs),
					checkInt64("snowflake_task.solo_task", "user_task_timeout_ms", stepThree.SoloTask.UserTaskTimeoutMs),
				),
			},
			{
				Config: taskConfig(initialState),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", true),
					checkBool("snowflake_task.child_task", "enabled", false),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql_statement", initialState.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql_statement", initialState.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", initialState.ChildTask.Comment),
					checkInt64("snowflake_task.root_task", "user_task_timeout_ms", stepOne.RootTask.UserTaskTimeoutMs),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schedule", initialState.RootTask.Schedule),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schedule", initialState.ChildTask.Schedule),
					// Terraform SDK is not able to differentiate if the
					// attribute has deleted or set to zero value.
					// ResourceData.GetChange returns the zero value of defined
					// type in schema as new the value. Provider handles 0 for
					// `user_task_timeout_ms` by unsetting the
					// USER_TASK_TIMEOUT_MS session variable.
					checkInt64("snowflake_task.solo_task", "user_task_timeout_ms", initialState.ChildTask.UserTaskTimeoutMs),
				),
			},
		},
	})
}

func taskConfig(settings *AccTaskTestSettings) string { //nolint
	config, err := template.New("task_acceptance_test_config").Parse(`
resource "snowflake_warehouse" "test_wh" {
	name = "{{ .WarehouseName }}"
}
resource "snowflake_task" "root_task" {
	name     	  = "{{ .RootTask.Name }}"
	database  	  = "terraform_test_database"
	schema   	  = "{{ .RootTask.Schema }}"
	warehouse 	  = snowflake_warehouse.test_wh.name
	sql_statement = "{{ .RootTask.SQL }}"
	enabled  	  = {{ .RootTask.Enabled }}
	schedule 	  = "{{ .RootTask.Schedule }}"
	{{ if .RootTask.UserTaskTimeoutMs }}
	user_task_timeout_ms = {{ .RootTask.UserTaskTimeoutMs }}
	{{- end }}

	{{ if .ChildTask.SessionParams }}
	session_parameters = {
	{{ range $key, $value :=  .RootTask.SessionParams}}
        {{ $key }} = "{{ $value }}",
	}
	{{- end }}
	{{- end }}
}
resource "snowflake_task" "child_task" {
	name     	  = "{{ .ChildTask.Name }}"
	database   	  = snowflake_task.root_task.database
	schema    	  = snowflake_task.root_task.schema
	warehouse 	  = snowflake_task.root_task.warehouse
	sql_statement = "{{ .ChildTask.SQL }}"
	enabled  	  = {{ .ChildTask.Enabled }}
	after    	  = [snowflake_task.root_task.name]
	comment 	  = "{{ .ChildTask.Comment }}"
	{{ if .ChildTask.UserTaskTimeoutMs }}
	user_task_timeout_ms = {{ .ChildTask.UserTaskTimeoutMs }}
	{{- end }}

	{{ if .ChildTask.SessionParams }}
	session_parameters = {
	{{ range $key, $value :=  .ChildTask.SessionParams}}
        {{ $key }} = "{{ $value }}",
	}
	{{- end }}
	{{- end }}
}
resource "snowflake_task" "solo_task" {
	name     	  = "{{ .SoloTask.Name }}"
	database  	  = "terraform_test_database"
	schema    	  = "{{ .SoloTask.Schema }}"
	warehouse 	  = snowflake_warehouse.test_wh.name
	sql_statement = "{{ .SoloTask.SQL }}"
	enabled  	  = {{ .SoloTask.Enabled }}
	when     	  = "{{ .SoloTask.When }}"
	{{ if .SoloTask.Schedule }}
	schedule    = "{{ .SoloTask.Schedule }}"
	{{- end }}

	{{ if .SoloTask.UserTaskTimeoutMs }}
	user_task_timeout_ms = {{ .SoloTask.UserTaskTimeoutMs }}
	{{- end }}

	{{ if .SoloTask.SessionParams }}
	session_parameters = {
	{{ range $key, $value :=  .SoloTask.SessionParams}}
        {{ $key }} = "{{ $value }}",
	}
	{{- end }}
	{{- end }}
}
	`)
	if err != nil {
		fmt.Println(err)
	}

	var result bytes.Buffer
	config.Execute(&result, settings) //nolint

	return result.String()
}

func TestAcc_Task_Managed(t *testing.T) {
	accName := "tst-terraform-" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	whName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskConfigManaged1(accName),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.managed_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "schedule", "5 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "user_task_managed_initial_warehouse_size", "XSMALL"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task_no_init", "user_task_managed_initial_warehouse_size", ""),
					resource.TestCheckResourceAttr("snowflake_task.managed_task_no_init", "session_parameters.TIMESTAMP_INPUT_FORMAT", "YYYY-MM-DD HH24"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "warehouse", ""),
				),
			},
			{
				Config: taskConfigManaged2(accName, whName),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.managed_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "schedule", "5 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "user_task_managed_initial_warehouse_size", ""),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "warehouse", whName),
				),
			},
			{
				Config: taskConfigManaged1(accName),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.managed_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "schedule", "5 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task_no_init", "session_parameters.TIMESTAMP_INPUT_FORMAT", "YYYY-MM-DD HH24"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task_no_init", "user_task_managed_initial_warehouse_size", ""),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "warehouse", ""),
				),
			},
			{
				Config: taskConfigManaged3(accName),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.managed_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "schedule", "5 MINUTE"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "user_task_managed_initial_warehouse_size", "SMALL"),
					resource.TestCheckResourceAttr("snowflake_task.managed_task", "warehouse", ""),
				),
			},
		},
	})
}

func taskConfigManaged1(name string) string {
	s := `
resource "snowflake_task" "managed_task" {
	name     	                             = "%s"
	database  	                             = "terraform_test_database"
	schema    	                             = "terraform_test_schema"
	sql_statement                            = "SELECT 1"
	enabled  	                             = true
	schedule                                 = "5 MINUTE"
    user_task_managed_initial_warehouse_size = "XSMALL"
}
resource "snowflake_task" "managed_task_no_init" {
	name     	  = "%s_no_init"
	database  	  = "terraform_test_database"
	schema    	  = "terraform_test_schema"
	sql_statement = "SELECT 1"
	enabled  	  = true
	schedule      = "5 MINUTE"
	session_parameters = {
		TIMESTAMP_INPUT_FORMAT = "YYYY-MM-DD HH24",
	}
}

`
	return fmt.Sprintf(s, name, name)
}

func taskConfigManaged2(name, whName string) string {
	s := `
resource "snowflake_warehouse" "test_wh" {
	name = "%s"
}

resource "snowflake_task" "managed_task" {
	name     	  = "%s"
	database  	  = "terraform_test_database"
	schema    	  = "terraform_test_schema"
	sql_statement = "SELECT 1"
	enabled  	  = true
	schedule      = "5 MINUTE"
	warehouse     = snowflake_warehouse.test_wh.name
}
`
	return fmt.Sprintf(s, whName, name)
}

func taskConfigManaged3(name string) string {
	s := `
resource "snowflake_task" "managed_task" {
	name     	                             = "%s"
	database  	                             = "terraform_test_database"
	schema    	                             = "terraform_test_schema"
	sql_statement                            = "SELECT 1"
	enabled  	                             = true
	schedule                                 = "5 MINUTE"
    user_task_managed_initial_warehouse_size = "SMALL"
}
`
	return fmt.Sprintf(s, name)
}

func TestAcc_Task_SwitchScheduled(t *testing.T) {
	accName := "tst-terraform-" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	taskRootName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskConfigManagedScheduled(accName, taskRootName),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.test_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schedule", "5 MINUTE"),
				),
			},
			{
				Config: taskConfigManagedScheduled2(accName, taskRootName),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.test_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schedule", ""),
				),
			},
			{
				Config: taskConfigManagedScheduled(accName, taskRootName),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.test_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schedule", "5 MINUTE"),
				),
			},
			{
				Config: taskConfigManagedScheduled3(accName, taskRootName),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.test_task", "enabled", false),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "sql_statement", "SELECT 1"),
					resource.TestCheckResourceAttr("snowflake_task.test_task", "schedule", ""),
				),
			},
		},
	})
}

func taskConfigManagedScheduled(name string, taskRootName string) string {
	s := `
resource "snowflake_task" "test_task_root" {
	name     	  = "%s"
	database  	  = "terraform_test_database"
	schema    	  = "terraform_test_schema"
	sql_statement = "SELECT 1"
	enabled  	  = true
	schedule      = "5 MINUTE"
}

resource "snowflake_task" "test_task" {
	name     	  = "%s"
	database  	  = "terraform_test_database"
	schema    	  = "terraform_test_schema"
	sql_statement = "SELECT 1"
	enabled  	  = true
	schedule      = "5 MINUTE"
}
`
	return fmt.Sprintf(s, taskRootName, name)
}

func taskConfigManagedScheduled2(name string, taskRootName string) string {
	s := `
resource "snowflake_task" "test_task_root" {
	name     	  = "%s"
	database  	  = "terraform_test_database"
	schema    	  = "terraform_test_schema"
	sql_statement = "SELECT 1"
	enabled  	  = true
	schedule      = "5 MINUTE"
}

resource "snowflake_task" "test_task" {
	name     	  = "%s"
	database  	  = "terraform_test_database"
	schema    	  = "terraform_test_schema"
	sql_statement = "SELECT 1"
	enabled  	  = true
	after         = [snowflake_task.test_task_root.name]
}
`
	return fmt.Sprintf(s, taskRootName, name)
}

func taskConfigManagedScheduled3(name string, taskRootName string) string {
	s := `
resource "snowflake_task" "test_task_root" {
	name     	  = "%s"
	database  	  = "terraform_test_database"
	schema    	  = "terraform_test_schema"
	sql_statement = "SELECT 1"
	enabled  	  = false
	schedule      = "5 MINUTE"
}

resource "snowflake_task" "test_task" {
	name     	  = "%s"
	database  	  = "terraform_test_database"
	schema    	  = "terraform_test_schema"
	sql_statement = "SELECT 1"
	enabled  	  = false
	after         = [snowflake_task.test_task_root.name]
}

`
	return fmt.Sprintf(s, taskRootName, name)
}

func checkInt64(name, key string, value int64) func(*terraform.State) error {
	return func(state *terraform.State) error {
		return resource.TestCheckResourceAttr(name, key, fmt.Sprintf("%v", value))(state)
	}
}
