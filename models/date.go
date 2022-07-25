package models

import "time"

type Date string

func DateFromTime(t time.Time) Date {
	return Date(t.Format("2006-01-02"))
}

// DateNow InLocalCurrency today's date
func DateNow() Date {
	return Date(time.Now().Format("2006-01-02"))
}

func (d Date) Before(other Date) bool {
	return d < other
}

func (d Date) BeforeOrEqual(other Date) bool {
	return d <= other
}

func (d Date) After(other Date) bool {
	return d > other
}

func (d Date) AfterOrEqual(other Date) bool {
	return d >= other
}

func (d Date) Between(start, end Date) bool {
	return d.AfterOrEqual(start) && d.BeforeOrEqual(end)
}

func (d Date) Year() string {
	return string(d)[:4]
}

func (d Date) Month() string {
	return string(d)[5:7]
}

func (d Date) Day() string {
	return string(d)[8:]
}

func (d Date) SetYear(year string) Date {
	return Date(year + string(d[4:]))
}

func (d Date) SetMonth(month string) Date {
	return Date(string(d[:4]) + "-" + month + string(d[7:]))
}

func (d Date) SetDay(day string) Date {
	return Date(string(d[:8]) + day)
}

func (d Date) AddYears(years int) Date {
	return DateFromTime(d.Time().AddDate(years, 0, 0))
}

func (d Date) AddMonths(months int) Date {
	return DateFromTime(d.Time().AddDate(0, months, 0))
}

func (d Date) AddDays(days int) Date {
	return DateFromTime(d.Time().AddDate(0, 0, days))
}

func (d Date) Time() time.Time {
	t, err := time.Parse("2006-01-02", string(d))
	if err != nil {
		panic(err)
	}
	return t
}
