package resources

import "golang.org/x/exp/slices"

// borrowed from https://github.com/terraform-providers/terraform-provider-aws/blob/master/aws/structure.go#L924:6

func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		if val, ok := v.(int); ok {
			vs = append(vs, val)
		}
	}
	return vs
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func expandStringListAllowEmpty(configured []interface{}) []string {
	// Allow empty values during expansion
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok {
			vs = append(vs, val)
		}
	}
	return vs
}

func expandObjectIdentifier(objectIdentifier interface{}) (string, string, string) {
	objectIdentifierMap := objectIdentifier.([]interface{})[0].(map[string]interface{})
	objectName := objectIdentifierMap["name"].(string)
	var objectSchema string
	if v := objectIdentifierMap["schema"]; v != nil {
		objectSchema = v.(string)
	}
	var objectDatabase string
	if v := objectIdentifierMap["database"]; v != nil {
		objectDatabase = v.(string)
	}
	return objectDatabase, objectSchema, objectName
}

// intersectionAAndNotB takes the intersection of set A and the intersection of not set B. A∩B′ in set notation.
func intersectionAAndNotB(setA []interface{}, setB []interface{}) []string {
	res := make([]string, 0)
	sliceA := expandStringList(setA)
	sliceB := expandStringList(setB)
	for _, s := range sliceA {
		if !slices.Contains(sliceB, s) {
			res = append(res, s)
		}
	}
	return res
}
