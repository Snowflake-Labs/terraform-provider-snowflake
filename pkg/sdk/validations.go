package sdk

func IsValidDataType(v string) bool {
	dt := DataTypeFromString(v)
	return dt != DataTypeUnknown
}
