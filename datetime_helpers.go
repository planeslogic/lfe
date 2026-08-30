package lfe

import "time"

const dateTimeUTCLayout = "2006-01-02 15:04:05"

// DateTimeUTC converts a time.Time to the LFE DateTime representation in UTC.
func DateTimeUTC(value time.Time) DateTime {
	value = value.UTC()
	return DateTime{
		Year: uint16(value.Year()), Month: uint8(value.Month()), Day: uint8(value.Day()),
		Hour: uint8(value.Hour()), Minute: uint8(value.Minute()), Second: uint8(value.Second()),
	}
}

// DateTimeUTCString parses YYYY-MM-DD HH:MM:SS as UTC.
func DateTimeUTCString(value string) (DateTime, error) {
	parsed, err := time.ParseInLocation(dateTimeUTCLayout, value, time.UTC)
	if err != nil {
		return DateTime{}, err
	}
	return DateTimeUTC(parsed), nil
}

func (v DateTime) cmBatchValue() batchValue {
	return batchValue{
		kind: 1, year: v.Year, month: v.Month, day: v.Day,
		hour: v.Hour, minute: v.Minute, second: v.Second,
	}
}
