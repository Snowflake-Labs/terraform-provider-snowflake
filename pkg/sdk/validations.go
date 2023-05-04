package sdk

func IsValidDataType(v string) bool {
	_, err := DataTypeFromString(v)
	return err == nil
}
