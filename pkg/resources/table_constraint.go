package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	snowflakeValidation "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var tableConstraintSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of constraint",
	},
	"type": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Type of constraint, one of 'UNIQUE', 'PRIMARY KEY', 'FOREIGN KEY', or 'NOT NULL'",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"UNIQUE", "PRIMARY KEY", "FOREIGN KEY", "NOT NULL"}, false),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(old, new)
		},
		StateFunc: func(val any) string {
			return strings.ToUpper(val.(string))
		},
	},
	"table_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Idenfifier for table to create constraint on. Must be of the form Note: format must follow: \"<db_name>\".\"<schema_name>\".\"<table_name>\" or \"<db_name>.<schema_name>.<table_name>\" or \"<db_name>|<schema_name>.<table_name>\" (snowflake_table.my_table.id)",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"columns": {
		Type:     schema.TypeList,
		MinItems: 1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		ForceNew:    true,
		Required:    true,
		Description: "Columns to use in constraint key",
	},
	"enforced": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Whether the constraint is enforced",
	},
	"deferrable": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     true,
		Description: "Whether the constraint is deferrable",
	},
	"initially": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Default:      "DEFERRED",
		Description:  "Whether the constraint is initially deferred or immediate",
		ValidateFunc: validation.StringInSlice([]string{"DEFERRED", "IMMEDIATE"}, true),
	},
	"enable": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     true,
		Description: "Specifies whether the constraint is enabled or disabled. These properties are provided for compatibility with Oracle.",
	},
	"validate": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Specifies whether to validate existing data on the table when a constraint is created. Only used in conjunction with the ENABLE property.",
	},
	"rely": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     true,
		Description: "Specifies whether a constraint in NOVALIDATE mode is taken into account during query rewrite.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Comment for the table constraint",
	},
	"foreign_key_properties": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Description: "Additional properties when type is set to foreign key. Not applicable for primary/unique keys",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"references": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: "The table and columns that the foreign key references. Not applicable for primary/unique keys",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"table_id": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Name of constraint",
							},
							"columns": {
								Type:     schema.TypeList,
								MinItems: 1,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Required:    true,
								Description: "Columns to use in foreign key reference",
							},
						},
					},
				},
				"match": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "FULL",
					ValidateFunc: validation.StringInSlice([]string{"FULL", "PARTIAL", "SIMPLE"}, true),
					Description:  "The match type for the foreign key. Not applicable for primary/unique keys",
				},
				"on_update": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "NO ACTION",
					ValidateFunc: validation.StringInSlice([]string{"NO ACTION", "CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT"}, true),
					Description:  "Specifies the action performed when the primary/unique key for the foreign key is updated. Not applicable for primary/unique keys",
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return strings.EqualFold(old, new)
					},
					StateFunc: func(val any) string {
						return strings.ToUpper(val.(string))
					},
				},
				"on_delete": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "NO ACTION",
					ValidateFunc: validation.StringInSlice([]string{"NO ACTION", "CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT"}, true),
					Description:  "Specifies the action performed when the primary/unique key for the foreign key is deleted. Not applicable for primary/unique keys",
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return strings.EqualFold(old, new)
					},
					StateFunc: func(val any) string {
						return strings.ToUpper(val.(string))
					},
				},
			},
		},
	},
}

