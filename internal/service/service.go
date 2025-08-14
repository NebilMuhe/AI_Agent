package service

import (
	"ai_agent/internal/constants/model/dto"
	"context"
	"time"
)

type AgentService interface {
	ProcessNaturalLanguageCommand(ctx context.Context, command string) (string, error)
	ScheduleMeeting(ctx context.Context, attendees []string,
		startTime time.Time, duration time.Duration, title string) error
	SendEmail(ctx context.Context, toEmail string, subject string,
		body string) error
	GetUpcomingEvents(ctx context.Context) ([]dto.Event, error)
	SendDailyReminder(ctx context.Context) error
}
