package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

const (
	maskingPolicyIDDelimiter = '|'
)

var maskingPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the masking policy; must be unique for the database and schema in which the masking policy is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the masking policy.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the masking policy.",
		ForceNew:    true,
	},
	"value_data_type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the data type to mask.",
		ForceNew:    true,
	},
	"masking_expression": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the SQL expression that transforms the data.",
	},
	"return_data_type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the data type to return.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the masking policy.",
	},
}

type maskingPolicyID struct {
	DatabaseName      string
	SchemaName        string
	MaskingPolicyName string
}

// String() takes in a maskingPolicyID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|MaskingPolicyName
func (mpi *maskingPolicyID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = maskingPolicyIDDelimiter
	dataIdentifiers := [][]string{{mpi.DatabaseName, mpi.SchemaName, mpi.MaskingPolicyName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strMaskingPolicyID := strings.TrimSpace(buf.String())
	return strMaskingPolicyID, nil
}

/// maskingPolicyIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|MaskingPolicyName
// and returns a maskingPolicyID object
func maskingPolicyIDFromString(stringID string) (*maskingPolicyID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = maskingPolicyIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per masking policy")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	maskingPolicyResult := &maskingPolicyID{
		DatabaseName:      lines[0][0],
		SchemaName:        lines[0][1],
		MaskingPolicyName: lines[0][2],
	}
	return maskingPolicyResult, nil
}

// MaskingPolicy returns a pointer to the resource representing a masking policy
func MaskingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: CreateMaskingPolicy,
		Read:   ReadMaskingPolicy,
		Update: UpdateMaskingPolicy,
		Delete: DeleteMaskingPolicy,

		Schema: maskingPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateMaskingPolicy implements schema.CreateFunc
func CreateMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	valueDataType := d.Get("value_data_type").(string)
	maskingExpression := d.Get("masking_expression").(string)
	returnDataType := d.Get("return_data_type").(string)

	builder := snowflake.MaskingPolicy(name, database, schema)

	builder.WithValueDataType(valueDataType)
	builder.WithMaskingExpression(maskingExpression)
	builder.WithReturnDataType(returnDataType)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	stmt := builder.Create()
	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating masking policy %v", name)
	}

	maskingPolicyID := &maskingPolicyID{
		DatabaseName:      database,
		SchemaName:        schema,
		MaskingPolicyName: name,
	}
	dataIDInput, err := maskingPolicyID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadMaskingPolicy(d, meta)
}

// ReadMaskingPolicy implements schema.ReadFunc
func ReadMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	maskingPolicyID, err := maskingPolicyIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := maskingPolicyID.DatabaseName
	schema := maskingPolicyID.SchemaName
	policyName := maskingPolicyID.MaskingPolicyName

	builder := snowflake.MaskingPolicy(policyName, dbName, schema)

	showSQL := builder.Show()

	row := snowflake.QueryRow(db, showSQL)

	s, err := snowflake.ScanMaskingPolicies(row)
	if err != nil {
		return err
	}

	err = d.Set("name", s.Name.String)
	if err != nil {
		return err
	}

	err = d.Set("database", s.DatabaseName.String)
	if err != nil {
		return err
	}

	err = d.Set("schema", s.SchemaName.String)
	if err != nil {
		return err
	}

	err = d.Set("comment", s.Comment.String)
	if err != nil {
		return err
	}

	descSQL := builder.Describe()
	rows, err := snowflake.Query(db, descSQL)
	if err != nil {
		return err
	}

	var (
		name       string
		signature  string
		returnType string
		body       string
	)
	for rows.Next() {
		err := rows.Scan(&name, &signature, &returnType, &body)
		if err != nil {
			return err
		}

		err = d.Set("masking_expression", body)
		if err != nil {
			return err
		}

		err = d.Set("return_data_type", returnType)
		if err != nil {
			return err
		}

		// format in database is `(VAL <data_type>)`
		valueDataType := strings.TrimSuffix(strings.Split(signature, " ")[1], ")")
		err = d.Set("value_data_type", valueDataType)
		if err != nil {
			return err
		}
	}

	return err
}

// UpdateMaskingPolicy implements schema.UpdateFunc
func UpdateMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	maskingPolicyID, err := maskingPolicyIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := maskingPolicyID.DatabaseName
	schema := maskingPolicyID.SchemaName
	policyName := maskingPolicyID.MaskingPolicyName

	builder := snowflake.MaskingPolicy(policyName, dbName, schema)

	if d.HasChange("comment") {
		comment := d.Get("comment")
		if c := comment.(string); c == "" {
			q := builder.RemoveComment()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting comment for masking policy on %v", d.Id())
			}
		} else {
			q := builder.ChangeComment(c)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error updating comment for masking policy on %v", d.Id())
			}
		}
	}

	if d.HasChange("masking_expression") {
		maskingExpression := d.Get("masking_expression")
		q := builder.ChangeMaskingExpression(maskingExpression.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating masking policy expression on %v", d.Id())
		}
	}

	return ReadMaskingPolicy(d, meta)
}

// DeleteMaskingPolicy implements schema.DeleteFunc
func DeleteMaskingPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	maskingPolicyID, err := maskingPolicyIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := maskingPolicyID.DatabaseName
	schema := maskingPolicyID.SchemaName
	policyName := maskingPolicyID.MaskingPolicyName

	q := snowflake.MaskingPolicy(policyName, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting masking policy %v", d.Id())
	}

	d.SetId("")

	return nil
}
