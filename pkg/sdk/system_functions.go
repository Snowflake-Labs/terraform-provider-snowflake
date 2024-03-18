package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

type SystemFunctions interface {
	GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) (string, error)
	PipeStatus(pipeId SchemaObjectIdentifier) (PipeExecutionState, error)
}

var _ SystemFunctions = (*systemFunctions)(nil)

type systemFunctions struct {
	client *Client
}

func (c *systemFunctions) GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) (string, error) {
	s := &struct {
		Tag string `db:"TAG"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$GET_TAG('%s', '%s', '%v') AS "TAG"`, tagID.FullyQualifiedName(), objectID.FullyQualifiedName(), objectType)
	err := c.client.queryOne(ctx, s, sql)
	if err != nil {
		return "", err
	}
	return s.Tag, nil
}

type PipeExecutionState string

const (
	FailingOverPipeExecutionState                           PipeExecutionState = "FAILING_OVER"
	PausedPipeExecutionState                                PipeExecutionState = "PAUSED"
	ReadOnlyPipeExecutionState                              PipeExecutionState = "READ_ONLY"
	RunningPipeExecutionState                               PipeExecutionState = "RUNNING"
	StoppedBySnowflakeAdminPipeExecutionState               PipeExecutionState = "STOPPED_BY_SNOWFLAKE_ADMIN"
	StoppedClonedPipeExecutionState                         PipeExecutionState = "STOPPED_CLONED"
	StoppedFeatureDisabledPipeExecutionState                PipeExecutionState = "STOPPED_FEATURE_DISABLED"
	StoppedStageAlteredPipeExecutionState                   PipeExecutionState = "STOPPED_STAGE_ALTERED"
	StoppedStageDroppedPipeExecutionState                   PipeExecutionState = "STOPPED_STAGE_DROPPED"
	StoppedFileFormatDroppedPipeExecutionState              PipeExecutionState = "STOPPED_FILE_FORMAT_DROPPED"
	StoppedNotificationIntegrationDroppedPipeExecutionState PipeExecutionState = "STOPPED_NOTIFICATION_INTEGRATION_DROPPED"
	StoppedMissingPipePipeExecutionState                    PipeExecutionState = "STOPPED_MISSING_PIPE"
	StoppedMissingTablePipeExecutionState                   PipeExecutionState = "STOPPED_MISSING_TABLE"
	StalledCompilationErrorPipeExecutionState               PipeExecutionState = "STALLED_COMPILATION_ERROR"
	StalledInitializationErrorPipeExecutionState            PipeExecutionState = "STALLED_INITIALIZATION_ERROR"
	StalledExecutionErrorPipeExecutionState                 PipeExecutionState = "STALLED_EXECUTION_ERROR"
	StalledInternalErrorPipeExecutionState                  PipeExecutionState = "STALLED_INTERNAL_ERROR"
	StalledStagePermissionErrorPipeExecutionState           PipeExecutionState = "STALLED_STAGE_PERMISSION_ERROR"
)

func (c *systemFunctions) PipeStatus(pipeId SchemaObjectIdentifier) (PipeExecutionState, error) {
	row := &struct {
		PipeStatus string `db:"PIPE_STATUS"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$PIPE_STATUS('%s') AS "PIPE_STATUS"`, pipeId.FullyQualifiedName())
	ctx := context.Background()

	err := c.client.queryOne(ctx, row, sql)
	if err != nil {
		return "", err
	}

	pipeStatus := make(map[string]any)
	err = json.Unmarshal([]byte(row.PipeStatus), pipeStatus)
	if err != nil {
		return "", err
	}

	return PipeExecutionState(pipeStatus["executionState"].(string)), nil
}
