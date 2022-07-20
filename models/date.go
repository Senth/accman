package models

import "time"

type Date string

// DateNow Get today's date
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
