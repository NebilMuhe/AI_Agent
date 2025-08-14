package dto

import "time"

type Event struct {
	Title     string
	Attendees []string
	StartTime time.Time
	EndTime   time.Time
}
