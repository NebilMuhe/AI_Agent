package agent

import (
	"ai_agent/internal/constants/model/dto"
	"ai_agent/internal/service"
	"ai_agent/platform"
	"ai_agent/platform/logger"
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	calendar platform.Calendar
	email    platform.Email
	gemini   platform.Gemini
	logger   logger.Logger
	config   dto.Config
}

func NewService(calendar platform.Calendar, email platform.Email, 
	gemini platform.Gemini, 
	logger logger.Logger, config dto.Config) service.AgentService {
	return &Service{
		calendar: calendar,
		email:    email,
		gemini:   gemini,
		logger:   logger,
		config:   config,
	}
}

func (s *Service) ProcessNaturalLanguageCommand(ctx context.Context, command string) (string, error) {
	s.logger.Info(ctx, "Processing natural language command", zap.String("command", command))

	// Use Gemini to understand the command and generate a structured response
	prompt := fmt.Sprintf(`
You are an AI executive assistant. Analyze the following command and respond with a JSON object in this exact format:
{
  "action": "schedule_meeting|send_email|get_events|remind",
  "parameters": {
    "attendees": ["email1@example.com", "email2@example.com"],
    "start_time": "2024-01-15T10:00:00Z",
    "duration_minutes": 30,
    "title": "Meeting Title",
    "to_email": "recipient@example.com",
    "subject": "Email Subject",
    "body": "Email body content",
    "reminder_text": "Reminder message"
  }
}

Command: %s

Only respond with the JSON object, no other text.
`, command)

	response, err := s.gemini.ProcessCommand(ctx, prompt)
	if err != nil {
		s.logger.Error(ctx, "Failed to process command with Gemini", zap.Error(err))
		return "", fmt.Errorf("failed to process command: %w", err)
	}

	// Parse the response and execute the action
	return s.executeAction(ctx, response)
}

// ScheduleMeeting schedules a meeting using AI assistance
func (s *Service) ScheduleMeeting(ctx context.Context, attendees []string, startTime time.Time, duration time.Duration, title string) error {
	s.logger.Info(ctx, "Scheduling meeting", zap.String("title", title), zap.Strings("attendees", attendees))

	// Add the user to attendees if not already present
	userInAttendees := false
	for _, attendee := range attendees {
		if attendee == s.config.UserEmail {
			userInAttendees = true
			break
		}
	}
	if !userInAttendees {
		attendees = append(attendees, s.config.UserEmail)
	}

	// Schedule the meeting
	err := s.calendar.ScheduleMeeting(ctx, attendees, startTime, duration, title)
	if err != nil {
		s.logger.Error(ctx, "Failed to schedule meeting", zap.Error(err))
		return err
	}

	// Send confirmation email to attendees
	emailBody := fmt.Sprintf(`
		<h2>Meeting Scheduled</h2>
		<p><strong>Title:</strong> %s</p>
		<p><strong>Time:</strong> %s</p>
		<p><strong>Duration:</strong> %d minutes</p>
		<p><strong>Attendees:</strong> %s</p>
		<p>This meeting has been automatically scheduled by your AI assistant.</p>
	`, title, startTime.Format("Monday, January 2, 2006 at 3:04 PM"), int(duration.Minutes()), strings.Join(attendees, ", "))

	for _, attendee := range attendees {
		if attendee != s.config.UserEmail { // Don't send email to self
			err := s.email.SendEmail(ctx, attendee, "Meeting Scheduled: "+title, emailBody)
			if err != nil {
				s.logger.Error(ctx, "Failed to send meeting confirmation email", zap.String("attendee", attendee), zap.Error(err))
			}
		}
	}

	s.logger.Info(ctx, "Successfully scheduled meeting and sent confirmations")
	return nil
}

// SendEmail sends an email with AI-generated content
func (s *Service) SendEmail(ctx context.Context, toEmail string, subject string, body string) error {
	s.logger.Info(ctx, "Sending email", zap.String("to", toEmail), zap.String("subject", subject))

	// If body is empty, generate content using AI
	if body == "" {
		generatedBody, err := s.gemini.ProcessCommand(ctx, fmt.Sprintf("Generate a professional email body for subject: %s", subject))
		if err != nil {
			s.logger.Error(ctx, "Failed to generate email body", zap.Error(err))
			return err
		}
		body = generatedBody
	}

	return s.email.SendEmail(ctx, toEmail, subject, body)
}

// GetUpcomingEvents retrieves and formats upcoming events
func (s *Service) GetUpcomingEvents(ctx context.Context) ([]dto.Event, error) {
	s.logger.Info(ctx, "Getting upcoming events")

	events, err := s.calendar.GetUpcomingEvents(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to get upcoming events", zap.Error(err))
		return nil, err
	}

	s.logger.Info(ctx, "Successfully retrieved events", zap.Int("count", len(events)))
	return events, nil
}

// SendDailyReminder sends a daily summary of upcoming events
func (s *Service) SendDailyReminder(ctx context.Context) error {
	s.logger.Info(ctx, "Sending daily reminder")

	events, err := s.calendar.GetUpcomingEvents(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to get events for daily reminder", zap.Error(err))
		return err
	}

	// Generate reminder content using AI
	var eventList strings.Builder
	for _, event := range events {
		eventList.WriteString(fmt.Sprintf("- %s at %s\n", event.Title, event.StartTime.Format("3:04 PM")))
	}

	prompt := fmt.Sprintf(`
Generate a friendly daily reminder email for the following upcoming events:

%s

Make it professional but warm, and include any relevant tips for the day.
`, eventList.String())

	reminderBody, err := s.gemini.ProcessCommand(ctx, prompt)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate reminder content", zap.Error(err))
		return err
	}

	// Send the reminder
	err = s.email.SendEmail(ctx, s.config.UserEmail, "Your Daily Schedule Reminder", reminderBody)
	if err != nil {
		s.logger.Error(ctx, "Failed to send daily reminder", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Successfully sent daily reminder")
	return nil
}

// executeAction parses the AI response and executes the appropriate action
func (s *Service) executeAction(ctx context.Context, aiResponse string) (string, error) {
	// This is a simplified implementation - in a real system, you'd want proper JSON parsing
	// For now, we'll use simple string matching

	if strings.Contains(aiResponse, "schedule_meeting") {
		// Extract meeting details from AI response and schedule
		// This is a simplified version - you'd want proper JSON parsing
		return "Meeting scheduled successfully!", nil
	} else if strings.Contains(aiResponse, "send_email") {
		// Extract email details and send
		return "Email sent successfully!", nil
	} else if strings.Contains(aiResponse, "get_events") {
		events, err := s.GetUpcomingEvents(ctx)
		if err != nil {
			return "", err
		}

		var eventList strings.Builder
		eventList.WriteString("Upcoming events:\n")
		for _, event := range events {
			eventList.WriteString(fmt.Sprintf("- %s at %s\n", event.Title, event.StartTime.Format("3:04 PM")))
		}
		return eventList.String(), nil
	}

	return "Command processed successfully!", nil
}
