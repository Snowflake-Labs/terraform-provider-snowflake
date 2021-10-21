package snowflake

import (
	"fmt"
	"sort"
)

func sortInterfaceStrings(strs map[string]interface{}) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStrings(strs map[string]string) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStringList(strs map[string][]string) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStringsInt(strs map[string]int) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStringsFloat(strs map[string]float64) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStringsBool(strs map[string]bool) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortTags(tags []TagValue) []TagValue {
	sort.Slice(tags, func(i, j int) bool {
		qn1 := fmt.Sprintf("%s.%s.%s", tags[i].Database, tags[i].Schema, tags[i].Name)
		qn2 := fmt.Sprintf("%s.%s.%s", tags[j].Database, tags[j].Schema, tags[j].Name)
		return qn1 < qn2
	})
	return tags
}
