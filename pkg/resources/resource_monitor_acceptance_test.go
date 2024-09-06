package resources_test

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	configvariable "github.com/hashicorp/terraform-plugin-testing/config"
	"regexp"
	"strings"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ResourceMonitor_Basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoCreditQuota().
						HasNotifyUsersLen(0).
						HasNoFrequency().
						HasNoStartTimestamp().
						HasNoEndTimestamp().
						HasTriggerLen(0),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(0).
						HasUsedCredits(0).
						HasRemainingCredits(0).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTime("").
						HasSuspendAt(0).
						HasSuspendImmediateAt(0).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			{
				ResourceName: "snowflake_resource_monitor.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedResourceMonitorResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("0").
						HasNotifyUsersLen(0).
						HasFrequencyString(string(sdk.FrequencyMonthly)).
						HasStartTimestampNotEmpty().
						HasEndTimestampString("").
						HasTriggerLen(0),
				),
			},
		},
	})
}

func TestAcc_ResourceMonitor_Complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name()).
		WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("JAN_CIESLAK"))).
		WithCreditQuota(10).
		WithFrequency(string(sdk.FrequencyWeekly)).
		WithStartTimestamp(time.Now().Add(time.Hour * 24 * 30).Format("2006-01-02 15:01")).
		WithEndTimestamp(time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02 15:01")).
		WithTriggerValue(configvariable.SetVariable(
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(100),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionNotify)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(110),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionNotify)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(120),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionSuspend)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(150),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionSuspendImmediate)),
			}),
		))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ResourceMonitor/complete"),
				ConfigVariables: configModel.ToConfigVariables(),
				Check: assert.AssertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("10").
						HasNotifyUsersLen(1).
						HasNotifyUser(0, "JAN_CIESLAK").
						HasFrequencyString(string(sdk.FrequencyWeekly)).
						HasStartTimestampString(time.Now().Add(time.Hour*24*30).Format("2006-01-02 15:01")).
						HasEndTimestampString(time.Now().Add(time.Hour*24*60).Format("2006-01-02 15:01")).
						HasTriggerLen(4).
						HasTrigger(0, 100, sdk.TriggerActionNotify).
						HasTrigger(1, 110, sdk.TriggerActionNotify).
						HasTrigger(2, 120, sdk.TriggerActionSuspend).
						HasTrigger(3, 150, sdk.TriggerActionSuspendImmediate),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(10).
						HasUsedCredits(0).
						HasRemainingCredits(10).
						HasLevel("").
						HasFrequency(sdk.FrequencyWeekly).
						HasStartTimeNotEmpty().
						HasEndTimeNotEmpty().
						HasSuspendAt(120).
						HasSuspendImmediateAt(150).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			{
				ResourceName:    "snowflake_resource_monitor.test",
				ImportState:     true,
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ResourceMonitor/complete"),
				ConfigVariables: configModel.ToConfigVariables(),
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedResourceMonitorResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("10").
						HasNotifyUsersLen(1).
						HasNotifyUser(0, "JAN_CIESLAK").
						HasFrequencyString(string(sdk.FrequencyWeekly)).
						HasStartTimestampNotEmpty().
						HasEndTimestampNotEmpty().
						HasTriggerLen(4).
						HasTrigger(0, 100, sdk.TriggerActionNotify).
						HasTrigger(1, 110, sdk.TriggerActionNotify).
						HasTrigger(2, 120, sdk.TriggerActionSuspend).
						HasTrigger(3, 150, sdk.TriggerActionSuspendImmediate),
				),
			},
		},
	})
}

