package resources_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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

func TestAcc_Alert(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	alertInitialState := &AccAlertTestSettings{ //nolint
		WarehouseName: acc.TestWarehouseName,
		DatabaseName:  acc.TestDatabaseName,
		Alert: &AlertSettings{
			Name:      id.Name(),
			Condition: "select 0 as c",
			Action:    "select 0 as c",
			Schema:    acc.TestSchemaName,
			Enabled:   true,
			Schedule:  5,
			Comment:   "dummy",
		},
	}

	// Changes: condition, action, comment, schedule.
	alertStepOne := &AccAlertTestSettings{ //nolint
		WarehouseName: acc.TestWarehouseName,
		DatabaseName:  acc.TestDatabaseName,
		Alert: &AlertSettings{
			Name:      id.Name(),
			Condition: "select 1 as c",
			Action:    "select 1 as c",
			Schema:    acc.TestSchemaName,
			Enabled:   true,
			Schedule:  15,
			Comment:   "test",
		},
	}

	// Changes: condition, action, comment, schedule.
	alertStepTwo := &AccAlertTestSettings{ //nolint
		WarehouseName: acc.TestWarehouseName,
		DatabaseName:  acc.TestDatabaseName,
		Alert: &AlertSettings{
			Name:      id.Name(),
			Condition: "select 2 as c",
			Action:    "select 2 as c",
			Schema:    acc.TestSchemaName,
			Enabled:   true,
			Schedule:  25,
			Comment:   "text",
		},
	}

	// Changes: condition, action, comment, schedule.
	alertStepThree := &AccAlertTestSettings{ //nolint
		WarehouseName: acc.TestWarehouseName,
		DatabaseName:  acc.TestDatabaseName,
		Alert: &AlertSettings{
			Name:      id.Name(),
			Condition: "select 2 as c",
			Action:    "select 2 as c",
			Schema:    acc.TestSchemaName,
			Enabled:   false,
			Schedule:  5,
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Alert),
		Steps: []resource.TestStep{
			{
				Config: alertConfig(alertInitialState),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_alert.test_alert", "enabled", alertInitialState.Alert.Enabled),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", acc.TestSchemaName),
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
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", acc.TestSchemaName),
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
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", acc.TestSchemaName),
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
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", acc.TestSchemaName),
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
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_alert.test_alert", "schema", acc.TestSchemaName),
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
resource "snowflake_alert" "test_alert" {
	name     	      = "{{ .Alert.Name }}"
	database  	      = "{{ .DatabaseName }}"
	schema   	      = "{{ .Alert.Schema }}"
	warehouse 	      = "{{ .WarehouseName }}"
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
