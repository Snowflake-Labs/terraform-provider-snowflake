package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
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
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// these are all equivalent as per https://docs.snowflake.com/en/sql-reference/data-types-text.html
			varcharType := []string{"VARCHAR(16777216)", "VARCHAR", "text", "string", "NVARCHAR", "NVARCHAR2", "CHAR VARYING", "NCHAR VARYING"}
			return slices.Contains(varcharType, new) && slices.Contains(varcharType, old)
		},
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
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// these are all equivalent as per https://docs.snowflake.com/en/sql-reference/data-types-text.html
			varcharType := []string{"VARCHAR(16777216)", "VARCHAR", "text", "string", "NVARCHAR", "NVARCHAR2", "CHAR VARYING", "NCHAR VARYING"}
			return slices.Contains(varcharType, new) && slices.Contains(varcharType, old)
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the masking policy.",
	},
	"qualified_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the qualified identifier for the masking policy.",
	},
}

type maskingPolicyID struct {
	DatabaseName      string
	SchemaName        string
	MaskingPolicyName string
}

// String() takes in a maskingPolicyID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|MaskingPolicyName.
func (mpi *maskingPolicyID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = maskingPolicyIDDelimiter
	dataIdentifiers := [][]string{{mpi.DatabaseName, mpi.SchemaName, mpi.MaskingPolicyName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strMaskingPolicyID := strings.TrimSpace(buf.String())
	return strMaskingPolicyID, nil
}

// / maskingPolicyIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|MaskingPolicyName
// and returns a maskingPolicyID object.
func maskingPolicyIDFromString(stringID string) (*maskingPolicyID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = maskingPolicyIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
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

// MaskingPolicy returns a pointer to the resource representing a masking policy.
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

// CreateMaskingPolicy implements schema.CreateFunc.
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
	if err := snowflake.Exec(db, stmt); err != nil {
		return fmt.Errorf("error creating masking policy %v err = %w", name, err)
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

// ReadMaskingPolicy implements schema.ReadFunc.
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

	if err := d.Set("name", s.Name.String); err != nil {
		return err
	}

	if err := d.Set("database", s.DatabaseName.String); err != nil {
		return err
	}

	if err := d.Set("schema", s.SchemaName.String); err != nil {
		return err
	}

	if err := d.Set("comment", s.Comment.String); err != nil {
		return err
	}

	if err := d.Set("qualified_name", builder.QualifiedName()); err != nil {
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
		if err := rows.Scan(&name, &signature, &returnType, &body); err != nil {
			return err
		}

		if err := d.Set("masking_expression", body); err != nil {
			return err
		}

		if err := d.Set("return_data_type", returnType); err != nil {
			return err
		}

		// format in database is `(VAL <data_type>)`
		valueDataType := strings.TrimSuffix(strings.Split(signature, " ")[1], ")")
		if err := d.Set("value_data_type", valueDataType); err != nil {
			return err
		}
	}

	return err
}

// UpdateMaskingPolicy implements schema.UpdateFunc.
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
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error unsetting comment for masking policy on %v err = %w", d.Id(), err)
			}
		} else {
			q := builder.ChangeComment(c)
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error updating comment for masking policy on %v err = %w", d.Id(), err)
			}
		}
	}

	if d.HasChange("masking_expression") {
		maskingExpression := d.Get("masking_expression")
		q := builder.ChangeMaskingExpression(maskingExpression.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating masking policy expression on %v err = %w", d.Id(), err)
		}
	}

	return ReadMaskingPolicy(d, meta)
}

// DeleteMaskingPolicy implements schema.DeleteFunc.
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
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting masking policy %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