func TestAcc_ResourceMonitor_Updates(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configModelNothingSet := model.ResourceMonitor("test", id.Name())

	configModelEverythingSet := model.ResourceMonitor("test", id.Name()).
		WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("JAN_CIESLAK"))).
		WithCreditQuota(10).
		WithFrequency(string(sdk.FrequencyWeekly)).
		WithStartTimestamp(time.Now().Add(time.Hour * 24 * 30).Format("2006-01-02 15:01")).
		WithEndTimestamp(time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02 15:01")).
		WithTriggerValue(configvariable.SetVariable(
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(100),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionNotify)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(110),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionNotify)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(120),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionSuspend)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(150),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionSuspendImmediate)),
			}),
		))

	configModelUpdated := model.ResourceMonitor("test", id.Name()).
		WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("JAN_CIESLAK"), configvariable.StringVariable("ARTUR_SAWICKI"))).
		WithCreditQuota(20).
		WithFrequency(string(sdk.FrequencyMonthly)).
		WithStartTimestamp(time.Now().Add(time.Hour * 24 * 40).Format("2006-01-02 15:01")).
		WithEndTimestamp(time.Now().Add(time.Hour * 24 * 70).Format("2006-01-02 15:01")).
		WithTriggerValue(configvariable.SetVariable(
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(110),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionNotify)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(120),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionNotify)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(130),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionSuspend)),
			}),
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(160),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionSuspendImmediate)),
			}),
		))

	configModelEverythingUnset := model.ResourceMonitor("test", id.Name()).
		WithTriggerValue(configvariable.SetVariable(
			configvariable.ObjectVariable(map[string]configvariable.Variable{
				"threshold":            configvariable.IntegerVariable(100),
				"on_threshold_reached": configvariable.StringVariable(string(sdk.TriggerActionNotify)),
			}),
		))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModelNothingSet),
				Check: assert.AssertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoCreditQuota().
						HasNotifyUsersLen(0).
						HasNoFrequency().
						HasNoStartTimestamp().
						HasNoEndTimestamp().
						HasTriggerLen(0),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(0).
						HasUsedCredits(0).
						HasRemainingCredits(0).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTime("").
						HasSuspendAt(0).
						HasSuspendImmediateAt(0).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			// Set
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ResourceMonitor/complete"),
				ConfigVariables: configModelEverythingSet.ToConfigVariables(),
				Check: assert.AssertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("10").
						HasNotifyUsersLen(1).
						HasNotifyUser(0, "JAN_CIESLAK").
						HasFrequencyString(string(sdk.FrequencyWeekly)).
						HasStartTimestampString(time.Now().Add(time.Hour*24*30).Format("2006-01-02 15:01")).
						HasEndTimestampString(time.Now().Add(time.Hour*24*60).Format("2006-01-02 15:01")).
						HasTriggerLen(4).
						HasTrigger(0, 100, sdk.TriggerActionNotify).
						HasTrigger(1, 110, sdk.TriggerActionNotify).
						HasTrigger(2, 120, sdk.TriggerActionSuspend).
						HasTrigger(3, 150, sdk.TriggerActionSuspendImmediate),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(10).
						HasUsedCredits(0).
						HasRemainingCredits(10).
						HasLevel("").
						HasFrequency(sdk.FrequencyWeekly).
						HasStartTimeNotEmpty().
						HasEndTimeNotEmpty().
						HasSuspendAt(120).
						HasSuspendImmediateAt(150).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			// Update
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ResourceMonitor/complete"),
				ConfigVariables: configModelUpdated.ToConfigVariables(),
				Check: assert.AssertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("20").
						HasNotifyUsersLen(2).
						HasNotifyUser(0, "ARTUR_SAWICKI").
						HasNotifyUser(1, "JAN_CIESLAK").
						HasFrequencyString(string(sdk.FrequencyMonthly)).
						HasStartTimestampString(time.Now().Add(time.Hour*24*40).Format("2006-01-02 15:01")).
						HasEndTimestampString(time.Now().Add(time.Hour*24*70).Format("2006-01-02 15:01")).
						HasTriggerLen(4).
						HasTrigger(0, 110, sdk.TriggerActionNotify).
						HasTrigger(1, 120, sdk.TriggerActionNotify).
						HasTrigger(2, 130, sdk.TriggerActionSuspend).
						HasTrigger(3, 160, sdk.TriggerActionSuspendImmediate),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(20).
						HasUsedCredits(0).
						HasRemainingCredits(20).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTimeNotEmpty().
						HasSuspendAt(130).
						HasSuspendImmediateAt(160).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			// Unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ResourceMonitor/only_triggers"),
				ConfigVariables: configModelEverythingUnset.ToConfigVariables(),
				Check: assert.AssertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("0").
						HasNotifyUsersLen(0).
						HasFrequencyString("").
						HasStartTimestampString("").
						HasEndTimestampString("").
						HasTriggerLen(1).
						HasTrigger(0, 100, sdk.TriggerActionNotify),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(0).
						HasUsedCredits(0).
						HasRemainingCredits(0).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTime("").
						HasSuspendAt(0).
						HasSuspendImmediateAt(0).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
		},
	})
}

