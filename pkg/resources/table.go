package resources

import (
	"context"
	"database/sql"
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

// TODO [SNOW-867235]: old implementation was quoting every column, SDK is not quoting them, therefore they are quoted here: decide if we quote columns or not
// TODO [SNOW-1031688]: move data manipulation logic to the SDK - SQL generation or builders part (e.g. different default types/identity)
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
						// TODO [SNOW-867235]: there is no such separation on SDK level. Should we keep it in V1?
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

type columnDefault struct {
	constant   *string
	expression *string
	sequence   *string
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

type column struct {
	name          string
	dataType      string
	nullable      bool
	_default      *columnDefault
	identity      *columnIdentity
	comment       string
	maskingPolicy string
}

type columns []column

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

	nameInQuotes := fmt.Sprintf(`"%v"`, snowflake.EscapeString(c["name"].(string)))
	request := sdk.NewTableColumnRequest(nameInQuotes, sdk.DataType(_type))

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

func toColumnConfig(descriptions []sdk.TableColumnDetails) []any {
	flattened := make([]any, 0)
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
			// TODO [SNOW-867240]: SHOW TABLE returns last part of id without double quotes... we have to quote it again. Move it to SDK.
			flat["masking_policy"] = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(*td.PolicyName).FullyQualifiedName()
		}

		identity := toColumnIdentityConfig(td)
		if identity != nil {
			flat["identity"] = []any{identity}
		} else {
			def := toColumnDefaultConfig(td)
			if def != nil {
				flat["default"] = []any{def}
			}
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
		// TODO [SNOW-867240]: SHOW TABLE returns last part of id without double quotes... we have to quote it again. Move it to SDK.
		sequenceIdRaw := strings.TrimSuffix(defaultRaw, ".NEXTVAL")
		def["sequence"] = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(sequenceIdRaw).FullyQualifiedName()
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
		start, err := strconv.Atoi(split[2])
		if err == nil {
			identity["start_num"] = start
		}
		step, err := strconv.Atoi(split[4])
		if err == nil {
			identity["step_num"] = step
		}

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
		createRequest.WithClusterBy(expandStringList(v.([]interface{})))
	}

	if v, ok := d.GetOk("primary_key"); ok {
		keysList := v.([]any)
		if len(keysList) > 0 {
			keys := expandStringList(keysList[0].(map[string]any)["keys"].([]interface{}))
			constraintRequest := sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypePrimaryKey).WithColumns(snowflake.QuoteStringList(keys))

			keyName, isPresent := keysList[0].(map[string]any)["name"]
			if isPresent && keyName != "" {
				constraintRequest.WithName(sdk.String(keyName.(string)))
			}
		}
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
		"name":            table.Name,
		"owner":           table.Owner,
		"database":        table.DatabaseName,
		"schema":          table.SchemaName,
		"comment":         table.Comment,
		"column":          toColumnConfig(tableDescription),
		"cluster_by":      table.GetClusterByKeys(),
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
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("name") {
		newName := d.Get("name").(string)

		newId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), newName)

		err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithNewName(&newId))
		if err != nil {
			return fmt.Errorf("error renaming table %v err = %w", d.Id(), err)
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	var runSetStatement bool
	var runUnsetStatement bool
	setRequest := sdk.NewTableSetRequest()
	unsetRequest := sdk.NewTableUnsetRequest()

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if comment == "" {
			runUnsetStatement = true
			unsetRequest.WithComment(true)
		} else {
			runSetStatement = true
			setRequest.WithComment(sdk.String(comment))
		}
	}

	if d.HasChange("change_tracking") {
		changeTracking := d.Get("change_tracking").(bool)
		runSetStatement = true
		setRequest.WithChangeTracking(sdk.Bool(changeTracking))
	}

	checkChangeForDataRetention := func(key string) {
		if d.HasChange(key) {
			dataRetentionDays := d.Get(key).(int)
			runSetStatement = true
			setRequest.WithDataRetentionTimeInDays(sdk.Int(dataRetentionDays))
		}
	}
	checkChangeForDataRetention("data_retention_days")
	checkChangeForDataRetention("data_retention_time_in_days")

	if runSetStatement {
		err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithSet(setRequest))
		if err != nil {
			return fmt.Errorf("error updating table: %w", err)
		}
	}

	if runUnsetStatement {
		err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithUnset(unsetRequest))
		if err != nil {
			return fmt.Errorf("error updating table: %w", err)
		}
	}

	if d.HasChange("cluster_by") {
		cb := expandStringList(d.Get("cluster_by").([]interface{}))

		if len(cb) != 0 {
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithClusteringAction(sdk.NewTableClusteringActionRequest().WithClusterBy(cb)))
			if err != nil {
				return fmt.Errorf("error updating table: %w", err)
			}
		} else {
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithClusteringAction(sdk.NewTableClusteringActionRequest().WithDropClusteringKey(sdk.Bool(true))))
			if err != nil {
				return fmt.Errorf("error updating table: %w", err)
			}
		}
	}

	if d.HasChange("column") {
		t, n := d.GetChange("column")
		removed, added, changed := getColumns(t).diffs(getColumns(n))

		if len(removed) > 0 {
			removedColumnNames := make([]string, len(removed))
			for i, r := range removed {
				removedColumnNames[i] = r.name
			}
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(sdk.NewTableColumnActionRequest().WithDropColumns(snowflake.QuoteStringList(removedColumnNames))))
			if err != nil {
				return fmt.Errorf("error updating table: %w", err)
			}
		}

		for _, cA := range added {
			addRequest := sdk.NewTableColumnAddActionRequest(fmt.Sprintf("\"%s\"", cA.name), sdk.DataType(cA.dataType)).
				WithInlineConstraint(sdk.NewTableColumnAddInlineConstraintRequest().WithNotNull(sdk.Bool(!cA.nullable)))

			if cA._default != nil {
				if cA._default._type() != "constant" {
					return fmt.Errorf("failed to add column %v => Only adding a column as a constant is supported by Snowflake", cA.name)
				}
				var expression string
				if strings.Contains(cA.dataType, "CHAR") || cA.dataType == "STRING" || cA.dataType == "TEXT" {
					expression = snowflake.EscapeSnowflakeString(*cA._default.constant)
				} else {
					expression = *cA._default.constant
				}
				addRequest.WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithExpression(sdk.String(expression)))
			}

			if cA.identity != nil {
				addRequest.WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(cA.identity.startNum, cA.identity.stepNum)))
			}

			if cA.maskingPolicy != "" {
				addRequest.WithMaskingPolicy(sdk.NewColumnMaskingPolicyRequest(sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(cA.maskingPolicy)))
			}

			if cA.comment != "" {
				addRequest.WithComment(sdk.String(cA.comment))
			}

			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(sdk.NewTableColumnActionRequest().WithAdd(addRequest)))
			if err != nil {
				return fmt.Errorf("error adding column: %w", err)
			}
		}
		for _, cA := range changed {
			if cA.changedDataType {
				err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*sdk.NewTableColumnAlterActionRequest(fmt.Sprintf("\"%s\"", cA.newColumn.name)).WithType(sdk.Pointer(sdk.DataType(cA.newColumn.dataType)))})))
				if err != nil {
					return fmt.Errorf("error changing property on %v: err %w", d.Id(), err)
				}
			}
			if cA.changedNullConstraint {
				nullabilityRequest := sdk.NewTableColumnNotNullConstraintRequest()
				if !cA.newColumn.nullable {
					nullabilityRequest.WithSet(sdk.Bool(true))
				} else {
					nullabilityRequest.WithDrop(sdk.Bool(true))
				}
				err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*sdk.NewTableColumnAlterActionRequest(fmt.Sprintf("\"%s\"", cA.newColumn.name)).WithNotNullConstraint(nullabilityRequest)})))
				if err != nil {
					return fmt.Errorf("error changing property on %v: err %w", d.Id(), err)
				}
			}
			if cA.dropedDefault {
				err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*sdk.NewTableColumnAlterActionRequest(fmt.Sprintf("\"%s\"", cA.newColumn.name)).WithDropDefault(sdk.Bool(true))})))
				if err != nil {
					return fmt.Errorf("error changing property on %v: err %w", d.Id(), err)
				}
			}
			if cA.changedComment {
				columnAlterActionRequest := sdk.NewTableColumnAlterActionRequest(fmt.Sprintf("\"%s\"", cA.newColumn.name))
				if cA.newColumn.comment == "" {
					columnAlterActionRequest.WithUnsetComment(sdk.Bool(true))
				} else {
					columnAlterActionRequest.WithComment(sdk.String(cA.newColumn.comment))
				}

				err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*columnAlterActionRequest})))
				if err != nil {
					return fmt.Errorf("error changing property on %v: err %w", d.Id(), err)
				}
			}
			if cA.changedMaskingPolicy {
				columnAction := sdk.NewTableColumnActionRequest()
				if strings.TrimSpace(cA.newColumn.maskingPolicy) == "" {
					columnAction.WithUnsetMaskingPolicy(sdk.NewTableColumnAlterUnsetMaskingPolicyActionRequest(fmt.Sprintf("\"%s\"", cA.newColumn.name)))
				} else {
					columnAction.WithSetMaskingPolicy(sdk.NewTableColumnAlterSetMaskingPolicyActionRequest(fmt.Sprintf("\"%s\"", cA.newColumn.name), sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(cA.newColumn.maskingPolicy), []string{}).WithForce(sdk.Bool(true)))
				}
				err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(columnAction))
				if err != nil {
					return fmt.Errorf("error changing property on %v: err %w", d.Id(), err)
				}
			}
		}
	}

	if d.HasChange("primary_key") {
		o, n := d.GetChange("primary_key")

		newKey := getPrimaryKey(n)
		oldKey := getPrimaryKey(o)

		if len(oldKey.keys) > 0 || len(newKey.keys) == 0 {
			// drop our pk if there was an old primary key, or pk has been removed
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithConstraintAction(
				sdk.NewTableConstraintActionRequest().
					WithDrop(sdk.NewTableConstraintDropActionRequest().WithPrimaryKey(sdk.Bool(true))),
			))
			if err != nil {
				return fmt.Errorf("error updating table: %w", err)
			}
		}

		if len(newKey.keys) > 0 {
			constraint := sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypePrimaryKey).WithColumns(snowflake.QuoteStringList(newKey.keys))
			if newKey.name != "" {
				constraint.WithName(sdk.String(newKey.name))
			}
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithConstraintAction(
				sdk.NewTableConstraintActionRequest().WithAdd(constraint),
			))
			if err != nil {
				return fmt.Errorf("error updating table: %w", err)
			}
		}
	}

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithUnsetTags(unsetTags))
			if err != nil {
				return fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err)
			}
		}

		if len(setTags) > 0 {
			tagAssociationRequests := make([]sdk.TagAssociationRequest, len(setTags))
			for i, t := range setTags {
				tagAssociationRequests[i] = *sdk.NewTagAssociationRequest(t.Name, t.Value)
			}
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithSetTags(tagAssociationRequests))
			if err != nil {
				return fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err)
			}
		}
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
