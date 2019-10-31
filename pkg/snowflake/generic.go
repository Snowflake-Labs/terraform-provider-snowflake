package snowflake

import (
	"fmt"
	"sort"
	"strings"
)

type EntityType string

const (
	DatabaseType        EntityType = "DATABASE"
	ManagedAccountType  EntityType = "MANAGED ACCOUNT"
	ResourceMonitorType EntityType = "RESOURCE MONITOR"
	RoleType            EntityType = "ROLE"
	ShareType           EntityType = "SHARE"
	UserType            EntityType = "USER"
	WarehouseType       EntityType = "WAREHOUSE"
)

type Builder struct {
	entityType EntityType
	name       string
}

func (b *Builder) Show() string {
	return fmt.Sprintf(`SHOW %sS LIKE '%s'`, b.entityType, b.name)
}

func (b *Builder) Drop() string {
	return fmt.Sprintf(`DROP %s "%s"`, b.entityType, b.name)
}

func (b *Builder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER %s "%s" RENAME TO "%s"`, b.entityType, b.name, newName)
}

type AlterPropertiesBuilder struct {
	name             string
	entityType       EntityType
	stringProperties map[string]string
	boolProperties   map[string]bool
	intProperties    map[string]int
	floatProperties  map[string]float64
}

func (b *Builder) Alter() *AlterPropertiesBuilder {
	return &AlterPropertiesBuilder{
		name:             b.name,
		entityType:       b.entityType,
		stringProperties: make(map[string]string),
		boolProperties:   make(map[string]bool),
		intProperties:    make(map[string]int),
		floatProperties:  make(map[string]float64),
	}
}

func (ab *AlterPropertiesBuilder) SetString(key, value string) {
	ab.stringProperties[key] = value
}

func (ab *AlterPropertiesBuilder) SetBool(key string, value bool) {
	ab.boolProperties[key] = value
}

func (ab *AlterPropertiesBuilder) SetInt(key string, value int) {
	ab.intProperties[key] = value
}

func (ab *AlterPropertiesBuilder) SetFloat(key string, value float64) {
	ab.floatProperties[key] = value
}

func (ab *AlterPropertiesBuilder) Statement() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`ALTER %s "%s" SET`, ab.entityType, ab.name)) // TODO handle error

	for k, v := range ab.stringProperties {
		sb.WriteString(fmt.Sprintf(" %s='%s'", strings.ToUpper(k), EscapeString(v)))
	}

	for k, v := range ab.boolProperties {
		sb.WriteString(fmt.Sprintf(" %s=%t", strings.ToUpper(k), v))
	}

	for k, v := range ab.intProperties {
		sb.WriteString(fmt.Sprintf(" %s=%d", strings.ToUpper(k), v))
	}

	for k, v := range ab.floatProperties {
		sb.WriteString(fmt.Sprintf(" %s=%.2f", strings.ToUpper(k), v))
	}

	return sb.String()
}

type CreateBuilder struct {
	name             string
	entityType       EntityType
	stringProperties map[string]string
	boolProperties   map[string]bool
	intProperties    map[string]int
	floatProperties  map[string]float64
}

func (b *Builder) Create() *CreateBuilder {
	return &CreateBuilder{
		name:             b.name,
		entityType:       b.entityType,
		stringProperties: make(map[string]string),
		boolProperties:   make(map[string]bool),
		intProperties:    make(map[string]int),
		floatProperties:  make(map[string]float64),
	}
}

func (b *CreateBuilder) SetString(key, value string) {
	b.stringProperties[key] = value
}

func (b *CreateBuilder) SetBool(key string, value bool) {
	b.boolProperties[key] = value
}

func (b *CreateBuilder) SetInt(key string, value int) {
	b.intProperties[key] = value
}

func (b *CreateBuilder) SetFloat(key string, value float64) {
	b.floatProperties[key] = value
}

func (b *CreateBuilder) Statement() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`CREATE %s "%s"`, b.entityType, b.name)) // TODO handle error

	sortedStringProperties := make([]string, 0)
	for k := range b.stringProperties {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)

	for _, k := range sortedStringProperties {
		sb.WriteString(fmt.Sprintf(" %s='%s'", strings.ToUpper(k), EscapeString(b.stringProperties[k])))
	}

	sortedBoolProperties := make([]string, 0)
	for k := range b.boolProperties {
		sortedBoolProperties = append(sortedBoolProperties, k)
	}
	sort.Strings(sortedBoolProperties)

	for _, k := range sortedBoolProperties {
		sb.WriteString(fmt.Sprintf(" %s=%t", strings.ToUpper(k), b.boolProperties[k]))
	}

	sortedIntProperties := make([]string, 0)
	for k := range b.intProperties {
		sortedIntProperties = append(sortedIntProperties, k)
	}
	sort.Strings(sortedIntProperties)

	for _, k := range sortedIntProperties {
		sb.WriteString(fmt.Sprintf(" %s=%d", strings.ToUpper(k), b.intProperties[k]))
	}

	sortedFloatProperties := make([]string, 0)
	for k := range b.floatProperties {
		sortedFloatProperties = append(sortedFloatProperties, k)
	}
	sort.Strings(sortedFloatProperties)

	for _, k := range sortedFloatProperties {
		sb.WriteString(fmt.Sprintf(" %s=%.2f", strings.ToUpper(k), b.floatProperties[k]))
	}

	return sb.String()
}
