package calendar

import (
	"ai_agent/internal/constants/model/dto"
	"ai_agent/platform"
	"ai_agent/platform/logger"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type simpleCalendar struct {
	config dto.Config
	client *http.Client
	logger logger.Logger
}

func InitSimpleCalendar(config dto.Config, logger logger.Logger) platform.Calendar {
	return &simpleCalendar{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
		logger: logger,
	}
}

// GetUpcomingEvents implements platform.Calendar.
func (c *simpleCalendar) GetUpcomingEvents(ctx context.Context) ([]dto.Event, error) {
	c.logger.Info(ctx, "Getting upcoming events (simple mode)")
	
	// Return mock events for demo purposes
	events := []dto.Event{
		{
			Title:     "Team Meeting",
			Attendees: []string{"team@example.com"},
			StartTime: time.Now().Add(2 * time.Hour),
			EndTime:   time.Now().Add(3 * time.Hour),
		},
		{
			Title:     "Client Call",
			Attendees: []string{"client@example.com"},
			StartTime: time.Now().Add(24 * time.Hour),
			EndTime:   time.Now().Add(25 * time.Hour),
		},
	}
	
	c.logger.Info(ctx, "Returning mock events", zap.Int("count", len(events)))
	return events, nil
}

// ScheduleMeeting implements platform.Calendar.
func (c *simpleCalendar) ScheduleMeeting(ctx context.Context, attendees []string, startTime time.Time, duration time.Duration, title string) error {
	c.logger.Info(ctx, "Scheduling meeting (simple mode)", 
		zap.String("title", title), 
		zap.Time("startTime", startTime), 
		zap.Strings("attendees", attendees))
	
	// In simple mode, just log the meeting details
	c.logger.Info(ctx, "Meeting would be scheduled", 
		zap.String("title", title),
		zap.String("start", startTime.Format("2006-01-02 15:04:05")),
		zap.String("duration", duration.String()),
		zap.Strings("attendees", attendees))
	
	// Return success (mock implementation)
	return nil
}
