package resources

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/jszwec/csvutil"
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

func readID(id string, data interface{}) error {
	reader := csv.NewReader(strings.NewReader(id))
	reader.Comma = delimiter

	header, err := csvutil.Header(data, "csv")
	if err != nil {
		return err
	}
	fmt.Printf("[DEBUG] header %#v\n", header)
	decoder, err := csvutil.NewDecoder(reader, header...)
	if err != nil {
		return err
	}

	return decoder.Decode(&data)
}

// grantID contains identifying elements that allow unique access privileges
type grantID struct {
	ResourceName string
	SchemaName   string
	ObjectName   string
	Privilege    string
	GrantOption  bool `csv:",omitempty"`
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
	result := &grantID{}
	err := readID(stringID, result)
	return result, err
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
	result := &schemaID{}
	err := readID(stringID, result)
	return result, err
}

type streamOnTableID struct {
	DatabaseName string
	SchemaName   string
	OnTableName  string
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

type schemaScopedID struct {
	Database string
	Schema   string
	Name     string
}

//String() takes in a schemaScopedID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|TaskName
func (t *schemaScopedID) String() (string, error) {
	return writeID([]string{t.Database, t.Schema, t.Name})
}

// taskIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TaskName
// and returns a taskID object
func idFromString(stringID string) (*schemaScopedID, error) {
	result := &schemaScopedID{}
	err := readID(stringID, result)
	return result, err
}
