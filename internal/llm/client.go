package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TransformRequest struct {
	Text      string  `json:"text"`
	Style     string  `json:"style"`
	Intensity float64 `json:"intensity"`
}

type TransformResponse struct {
	Original    string  `json:"original"`
	Transformed string  `json:"transformed"`
	Style       string  `json:"style"`
	Intensity   float64 `json:"intensity"`
	Score       float64 `json:"compliance_score"`
	Retries     int     `json:"retries"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "http://localhost:8080",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Transform(text, style string, intensity float64) (*TransformResponse, error) {
	req := TransformRequest{
		Text:      text,
		Style:     style,
		Intensity: intensity,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/transform", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("style-engine request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("style-engine returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result TransformResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}
