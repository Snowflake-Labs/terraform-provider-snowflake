package resources

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	tableIDDelimiter = '|'
)

var tableSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the table; must be unique for the database and schema in which the table is created.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the table.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the table.",
	},
	"cluster_by": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "A list of one or more table columns/expressions to be used as clustering key(s) for the table",
	},

	"column": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Description: "Definitions of a column to create in the table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column name",
				},
				"type": {
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
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Defines the column default value; note due to limitations of Snowflake's ALTER TABLE ADD/MODIFY COLUMN updates to default will not be applied",
					MinItems:    1,
					MaxItems:    1,
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
							"sequence": {
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
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Defines the identity start/step values for a column. **Note** Identity/default are mutually exclusive.",
					MinItems:    1,
					MaxItems:    1,
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
				"masking_policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "Masking policy to apply on column. It has to be a fully qualified name.",
				},
			},
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the table.",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the role that owns the table.",
	},
	"primary_key": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Definitions of primary key constraint to create on table",
		Deprecated:  "Use snowflake_table_constraint instead",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Name of constraint",
				},
				"keys": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Required:    true,
					Description: "Columns to use in primary key",
				},
			},
		},
	},
	"data_retention_days": {
		Type:          schema.TypeInt,
		Optional:      true,
		Description:   "Specifies the retention period for the table so that Time Travel actions (SELECT, CLONE, UNDROP) can be performed on historical data in the table. Default value is 1, if you wish to inherit the parent schema setting then pass in the schema attribute to this argument.",
		ValidateFunc:  validation.IntBetween(0, 90),
		Deprecated:    "Use data_retention_time_in_days attribute instead",
		ConflictsWith: []string{"data_retention_time_in_days"},
	},
	"data_retention_time_in_days": {
		Type:          schema.TypeInt,
		Optional:      true,
		Description:   "Specifies the retention period for the table so that Time Travel actions (SELECT, CLONE, UNDROP) can be performed on historical data in the table. Default value is 1, if you wish to inherit the parent schema setting then pass in the schema attribute to this argument.",
		ValidateFunc:  validation.IntBetween(0, 90),
		Deprecated:    "Use snowflake_object_parameter instead",
		ConflictsWith: []string{"data_retention_days"},
	},
	"change_tracking": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies whether to enable change tracking on the table. Default false.",
	},
	"qualified_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Qualified name of the table.",
	},
	"tag": tagReferenceSchema,
}

