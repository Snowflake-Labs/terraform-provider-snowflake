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
	rowAccessPolicyIDDelimiter = '|'
)

var rowAccessPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the row access policy; must be unique for the database and schema in which the row access policy is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the row access policy.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the row access policy.",
		ForceNew:    true,
	},
	"signature": {
		Type:        schema.TypeMap,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		ForceNew:    true,
		Description: "Specifies signature (arguments) for the row access policy (uppercase and sorted to avoid recreation of resource). A signature specifies a set of attributes that must be considered to determine whether the row is accessible. The attribute values come from the database object (e.g. table or view) to be protected by the row access policy.",
		//Implement DiffSuppressFunc after https://github.com/hashicorp/terraform-plugin-sdk/issues/477 is solved
	},
	"row_access_expression": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the SQL expression. The expression can be any boolean-valued SQL expression.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the row access policy.",
	},
}

type rowAccessPolicyID struct {
	DatabaseName        string
	SchemaName          string
	RowAccessPolicyName string
}

// String() takes in a rowAccessPolicyID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|RowAccessPolicyName
func (rapi *rowAccessPolicyID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = rowAccessPolicyIDDelimiter
	dataIdentifiers := [][]string{{rapi.DatabaseName, rapi.SchemaName, rapi.RowAccessPolicyName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strRowAccessPolicyID := strings.TrimSpace(buf.String())
	return strRowAccessPolicyID, nil
}

/// rowAccessPolicyIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|RowAccessPolicyName
// and returns a rowAccessPolicyID object
func rowAccessPolicyIDFromString(stringID string) (*rowAccessPolicyID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = rowAccessPolicyIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per row access policy")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	rowAccessPolicyResult := &rowAccessPolicyID{
		DatabaseName:        lines[0][0],
		SchemaName:          lines[0][1],
		RowAccessPolicyName: lines[0][2],
	}
	return rowAccessPolicyResult, nil
}

// RowAccessPolicy returns a pointer to the resource representing a row access policy
func RowAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: CreateRowAccessPolicy,
		Read:   ReadRowAccessPolicy,
		Update: UpdateRowAccessPolicy,
		Delete: DeleteRowAccessPolicy,

		Schema: rowAccessPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateRowAccessPolicy implements schema.CreateFunc
func CreateRowAccessPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	signature := d.Get("signature").(map[string]interface{})
	rowAccessExpression := d.Get("row_access_expression").(string)

	builder := snowflake.RowAccessPolicy(name, database, schema)

	builder.WithSignature(signature)
	builder.WithRowAccessExpression(rowAccessExpression)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	stmt := builder.Create()
	err := snowflake.Exec(db, stmt)
	if err != nil {
		return errors.Wrapf(err, "error creating row access policy %v", name)
	}

	rowAccessPolicyID := &rowAccessPolicyID{
		DatabaseName:        database,
		SchemaName:          schema,
		RowAccessPolicyName: name,
	}
	dataIDInput, err := rowAccessPolicyID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadRowAccessPolicy(d, meta)
}

// ReadRowAccessPolicy implements schema.ReadFunc
func ReadRowAccessPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	rowAccessPolicyID, err := rowAccessPolicyIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := rowAccessPolicyID.DatabaseName
	schema := rowAccessPolicyID.SchemaName
	policyName := rowAccessPolicyID.RowAccessPolicyName

	builder := snowflake.RowAccessPolicy(policyName, dbName, schema)

	showSQL := builder.Show()

	row := snowflake.QueryRow(db, showSQL)

	s, err := snowflake.ScanRowAccessPolicies(row)
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

		err = d.Set("row_access_expression", body)
		if err != nil {
			return err
		}

		err = d.Set("signature", ParseSignature(signature))
		if err != nil {
			return err
		}
	}

	return err
}

// UpdateRowAccessPolicy implements schema.UpdateFunc
func UpdateRowAccessPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	rowAccessPolicyID, err := rowAccessPolicyIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := rowAccessPolicyID.DatabaseName
	schema := rowAccessPolicyID.SchemaName
	policyName := rowAccessPolicyID.RowAccessPolicyName

	builder := snowflake.RowAccessPolicy(policyName, dbName, schema)

	if d.HasChange("comment") {
		comment := d.Get("comment")
		if c := comment.(string); c == "" {
			q := builder.RemoveComment()
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error unsetting comment for row access policy on %v", d.Id())
			}
		} else {
			q := builder.ChangeComment(c)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error updating comment for row access policy on %v", d.Id())
			}
		}
	}

	if d.HasChange("row_access_expression") {
		rowAccessExpression := d.Get("row_access_expression")
		q := builder.ChangeRowAccessExpression(rowAccessExpression.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating row access policy expression on %v", d.Id())
		}
	}

	return ReadRowAccessPolicy(d, meta)
}

// DeleteRowAccessPolicy implements schema.DeleteFunc
func DeleteRowAccessPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	rowAccessPolicyID, err := rowAccessPolicyIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := rowAccessPolicyID.DatabaseName
	schema := rowAccessPolicyID.SchemaName
	policyName := rowAccessPolicyID.RowAccessPolicyName

	q := snowflake.RowAccessPolicy(policyName, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting row access policy %v", d.Id())
	}

	d.SetId("")

	return nil
}

func ParseSignature(signature string) map[string]interface{} {
	// Format in database is `(column <data_type>)`
	plainSignature := strings.ReplaceAll(signature, "(", "")
	plainSignature = strings.ReplaceAll(plainSignature, ")", "")
	signatureParts := strings.Split(plainSignature, ", ")
	signatureMap := map[string]interface{}{}

	for _, e := range signatureParts {
		parts := strings.Split(e, " ")
		signatureMap[parts[0]] = parts[1]
	}

	return signatureMap
}
