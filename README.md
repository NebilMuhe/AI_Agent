# AI Executive Assistant

An intelligent AI-powered executive assistant that automates scheduling, email management, and task reminders using natural language processing.

## Features

- **Natural Language Processing**: Understand and execute commands in plain English
- **Meeting Scheduling**: Automatically schedule meetings with Google Calendar integration
- **Email Management**: Send emails with AI-generated content using SendGrid
- **Daily Reminders**: Automated daily schedule summaries
- **RESTful API**: Easy integration with existing systems

## Architecture

The application follows a clean architecture pattern with the following layers:

- **Platform Layer**: External service integrations (Google Calendar, SendGrid, Gemini AI)
- **Service Layer**: Business logic and orchestration
- **Handler Layer**: HTTP request handling and response formatting
- **DTO Layer**: Data transfer objects and configuration

## Prerequisites

- Go 1.21 or higher
- Google Calendar API key
- SendGrid API key
- Gemini AI API key

## Setup Instructions

### 1. Get API Keys

#### Google Calendar API
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable the Google Calendar API
4. Create credentials (API Key)
5. Set up OAuth 2.0 for calendar access

#### SendGrid API
1. Sign up for a free SendGrid account
2. Navigate to Settings > API Keys
3. Create a new API key with "Mail Send" permissions
4. Verify your sender email address

#### Gemini AI API
1. Go to [Google AI Studio](https://ai.google.dev/)
2. Sign in with your Google account
3. Click "Get API Key" and create a new key
4. The free tier allows 15 requests/minute

### 2. Environment Variables

Create a `.env` file in the project root:

```bash
# API Keys
GOOGLE_CALENDAR_API_KEY=your_google_calendar_api_key
SENDGRID_API_KEY=your_sendgrid_api_key
GEMINI_API_KEY=your_gemini_api_key

# Server Configuration
SERVER_PORT=8080
LOG_LEVEL=info

# Email Configuration
FROM_EMAIL=your_verified_sender@example.com
FROM_NAME=AI Executive Assistant
USER_EMAIL=your_email@example.com

# Calendar Configuration
TIMEZONE=America/New_York
DAILY_REMINDER_TIME=09:00
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Process Natural Language Command
**POST** `/api/command`

Process natural language commands like "Schedule a meeting with John tomorrow at 2 PM"

```json
{
  "command": "Schedule a meeting with john@example.com tomorrow at 2 PM for 30 minutes to discuss project updates"
}
```

Response:
```json
{
  "result": "Meeting scheduled successfully!",
  "error": null
}
```

### 2. Schedule Meeting
**POST** `/api/schedule`

Schedule a meeting with specific details

```json
{
  "attendees": ["john@example.com", "jane@example.com"],
  "start_time": "2024-01-15T14:00:00Z",
  "duration_minutes": 30,
  "title": "Project Update Meeting"
}
```

### 3. Send Email
**POST** `/api/email`

Send an email with optional AI-generated content

```json
{
  "to_email": "recipient@example.com",
  "subject": "Follow-up on our meeting",
  "body": "Optional email body (leave empty for AI generation)"
}
```

### 4. Get Upcoming Events
**GET** `/api/events`

Retrieve upcoming calendar events

Response:
```json
{
  "events": [
    {
      "title": "Team Meeting",
      "attendees": ["john@example.com"],
      "start_time": "2024-01-15T10:00:00Z",
      "end_time": "2024-01-15T11:00:00Z"
    }
  ],
  "error": null
}
```

### 5. Send Daily Reminder
**POST** `/api/reminder`

Trigger a daily reminder email with upcoming events

### 6. Health Check
**GET** `/health`

Check if the service is running

## Example Usage

### Scheduling a Meeting
```bash
curl -X POST http://localhost:8080/api/command \
  -H "Content-Type: application/json" \
  -d '{
    "command": "Schedule a meeting with john@example.com and jane@example.com tomorrow at 3 PM for 1 hour to discuss Q1 planning"
  }'
```

### Sending an Email
```bash
curl -X POST http://localhost:8080/api/email \
  -H "Content-Type: application/json" \
  -d '{
    "to_email": "client@example.com",
    "subject": "Project Status Update",
    "body": ""
  }'
```

### Getting Events
```bash
curl -X GET http://localhost:8080/api/events
```

## Free Tier Limits

- **Google Calendar API**: 1,000,000 requests/day
- **SendGrid**: 100 emails/day (free tier)
- **Gemini AI**: 15 requests/minute, 1,500 requests/day
- **GCP Always Free Tier**: 2 million Cloud Function invocations/month

## Deployment

### Local Development
```bash
go run cmd/main.go
```

### Docker Deployment
```bash
# Build the image
docker build -t ai-executive-assistant .

# Run the container
docker run -p 8080:8080 --env-file .env ai-executive-assistant
```

### Cloud Run Deployment (GCP Free Tier)
```bash
# Build and deploy to Cloud Run
gcloud run deploy ai-executive-assistant \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "GOOGLE_CALENDAR_API_KEY=$GOOGLE_CALENDAR_API_KEY,SENDGRID_API_KEY=$SENDGRID_API_KEY,GEMINI_API_KEY=$GEMINI_API_KEY"
```

## Project Structure

```
AI_Agent/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── constants/
│   │   └── model/
│   │       ├── dto/            # Data transfer objects
│   │       └── response/       # Response models
│   ├── handler/                # HTTP handlers
│   └── service/                # Business logic
├── platform/                   # External service integrations
│   ├── calendar/               # Google Calendar integration
│   ├── email/                  # SendGrid integration
│   ├── gemini/                 # Gemini AI integration
│   └── logger/                 # Logging
├── go.mod                      # Go module file
├── go.sum                      # Go module checksums
└── README.md                   # This file
```