func Table() *schema.Resource {
	return &schema.Resource{
		Create: CreateTable,
		Read:   ReadTable,
		Update: UpdateTable,
		Delete: DeleteTable,

		Schema: tableSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

type tableID struct {
	DatabaseName string
	SchemaName   string
	TableName    string
}

// String() takes in a tableID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|TableName.
func (si *tableID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = tableIDDelimiter
	dataIdentifiers := [][]string{{si.DatabaseName, si.SchemaName, si.TableName}}
	if err := csvWriter.WriteAll(dataIdentifiers); err != nil {
		return "", err
	}
	strTableID := strings.TrimSpace(buf.String())
	return strTableID, nil
}

// tableIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TableName
// and returns a tableID object.
func tableIDFromString(stringID string) (*tableID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = tableIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
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

type columnDefault struct {
	constant   *string
	expression *string
	sequence   *string
}

func (cd *columnDefault) toSnowflakeColumnDefault() *snowflake.ColumnDefault {
	if cd.constant != nil {
		return snowflake.NewColumnDefaultWithConstant(*cd.constant)
	}

	if cd.expression != nil {
		return snowflake.NewColumnDefaultWithExpression(*cd.expression)
	}

	if cd.sequence != nil {
		return snowflake.NewColumnDefaultWithSequence(*cd.sequence)
	}

	return nil
}

func (cd *columnDefault) _type() string {
	if cd.constant != nil {
		return "constant"
	}

	if cd.expression != nil {
		return "expression"
	}

	if cd.sequence != nil {
		return "sequence"
	}

	return "unknown"
}

type columnIdentity struct {
	startNum int
	stepNum  int
}

func (identity *columnIdentity) toSnowflakeColumnIdentity() *snowflake.ColumnIdentity {
	snowIdentity := snowflake.ColumnIdentity{}
	return snowIdentity.WithStartNum(identity.startNum).WithStep(identity.stepNum)
}

type column struct {
	name          string
	dataType      string
	nullable      bool
	_default      *columnDefault
	identity      *columnIdentity
	comment       string
	maskingPolicy string
}

func (c column) toSnowflakeColumn() snowflake.Column {
	sC := &snowflake.Column{}

	if c._default != nil {
		sC = sC.WithDefault(c._default.toSnowflakeColumnDefault())
	}

	if c.identity != nil {
		sC = sC.WithIdentity(c.identity.toSnowflakeColumnIdentity())
	}

	return *sC.WithName(c.name).
		WithType(c.dataType).
		WithNullable(c.nullable).
		WithComment(c.comment).
		WithMaskingPolicy(c.maskingPolicy)
}

type columns []column

func (c columns) toSnowflakeColumns() []snowflake.Column {
	sC := make([]snowflake.Column, len(c))
	for i, col := range c {
		sC[i] = col.toSnowflakeColumn()
	}
	return sC
}

type changedColumns []changedColumn

type changedColumn struct {
	newColumn             column // our new column
	changedDataType       bool
	changedNullConstraint bool
	dropedDefault         bool
	changedComment        bool
	changedMaskingPolicy  bool
}

func (c columns) getChangedColumnProperties(new columns) (changed changedColumns) {
	changed = changedColumns{}
	for _, cO := range c {
		for _, cN := range new {
			changeColumn := changedColumn{cN, false, false, false, false, false}
			if cO.name == cN.name && cO.dataType != cN.dataType {
				changeColumn.changedDataType = true
			}
			if cO.name == cN.name && cO.nullable != cN.nullable {
				changeColumn.changedNullConstraint = true
			}
			if cO.name == cN.name && cO._default != nil && cN._default == nil {
				changeColumn.dropedDefault = true
			}

			if cO.name == cN.name && cO.comment != cN.comment {
				changeColumn.changedComment = true
			}

			if cO.name == cN.name && cO.maskingPolicy != cN.maskingPolicy {
				changeColumn.changedMaskingPolicy = true
			}

			changed = append(changed, changeColumn)
		}
	}
	return
}

func (c columns) diffs(new columns) (removed columns, added columns, changed changedColumns) {
	return c.getNewIn(new), new.getNewIn(c), c.getChangedColumnProperties(new)
}

func getColumnDefault(def map[string]interface{}) *columnDefault {
	if c, ok := def["constant"]; ok {
		if constant, ok := c.(string); ok && len(constant) > 0 {
			return &columnDefault{
				constant: &constant,
			}
		}
	}

	if e, ok := def["expression"]; ok {
		if expr, ok := e.(string); ok && len(expr) > 0 {
			return &columnDefault{
				expression: &expr,
			}
		}
	}

	if s, ok := def["sequence"]; ok {
		if seq := s.(string); ok && len(seq) > 0 {
			return &columnDefault{
				sequence: &seq,
			}
		}
	}

	return nil
}

func getColumnIdentity(identity map[string]interface{}) *columnIdentity {
	if len(identity) > 0 {
		startNum := identity["start_num"].(int)
		stepNum := identity["step_num"].(int)
		return &columnIdentity{startNum, stepNum}
	}

	return nil
}

func getColumn(from interface{}) (to column) {
	c := from.(map[string]interface{})
	var cd *columnDefault
	var id *columnIdentity

	_default := c["default"].([]interface{})
	identity := c["identity"].([]interface{})

	if len(_default) == 1 {
		cd = getColumnDefault(_default[0].(map[string]interface{}))
	}
	if len(identity) == 1 {
		id = getColumnIdentity(identity[0].(map[string]interface{}))
	}

	return column{
		name:          c["name"].(string),
		dataType:      c["type"].(string),
		nullable:      c["nullable"].(bool),
		_default:      cd,
		identity:      id,
		comment:       c["comment"].(string),
		maskingPolicy: c["masking_policy"].(string),
	}
}

func getColumns(from interface{}) (to columns) {
	cols := from.([]interface{})
	to = make(columns, len(cols))
	for i, c := range cols {
		to[i] = getColumn(c)
	}
	return to
}

func getTableColumnRequest(from interface{}) *sdk.TableColumnRequest {
	c := from.(map[string]interface{})
	_type := c["type"].(string)

	// TODO [SNOW-884959]: old implementation was quoting the column names - should we leave it?
	nameInQuotes := fmt.Sprintf(`"%v"`, snowflake.EscapeString(c["name"].(string)))
	request := sdk.NewTableColumnRequest(nameInQuotes, sdk.DataType(_type))

	// TODO [SNOW-884959]: move each default possibility logic to request builder/SDK
	_default := c["default"].([]interface{})
	var expression string
	if len(_default) == 1 {
		if c, ok := _default[0].(map[string]interface{})["constant"]; ok {
			if constant, ok := c.(string); ok && len(constant) > 0 {
				if strings.Contains(_type, "CHAR") || _type == "STRING" || _type == "TEXT" {
					expression = snowflake.EscapeSnowflakeString(constant)
				} else {
					expression = constant
				}

			}
		}

		if e, ok := _default[0].(map[string]interface{})["expression"]; ok {
			if expr, ok := e.(string); ok && len(expr) > 0 {
				expression = expr
			}
		}

		if s, ok := _default[0].(map[string]interface{})["sequence"]; ok {
			if seq := s.(string); ok && len(seq) > 0 {
				expression = fmt.Sprintf(`%v.NEXTVAL`, seq)
			}
		}
		request.WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithExpression(sdk.String(expression)))
	}

	identity := c["identity"].([]interface{})
	if len(identity) == 1 {
		identityProp := identity[0].(map[string]interface{})
		startNum := identityProp["start_num"].(int)
		stepNum := identityProp["step_num"].(int)
		request.WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(startNum, stepNum)))
	}

	maskingPolicy := c["masking_policy"].(string)
	if maskingPolicy != "" {
		request.WithMaskingPolicy(sdk.NewColumnMaskingPolicyRequest(sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(maskingPolicy)))
	}

	return request.
		WithNotNull(sdk.Bool(!c["nullable"].(bool))).
		WithComment(sdk.String(c["comment"].(string)))
}

