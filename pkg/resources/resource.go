package resources

type ResourceValueSetter interface {
	Set(string, any) error
}
