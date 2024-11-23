package schema

// StringToDataType Convert string to DataType
func StringToDataType(value string) (DataType, bool) {
	switch value {
	case string(DataTypeString):
		return DataTypeString, true
	case string(DataTypeNumber):
		return DataTypeNumber, true
	case string(DataTypeInteger):
		return DataTypeInteger, true
	case string(DataTypeObject):
		return DataTypeObject, true
	case string(DataTypeArray):
		return DataTypeArray, true
	case string(DataTypeBoolean):
		return DataTypeBoolean, true
	case string(DataTypeNull):
		return DataTypeNull, true
	default:
		return "", false
	}
}

// StringToDataFormat Convert string to DataFormat
func StringToDataFormat(value string) (DataFormat, bool) {
	switch value {
	case string(DataFormatDateTime):
		return DataFormatDateTime, true
	case string(DataFormatTime):
		return DataFormatTime, true
	case string(DataFormatDate):
		return DataFormatDate, true
	case string(DataFormatDuration):
		return DataFormatDuration, true
	default:
		return "", false
	}
}
