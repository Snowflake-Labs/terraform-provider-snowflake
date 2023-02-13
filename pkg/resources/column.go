package resources

import (
	"database/sql"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
	"strings"
)

var columnParameterSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Column name",
	},
	"table_identifier": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Description: "Specifies the table identifier for the column.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Name of the table to set the parameter for.",
				},
				"database": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the database that the table was created in.",
				},
				"schema": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the schema that the table was created in.",
				},
			},
		},
	},
	"data_type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Column type, e.g. VARIANT",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// these are all equivalent as per https://docs.snowflake.com/en/sql-reference/data-types-text.html
			varcharType := []string{"VARCHAR(16777216)", "VARCHAR", "text", "string", "NVARCHAR", "NVARCHAR2", "CHAR VARYING", "NCHAR VARYING"}
			return slices.Contains(varcharType, new) && slices.Contains(varcharType, old)
		},
	},
	"nullable": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Whether this column can contain null values. **Note**: Depending on your Snowflake version, the default value will not suffice if this column is used in a primary key constraint.",
	},
	"default": {
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "Defines the column default value; note due to limitations of Snowflake's ALTER TABLE ADD/MODIFY COLUMN updates to default will not be applied",
		MinItems:      1,
		MaxItems:      1,
		ConflictsWith: []string{".identity"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"constant": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The default constant value for the column",
					// ConflictsWith: []string{".expression", ".sequence"}, - can't use, nor ExactlyOneOf due to column type being TypeList
				},
				"expression": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The default expression value for the column",
					// ConflictsWith: []string{".constant", ".sequence"}, - can't use, nor ExactlyOneOf due to column type being TypeList
				},
				"sequence_identifier": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The default sequence to use for the column",
					// ConflictsWith: []string{".constant", ".expression"}, - can't use, nor ExactlyOneOf due to column type being TypeList
				},
			},
		},
	},
	/*Note: Identity and default are mutually exclusive. From what I can tell we can't enforce this here
	the snowflake query will error so we can defer enforcement to there.
	*/
	"identity": {
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "Defines the identity start/step values for a column. **Note** Identity/default are mutually exclusive.",
		MinItems:      1,
		MaxItems:      1,
		ConflictsWith: []string{".default"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"start_num": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "The number to start incrementing at.",
					Default:     1,
				},
				"step_num": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Step size to increment by.",
					Default:     1,
				},
			},
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Column comment",
	},
	//"masking_policy": {
	//	Type:        schema.TypeString,
	//	Optional:    true,
	//	Default:     "",
	//	Description: "Masking policy to apply on column",
	//},
}

func ColumnParameter() *schema.Resource {
	return &schema.Resource{
		Create: CreateColumnParameter,
		Read:   ReadColumnParameter,
		Update: UpdateColumnParameter,
		Delete: DeleteColumnParameter,

		Schema: columnParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateColumnParameter implements schema.CreateFunc.
func CreateColumnParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	tableObjectID := d.Get("table_identifier").(snowflake.ObjectIdentifier)
	dataType := d.Get("data_type").(snowflake.ColumnDataType)

	nullable := d.Get("nullable").(bool)
	defaultConstant, defaultConstantBool := d.GetOk("default.constant")
	defaultExpression, defaultExpressionBool := d.GetOk("default.expression")
	defaultIdentifier, defaultIdentifierBool := d.GetOk("default.sequence_identifier")
	defaultBool := defaultConstantBool || defaultExpressionBool || defaultIdentifierBool

	identity, identityBool := d.GetOk("identity")

	comment := d.Get("comment").(string)

	q := snowflake.NewCreateTableColumnBuilder(name, &tableObjectID, &dataType, db)

	q.WithNullable(nullable)

	if defaultBool {
		var default_ snowflake.ColumnDefault
		if defaultConstantBool {
			default_ = defaultConstant.(snowflake.ColumnDefault)
		} else if defaultExpressionBool {
			default_ = defaultExpression.(snowflake.ColumnDefault)
		} else if defaultIdentifierBool {
			default_ = defaultIdentifier.(snowflake.ColumnDefault)
		}
		q.WithDefault(&default_)
	}

	if identityBool {
		ident := identity.(snowflake.Identity)
		q.WithIdentity(&ident)
	}

	q.WithComment(comment)

	err := q.Create()
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s.%s.%s.%s", tableObjectID.Database, tableObjectID.Schema, tableObjectID.Name, name))
	return nil
}

// ReadColumnParameter implements schema.ReadFunc.
func ReadColumnParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	tableObjectID := d.Get("table_identifier").(snowflake.ObjectIdentifier)
	var err error

	q := snowflake.NewCreateTableColumnBuilder(name, &tableObjectID, nil, db)

	show, err := q.Show()
	if err != nil {
		return err
	}

	if err = d.Set("data_type", show.DataType); err != nil {
		return err
	}

	if err = d.Set("nullable", show.Null); err != nil {
		return err
	}

	// TODO default stuff

	if err = d.Set("comment", show.Comment); err != nil {
		return err
	}

	autoIncrementNullString := show.Autoincrement

	if autoIncrementNullString.Valid {
		autoIncrementString := autoIncrementNullString.String
		split := strings.Split(autoIncrementString, " INCREMENT ")
		splitLeft := split[0]
		stepString := split[1]
		startString := strings.Split(splitLeft, "IDENTITY START ")[0]

		if err = d.Set("identity.start_num", startString); err != nil {
			return err
		}
		if err = d.Set("identity.step_num", stepString); err != nil {
			return err
		}
	} else {
		if err = d.Set("identity.start_num", nil); err != nil {
			return err
		}
		if err = d.Set("identity.step_num", nil); err != nil {
			return err
		}
	}

	return nil
}

// UpdateColumnParameter implements schema.UpdateFunc.
func UpdateColumnParameter(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// DeleteColumnParameter implements schema.DeleteFunc.
func DeleteColumnParameter(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	tableObjectID := d.Get("table_identifier").(snowflake.ObjectIdentifier)

	q := snowflake.NewCreateTableColumnBuilder(name, &tableObjectID, nil, db)
	if err := q.Drop(); err != nil {
		return fmt.Errorf("error deleting Table column %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}
