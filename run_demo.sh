#!/bin/bash

echo "🚀 Starting AI Executive Assistant in DEMO MODE"
echo "================================================"
echo ""
echo "This will start the server without requiring API keys."
echo "The application will run in demo mode and show appropriate messages."
echo ""
echo "To get full functionality, set the following environment variables:"
echo "  - GOOGLE_CALENDAR_API_KEY"
echo "  - SENDGRID_API_KEY"
echo "  - GEMINI_API_KEY"
echo ""
echo "Starting server on http://localhost:8080"
echo "Visit http://localhost:8080/demo for API information"
echo ""

# Run the application
go run cmd/main.go