// TestAcc_ResourceMonitor_issue2167 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2167 issue.
// Second step is purposely error, because tests TestAcc_ResourceMonitorUpdateNotifyUsers and TestAcc_ResourceMonitorNotifyUsers are still skipped.
// It can be fixed with them.
func TestAcc_ResourceMonitor_issue2167(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	configNoUsers := model.ResourceMonitor("test", id.Name()).WithNotifyUsersValue(configvariable.SetVariable())
	configWithNonExistingUser := model.ResourceMonitor("test", id.Name()).WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("non_existing_user")))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configNoUsers),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", id.Name()),
				),
			},
			{
				Config:      config.FromModel(t, configWithNonExistingUser),
				ExpectError: regexp.MustCompile(`.*090268 \(22023\): User non_existing_user does not exist.*`),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1990 is fixed
func TestAcc_ResourceMonitor_Issue1990_RemovingResourceMonitorOutsideOfTerraform(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			// Create resource monitor
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.69.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: config.FromModel(t, configModel),
			},
			// Same configuration, but we drop resource monitor externally
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.69.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				PreConfig: func() {
					acc.TestClient().ResourceMonitor.DropResourceMonitorFunc(t, id)()
				},
				Config:      config.FromModel(t, configModel),
				ExpectError: regexp.MustCompile("object does not exist or not authorized"),
			},
			// Same configuration, but it's the last version where it's still not working
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.95.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config:      config.FromModel(t, configModel),
				ExpectError: regexp.MustCompile("object does not exist or not authorized"),
			},
			// Same configuration, but it's the latest version of the provider (0.96.0 and above)
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModel(t, configModel),
			},
		},
	})
}

// TODO: Timestamp issues
// TODO: Reference related issues
func TestAcc_ResourceMonitor_Issue_TimestampInfinitePlan(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name())

	// Steps
	// old version
	// - create with and without timestamps
	// - different formats
	// - same format as in Snowflake
	// new version
	// - create with and without timestamps
	// - different formats
	// - same format as in Snowflake

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			// Create resource monitor
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.69.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: config.FromModel(t, configModel),
			},
			// Same configuration, but we drop resource monitor externally
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.69.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				PreConfig: func() {
					acc.TestClient().ResourceMonitor.DropResourceMonitorFunc(t, id)()
				},
				Config:      config.FromModel(t, configModel),
				ExpectError: regexp.MustCompile("object does not exist or not authorized"),
			},
			// Same configuration, but it's the last version where it's still not working
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.95.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config:      config.FromModel(t, configModel),
				ExpectError: regexp.MustCompile("object does not exist or not authorized"),
			},
			// Same configuration, but it's the latest version of the provider (0.96.0 and above)
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModel(t, configModel),
			},
		},
	})
}

// TODO: Issue #1500 (creating and altering resource monitor with only triggers)
// - On create we have required_with validation, so it's not possible
// - On update we can set e.g. credit_quota to the same (or new) value and it will work.

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1500 is fixed
func TestAcc_ResourceMonitor_Issue1500_CreatingAndAlteringResourceMonitorWithOnlyTriggers(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	//configModel := model.ResourceMonitor("test", id.Name())
	triggers := []map[string]any{
		{
			"threshold":            100,
			"on_threshold_reached": string(sdk.TriggerActionNotify),
		},
		{
			"threshold":            120,
			"on_threshold_reached": string(sdk.TriggerActionNotify),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			// Create resource monitor with only triggers (old version)
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.55.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config:      resourceMonitorConfigWithOnlyTriggers(t, id.Name(), triggers),
				ExpectError: regexp.MustCompile("LULULULUL"),
			},
		},
	})
}

func resourceMonitorConfigWithOnlyTriggers(t *testing.T, name string, triggers []map[string]any) string {
	t.Helper()

	triggersBuilder := new(strings.Builder)
	for _, trigger := range triggers {
		triggersBuilder.WriteString(fmt.Sprintf(`
trigger {
	threshold = %d
	on_threshold_reached = "%s"
}
`, trigger["threshold"], trigger["on_threshold_reached"]))
	}

	return fmt.Sprintf(`
resource "snowflake_resource_monitor" "test" {
  name     = "%[1]s"

  %[2]s
}
`, name, triggersBuilder.String())
}