func getTableColumnRequests(from interface{}) []sdk.TableColumnRequest {
	cols := from.([]interface{})
	to := make([]sdk.TableColumnRequest, len(cols))
	for i, c := range cols {
		to[i] = *getTableColumnRequest(c)
	}
	return to
}

type primarykey struct {
	name string
	keys []string
}

func getPrimaryKey(from interface{}) (to primarykey) {
	pk := from.([]interface{})
	to = primarykey{}
	if len(pk) > 0 {
		pkDetails := pk[0].(map[string]interface{})
		to.name = pkDetails["name"].(string)
		to.keys = expandStringList(pkDetails["keys"].([]interface{}))
		return to
	}
	return to
}

func (pk primarykey) toSnowflakePrimaryKey() snowflake.PrimaryKey {
	snowPk := snowflake.PrimaryKey{}
	return *snowPk.WithName(pk.name).WithKeys(pk.keys)
}

func toColumnConfig(descriptions []sdk.TableColumnDetails) []any {
	var flattened []any
	for _, td := range descriptions {
		if td.Kind != "COLUMN" {
			continue
		}

		flat := map[string]any{}
		flat["name"] = td.Name
		flat["type"] = string(td.Type)
		flat["nullable"] = td.IsNullable

		if td.Comment != nil {
			flat["comment"] = *td.Comment
		}

		if td.PolicyName != nil {
			flat["masking_policy"] = *td.PolicyName
		}

		def := toColumnDefaultConfig(td)
		if def != nil {
			flat["default"] = []any{def}
		}

		identity := toColumnIdentityConfig(td)
		if identity != nil {
			flat["identity"] = []any{identity}
		}

		flattened = append(flattened, flat)
	}
	return flattened
}

func toColumnDefaultConfig(td sdk.TableColumnDetails) map[string]any {
	if td.Default == nil {
		return nil
	}

	defaultRaw := *td.Default
	def := map[string]any{}
	if strings.HasSuffix(defaultRaw, ".NEXTVAL") {
		def["sequence"] = strings.TrimSuffix(defaultRaw, ".NEXTVAL")
		return def
	}

	if strings.Contains(defaultRaw, "(") && strings.Contains(defaultRaw, ")") {
		def["expression"] = defaultRaw
		return def
	}

	columnType := strings.ToUpper(string(td.Type))
	if strings.Contains(columnType, "CHAR") || columnType == "STRING" || columnType == "TEXT" {
		def["constant"] = snowflake.UnescapeSnowflakeString(defaultRaw)
		return def
	}

	if toColumnIdentityConfig(td) != nil {
		/*
			Identity/autoincrement information is stored in the same column as default information. We want to handle the identity separate so will return nil
			here if identity information is present. Default/identity are mutually exclusive
		*/
		return nil
	}

	def["constant"] = defaultRaw
	return def
}

