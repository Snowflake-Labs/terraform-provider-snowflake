package resources

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NormalizeAndCompare[T comparable](normalize func(string) (T, error)) schema.SchemaDiffSuppressFunc {
	return func(_, oldValue, newValue string, _ *schema.ResourceData) bool {
		oldNormalized, err := normalize(oldValue)
		if err != nil {
			return false
		}
		newNormalized, err := normalize(newValue)
		if err != nil {
			return false
		}
		return oldNormalized == newNormalized
	}
}

// IgnoreAfterCreation should be used to ignore changes to the given attribute post creation.
func IgnoreAfterCreation(_, _, _ string, d *schema.ResourceData) bool {
	// For new resources always show the diff and in every other case we do not want to use this attribute
	return d.Id() != ""
}

func IgnoreChangeToCurrentSnowflakeValue(keyInShowOutput string) schema.SchemaDiffSuppressFunc {
	return func(_, _, new string, d *schema.ResourceData) bool {
		if d.Id() == "" {
			return false
		}

		if showOutput, ok := d.GetOk(showOutputAttributeName); ok {
			showOutputList := showOutput.([]any)
			if len(showOutputList) == 1 {
				result := showOutputList[0].(map[string]any)
				log.Printf("[DEBUG] IgnoreChangeToCurrentSnowflakeValue: value for key %s is %v, new value is %s, comparison result is: %t", keyInShowOutput, result[keyInShowOutput], new, new == fmt.Sprintf("%v", result[keyInShowOutput]))
				if new == fmt.Sprintf("%v", result[keyInShowOutput]) {
					return true
				}
			}
		}
		return false
	}
}

func SuppressIfAny(diffSuppressFunctions ...schema.SchemaDiffSuppressFunc) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		var suppress bool
		for _, f := range diffSuppressFunctions {
			suppress = suppress || f(k, old, new, d)
		}
		return suppress
	}
}
