package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient 与本地 Ollama 模型交互
type OllamaClient struct {
	BaseURL    string
	Model      string
	EmbedModel string
	client     *http.Client
}

// NewOllamaClient 创建客户端
func NewOllamaClient(baseURL, model, embedModel string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2:latest"
	}
	if embedModel == "" {
		embedModel = model
	}
	return &OllamaClient{
		BaseURL:    baseURL,
		Model:      model,
		EmbedModel: embedModel,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Error    string `json:"error"`
}

// Generate 调用Llama生成内容
func (c *OllamaClient) Generate(prompt string) (string, error) {
	req := ollamaGenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	body, _ := json.Marshal(req)
	resp, err := c.client.Post(fmt.Sprintf("%s/api/generate", c.BaseURL), "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("ollama generate error: %s", string(data))
	}

	var result ollamaGenerateResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("ollama generate error: %s", result.Error)
	}

	return result.Response, nil
}

type ollamaEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
	Error     string    `json:"error"`
}

// Embed 生成文本向量
func (c *OllamaClient) Embed(text string) ([]float32, error) {
	req := ollamaEmbedRequest{
		Model: c.EmbedModel,
		Input: text,
	}

	body, _ := json.Marshal(req)
	resp, err := c.client.Post(fmt.Sprintf("%s/api/embeddings", c.BaseURL), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ollama embed error: %s", string(data))
	}

	var result ollamaEmbedResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, fmt.Errorf("ollama embed error: %s", result.Error)
	}

	return result.Embedding, nil
}

