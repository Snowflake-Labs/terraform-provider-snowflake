package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var sequenceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the name for the sequence.",
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
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the sequence. Don't use the | character.",
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

	return ReadSequence(d, meta)
}

// ReadSequence implements schema.ReadFunc
func ReadSequence(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)

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

	i, err = strconv.ParseInt(sequence.NextValue.String, 10, 64)
	if err != nil {
		return err
	}

	err = d.Set("next_value", i)
	if err != nil {
		return err
	}

	err = d.Set("fully_qualified_name", seq.Address())
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf(`%v|%v|%v`, sequence.DBName.String, sequence.SchemaName.String, sequence.Name.String))
	if err != nil {
		return err
	}

	return err
}

func UpdateSequence(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)
	next := d.Get("next_value").(int)

	DeleteSequence(d, meta)

	sq := snowflake.Sequence(name, database, schema)

	if i, ok := d.GetOk("increment"); ok {
		sq.WithIncrement(i.(int))
	}

	if v, ok := d.GetOk("comment"); ok {
		sq.WithComment(v.(string))
	}

	sq.WithStart(next)

	err := snowflake.Exec(db, sq.Create())
	if err != nil {
		return errors.Wrapf(err, "error creating sequence")
	}

	return ReadSequence(d, meta)
}

func DeleteSequence(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)

	stmt := snowflake.Sequence(name, database, schema).Drop()

	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error dropping sequence %s", name)
	}

	d.SetId("")
	return nil
}
