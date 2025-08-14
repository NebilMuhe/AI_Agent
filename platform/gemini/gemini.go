package gemini

import (
	"ai_agent/internal/constants/model/dto"
	"ai_agent/platform"
	"ai_agent/platform/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type gemini struct {
	config dto.Config
	client *http.Client
	logger logger.Logger
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

func InitGemini(config dto.Config, logger logger.Logger) platform.Gemini {
	return &gemini{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
		logger: logger,
	}
}

func (g *gemini) ProcessCommand(ctx context.Context, command string) (string, error) {
	g.logger.Info(ctx, "Executing Gemini command", zap.String("command", command))

	// Prepare the request data
	requestData := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: command},
				},
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		g.logger.Error(ctx, "Failed to marshal request data", zap.Error(err))
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Build the API URL
	fmt.Println("apikey=============================================", g.config.GeminiAPIKey)
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s", g.config.GeminiAPIKey)

	// Create the request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		g.logger.Error(ctx, "Failed to create request", zap.Error(err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := g.client.Do(req)
	if err != nil {
		g.logger.Error(ctx, "Failed to execute Gemini command", zap.Error(err))
		return "", fmt.Errorf("failed to execute Gemini command: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		g.logger.Error(ctx, "Gemini API returned error status", zap.Int("status", resp.StatusCode))
		return "", fmt.Errorf("gemini API returned status: %d", resp.StatusCode)
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		g.logger.Error(ctx, "Failed to decode response", zap.Error(err))
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		g.logger.Error(ctx, "No response from Gemini")
		return "", fmt.Errorf("no response from Gemini")
	}

	response := geminiResp.Candidates[0].Content.Parts[0].Text
	g.logger.Info(ctx, "Successfully executed Gemini command", zap.Int("response_length", len(response)))

	return response, nil
}