func TableConstraint() *schema.Resource {
	return &schema.Resource{
		Create: CreateTableConstraint,
		Read:   ReadTableConstraint,
		Update: UpdateTableConstraint,
		Delete: DeleteTableConstraint,

		Schema: tableConstraintSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

type tableConstraintID struct {
	name           string
	constraintType string
	tableID        string
}

func (v *tableConstraintID) String() string {
	return fmt.Sprintf("%s❄️%s❄️%s", v.name, v.constraintType, v.tableID)
}

func (v *tableConstraintID) parse(s string) {
	parts := strings.Split(s, "❄️")
	v.name = parts[0]
	v.constraintType = parts[1]
	v.tableID = parts[2]
}

// CreateTableConstraint implements schema.CreateFunc.
func CreateTableConstraint(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)
	constraintType := d.Get("type").(string)
	tableID := d.Get("table_id").(string)

	formattedTableID := snowflakeValidation.ParseAndFormatFullyQualifiedObectID(tableID)
	builder := snowflake.NewTableConstraintBuilder(name, constraintType, formattedTableID)

	cc := d.Get("columns").([]interface{})
	columns := make([]string, 0, len(cc))
	for _, c := range cc {
		columns = append(columns, c.(string))
	}
	builder.WithColumns(columns)

	// set optionals
	if v, ok := d.GetOk("enforced"); ok {
		builder.WithEnforced(v.(bool))
	}

	if v, ok := d.GetOk("deferrable"); ok {
		builder.WithDeferrable(v.(bool))
	}

	if v, ok := d.GetOk("initially"); ok {
		builder.WithInitially(v.(string))
	}

	if v, ok := d.GetOk("enable"); ok {
		builder.WithEnable(v.(bool))
	}

	if v, ok := d.GetOk("validate"); ok {
		builder.WithValidate(v.(bool))
	}

	if v, ok := d.GetOk("rely"); ok {
		builder.WithRely(v.(bool))
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	// set foreign key specific settings
	if v, ok := d.GetOk("foreign_key_properties"); ok {
		foreignKeyProperties := v.([]interface{})[0].(map[string]interface{})
		builder.WithMatch(foreignKeyProperties["match"].(string))
		builder.WithUpdate(foreignKeyProperties["on_update"].(string))
		builder.WithDelete(foreignKeyProperties["on_delete"].(string))
		references := foreignKeyProperties["references"].([]interface{})[0].(map[string]interface{})
		fkTableID := references["table_id"].(string)
		formattedFkTableID := snowflakeValidation.ParseAndFormatFullyQualifiedObectID(fkTableID)
		builder.WithReferenceTableID(formattedFkTableID)
		log.Printf("reference table id : %s", formattedFkTableID)
		cols := references["columns"].([]interface{})
		var fkColumns []string
		for _, c := range cols {
			fkColumns = append(fkColumns, c.(string))
		}
		builder.WithReferenceColumns(fkColumns)
	}

	stmt := builder.Create()
	log.Printf("[DEBUG] create table constraint statement: %v\n", stmt)
	result, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error creating table constraint %v err = %w", name, err)
	}
	log.Printf("[DEBUG] result: %v\n", result)

	tc := tableConstraintID{
		name,
		constraintType,
		tableID,
	}
	d.SetId(tc.String())

	return ReadTableConstraint(d, meta)
}

// ReadTableConstraint implements schema.ReadFunc.
func ReadTableConstraint(_ *schema.ResourceData, _ interface{}) error {
	// commenting this out since it requires an active warehouse to be set which may not be intuitive.
	// also it takes a while for the database to reflect changes. Would likely need to add a validation
	// step like in tag association. People don't like waiting 40 minutes for Terraform to run.

	/*db := meta.(*sql.DB)
	tc := tableConstraintID{}
	tc.parse(d.Id())
	databaseName, schemaName, tableName := snowflakeValidation.ParseFullyQualifiedObjectID(tc.tableID)

	// just need to check to make sure it exists
	_, err := snowflake.ShowTableConstraint(tc.name, databaseName, schemaName, tableName, db)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error reading table constraint %v", tc.String()))
	}*/

	return nil
}

// UpdateTableConstraint implements schema.UpdateFunc.
func UpdateTableConstraint(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tc := tableConstraintID{}
	tc.parse(d.Id())
	formattedTableID := snowflakeValidation.ParseAndFormatFullyQualifiedObectID(tc.tableID)

	builder := snowflake.NewTableConstraintBuilder(tc.name, tc.constraintType, formattedTableID)

	/* "unsupported feature comment error message"
	if d.HasChange("comment") {
		_, new := d.GetChange("comment")
		_, err := db.Exec(builder.SetComment(new.(string)))
		if err != nil {
			return fmt.Errorf("error setting comment for table constraint %v", tc.name)
		}
	}*/

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		_, err := db.Exec(builder.Rename(n.(string)))
		if err != nil {
			return fmt.Errorf("error renaming table constraint %v err = %w", tc.name, err)
		}
	}

	return ReadTableConstraint(d, meta)
}

// DeleteTableConstraint implements schema.DeleteFunc.
func DeleteTableConstraint(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	tc := tableConstraintID{}
	tc.parse(d.Id())
	formattedTableID := snowflakeValidation.ParseAndFormatFullyQualifiedObectID(tc.tableID)
	builder := snowflake.NewTableConstraintBuilder(tc.name, tc.constraintType, formattedTableID)
	cc := d.Get("columns").([]interface{})
	columns := make([]string, 0, len(cc))
	for _, c := range cc {
		columns = append(columns, c.(string))
	}
	builder.WithColumns(columns)

	stmt := builder.Drop()
	_, err := db.Exec(stmt)
	if err != nil {
		// if the table constraint does not exist, then remove from state file
		if strings.Contains(err.Error(), "does not exist") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error deleting table constraint %v err = %w", tc.name, err)
	}

	d.SetId("")
	return nil
}
