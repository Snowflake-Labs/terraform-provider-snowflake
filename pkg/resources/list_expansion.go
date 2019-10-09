package resources

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

type stringSet map[string]struct{}

func (ss stringSet) Difference(other stringSet) stringSet {
	difference := map[string]struct{}{}

	for k, v := range ss {
		if _, ok := other[k]; ok {
			continue
		}
		difference[k] = v
	}
	return difference
}

func (ss stringSet) List() []string {
	ls := make([]string, len(ss))
	for k := range ss {
		ls = append(ls, k)
	}
	return ls
}

func createStringSet(entities []interface{}) stringSet {
	set := stringSet{}

	stringEntities := expandStringList(entities)
	for _, entity := range stringEntities {
		set[entity] = struct{}{}
	}

	return set
}
