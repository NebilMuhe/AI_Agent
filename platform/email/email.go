package email

import (
	"ai_agent/internal/constants/model/dto"
	"ai_agent/platform"
	"ai_agent/platform/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
)

type email struct {
	config dto.Config
	client *http.Client
	logger logger.Logger
}

type SendGridEmail struct {
	Personalizations []Personalization `json:"personalizations"`
	From             From              `json:"from"`
	Subject          string            `json:"subject"`
	Content          []Content         `json:"content"`
}

type Personalization struct {
	To []To `json:"to"`
}

type To struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type From struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func InitEmail(config dto.Config, logger logger.Logger) platform.Email {
	return &email{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
		logger: logger,
	}
}

// SendEmail implements platform.Email.
func (e *email) SendEmail(ctx context.Context, toEmail string, subject string, body string) error {
	e.logger.Info(ctx, "Sending email", zap.String("to", toEmail), zap.String("subject", subject))

	// Prepare the email data
	emailData := SendGridEmail{
		Personalizations: []Personalization{
			{
				To: []To{
					{Email: toEmail},
				},
			},
		},
		From: From{
			Email: e.config.FromEmail,
			Name:  e.config.FromName,
		},
		Subject: subject,
		Content: []Content{
			{
				Type:  "text/html",
				Value: body,
			},
			{
				Type:  "text/plain",
				Value: stripHTML(body),
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(emailData)
	if err != nil {
		e.logger.Error(ctx, "Failed to marshal email data", zap.Error(err))
		return fmt.Errorf("failed to marshal email data: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonData))
	if err != nil {
		e.logger.Error(ctx, "Failed to create request", zap.Error(err))
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+e.config.SendGridAPIKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := e.client.Do(req)
	if err != nil {
		e.logger.Error(ctx, "Failed to send email", zap.Error(err))
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		e.logger.Error(ctx, "SendGrid API returned error status", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("sendgrid API returned status: %d", resp.StatusCode)
	}

	e.logger.Info(ctx, "Successfully sent email", zap.String("to", toEmail), zap.String("subject", subject))
	return nil
}

// stripHTML removes HTML tags from text for plain text version
func stripHTML(html string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, "")

	// Clean up extra whitespace
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text
}