func toColumnIdentityConfig(td sdk.TableColumnDetails) map[string]any {
	// if autoincrement is used this is reflected back IDENTITY START 1 INCREMENT 1
	if td.Default == nil {
		return nil
	}

	defaultRaw := *td.Default

	if strings.Contains(defaultRaw, "IDENTITY") {
		identity := map[string]any{}

		split := strings.Split(defaultRaw, " ")
		start, _ := strconv.Atoi(split[2])
		step, _ := strconv.Atoi(split[4])

		identity["start_num"] = start
		identity["step_num"] = step

		return identity
	}
	return nil
}

// CreateTable implements schema.CreateFunc.
func CreateTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	tableColumnRequests := getTableColumnRequests(d.Get("column").([]interface{}))

	createRequest := sdk.NewCreateTableRequest(id, tableColumnRequests)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	if v, ok := d.GetOk("cluster_by"); ok {
		// TODO [SNOW-884959]: in old implementation LINEAR was wrapping the list. Is it needed?
		createRequest.WithClusterBy(expandStringList(v.([]interface{})))
	}

	if v, ok := d.GetOk("primary_key"); ok {
		pk := getPrimaryKey(v.([]interface{}))
		// TODO [SNOW-884959]: do we need quoteStringList?
		// TODO [SNOW-884959]: change name to optional
		sdk.NewOutOfLineConstraintRequest("TODO - optional", sdk.ColumnConstraintTypePrimaryKey).WithColumns(snowflake.QuoteStringList(pk.keys))
	}

	if v, ok := d.GetOk("data_retention_days"); ok {
		createRequest.WithDataRetentionTimeInDays(sdk.Int(v.(int)))
	} else if v, ok := d.GetOk("data_retention_time_in_days"); ok {
		createRequest.WithDataRetentionTimeInDays(sdk.Int(v.(int)))
	}

	if v, ok := d.GetOk("change_tracking"); ok {
		createRequest.WithChangeTracking(sdk.Bool(v.(bool)))
	}

	var tagAssociationRequests []sdk.TagAssociationRequest
	if _, ok := d.GetOk("tag"); ok {
		tagAssociations := getPropertyTags(d, "tag")
		tagAssociationRequests = make([]sdk.TagAssociationRequest, len(tagAssociations))
		for i, t := range tagAssociations {
			tagAssociationRequests[i] = *sdk.NewTagAssociationRequest(t.Name, t.Value)
		}
		createRequest.WithTags(tagAssociationRequests)
	}

	err := client.Tables.Create(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("error creating table %v err = %w", name, err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadTable(d, meta)
}

// ReadTable implements schema.ReadFunc.
func ReadTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	table, err := client.Tables.ShowByID(ctx, id)
	if err != nil {
		log.Printf("[DEBUG] table (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	tableDescription, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(id))
	if err != nil {
		return err
	}

	// Set the relevant data in the state
	toSet := map[string]interface{}{
		"name":       table.Name,
		"owner":      table.Owner,
		"database":   table.DatabaseName,
		"schema":     table.SchemaName,
		"comment":    table.Comment,
		"column":     toColumnConfig(tableDescription),
		"cluster_by": table.GetClusterByKeys(),
		// TODO [SNOW-884959]: SHOW PRIMARY KEYS IN TABLE? It was deprecated when table_constraint resource was introduced; should it be set here or not?
		// "primary_key":         snowflake.FlattenTablePrimaryKey(pkDescription),
		"change_tracking": table.ChangeTracking,
		"qualified_name":  id.FullyQualifiedName(),
	}
	var dataRetentionKey string
	if _, ok := d.GetOk("data_retention_time_in_days"); ok {
		dataRetentionKey = "data_retention_time_in_days"
	} else if _, ok := d.GetOk("data_retention_days"); ok {
		dataRetentionKey = "data_retention_days"
	}
	if dataRetentionKey != "" {
		toSet[dataRetentionKey] = table.RetentionTime
	}

	for key, val := range toSet {
		if err := d.Set(key, val); err != nil { // lintignore:R001
			return err
		}
	}
	return nil
}

// UpdateTable implements schema.UpdateFunc.
func UpdateTable(d *schema.ResourceData, meta interface{}) error {
	tid, err := tableIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := tid.DatabaseName
	schema := tid.SchemaName
	tableName := tid.TableName

	builder := snowflake.NewTableBuilder(tableName, dbName, schema)

	db := meta.(*sql.DB)
	if d.HasChange("name") {
		name := d.Get("name")
		q := builder.Rename(name.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating table name on %v", d.Id())
		}
		tableID := &tableID{
			DatabaseName: dbName,
			SchemaName:   schema,
			TableName:    name.(string),
		}
		dataIDInput, err := tableID.String()
		if err != nil {
			return err
		}
		d.SetId(dataIDInput)
	}
	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating table comment on %v", d.Id())
		}
	}
	if d.HasChange("column") {
		t, n := d.GetChange("column")
		removed, added, changed := getColumns(t).diffs(getColumns(n))
		for _, cA := range removed {
			q := builder.DropColumn(cA.name)
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error dropping column on %v", d.Id())
			}
		}
		for _, cA := range added {
			var q string

			if cA.identity == nil && cA._default == nil { //nolint:gocritic
				q = builder.AddColumn(cA.name, cA.dataType, cA.nullable, nil, nil, cA.comment, cA.maskingPolicy)
			} else if cA.identity != nil {
				q = builder.AddColumn(cA.name, cA.dataType, cA.nullable, nil, cA.identity.toSnowflakeColumnIdentity(), cA.comment, cA.maskingPolicy)
			} else {
				if cA._default._type() != "constant" {
					return fmt.Errorf("failed to add column %v => Only adding a column as a constant is supported by Snowflake", cA.name)
				}

				q = builder.AddColumn(cA.name, cA.dataType, cA.nullable, cA._default.toSnowflakeColumnDefault(), nil, cA.comment, cA.maskingPolicy)
			}

			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error adding column on %v", d.Id())
			}
		}
		for _, cA := range changed {
			if cA.changedDataType {
				q := builder.ChangeColumnType(cA.newColumn.name, cA.newColumn.dataType)
				if err := snowflake.Exec(db, q); err != nil {
					return fmt.Errorf("error changing property on %v", d.Id())
				}
			}
			if cA.changedNullConstraint {
				q := builder.ChangeNullConstraint(cA.newColumn.name, cA.newColumn.nullable)
				if err := snowflake.Exec(db, q); err != nil {
					return fmt.Errorf("error changing property on %v", d.Id())
				}
			}
			if cA.dropedDefault {
				q := builder.DropColumnDefault(cA.newColumn.name)
				if err := snowflake.Exec(db, q); err != nil {
					return fmt.Errorf("error changing property on %v", d.Id())
				}
			}
			if cA.changedComment {
				q := builder.ChangeColumnComment(cA.newColumn.name, cA.newColumn.comment)
				if err := snowflake.Exec(db, q); err != nil {
					return fmt.Errorf("error changing property on %v", d.Id())
				}
			}
			if cA.changedMaskingPolicy {
				q := builder.ChangeColumnMaskingPolicy(cA.newColumn.name, cA.newColumn.maskingPolicy)
				if err := snowflake.Exec(db, q); err != nil {
					return fmt.Errorf("error changing property on %v", d.Id())
				}
			}
		}
	}
	if d.HasChange("cluster_by") {
		cb := expandStringList(d.Get("cluster_by").([]interface{}))

		var q string
		if len(cb) != 0 {
			builder.WithClustering(cb)
			q = builder.ChangeClusterBy(builder.GetClusterKeyString())
		} else {
			q = builder.DropClustering()
		}

		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating table clustering on %v", d.Id())
		}
	}
	if d.HasChange("primary_key") {
		opk, npk := d.GetChange("primary_key")

		newpk := getPrimaryKey(npk)
		oldpk := getPrimaryKey(opk)

		if len(oldpk.keys) > 0 || len(newpk.keys) == 0 {
			// drop our pk if there was an old primary key, or pk has been removed
			q := builder.DropPrimaryKey()
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error changing primary key first on %v", d.Id())
			}
		}

		if len(newpk.keys) > 0 {
			// add our new pk
			q := builder.ChangePrimaryKey(newpk.toSnowflakePrimaryKey())
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error changing property on %v", d.Id())
			}
		}
	}
	updateDataRetention := func(key string) error {
		if d.HasChange(key) {
			ndr := d.Get(key)
			q := builder.ChangeDataRetention(ndr.(int))
			if err := snowflake.Exec(db, q); err != nil {
				return fmt.Errorf("error changing property on %v", d.Id())
			}
		}
		return nil
	}
	err = updateDataRetention("data_retention_days")
	if err != nil {
		return err
	}
	err = updateDataRetention("data_retention_time_in_days")
	if err != nil {
		return err
	}
	if d.HasChange("change_tracking") {
		nct := d.Get("change_tracking")

		q := builder.ChangeChangeTracking(nct.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error changing property on %v", d.Id())
		}
	}
	tagChangeErr := handleTagChanges(db, d, builder)
	if tagChangeErr != nil {
		return tagChangeErr
	}

	return ReadTable(d, meta)
}

// DeleteTable implements schema.DeleteFunc.
func DeleteTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(id))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
