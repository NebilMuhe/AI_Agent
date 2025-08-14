#!/bin/bash

echo "🧪 Testing AI Executive Assistant API Endpoints"
echo "=============================================="
echo ""

# Base URL
BASE_URL="http://localhost:8080"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to test endpoint
test_endpoint() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    
    echo -e "${BLUE}Testing: ${name}${NC}"
    echo -e "${YELLOW}curl -X ${method} ${BASE_URL}${endpoint}${NC}"
    
    if [ -n "$data" ]; then
        echo -e "${YELLOW}Data: ${data}${NC}"
        response=$(curl -s -X ${method} \
            -H "Content-Type: application/json" \
            -d "${data}" \
            "${BASE_URL}${endpoint}")
    else
        response=$(curl -s -X ${method} "${BASE_URL}${endpoint}")
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Success${NC}"
        echo "Response: $response"
    else
        echo -e "${RED}❌ Failed${NC}"
    fi
    echo ""
}

# Test 1: Health Check
test_endpoint "Health Check" "GET" "/health"

# Test 2: Demo Info
test_endpoint "Demo Info" "GET" "/demo"

# Test 3: Get Events
test_endpoint "Get Events" "GET" "/api/events"

# Test 4: Process Natural Language Command
test_endpoint "Natural Language Command" "POST" "/api/command" '{"command": "What meetings do I have today?"}'

# Test 5: Schedule Meeting
test_endpoint "Schedule Meeting" "POST" "/api/schedule" '{
  "attendees": ["test@example.com"],
  "start_time": "2024-12-15T14:00:00Z",
  "duration_minutes": 30,
  "title": "Test Meeting"
}'

# Test 6: Send Email
test_endpoint "Send Email" "POST" "/api/email" '{
  "to_email": "test@example.com",
  "subject": "Test Email from AI Assistant",
  "body": "This is a test email sent by the AI Executive Assistant."
}'

# Test 7: Send Daily Reminder
test_endpoint "Send Daily Reminder" "POST" "/api/reminder"

echo -e "${GREEN}🎉 All API tests completed!${NC}"
echo ""
echo -e "${YELLOW}📋 Individual curl commands for manual testing:${NC}"
echo "=============================================="
echo ""
echo "1. Health Check:"
echo "curl -X GET http://localhost:8080/health"
echo ""
echo "2. Demo Info:"
echo "curl -X GET http://localhost:8080/demo"
echo ""
echo "3. Get Events:"
echo "curl -X GET http://localhost:8080/api/events"
echo ""
echo "4. Natural Language Command:"
echo 'curl -X POST http://localhost:8080/api/command -H "Content-Type: application/json" -d '"'"'{"command": "What meetings do I have today?"}'"'"
echo ""
echo "5. Schedule Meeting:"
echo 'curl -X POST http://localhost:8080/api/schedule -H "Content-Type: application/json" -d '"'"'{"attendees": ["test@example.com"], "start_time": "2024-12-15T14:00:00Z", "duration_minutes": 30, "title": "Test Meeting"}'"'"
echo ""
echo "6. Send Email:"
echo 'curl -X POST http://localhost:8080/api/email -H "Content-Type: application/json" -d '"'"'{"to_email": "test@example.com", "subject": "Test Email", "body": "Test message"}'"'"
echo ""
echo "7. Send Daily Reminder:"
echo "curl -X POST http://localhost:8080/api/reminder"

