package config

// ProviderModel is the base interface all of our provider config models will implement.
// To allow easy implementation, ProviderModelMeta can be embedded inside the struct (and the struct will automatically implement it).
type ProviderModel interface {
	ProviderName() string
	Alias() string
}

type ProviderModelMeta struct {
	name  string
	alias string
}

func DefaultProviderMeta(name string) *ProviderModelMeta {
	return &ProviderModelMeta{name: name}
}

func ProviderMeta(name string, alias string) *ProviderModelMeta {
	return &ProviderModelMeta{name: name, alias: alias}
}

func (m *ProviderModelMeta) ProviderName() string {
	return m.name
}

func (m *ProviderModelMeta) Alias() string {
	return m.alias
}
