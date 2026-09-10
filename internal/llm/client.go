package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	defaultOpenAIEndpoint = "https://api.openai.com/v1/chat/completions"
	defaultOpenRouterEndpoint = "https://openrouter.ai/api/v1/chat/completions"
	defaultModel = "gpt-4o-mini"
	openRouterModel = "openai/gpt-4o-mini"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Client struct {
	APIKey       string
	Model        string
	Endpoint     string
	IsOpenRouter bool
}

func NewClient() *Client {
	// 1. Check OPENROUTER_API_KEY first
	key := os.Getenv("OPENROUTER_API_KEY")
	isOpenRouter := key != ""
	
	if isOpenRouter {
		// Use OpenRouter defaults
		model := os.Getenv("STYLE_ENGINE_MODEL")
		if model == "" {
			model = openRouterModel
		}
		endpoint := os.Getenv("STYLE_ENGINE_ENDPOINT")
		if endpoint == "" {
			endpoint = defaultOpenRouterEndpoint
		}
		return &Client{
			APIKey:       key,
			Model:        model,
			Endpoint:     endpoint,
			IsOpenRouter: true,
		}
	}

	// 2. Fall back to OPENAI_API_KEY
	key = os.Getenv("OPENAI_API_KEY")
	if key == "" {
		key = os.Getenv("STYLE_ENGINE_API_KEY")
	}

	model := os.Getenv("STYLE_ENGINE_MODEL")
	if model == "" {
		model = defaultModel
	}
	endpoint := os.Getenv("STYLE_ENGINE_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultOpenAIEndpoint
	}
	return &Client{
		APIKey:       key,
		Model:        model,
		Endpoint:     endpoint,
		IsOpenRouter: false,
	}
}

func (c *Client) Transform(systemPrompt, userPrompt string) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("no API key set - export OPENROUTER_API_KEY, OPENAI_API_KEY, or STYLE_ENGINE_API_KEY")
	}

	reqBody := Request{
		Model: c.Model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	if c.IsOpenRouter {
		httpReq.Header.Set("HTTP-Referer", "https://style-engine.local")
		httpReq.Header.Set("X-Title", "Style Engine")
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result Response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in API response")
	}

	return result.Choices[0].Message.Content, nil
}
