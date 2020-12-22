package resources

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

const (
	delimiter         = '|'
	streamOndelimiter = '.'
)

func writeID(in []string) (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = delimiter
	err := csvWriter.WriteAll([][]string{in})
	if err != nil {
		return "", err
	}
	strGrantID := strings.TrimSpace(buf.String())
	return strGrantID, nil
}

// grantID contains identifying elements that allow unique access privileges
type grantID struct {
	ResourceName string
	SchemaName   string
	ObjectName   string
	Privilege    string
	GrantOption  bool
}

// String() takes in a grantID object and returns a pipe-delimited string:
// resourceName|schemaName|ObjectName|Privilege|GrantOption
func (gi *grantID) String() (string, error) {
	grantOption := fmt.Sprintf("%v", gi.GrantOption)
	dataIdentifiers := []string{gi.ResourceName, gi.SchemaName, gi.ObjectName, gi.Privilege, grantOption}
	return writeID(dataIdentifiers)
}

// grantIDFromString() takes in a pipe-delimited string: resourceName|schemaName|ObjectName|Privilege
// and returns a grantID object
func grantIDFromString(stringID string) (*grantID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = delimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per grant")
	}
	if len(lines[0]) != 4 && len(lines[0]) != 5 {
		return nil, fmt.Errorf("4 or 5 fields allowed")
	}

	grantOption := false
	if len(lines[0]) == 5 && lines[0][4] == "true" {
		grantOption = true
	}

	grantResult := &grantID{
		ResourceName: lines[0][0],
		SchemaName:   lines[0][1],
		ObjectName:   lines[0][2],
		Privilege:    lines[0][3],
		GrantOption:  grantOption,
	}
	return grantResult, nil
}

type pipeID struct {
	DatabaseName string
	SchemaName   string
	PipeName     string
}

//String() takes in a pipeID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|PipeName
func (si *pipeID) String() (string, error) {
	dataIdentifiers := []string{si.DatabaseName, si.SchemaName, si.PipeName}
	return writeID(dataIdentifiers)
}

// pipeIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|PipeName
// and returns a pipeID object
func pipeIDFromString(stringID string) (*pipeID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = delimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per pipe")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	pipeResult := &pipeID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		PipeName:     lines[0][2],
	}
	return pipeResult, nil
}

type schemaID struct {
	DatabaseName string
	SchemaName   string
}

// String() takes in a schemaID object and returns a pipe-delimited string:
// DatabaseName|schemaName
func (si *schemaID) String() (string, error) {
	dataIdentifiers := []string{si.DatabaseName, si.SchemaName}
	return writeID(dataIdentifiers)
}

// schemaIDFromString() takes in a pipe-delimited string: DatabaseName|schemaName
// and returns a schemaID object
func schemaIDFromString(stringID string) (*schemaID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = delimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per schema")
	}
	if len(lines[0]) != 2 {
		return nil, fmt.Errorf("2 fields allowed")
	}

	schemaResult := &schemaID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
	}
	return schemaResult, nil
}

type stageID struct {
	DatabaseName string
	SchemaName   string
	StageName    string
}

// String() takes in a stageID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|StageName
func (si *stageID) String() (string, error) {
	dataIdentifiers := []string{si.DatabaseName, si.SchemaName, si.StageName}
	return writeID(dataIdentifiers)
}

// stageIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|StageName
// and returns a stageID object
func stageIDFromString(stringID string) (*stageID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = delimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per stage")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	stageResult := &stageID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		StageName:    lines[0][2],
	}
	return stageResult, nil
}

type streamID struct {
	DatabaseName string
	SchemaName   string
	StreamName   string
}

type streamOnTableID struct {
	DatabaseName string
	SchemaName   string
	OnTableName  string
}

//String() takes in a streamID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|StreamName
func (si *streamID) String() (string, error) {
	dataIdentifiers := []string{si.DatabaseName, si.SchemaName, si.StreamName}
	return writeID(dataIdentifiers)
}

// streamIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|StreamName
// and returns a streamID object
func streamIDFromString(stringID string) (*streamID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = delimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	streamResult := &streamID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		StreamName:   lines[0][2],
	}
	return streamResult, nil
}

// streamOnTableIDFromString() takes in a dot-delimited string: DatabaseName.SchemaName.TableName
// and returns a streamOnTableID object
func streamOnTableIDFromString(stringID string) (*streamOnTableID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = streamOndelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		//return nil, fmt.Errorf("on table format: database_name.schema_name.target_table_name")
		return nil, fmt.Errorf("invalid format for on_table: %v , expected: <database_name.schema_name.target_table_name>", strings.Join(lines[0], "."))
	}

	streamOnTableResult := &streamOnTableID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		OnTableName:  lines[0][2],
	}
	return streamOnTableResult, nil
}

type tableID struct {
	DatabaseName string
	SchemaName   string
	TableName    string
}

//String() takes in a tableID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|TableName
func (si *tableID) String() (string, error) {
	dataIdentifiers := []string{si.DatabaseName, si.SchemaName, si.TableName}
	return writeID(dataIdentifiers)
}

// tableIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TableName
// and returns a tableID object
func tableIDFromString(stringID string) (*tableID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = delimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	tableResult := &tableID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		TableName:    lines[0][2],
	}
	return tableResult, nil
}

type taskID struct {
	DatabaseName string
	SchemaName   string
	TaskName     string
}

//String() takes in a taskID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|TaskName
func (t *taskID) String() (string, error) {
	dataIdentifiers := []string{t.DatabaseName, t.SchemaName, t.TaskName}
	return writeID(dataIdentifiers)
}

// taskIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TaskName
// and returns a taskID object
func taskIDFromString(stringID string) (*taskID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = delimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per task")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	taskResult := &taskID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		TaskName:     lines[0][2],
	}
	return taskResult, nil
}
