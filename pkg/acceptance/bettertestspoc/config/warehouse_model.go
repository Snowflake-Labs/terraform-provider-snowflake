package config

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

// TODO: add possibility to have reference to another object (e.g. WithResourceMonitorReference); new config.Variable impl?
// TODO: add possibility to have depends_on to other resources (in meta?)
// TODO: add a convenience method to use multiple configs from multiple models
type WarehouseModel struct {
	Name                            config.Variable `json:"name,omitempty"`
	WarehouseType                   config.Variable `json:"warehouse_type,omitempty"`
	WarehouseSize                   config.Variable `json:"warehouse_size,omitempty"`
	MaxClusterCount                 config.Variable `json:"max_cluster_count,omitempty"`
	MinClusterCount                 config.Variable `json:"min_cluster_count,omitempty"`
	ScalingPolicy                   config.Variable `json:"scaling_policy,omitempty"`
	AutoSuspend                     config.Variable `json:"auto_suspend,omitempty"`
	AutoResume                      config.Variable `json:"auto_resume,omitempty"`
	InitiallySuspended              config.Variable `json:"initially_suspended,omitempty"`
	ResourceMonitor                 config.Variable `json:"resource_monitor,omitempty"`
	Comment                         config.Variable `json:"comment,omitempty"`
	EnableQueryAcceleration         config.Variable `json:"enable_query_acceleration,omitempty"`
	QueryAccelerationMaxScaleFactor config.Variable `json:"query_acceleration_max_scale_factor,omitempty"`

	MaxConcurrencyLevel             config.Variable `json:"max_concurrency_level,omitempty"`
	StatementQueuedTimeoutInSeconds config.Variable `json:"statement_queued_timeout_in_seconds,omitempty"`
	StatementTimeoutInSeconds       config.Variable `json:"statement_timeout_in_seconds,omitempty"`

	*ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func NewWarehouseModel(
	resourceName string,
	name string,
) *WarehouseModel {
	m := &WarehouseModel{ResourceModelMeta: Meta(resourceName, resources.Warehouse)}
	m.WithName(name)
	return m
}

func NewDefaultWarehouseModel(
	name string,
) *WarehouseModel {
	m := &WarehouseModel{ResourceModelMeta: DefaultMeta(resources.Warehouse)}
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
