package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

const (
	sequenceIDDelimiter = '|'
)

var sequenceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the name for the sequence.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies a comment for the sequence.",
	},
	"increment": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     1,
		Description: "The amount the sequence will increase by each time it is used",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the sequence. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the sequence. Don't use the | character.",
		ForceNew:    true,
	},
	"next_value": {
		Type:        schema.TypeInt,
		Description: "The next value the sequence will provide.",
		Computed:    true,
	},
	"fully_qualified_name": {
		Type:        schema.TypeString,
		Description: "The fully qualified name of the sequence.",
		Computed:    true,
	},
}

var sequenceProperties = []string{"comment", "data_retention_time_in_days"}

type sequenceID struct {
	DatabaseName string
	SchemaName   string
	SequenceName string
}

//String() takes in a sequenceID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|SequenceName
func (si *sequenceID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = pipeIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.SequenceName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strSequenceID := strings.TrimSpace(buf.String())
	return strSequenceID, nil
}

// Sequence returns a pointer to the resource representing a sequence
func Sequence() *schema.Resource {
	return &schema.Resource{
		Create: CreateSequence,
		Read:   ReadSequence,
		Delete: DeleteSequence,
		Update: UpdateSequence,

		Schema: sequenceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSequence implements schema.CreateFunc
func CreateSequence(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)

	sq := snowflake.Sequence(name, database, schema)

	if i, ok := d.GetOk("increment"); ok {
		sq.WithIncrement(i.(int))
	}

	if v, ok := d.GetOk("comment"); ok {
		sq.WithComment(v.(string))
	}

	err := snowflake.Exec(db, sq.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating sequence")
	}

	sequenceID := &sequenceID{
		DatabaseName: database,
		SchemaName:   schema,
		SequenceName: name,
	}

	dataIDInput, err := sequenceID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadSequence(d, meta)
}

// ReadSequence implements schema.ReadFunc
func ReadSequence(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	sequenceID, err := sequenceIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := sequenceID.DatabaseName
	schema := sequenceID.SchemaName
	name := sequenceID.SequenceName

	seq := snowflake.Sequence(name, database, schema)
	stmt := seq.Show()
	row := snowflake.QueryRow(db, stmt)

	sequence, err := snowflake.ScanSequence(row)

	if err != nil {
		if err == sql.ErrNoRows {
			// If not found, mark resource to be removed from statefile during apply or refresh
			log.Printf("[DEBUG] sequence (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return errors.Wrap(err, "unable to scan row for SHOW SEQUENCES")
	}

	err = d.Set("name", sequence.Name.String)
	if err != nil {
		return err
	}

	err = d.Set("schema", sequence.SchemaName.String)
	if err != nil {
		return err
	}

	err = d.Set("database", sequence.DBName.String)
	if err != nil {
		return err
	}

	err = d.Set("comment", sequence.Comment.String)
	if err != nil {
		return err
	}

	i, err := strconv.ParseInt(sequence.Increment.String, 10, 64)
	if err != nil {
		return err
	}

	err = d.Set("increment", i)
	if err != nil {
		return err
	}

	n, err := strconv.ParseInt(sequence.NextValue.String, 10, 64)
	if err != nil {
		return err
	}

	err = d.Set("next_value", n)
	if err != nil {
		return err
	}

	err = d.Set("fully_qualified_name", seq.Address())
	if err != nil {
		return err
	}

	return nil
}

func UpdateSequence(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	sequenceID, err := sequenceIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := sequenceID.DatabaseName
	schema := sequenceID.SchemaName
	name := sequenceID.SequenceName

	sq := snowflake.Sequence(name, database, schema)
	stmt := sq.Show()
	row := snowflake.QueryRow(db, stmt)

	sequence, err := snowflake.ScanSequence(row)
	DeleteSequence(d, meta)

	if i, ok := d.GetOk("increment"); ok {
		sq.WithIncrement(i.(int))
	}

	if v, ok := d.GetOk("comment"); ok {
		sq.WithComment(v.(string))
	}

	nextValue, err := strconv.Atoi(sequence.NextValue.String)
	if err != nil {
		return err
	}

	err = d.Set("next_value", nextValue)
	if err != nil {
		return err
	}

	sq.WithStart(nextValue)

	err = snowflake.Exec(db, sq.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating sequence")
	}

	return ReadSequence(d, meta)
}

func DeleteSequence(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	sequenceID, err := sequenceIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := sequenceID.DatabaseName
	schema := sequenceID.SchemaName
	name := sequenceID.SequenceName

	stmt := snowflake.Sequence(name, database, schema).Drop()

	err = snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error dropping sequence %s", name)
	}

	d.SetId("")
	return nil
}

// sequenceIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|PipeName
// and returns a sequenceID object
func sequenceIDFromString(stringID string) (*sequenceID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = sequenceIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per sequence")
	}

	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	sequenceResult := &sequenceID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		SequenceName: lines[0][2],
	}

	return sequenceResult, nil
}
