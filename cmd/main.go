package main

import (
	"ai_agent/internal/constants/model/dto"
	"ai_agent/internal/constants/model/response"
	agentHandler "ai_agent/internal/handler/agent"
	"ai_agent/internal/service/agent"
	"ai_agent/platform"
	"ai_agent/platform/calendar"
	"ai_agent/platform/email"
	"ai_agent/platform/gemini"
	"ai_agent/platform/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Load configuration
	config := loadConfig()

	// Initialize logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()

	logger := logger.InitLogger(zapLogger)
	logger.Info(context.Background(), "Starting AI Executive Assistant")

	// Initialize platform services
	var calendarService platform.Calendar
	if config.GoogleCalendarAPIKey != "" && config.GoogleCalendarAPIKey != "your_google_calendar_api_key_here" && len(config.GoogleCalendarAPIKey) > 30 {
		calendarService = calendar.InitCalendar(config, logger)
	} else {
		log.Println("⚠️  Using Simple Calendar Mode (no valid Google Calendar API key)")
		calendarService = calendar.InitSimpleCalendar(config, logger)
	}

	// Initialize email service (SendGrid or Gmail)
	var emailService platform.Email
	if config.SendGridAPIKey != "" {
		emailService = email.InitEmail(config, logger)
	} else if config.GmailAppPassword != "" {
		emailService = email.InitGmail(config, logger)
	} else {
		// Use a mock email service for demo mode
		emailService = email.InitEmail(config, logger) // This will show errors but won't crash
	}

	geminiService := gemini.InitGemini(config, logger)

	// Initialize business service
	service := agent.NewService(calendarService, emailService, geminiService, logger, config)

	// Initialize HTTP handler
	handler := agentHandler.NewHandler(service, logger)

	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/command", handler.ProcessCommand)
	mux.HandleFunc("POST /api/schedule", handler.ScheduleMeeting)
	mux.HandleFunc("POST /api/email", handler.SendEmail)
	mux.HandleFunc("GET /api/events", handler.GetEvents)
	mux.HandleFunc("POST /api/reminder", handler.SendDailyReminder)

	// Add health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.SendSuccessResponse(w, http.StatusOK, "AI Executive Assistant is running")
	})

	// Add demo endpoint for testing without API keys
	mux.HandleFunc("GET /demo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "AI Executive Assistant Demo Mode",
			"status": "running",
			"endpoints": {
				"health": "GET /health",
				"command": "POST /api/command",
				"schedule": "POST /api/schedule", 
				"email": "POST /api/email",
				"events": "GET /api/events",
				"reminder": "POST /api/reminder"
			},
			"note": "Set API keys in environment variables for full functionality"
		}`))
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting HTTP server", zap.String("port", config.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "Server error", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "Server forced to shutdown", zap.Error(err))
	}

	logger.Info(context.Background(), "Server exited")
}

func loadConfig() dto.Config {
	config := dto.Config{
		GoogleCalendarAPIKey:   getEnv("GOOGLE_CALENDAR_API_KEY", ""),
		SendGridAPIKey:         getEnv("SENDGRID_API_KEY", ""),
		GeminiAPIKey:           getEnv("GEMINI_API_KEY", ""),
		GmailAppPassword:       getEnv("GMAIL_APP_PASSWORD", ""),
		ServerPort:             getEnv("SERVER_PORT", "8080"),
		FromEmail:              getEnv("FROM_EMAIL", "assistant@example.com"),
		FromName:               getEnv("FROM_NAME", "AI Executive Assistant"),
		UserEmail:              getEnv("USER_EMAIL", "user@example.com"),
		TimeZone:               getEnv("TIMEZONE", "UTC"),
		DailyReminderTime:      getEnv("DAILY_REMINDER_TIME", "09:00"),
		MeetingReminderMinutes: 15,
	}

	// Check if we're in demo mode (no API keys provided)
	if config.GoogleCalendarAPIKey == "" && config.SendGridAPIKey == "" && config.GeminiAPIKey == "" && config.GmailAppPassword == "" {
		log.Println("⚠️  Running in DEMO MODE - No API keys provided")
		log.Println("   Set the following environment variables for full functionality:")
		log.Println("   - GOOGLE_CALENDAR_API_KEY")
		log.Println("   - SENDGRID_API_KEY (or GMAIL_APP_PASSWORD for Gmail SMTP)")
		log.Println("   - GEMINI_API_KEY")
		log.Println("   Visit http://localhost:8080/demo for more info")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
