package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Test the API endpoints
func main() {
	baseURL := "http://localhost:8080"

	// Test 1: Health check
	fmt.Println("=== Testing Health Check ===")
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Health check response: %s\n", string(body))

	// Test 2: Get upcoming events
	fmt.Println("\n=== Testing Get Events ===")
	resp, err = http.Get(baseURL + "/api/events")
	if err != nil {
		fmt.Printf("Get events failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Events response: %s\n", string(body))

	// Test 3: Send email
	fmt.Println("\n=== Testing Send Email ===")
	emailData := map[string]interface{}{
		"to_email": "test@example.com",
		"subject":  "Test Email from AI Assistant",
		"body":     "This is a test email sent by the AI Executive Assistant.",
	}
	
	jsonData, _ := json.Marshal(emailData)
	resp, err = http.Post(baseURL+"/api/email", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Send email failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Send email response: %s\n", string(body))

	// Test 4: Natural language command
	fmt.Println("\n=== Testing Natural Language Command ===")
	commandData := map[string]interface{}{
		"command": "What meetings do I have today?",
	}
	
	jsonData, _ = json.Marshal(commandData)
	resp, err = http.Post(baseURL+"/api/command", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Process command failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Command response: %s\n", string(body))

	// Test 5: Schedule meeting
	fmt.Println("\n=== Testing Schedule Meeting ===")
	meetingData := map[string]interface{}{
		"attendees":      []string{"test@example.com"},
		"start_time":     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"duration_minutes": 30,
		"title":          "Test Meeting",
	}
	
	jsonData, _ = json.Marshal(meetingData)
	resp, err = http.Post(baseURL+"/api/schedule", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Schedule meeting failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Schedule meeting response: %s\n", string(body))

	fmt.Println("\n=== All tests completed ===")
}
