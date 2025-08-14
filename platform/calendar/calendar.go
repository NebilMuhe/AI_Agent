package calendar

import (
	"ai_agent/internal/constants/model/dto"
	"ai_agent/platform"
	"ai_agent/platform/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type calendar struct {
	config dto.Config
	client *http.Client
	logger logger.Logger
}

type GoogleCalendarEvent struct {
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	Start       struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"end"`
	Attendees []struct {
		Email string `json:"email"`
	} `json:"attendees,omitempty"`
}

type GoogleCalendarResponse struct {
	Items []GoogleCalendarEvent `json:"items"`
}

func InitCalendar(config dto.Config, logger logger.Logger) platform.Calendar {
	return &calendar{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
		logger: logger,
	}
}

func (c *calendar) GetUpcomingEvents(ctx context.Context) ([]dto.Event, error) {
	c.logger.Info(ctx,"Fetching upcoming events from Google Calendar")

	// Build the API URL
	baseURL := "https://www.googleapis.com/calendar/v3/calendars/primary/events"
	params := url.Values{}
	params.Add("key", c.config.GoogleCalendarAPIKey)
	params.Add("timeMin", time.Now().Format(time.RFC3339))
	params.Add("timeMax", time.Now().AddDate(0, 0, 7).Format(time.RFC3339)) // Next 7 days
	params.Add("singleEvents", "true")
	params.Add("orderBy", "startTime")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		c.logger.Error(ctx, "Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(ctx, "Failed to fetch events", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error(ctx, "Calendar API returned error status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("calendar API returned status: %d", resp.StatusCode)
	}

	var calendarResp GoogleCalendarResponse
	if err := json.NewDecoder(resp.Body).Decode(&calendarResp); err != nil {
		c.logger.Error(ctx, "Failed to decode response", zap.Error(err))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our Event format
	events := make([]dto.Event, 0, len(calendarResp.Items))
	for _, item := range calendarResp.Items {
		startTime, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, item.End.DateTime)

		attendees := make([]string, 0, len(item.Attendees))
		for _, attendee := range item.Attendees {
			attendees = append(attendees, attendee.Email)
		}

		events = append(events, dto.Event{
			Title:     item.Summary,
			Attendees: attendees,
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	c.logger.Info(ctx, "Successfully fetched events", zap.Int("count", len(events)))
	return events, nil
}

// ScheduleMeeting implements platform.Calendar.
func (c *calendar) ScheduleMeeting(ctx context.Context, attendees []string, startTime time.Time, duration time.Duration, title string) error {
	c.logger.Info(ctx, "Scheduling meeting", zap.String("title", title), zap.Time("startTime", startTime), zap.Strings("attendees", attendees))

	endTime := startTime.Add(duration)

	// Prepare the event data
	event := GoogleCalendarEvent{
		Summary: title,
		Start: struct {
			DateTime string `json:"dateTime"`
			TimeZone string `json:"timeZone"`
		}{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: c.config.TimeZone,
		},
		End: struct {
			DateTime string `json:"dateTime"`
			TimeZone string `json:"timeZone"`
		}{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: c.config.TimeZone,
		},
	}

	// Add attendees
	for _, attendee := range attendees {
		event.Attendees = append(event.Attendees, struct {
			Email string `json:"email"`
		}{Email: attendee})
	}

	// Convert to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		c.logger.Error(ctx, "Failed to marshal event data", zap.Error(err))
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Build the API URL
	baseURL := "https://www.googleapis.com/calendar/v3/calendars/primary/events"
	params := url.Values{}
	fmt.Println("===================================",c.config.GoogleCalendarAPIKey)
	params.Add("key", c.config.GoogleCalendarAPIKey)
	params.Add("sendUpdates", "all") // Send invitations to attendees

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create the request
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(eventData))
	if err != nil {
		c.logger.Error(ctx, "Failed to create request", zap.Error(err))
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(ctx, "Failed to schedule meeting", zap.Error(err))
		return fmt.Errorf("failed to schedule meeting: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.logger.Error(ctx, "Calendar API returned error status", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("calendar API returned status: %d", resp.StatusCode)
	}

	c.logger.Info(ctx, "Successfully scheduled meeting", zap.String("title", title))
	return nil
}
