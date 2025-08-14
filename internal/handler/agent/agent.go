package agent

import (
	"ai_agent/internal/constants/model/dto"
	"ai_agent/internal/service"
	"ai_agent/internal/handler"
	"ai_agent/platform/logger"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type agentHandler struct {
	service service.AgentService
	logger  logger.Logger
}

func NewHandler(service service.AgentService, logger logger.Logger) handler.Agent {
	return &agentHandler{
		service: service,
		logger:  logger,
	}
}

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

type MeetingRequest struct {
	Attendees []string `json:"attendees"`
	StartTime string   `json:"start_time"`
	Duration  int      `json:"duration_minutes"`
	Title     string   `json:"title"`
}

type EmailRequest struct {
	ToEmail string `json:"to_email"`
	Subject string `json:"subject"`
	Body    string `json:"body,omitempty"`
}

type EventsResponse struct {
	Events []dto.Event `json:"events"`
	Error  string      `json:"error,omitempty"`
}

// ProcessCommand handles natural language commands
func (h *agentHandler) ProcessCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "Failed to decode command request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.ProcessNaturalLanguageCommand(r.Context(), req.Command)

	response := CommandResponse{
		Result: result,
	}
	if err != nil {
		response.Error = err.Error()
		h.logger.Error(r.Context(), "Failed to process command", zap.Error(err))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ScheduleMeeting handles meeting scheduling requests
func (h *agentHandler) ScheduleMeeting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MeetingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "Failed to decode meeting request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		h.logger.Error(r.Context(), "Failed to parse start time", zap.Error(err))
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	duration := time.Duration(req.Duration) * time.Minute
	err = h.service.ScheduleMeeting(r.Context(), req.Attendees, startTime, duration, req.Title)

	response := CommandResponse{}
	if err != nil {
		response.Error = err.Error()
		h.logger.Error(r.Context(), "Failed to schedule meeting", zap.Error(err))
	} else {
		response.Result = "Meeting scheduled successfully!"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SendEmail handles email sending requests
func (h *agentHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "Failed to decode email request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.SendEmail(r.Context(), req.ToEmail, req.Subject, req.Body)

	response := CommandResponse{}
	if err != nil {
		response.Error = err.Error()
		h.logger.Error(r.Context(), "Failed to send email", zap.Error(err))
	} else {
		response.Result = "Email sent successfully!"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEvents retrieves upcoming events
func (h *agentHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	events, err := h.service.GetUpcomingEvents(r.Context())

	response := EventsResponse{
		Events: events,
	}
	if err != nil {
		response.Error = err.Error()
		h.logger.Error(r.Context(), "Failed to get events", zap.Error(err))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SendDailyReminder triggers a daily reminder
func (h *agentHandler) SendDailyReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.service.SendDailyReminder(r.Context())

	response := CommandResponse{}
	if err != nil {
		response.Error = err.Error()
		h.logger.Error(r.Context(), "Failed to send daily reminder", zap.Error(err))
	} else {
		response.Result = "Daily reminder sent successfully!"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}