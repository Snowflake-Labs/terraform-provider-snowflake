package resources

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTypeValidateFunc(val interface{}, _ string) (warns []string, errs []error) {
	if ok := sdk.IsValidDataType(val.(string)); !ok {
		errs = append(errs, fmt.Errorf("%v is not a valid data type", val))
	}
	return
}

func dataTypeDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	oldDT, err := sdk.ToDataType(old)
	if err != nil {
		return false
	}
	newDT, err := sdk.ToDataType(new)
	if err != nil {
		return false
	}
	return oldDT == newDT
}

// DataTypeIssue3007DiffSuppressFunc is a temporary solution to handle data type suppression problems.
// Currently, it handles only number and text data types.
// It falls back to Snowflake defaults for arguments if no arguments were provided for the data type.
// TODO [SNOW-1348103 or SNOW-1348106]: visit with functions and procedures rework
func DataTypeIssue3007DiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	oldDataType, err := sdk.ToDataType(old)
	if err != nil {
		return false
	}
	newDataType, err := sdk.ToDataType(new)
	if err != nil {
		return false
	}
	if oldDataType != newDataType {
		return false
	}
	switch v := oldDataType; v {
	case sdk.DataTypeNumber:
		logging.DebugLogger.Printf("[DEBUG] DataTypeIssue3007DiffSuppressFunc: Handling number data type diff suppression")
		oldPrecision, oldScale := sdk.ParseNumberDataTypeRaw(old)
		newPrecision, newScale := sdk.ParseNumberDataTypeRaw(new)
		return oldPrecision == newPrecision && oldScale == newScale
	case sdk.DataTypeVARCHAR:
		logging.DebugLogger.Printf("[DEBUG] DataTypeIssue3007DiffSuppressFunc: Handling text data type diff suppression")
		oldLength := sdk.ParseVarcharDataTypeRaw(old)
		newLength := sdk.ParseVarcharDataTypeRaw(new)
		return oldLength == newLength
	default:
		logging.DebugLogger.Printf("[DEBUG] DataTypeIssue3007DiffSuppressFunc: Diff suppression for %s can't be currently handled", v)
	}
	return true
}

func ignoreTrimSpaceSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func ignoreCaseSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

func ignoreCaseAndTrimSpaceSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(strings.TrimSpace(old), strings.TrimSpace(new))
}

func getTagObjectIdentifier(obj map[string]any) sdk.ObjectIdentifier {
	database := obj["database"].(string)
	schema := obj["schema"].(string)
	name := obj["name"].(string)
	switch {
	case schema != "":
		return sdk.NewSchemaObjectIdentifier(database, schema, name)
	case database != "":
		return sdk.NewDatabaseObjectIdentifier(database, name)
	default:
		return sdk.NewAccountObjectIdentifier(name)
	}
}

func getPropertyTags(d *schema.ResourceData, key string) []sdk.TagAssociation {
	if from, ok := d.GetOk(key); ok {
		return getTagsFromList(from.([]any))
	}
	return nil
}

func getTagsFromList(tags []any) []sdk.TagAssociation {
	to := make([]sdk.TagAssociation, len(tags))
	for i, t := range tags {
		v := t.(map[string]any)
		to[i] = sdk.TagAssociation{
			Name:  getTagObjectIdentifier(v),
			Value: v["value"].(string),
		}
	}
	return to
}

func GetTagsDiff(d *schema.ResourceData, key string) (unsetTags []sdk.ObjectIdentifier, setTags []sdk.TagAssociation) {
	o, n := d.GetChange(key)
	removed, added, changed := getTags(o).diffs(getTags(n))

	unsetTags = make([]sdk.ObjectIdentifier, len(removed))
	for i, t := range removed {
		unsetTags[i] = sdk.NewSchemaObjectIdentifier(t.database, t.schema, t.name)
	}

	setTags = make([]sdk.TagAssociation, len(added)+len(changed))
	for i, t := range added {
		setTags[i] = sdk.TagAssociation{
			Name:  sdk.NewSchemaObjectIdentifier(t.database, t.schema, t.name),
			Value: t.value,
		}
	}
	for i, t := range changed {
		setTags[len(added)+i] = sdk.TagAssociation{
			Name:  sdk.NewSchemaObjectIdentifier(t.database, t.schema, t.name),
			Value: t.value,
		}
	}

	return unsetTags, setTags
}

func GetPropertyAsPointer[T any](d *schema.ResourceData, property string) *T {
	value, ok := d.GetOk(property)
	if !ok {
		return nil
	}
	typedValue, ok := value.(T)
	if !ok {
		return nil
	}
	return &typedValue
}

func GetConfigPropertyAsPointerAllowingZeroValue[T any](d *schema.ResourceData, property string) *T {
	if d.GetRawConfig().AsValueMap()[property].IsNull() {
		return nil
	}
	value := d.Get(property)
	typedValue, ok := value.(T)
	if !ok {
		return nil
	}
	return &typedValue
}

func GetPropertyOfFirstNestedObjectByValueKey[T any](d *schema.ResourceData, propertyKey string) (*T, error) {
	return GetPropertyOfFirstNestedObjectByKey[T](d, propertyKey, "value")
}

