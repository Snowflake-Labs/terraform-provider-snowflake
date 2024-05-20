package generator

type Config struct {
	Values map[string]any
}

func NewConfig() *Config {
	return &Config{
		Values: make(map[string]any),
	}
}

type Pattern struct {
	Field   string
	Pattern string
}

func (c *Config) WithLike(patterns []Pattern) *Config {
	c.Values["Like"] = patterns
	return c
}

func (c *Config) WithIn(patterns []Pattern) *Config {
	c.Values["In"] = patterns
	return c
}
