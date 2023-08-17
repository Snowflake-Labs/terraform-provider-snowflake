package sdk

func findOne[T any](collection []T, condition func(T) bool) (*T, error) {
	for _, o := range collection {
		if condition(o) {
			return &o, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}
