package snowflake

import (
	"fmt"
	"sort"
	"strings"
)

type UserBuilder struct {
	name string
}

func User(name string) *UserBuilder {
	return &UserBuilder{name: name}
}

func (b *UserBuilder) Show() string {
	return fmt.Sprintf(`SHOW USERS LIKE '%s'`, b.name)
}

func (b *UserBuilder) Drop() string {
	return fmt.Sprintf(`DROP USER "%s"`, b.name)
}

func (b *UserBuilder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER USER "%s" RENAME TO "%s"`, b.name, newName)
}

type UserAlterPropertiesBuilder struct {
	name             string
	stringProperties map[string]string
	boolProperties   map[string]bool
}

func (b *UserBuilder) Alter() *UserAlterPropertiesBuilder {
	return &UserAlterPropertiesBuilder{
		name:             b.name,
		stringProperties: make(map[string]string),
		boolProperties:   make(map[string]bool),
	}
}

func (ab *UserAlterPropertiesBuilder) SetString(key, value string) {
	ab.stringProperties[key] = value
}

func (ab *UserAlterPropertiesBuilder) SetBool(key string, value bool) {
	ab.boolProperties[key] = value
}

func (ab *UserAlterPropertiesBuilder) Statement() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`ALTER USER "%s" SET`, ab.name)) // TODO handle error

	for k, v := range ab.stringProperties {
		sb.WriteString(fmt.Sprintf(" %s='%s'", strings.ToUpper(k), EscapeString(v)))
	}

	for k, v := range ab.boolProperties {
		sb.WriteString(fmt.Sprintf(" %s=%t", strings.ToUpper(k), v))
	}

	return sb.String()
}

type UserCreateBuilder struct {
	name             string
	stringProperties map[string]string
	boolProperties   map[string]bool
}

func (b *UserBuilder) Create() *UserCreateBuilder {
	return &UserCreateBuilder{
		name:             b.name,
		stringProperties: make(map[string]string),
		boolProperties:   make(map[string]bool),
	}
}

func (b *UserCreateBuilder) SetString(key, value string) {
	b.stringProperties[key] = value
}

func (b *UserCreateBuilder) SetBool(key string, value bool) {
	b.boolProperties[key] = value
}

func (b *UserCreateBuilder) Statement() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`CREATE USER "%s"`, b.name)) // TODO handle error

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

	return sb.String()
}
