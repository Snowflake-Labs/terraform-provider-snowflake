package resources

import (
	"bytes"
	"text/template"
)

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

func expandStringListToStorageLocations(configured []interface{}) string {
	list := expandStringList(configured)

	t, err := template.New("StorageLocations").Parse(`({{ range $i, $v := .}}{{ if $i }}, {{ end }}'{{ $v }}'{{ end }})`)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer

	if err := t.Execute(&buf, list); err != nil {
		return ""
	}

	return buf.String()
}
