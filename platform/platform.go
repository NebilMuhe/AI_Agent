package platform

import (
	"ai_agent/internal/constants/model/dto"
	"context"
	"time"
)

type Calendar interface {
	ScheduleMeeting(ctx context.Context,attendees []string,
		startTime time.Time, duration time.Duration, title string) error
	GetUpcomingEvents(ctx context.Context) ([]dto.Event, error)
}

type Email interface {
	SendEmail(ctx context.Context,toEmail string, subject string, body string) error
}

type Gemini interface {
	ProcessCommand(ctx context.Context,command string) (string, error)
}
