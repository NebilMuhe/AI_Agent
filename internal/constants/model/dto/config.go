package dto

type Config struct {
	GoogleCalendarAPIKey string
	SendGridAPIKey       string
	GeminiAPIKey         string
	ServerPort           string
	GmailAppPassword     string

	GoogleCalendarURL string
	SendGridURL       string
	GeminiURL         string

	FromEmail string
	FromName  string
	UserEmail string

	CalendarID string
	TimeZone   string

	DailyReminderTime      string
	MeetingReminderMinutes int
}
