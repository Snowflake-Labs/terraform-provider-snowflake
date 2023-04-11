package snowflake

type BaseManager struct {
	genericBuilder NewBuilder
}

func (m *BaseManager) Ok(_ interface{}, ok bool) bool {
	return ok
}
