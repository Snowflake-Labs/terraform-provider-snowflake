package snowflake

type BaseManager struct {
	sqlBuilder SQLBuilder
}

func (m *BaseManager) Ok(_ interface{}, ok bool) bool {
	return ok
}
