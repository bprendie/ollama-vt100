// ollama/client.go

package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ... Client struct and NewClient function are unchanged ...
type Client struct {
	BaseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &Client{
		BaseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// ... ListModelsResponse and Model structs are unchanged ...
type ListModelsResponse struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
}

// --- NEW Structs for Chat API ---

// Message represents a single message in a chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the payload sent to the /api/chat endpoint.
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Options  struct {
		Temperature float32 `json:"temperature"`
		NumCtx      int     `json:"num_ctx"`
	} `json:"options"`
}

// ChatResponse is the structure of a single JSON object in a streaming response.
type ChatResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

// ListModels remains the same
func (c *Client) ListModels() ([]string, error) {
	// ... implementation is unchanged ...
	resp, err := c.httpClient.Get(c.BaseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("could not connect to Ollama: %w. \nIs Ollama running?", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned a non-200 status: %s", resp.Status)
	}
	var listResponse ListModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, fmt.Errorf("failed to decode ollama response: %w", err)
	}
	var modelNames []string
	for _, model := range listResponse.Models {
		modelNames = append(modelNames, model.Name)
	}
	return modelNames, nil
}

// NEW Chat function
// Chat sends a request to Ollama and streams the response.
// It calls the onChunk function for each piece of the response.
func (c *Client) Chat(req ChatRequest, onChunk func(string)) error {
	req.Stream = true // Always stream

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal chat request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/api/chat", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request to ollama: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var chatResp ChatResponse
		if err := decoder.Decode(&chatResp); err == io.EOF {
			break // End of stream
		} else if err != nil {
			return fmt.Errorf("failed to decode streaming response: %w", err)
		}

		onChunk(chatResp.Message.Content)

		if chatResp.Done {
			break
		}
	}
	return nil
}
