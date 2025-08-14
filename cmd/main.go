package main

import (
	"ai_agent/internal/constants/model/dto"
	agentHandler "ai_agent/internal/handler/agent"
	"ai_agent/internal/service/agent"
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

	"go.uber.org/zap"
)

func main() {
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
	calendarService := calendar.InitCalendar(config, logger)
	emailService := email.InitEmail(config, logger)
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
	mux.HandleFunc("POST /api/events", handler.GetEvents)
	mux.HandleFunc("POST /api/reminder", handler.SendDailyReminder)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("AI Executive Assistant is running"))
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
		ServerPort:             getEnv("SERVER_PORT", "8080"),
		FromEmail:              getEnv("FROM_EMAIL", "assistant@example.com"),
		FromName:               getEnv("FROM_NAME", "AI Executive Assistant"),
		UserEmail:              getEnv("USER_EMAIL", "user@example.com"),
		TimeZone:               getEnv("TIMEZONE", "UTC"),
		DailyReminderTime:      getEnv("DAILY_REMINDER_TIME", "09:00"),
		MeetingReminderMinutes: 15,
	}

	// Validate required configuration
	if config.GoogleCalendarAPIKey == "" {
		log.Fatal("GOOGLE_CALENDAR_API_KEY is required")
	}
	if config.SendGridAPIKey == "" {
		log.Fatal("SENDGRID_API_KEY is required")
	}
	if config.GeminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY is required")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
