package sdk

type DateFormat string

const (
	DateFormatISO   DateFormat = "DateFormatISO"
	DateFormatOther DateFormat = "DateFormatOther"
	DateFormatAny   DateFormat = "DateFormatAny"
)

func GetValidDateFormats(f DateFormat, includeAuto bool) []string {
	ISODateFormats := []string{
		"YYYY-MM-DD",
	}
	OtherDateFormats := []string{
		"DD-MON-YYYY",
		"MM/DD/YYYY",
	}
	AnyDateFormats := make([]string, 0, len(ISODateFormats)+len(OtherDateFormats))
	AnyDateFormats = append(AnyDateFormats, ISODateFormats...)
	AnyDateFormats = append(AnyDateFormats, OtherDateFormats...)
	var formats []string

	switch f {
	case DateFormatISO:
		formats = ISODateFormats
	case DateFormatOther:
		formats = OtherDateFormats
	case DateFormatAny:
		formats = AnyDateFormats
	}
	if includeAuto {
		formats = append(formats, "auto")
	}
	return formats
}

type TimeFormat string

const (
	TimeFormatISO TimeFormat = "TimeFormatISO"
	TimeFormatRFC TimeFormat = "TimeFormatRFC"
	TimeFormatAny TimeFormat = "TimeFormatAny"
)

func getValidTimeFormats(f TimeFormat, includeAuto bool) []string {
	ISOTimeFormats := []string{
		"HH24:MI:SS.FFTZH:TZM",
		"HH24:MI:SS.FF",
		"HH24:MI:SS",
		"HH24:MI",
	}
	RFCTimeFormats := []string{
		"HH12:MI:SS.FF AM",
		"HH12:MI:SS AM",
		"HH12:MI AM",
	}
	AnyTimeFormats := make([]string, 0, len(ISOTimeFormats)+len(RFCTimeFormats))
	AnyTimeFormats = append(AnyTimeFormats, ISOTimeFormats...)
	AnyTimeFormats = append(AnyTimeFormats, RFCTimeFormats...)
	var formats []string

	switch f {
	case TimeFormatISO:
		formats = ISOTimeFormats
	case TimeFormatRFC:
		formats = RFCTimeFormats
	case TimeFormatAny:
		formats = AnyTimeFormats
	}
	if includeAuto {
		formats = append(formats, "auto")
	}
	return formats
}

type TimeStampFormat string

const (
	TimeStampFormatISO   TimeStampFormat = "TimeStampFormatISO"
	TimeStampFormatRFC   TimeStampFormat = "TimeStampFormatRFC"
	TimeStampFormatOther TimeStampFormat = "TimeStampFormatOther"
	TimeStampFormatAny   TimeStampFormat = "TimeStampFormatAny"
)

func getValidTimeStampFormats(f TimeStampFormat, includeAuto bool) []string {
	ISOTimeStampFormats := []string{
		"YYYY-MM-DD\"T\"HH24:MI:SS.FFTZH:TZM",
		"YYYY-MM-DD HH24:MI:SS.FFTZH:TZM",
		"YYYY-MM-DD HH24:MI:SS.FFTZH",
		"YYYY-MM-DD HH24:MI:SS.FF TZH:TZM",
		"YYYY-MM-DD HH24:MI:SS.FF TZHTZM",
		"YYYY-MM-DD HH24:MI:SS TZH:TZM",
		"YYYY-MM-DD HH24:MI:SS TZHTZM",
		"YYYY-MM-DD\"T\"HH24:MI:SS.FF",
		"YYYY-MM-DD HH24:MI:SS.FF",
		"YYYY-MM-DD\"T\"HH24:MI:SS",
		"YYYY-MM-DD HH24:MI:SS",
		"YYYY-MM-DD\"T\"HH24:MI",
		"YYYY-MM-DD HH24:MI",
		"YYYY-MM-DD\"T\"HH24",
		"YYYY-MM-DD HH24",
		"YYYY-MM-DD\"T\"HH24:MI:SSTZH:TZM",
		"YYYY-MM-DD HH24:MI:SSTZH:TZM",
		"YYYY-MM-DD HH24:MI:SSTZH",
		"YYYY-MM-DD\"T\"HH24:MITZH:TZM",
		"YYYY-MM-DD HH24:MITZH:TZM",
	}
	RFCTimeStampFormats := []string{
		"DY, DD MON YYYY HH24:MI:SS TZHTZM",
		"DY, DD MON YYYY HH24:MI:SS.FF TZHTZM",
		"DY, DD MON YYYY HH12:MI:SS AM TZHTZM",
		"DY, DD MON YYYY HH12:MI:SS.FF AM TZHTZM",
		"DY, DD MON YYYY HH24:MI:SS",
		"DY, DD MON YYYY HH24:MI:SS.FF",
		"DY, DD MON YYYY HH12:MI:SS AM",
		"DY, DD MON YYYY HH12:MI:SS.FF AM",
	}
	OtherTimeStampFormats := []string{
		"MM/DD/YYYY HH24:MI:SS",
		"DY MON DD HH24:MI:SS TZHTZM YYYY",
	}
	AnyTimeStampFormats := make([]string, 0, len(ISOTimeStampFormats)+len(RFCTimeStampFormats)+len(OtherTimeStampFormats))
	AnyTimeStampFormats = append(AnyTimeStampFormats, ISOTimeStampFormats...)
	AnyTimeStampFormats = append(AnyTimeStampFormats, RFCTimeStampFormats...)
	AnyTimeStampFormats = append(AnyTimeStampFormats, OtherTimeStampFormats...)

	var formats []string
	switch f {
	case TimeStampFormatISO:
		formats = ISOTimeStampFormats
	case TimeStampFormatRFC:
		formats = RFCTimeStampFormats
	case TimeStampFormatOther:
		formats = OtherTimeStampFormats
	case TimeStampFormatAny:
		formats = AnyTimeStampFormats
	}

	if includeAuto {
		formats = append(formats, "auto")
	}
	return formats
}
