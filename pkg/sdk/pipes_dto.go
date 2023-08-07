package sdk

import (
	"bytes"

	"github.com/stretchr/testify/assert"
)

//go:generate go run ./dto-builder-generator/main.go

type CreatePipeRequest struct {
	orReplace               bool
	ifNotExists             bool
	name                    SchemaObjectIdentifier // required
	autoIngest              bool
	errorIntegration        string
	awsSnsTopic             string
	integration             string
	comment                 string
	copyStatement           string // required
	exampleOfImport         bytes.Buffer
	exampleOfExternalImport assert.Assertions
}

type AlterPipeRequest struct {
	ifExists  bool
	name      SchemaObjectIdentifier // required
	set       PipeSetRequest
	unset     PipeUnsetRequest
	setTags   PipeSetTagsRequest
	unsetTags PipeUnsetTagsRequest
	refresh   PipeRefreshRequest
}

type PipeSetRequest struct {
	errorIntegration    string
	pipeExecutionPaused bool
	comment             string
}

type PipeUnsetRequest struct {
	pipeExecutionPaused bool
	comment             bool
}

type PipeSetTagsRequest struct {
	tag []TagAssociation // required
}

type PipeUnsetTagsRequest struct {
	tag []ObjectIdentifier // required
}

type PipeRefreshRequest struct {
	prefix        string
	modifiedAfter string
}
