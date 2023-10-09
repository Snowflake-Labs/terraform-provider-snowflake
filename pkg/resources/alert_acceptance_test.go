package resources_test

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"text/template"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type (
	AccAlertTestSettings struct {
		WarehouseName string
		DatabaseName  string
		Alert         *AlertSettings
	}

	AlertSettings struct {
		Name      string
		Enabled   bool
		Schema    string
		Condition string
		Action    string
		Schedule  int
		Comment   string
	}
)

var (
	warehouseName = "wh_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseName  = "db_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName    = "PUBLIC"
	alertName     = "a_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	alertInitialState = &AccAlertTestSettings{ //nolint
		WarehouseName: warehouseName,
		DatabaseName:  databaseName,

		Alert: &AlertSettings{
			Name:      alertName,
			Schema:    schemaName,
			Condition: "select 0 as c",
			Action:    "select 0 as c",
			Enabled:   true,
			Schedule:  5,
			Comment:   "dummy",
		},
	}

	// Changes: condition, action, comment, schedule.
	alertStepOne = &AccAlertTestSettings{ //nolint
		WarehouseName: warehouseName,
		DatabaseName:  databaseName,

		Alert: &AlertSettings{
			Name:      alertName,
			Schema:    schemaName,
			Condition: "select 1 as c",
			Action:    "select 1 as c",
			Enabled:   true,
			Schedule:  15,
			Comment:   "test",
		},
	}

	// Changes: condition, action, comment, schedule.
	alertStepTwo = &AccAlertTestSettings{ //nolint
		WarehouseName: warehouseName,
		DatabaseName:  databaseName,

		Alert: &AlertSettings{
			Name:      alertName,
			Schema:    schemaName,
			Condition: "select 2 as c",
			Action:    "select 2 as c",
			Enabled:   true,
			Schedule:  25,
			Comment:   "text",
		},
	}

	// Changes: condition, action, comment, schedule.
	alertStepThree = &AccAlertTestSettings{ //nolint
		WarehouseName: warehouseName,
		DatabaseName:  databaseName,

		Alert: &AlertSettings{
			Name:      alertName,
			Schema:    schemaName,
			Condition: "select 2 as c",
			Action:    "select 2 as c",
			Enabled:   false,
			Schedule:  5,
		},
	}
)

func TestAcc_Alert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: alertConfig(alertInitialState),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_alert.test_alert", "enabled", alertInitialState.Alert.Enabled),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", alertName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", alertInitialState.Alert.Condition),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", alertInitialState.Alert.Action),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", alertInitialState.Alert.Comment),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(alertInitialState.Alert.Schedule)),
				),
			},
			{
				Config: alertConfig(alertStepOne),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_alert.test_alert", "enabled", alertStepOne.Alert.Enabled),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", alertName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", alertStepOne.Alert.Condition),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", alertStepOne.Alert.Action),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", alertStepOne.Alert.Comment),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(alertStepOne.Alert.Schedule)),
				),
			},
			{
				Config: alertConfig(alertStepTwo),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_alert.test_alert", "enabled", alertStepTwo.Alert.Enabled),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", alertName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", alertStepTwo.Alert.Condition),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", alertStepTwo.Alert.Action),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", alertStepTwo.Alert.Comment),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(alertStepTwo.Alert.Schedule)),
				),
			},
			{
				Config: alertConfig(alertStepThree),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_alert.test_alert", "enabled", alertStepThree.Alert.Enabled),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", alertName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", alertStepThree.Alert.Condition),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", alertStepThree.Alert.Action),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", alertStepThree.Alert.Comment),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(alertStepThree.Alert.Schedule)),
				),
			},
			{
				Config: alertConfig(alertInitialState),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_alert.test_alert", "enabled", alertInitialState.Alert.Enabled),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", alertName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", schemaName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "condition", alertInitialState.Alert.Condition),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "action", alertInitialState.Alert.Action),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "comment", alertInitialState.Alert.Comment),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "alert_schedule.0.interval", strconv.Itoa(alertInitialState.Alert.Schedule)),
				),
			},
		},
	})
}

func alertConfig(settings *AccAlertTestSettings) string { //nolint
	config, err := template.New("alert_acceptance_test_config").Parse(`
resource "snowflake_warehouse" "test_wh" {
	name = "{{ .WarehouseName }}"
}
resource "snowflake_database" "test_db" {
	name = "{{ .DatabaseName }}"
}
resource "snowflake_alert" "test_alert" {
	name     	      = "{{ .Alert.Name }}"
	database  	      = snowflake_database.test_db.name
	schema   	      = "{{ .Alert.Schema }}"
	warehouse 	      = snowflake_warehouse.test_wh.name
	alert_schedule 	  {
		interval = "{{ .Alert.Schedule }}"
	}
	condition         = "{{ .Alert.Condition }}"
	action            = "{{ .Alert.Action }}"
	enabled  	      = {{ .Alert.Enabled }}
	comment           = "{{ .Alert.Comment }}"
}
	`)
	if err != nil {
		fmt.Println(err)
	}

	var result bytes.Buffer
	err = config.Execute(&result, settings) //nolint
	if err != nil {
		fmt.Println(err)
	}
	return result.String()
}
