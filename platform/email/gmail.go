package email

import (
	"ai_agent/internal/constants/model/dto"
	"ai_agent/platform"
	"ai_agent/platform/logger"
	"context"
	"fmt"
	"net/smtp"

	"go.uber.org/zap"
)

type gmail struct {
	config dto.Config
	logger logger.Logger
}

func InitGmail(config dto.Config, logger logger.Logger) platform.Email {
	return &gmail{
		config: config,
		logger: logger,
	}
}

// SendEmail implements platform.Email using Gmail SMTP
func (g *gmail) SendEmail(ctx context.Context, toEmail string, subject string, body string) error {
	g.logger.Info(ctx, "Sending email via Gmail SMTP", zap.String("to", toEmail), zap.String("subject", subject))

	// Gmail SMTP configuration
	from := g.config.FromEmail
	password := g.config.GmailAppPassword // You'll need to set this in config

	// Gmail SMTP server
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Create message with proper MIME headers for HTML
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", from, toEmail, subject, body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, []byte(message))
	if err != nil {
		g.logger.Error(ctx, "Failed to send email via Gmail", zap.Error(err))
		return fmt.Errorf("failed to send email: %w", err)
	}

	g.logger.Info(ctx, "Successfully sent email via Gmail", zap.String("to", toEmail), zap.String("subject", subject))
	return nil
}