// GetPropertyOfFirstNestedObjectByKey should be used for single objects defined in the Terraform schema as
// schema.TypeList with MaxItems set to one and inner schema with single value. To easily retrieve
// the inner value, you can specify the top-level property with propertyKey and the nested value with nestedValueKey.
func GetPropertyOfFirstNestedObjectByKey[T any](d *schema.ResourceData, propertyKey string, nestedValueKey string) (*T, error) {
	value, ok := d.GetOk(propertyKey)
	if !ok {
		return nil, fmt.Errorf("nested property %s not found", propertyKey)
	}

	typedValue, ok := value.([]any)
	if !ok || len(typedValue) != 1 {
		return nil, fmt.Errorf("nested property %s is not an array or has incorrect number of values: %d, expected: 1", propertyKey, len(typedValue))
	}

	typedNestedMap, ok := typedValue[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("nested property %s is not of type map[string]any, got: %T", propertyKey, typedValue[0])
	}

	_, ok = typedNestedMap[nestedValueKey]
	if !ok {
		return nil, fmt.Errorf("nested value key %s couldn't be found in the nested property map %s", nestedValueKey, propertyKey)
	}

	typedNestedValue, ok := typedNestedMap[nestedValueKey].(T)
	if !ok {
		return nil, fmt.Errorf("nested property %s.%s is not of type %T, got: %T", propertyKey, nestedValueKey, *new(T), typedNestedMap[nestedValueKey])
	}

	return &typedNestedValue, nil
}

func SetPropertyOfFirstNestedObjectByValueKey[T any](d *schema.ResourceData, propertyKey string, value T) error {
	return SetPropertyOfFirstNestedObjectByKey[T](d, propertyKey, "value", value)
}

// SetPropertyOfFirstNestedObjectByKey should be used for single objects defined in the Terraform schema as
// schema.TypeList with MaxItems set to one and inner schema with single value. To easily set
// the inner value, you can specify top-level property with propertyKey, nested value with nestedValueKey and value at the end.
func SetPropertyOfFirstNestedObjectByKey[T any](d *schema.ResourceData, propertyKey string, nestedValueKey string, value T) error {
	return d.Set(propertyKey, []any{
		map[string]any{
			nestedValueKey: value,
		},
	})
}

type tags []tag

func (t tags) toSnowflakeTagValues() []snowflake.TagValue {
	sT := make([]snowflake.TagValue, len(t))
	for i, tag := range t {
		sT[i] = tag.toSnowflakeTagValue()
	}
	return sT
}

func (t tag) toSnowflakeTagValue() snowflake.TagValue {
	return snowflake.TagValue{
		Name:     t.name,
		Value:    t.value,
		Database: t.database,
		Schema:   t.schema,
	}
}

func (t tags) getNewIn(new tags) (added tags) {
	added = tags{}
	for _, t0 := range t {
		found := false
		for _, cN := range new {
			if t0.name == cN.name {
				found = true
				break
			}
		}
		if !found {
			added = append(added, t0)
		}
	}
	return
}

func (t tags) getChangedTagProperties(new tags) (changed tags) {
	changed = tags{}
	for _, t0 := range t {
		for _, tN := range new {
			if t0.name == tN.name && t0.value != tN.value {
				changed = append(changed, tN)
			}
		}
	}
	return
}

func (t tags) diffs(new tags) (removed tags, added tags, changed tags) {
	return t.getNewIn(new), new.getNewIn(t), t.getChangedTagProperties(new)
}

func (t columns) getNewIn(new columns) (added columns) {
	added = columns{}
	for _, cO := range t {
		found := false
		for _, cN := range new {
			if cO.name == cN.name {
				found = true
				break
			}
		}
		if !found {
			added = append(added, cO)
		}
	}
	return
}

type tag struct {
	name     string
	value    string
	database string
	schema   string
}

func getTags(from interface{}) (to tags) {
	tags := from.([]interface{})
	to = make([]tag, len(tags))
	for i, t := range tags {
		v := t.(map[string]interface{})
		to[i] = tag{
			name:     v["name"].(string),
			value:    v["value"].(string),
			database: v["database"].(string),
			schema:   v["schema"].(string),
		}
	}
	return to
}

// TODO(SNOW-1479870): Test
// JoinDiags iterates through passed diag.Diagnostics and joins them into one diag.Diagnostics.
// If none of the passed diagnostics contained any element a nil reference will be returned.
func JoinDiags(diagnostics ...diag.Diagnostics) diag.Diagnostics {
	var result diag.Diagnostics
	for _, diagnostic := range diagnostics {
		if len(diagnostic) > 0 {
			result = append(result, diagnostic...)
		}
	}
	return result
}

// ListDiff compares two lists (before and after), then compares and returns two lists that include
// added and removed items between those lists.
func ListDiff[T comparable](beforeList []T, afterList []T) (added []T, removed []T) {
	added, removed, _ = ListDiffWithCommonItems(beforeList, afterList)
	return
}

// ListDiffWithCommonItems compares two lists (before and after), then compares and returns three lists that include
// added, removed and common items between those lists.
func ListDiffWithCommonItems[T comparable](beforeList []T, afterList []T) (added []T, removed []T, common []T) {
	added = make([]T, 0)
	removed = make([]T, 0)
	common = make([]T, 0)

	for _, beforeItem := range beforeList {
		if !slices.Contains(afterList, beforeItem) {
			removed = append(removed, beforeItem)
		} else {
			common = append(common, beforeItem)
		}
	}

	for _, afterItem := range afterList {
		if !slices.Contains(beforeList, afterItem) {
			added = append(added, afterItem)
		}
	}

	return added, removed, common
}

// parseSchemaObjectIdentifierSet is a helper function to parse a given schema object identifier list from ResourceData.
func parseSchemaObjectIdentifierSet(v any) ([]sdk.SchemaObjectIdentifier, error) {
	idsRaw := expandStringList(v.(*schema.Set).List())
	ids := make([]sdk.SchemaObjectIdentifier, len(idsRaw))
	for i, idRaw := range idsRaw {
		id, err := sdk.ParseSchemaObjectIdentifier(idRaw)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}
