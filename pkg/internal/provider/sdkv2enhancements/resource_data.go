package sdkv2enhancements

import (
	"reflect"
	"unsafe"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// CreateResourceDataFromResourceDiff allows to create schema.ResourceData out of schema.ResourceDiff.
// Unfortunately, schema.resourceDiffer is an unexported interface, and functions like schema.SchemaDiffSuppressFunc do not use the interface but a concrete implementation.
// One use case in which we needed to have schema.ResourceData from schema.ResourceDiff was to run schema.SchemaDiffSuppressFunc from inside schema.CustomizeDiffFunc.
// This implementation uses:
// - schema.InternalMap that exposes hidden schema.schemaMap (a wrapper over map[string]*schema.Schema)
// - (m schemaMap) Data method allowing to create schema.ResourceData from terraform.InstanceState and terraform.InstanceDiff
// - terraform.InstanceState and terraform.InstanceDiff are unexported in schema.ResourceDiff, so we get them using reflection
func CreateResourceDataFromResourceDiff(resourceSchema schema.InternalMap, diff *schema.ResourceDiff) (*schema.ResourceData, bool) {
	unexportedState := reflect.ValueOf(diff).Elem().FieldByName("state")
	stateFromResourceDiff := reflect.NewAt(unexportedState.Type(), unsafe.Pointer(unexportedState.UnsafeAddr())).Elem().Interface()
	unexportedDiff := reflect.ValueOf(diff).Elem().FieldByName("diff")
	diffFroResourceDif := reflect.NewAt(unexportedDiff.Type(), unsafe.Pointer(unexportedDiff.UnsafeAddr())).Elem().Interface()
	castState, ok := stateFromResourceDiff.(*terraform.InstanceState)
	if !ok {
		return nil, false
	}
	castDiff, ok := diffFroResourceDif.(*terraform.InstanceDiff)
	if !ok {
		return nil, false
	}
	resourceData, err := resourceSchema.Data(castState, castDiff)
	if err != nil {
		return nil, false
	}
	return resourceData, true
}
