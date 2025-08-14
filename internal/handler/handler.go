package handler

import "net/http"

type Agent interface {
	ProcessCommand(w http.ResponseWriter, r *http.Request)
	ScheduleMeeting(w http.ResponseWriter, r *http.Request)
	SendEmail(w http.ResponseWriter, r *http.Request)
	GetEvents(w http.ResponseWriter, r *http.Request)
	SendDailyReminder(w http.ResponseWriter, r *http.Request)
}
