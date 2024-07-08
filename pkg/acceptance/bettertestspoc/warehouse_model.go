package bettertestspoc

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

type WarehouseModel struct {
	Name                            config.Variable
	WarehouseType                   config.Variable
	WarehouseSize                   config.Variable
	MaxClusterCount                 config.Variable
	MinClusterCount                 config.Variable
	ScalingPolicy                   config.Variable
	AutoSuspend                     config.Variable
	AutoResume                      config.Variable
	InitiallySuspended              config.Variable
	ResourceMonitor                 config.Variable
	Comment                         config.Variable
	EnableQueryAcceleration         config.Variable
	QueryAccelerationMaxScaleFactor config.Variable

	MaxConcurrencyLevel             config.Variable
	StatementQueuedTimeoutInSeconds config.Variable
	StatementTimeoutInSeconds       config.Variable
}

///////////////////////////////////
// Basic builder (only required) //
///////////////////////////////////

func NewWarehouseModel(
	name string,
) *WarehouseModel {
	m := &WarehouseModel{}
	m.WithName(name)
	return m
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (m *WarehouseModel) WithName(name string) *WarehouseModel {
	m.Name = config.StringVariable(name)
	return m
}

func (m *WarehouseModel) WithWarehouseType(warehouseType sdk.WarehouseType) *WarehouseModel {
	m.WarehouseType = config.StringVariable(string(warehouseType))
	return m
}

func (m *WarehouseModel) WithWarehouseSize(warehouseSize sdk.WarehouseSize) *WarehouseModel {
	m.WarehouseSize = config.StringVariable(string(warehouseSize))
	return m
}

func (m *WarehouseModel) WithMaxClusterCount(maxClusterCount int) *WarehouseModel {
	m.MaxClusterCount = config.IntegerVariable(maxClusterCount)
	return m
}

func (m *WarehouseModel) WithMinClusterCount(minClusterCount int) *WarehouseModel {
	m.MinClusterCount = config.IntegerVariable(minClusterCount)
	return m
}

func (m *WarehouseModel) WithScalingPolicy(scalingPolicy sdk.ScalingPolicy) *WarehouseModel {
	m.ScalingPolicy = config.StringVariable(string(scalingPolicy))
	return m
}

func (m *WarehouseModel) WithAutoSuspend(autoSuspend int) *WarehouseModel {
	m.AutoSuspend = config.IntegerVariable(autoSuspend)
	return m
}

func (m *WarehouseModel) WithAutoResume(autoResume bool) *WarehouseModel {
	m.AutoResume = config.BoolVariable(autoResume)
	return m
}

func (m *WarehouseModel) WithInitiallySuspended(initiallySuspended bool) *WarehouseModel {
	m.InitiallySuspended = config.BoolVariable(initiallySuspended)
	return m
}

func (m *WarehouseModel) WithResourceMonitor(resourceMonitor sdk.AccountObjectIdentifier) *WarehouseModel {
	m.ResourceMonitor = config.StringVariable(resourceMonitor.Name())
	return m
}

func (m *WarehouseModel) WithComment(comment string) *WarehouseModel {
	m.Comment = config.StringVariable(comment)
	return m
}

func (m *WarehouseModel) WithEnableQueryAcceleration(enableQueryAcceleration bool) *WarehouseModel {
	m.EnableQueryAcceleration = config.BoolVariable(enableQueryAcceleration)
	return m
}

func (m *WarehouseModel) WithQueryAccelerationMaxScaleFactor(queryAccelerationMaxScaleFactor int) *WarehouseModel {
	m.QueryAccelerationMaxScaleFactor = config.IntegerVariable(queryAccelerationMaxScaleFactor)
	return m
}

func (m *WarehouseModel) WithMaxConcurrencyLevel(maxConcurrencyLevel int) *WarehouseModel {
	m.MaxConcurrencyLevel = config.IntegerVariable(maxConcurrencyLevel)
	return m
}

func (m *WarehouseModel) WithStatementQueuedTimeoutInSeconds(statementQueuedTimeoutInSeconds int) *WarehouseModel {
	m.StatementQueuedTimeoutInSeconds = config.IntegerVariable(statementQueuedTimeoutInSeconds)
	return m
}

func (m *WarehouseModel) WithStatementTimeoutInSeconds(statementTimeoutInSeconds int) *WarehouseModel {
	m.StatementTimeoutInSeconds = config.IntegerVariable(statementTimeoutInSeconds)
	return m
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (m *WarehouseModel) WithNameValue(value config.Variable) *WarehouseModel {
	m.Name = value
	return m
}

func (m *WarehouseModel) WithWarehouseTypeValue(value config.Variable) *WarehouseModel {
	m.WarehouseType = value
	return m
}

func (m *WarehouseModel) WithWarehouseSizeValue(value config.Variable) *WarehouseModel {
	m.WarehouseSize = value
	return m
}

func (m *WarehouseModel) WithMaxClusterCountValue(value config.Variable) *WarehouseModel {
	m.MaxClusterCount = value
	return m
}

func (m *WarehouseModel) WithMinClusterCountValue(value config.Variable) *WarehouseModel {
	m.MinClusterCount = value
	return m
}

func (m *WarehouseModel) WithScalingPolicyValue(value config.Variable) *WarehouseModel {
	m.ScalingPolicy = value
	return m
}

func (m *WarehouseModel) WithAutoSuspendValue(value config.Variable) *WarehouseModel {
	m.AutoSuspend = value
	return m
}

func (m *WarehouseModel) WithAutoResumeValue(value config.Variable) *WarehouseModel {
	m.AutoResume = value
	return m
}

func (m *WarehouseModel) WithInitiallySuspendedValue(value config.Variable) *WarehouseModel {
	m.InitiallySuspended = value
	return m
}

func (m *WarehouseModel) WithResourceMonitorValue(value config.Variable) *WarehouseModel {
	m.ResourceMonitor = value
	return m
}

func (m *WarehouseModel) WithCommentValue(value config.Variable) *WarehouseModel {
	m.Comment = value
	return m
}

func (m *WarehouseModel) WithEnableQueryAccelerationValue(value config.Variable) *WarehouseModel {
	m.EnableQueryAcceleration = value
	return m
}

func (m *WarehouseModel) WithQueryAccelerationMaxScaleFactorValue(value config.Variable) *WarehouseModel {
	m.QueryAccelerationMaxScaleFactor = value
	return m
}

func (m *WarehouseModel) WithMaxConcurrencyLevelValue(value config.Variable) *WarehouseModel {
	m.MaxConcurrencyLevel = value
	return m
}

func (m *WarehouseModel) WithStatementQueuedTimeoutInSecondsValue(value config.Variable) *WarehouseModel {
	m.StatementQueuedTimeoutInSeconds = value
	return m
}

func (m *WarehouseModel) WithStatementTimeoutInSecondsValue(value config.Variable) *WarehouseModel {
	m.StatementTimeoutInSeconds = value
	return m
}
