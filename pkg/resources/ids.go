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

func readID(id string, minFields, maxFields int) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(id))
	reader.Comma = delimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("expecting 1 line")
	}
	if len(lines[0]) < minFields || len(lines[0]) > maxFields {
		if minFields == maxFields {
			return nil, fmt.Errorf("%d fields allowed", minFields)
		} else {
			return nil, fmt.Errorf("between %d and %d fields allowed", minFields, maxFields)
		}
	}
	return lines[0], nil
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
	return writeID([]string{gi.ResourceName, gi.SchemaName, gi.ObjectName, gi.Privilege, grantOption})
}

// grantIDFromString() takes in a pipe-delimited string: resourceName|schemaName|ObjectName|Privilege
// and returns a grantID object
func grantIDFromString(stringID string) (*grantID, error) {
	row, err := readID(stringID, 4, 5)
	if err != nil {
		return nil, err
	}

	grantOption := false
	if len(row) == 5 && row[4] == "true" {
		grantOption = true
	}

	grantResult := &grantID{
		ResourceName: row[0],
		SchemaName:   row[1],
		ObjectName:   row[2],
		Privilege:    row[3],
		GrantOption:  grantOption,
	}
	return grantResult, nil
}

type pipeID struct {
	Database string
	Schema   string
	Name     string
}

//String() takes in a pipeID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|PipeName
func (si *pipeID) String() (string, error) {
	return writeID([]string{si.Database, si.Schema, si.Name})
}

// pipeIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|PipeName
// and returns a pipeID object
func pipeIDFromString(stringID string) (*pipeID, error) {
	row, err := readID(stringID, 3, 3)
	if err != nil {
		return nil, err
	}

	pipeResult := &pipeID{
		Database: row[0],
		Schema:   row[1],
		Name:     row[2],
	}
	return pipeResult, nil
}

type schemaID struct {
	Database string
	Name     string
}

// String() takes in a schemaID object and returns a pipe-delimited string:
// DatabaseName|schemaName
func (si *schemaID) String() (string, error) {
	return writeID([]string{si.Database, si.Name})
}

// schemaIDFromString() takes in a pipe-delimited string: DatabaseName|schemaName
// and returns a schemaID object
func schemaIDFromString(stringID string) (*schemaID, error) {
	row, err := readID(stringID, 2, 2)
	if err != nil {
		return nil, err
	}

	schemaResult := &schemaID{
		Database: row[0],
		Name:     row[1],
	}
	return schemaResult, nil
}

type stageID struct {
	Database string
	Schema   string
	Name     string
}

// String() takes in a stageID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|StageName
func (si *stageID) String() (string, error) {
	return writeID([]string{si.Database, si.Schema, si.Name})
}

// stageIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|StageName
// and returns a stageID object
func stageIDFromString(stringID string) (*stageID, error) {
	row, err := readID(stringID, 3, 3)
	if err != nil {
		return nil, err
	}

	stageResult := &stageID{
		Database: row[0],
		Schema:   row[1],
		Name:     row[2],
	}
	return stageResult, nil
}

type streamID struct {
	Database string
	Schema   string
	Name     string
}

type streamOnTableID struct {
	DatabaseName string
	SchemaName   string
	OnTableName  string
}

//String() takes in a streamID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|StreamName
func (si *streamID) String() (string, error) {
	return writeID([]string{si.Database, si.Schema, si.Name})
}

// streamIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|StreamName
// and returns a streamID object
func streamIDFromString(stringID string) (*streamID, error) {
	row, err := readID(stringID, 3, 3)
	if err != nil {
		return nil, err
	}

	streamResult := &streamID{
		Database: row[0],
		Schema:   row[1],
		Name:     row[2],
	}
	return streamResult, nil
}

// streamOnTableIDFromString() takes in a dot-delimited string: DatabaseName.SchemaName.TableName
// and returns a streamOnTableID object
func streamOnTableIDFromString(stringID string) (*streamOnTableID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	// TODO switch this to delimter, requires state transition
	reader.Comma = streamOndelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("expecting 1 line")
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
	Database string
	Schema   string
	Name     string
}

//String() takes in a tableID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|TableName
func (si *tableID) String() (string, error) {
	return writeID([]string{si.Database, si.Schema, si.Name})
}

// tableIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TableName
// and returns a tableID object
func tableIDFromString(stringID string) (*tableID, error) {
	row, err := readID(stringID, 3, 3)
	if err != nil {
		return nil, err
	}

	tableResult := &tableID{
		Database: row[0],
		Schema:   row[1],
		Name:     row[2],
	}
	return tableResult, nil
}

type taskID struct {
	Database string
	Schema   string
	Name     string
}

//String() takes in a taskID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|TaskName
func (t *taskID) String() (string, error) {
	return writeID([]string{t.Database, t.Schema, t.Name})
}

// taskIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TaskName
// and returns a taskID object
func taskIDFromString(stringID string) (*taskID, error) {
	row, err := readID(stringID, 3, 3)
	if err != nil {
		return nil, err
	}

	taskResult := &taskID{
		Database: row[0],
		Schema:   row[1],
		Name:     row[2],
	}
	return taskResult, nil
}
