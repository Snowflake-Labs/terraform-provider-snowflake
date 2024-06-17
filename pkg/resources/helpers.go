package resources

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func isOk(_ interface{}, ok bool) bool {
	return ok
}

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

func ignoreTrimSpaceSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func ignoreCaseSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

func ignoreCaseAndTrimSpaceSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(strings.TrimSpace(old), strings.TrimSpace(new))
}

func setIntProperty(d *schema.ResourceData, key string, property *sdk.IntProperty) error {
	if property != nil && property.Value != nil {
		if err := d.Set(key, *property.Value); err != nil {
			return err
		}
	}
	return nil
}

func setStringProperty(d *schema.ResourceData, key string, property *sdk.StringProperty) error {
	if property != nil {
		if err := d.Set(key, property.Value); err != nil {
			return err
		}
	}
	return nil
}

func setBoolProperty(d *schema.ResourceData, key string, property *sdk.BoolProperty) error {
	if property != nil {
		if err := d.Set(key, property.Value); err != nil {
			return err
		}
	}
	return nil
}

func getTagObjectIdentifier(v map[string]any) sdk.ObjectIdentifier {
	if _, ok := v["database"]; ok {
		if _, ok := v["schema"]; ok {
			return sdk.NewSchemaObjectIdentifier(v["database"].(string), v["schema"].(string), v["name"].(string))
		}
		return sdk.NewDatabaseObjectIdentifier(v["database"].(string), v["name"].(string))
	}
	return sdk.NewAccountObjectIdentifier(v["name"].(string))
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
// MergeMaps takes any number of maps (of the same type) and concatinates them.
// In case of key collision, the value will be selected from the map that is provided
// later in the src function parameter.
func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
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

// ListDiff Compares two lists (before and after), then compares and returns two lists that include
// added and removed items between those lists.
func ListDiff[T comparable](beforeList []T, afterList []T) (added []T, removed []T) {
	added = make([]T, 0)
	removed = make([]T, 0)

	for _, privilegeBeforeChange := range beforeList {
		if !slices.Contains(afterList, privilegeBeforeChange) {
			removed = append(removed, privilegeBeforeChange)
		}
	}

	for _, privilegeAfterChange := range afterList {
		if !slices.Contains(beforeList, privilegeAfterChange) {
			added = append(added, privilegeAfterChange)
		}
	}

	return added, removed
}
