package config

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
