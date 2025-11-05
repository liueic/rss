package summarizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	maxContentLen  = 8000 // 限制输入内容长度，避免超出模型限制
)

type Summarizer struct {
	apiEndpoint string
	apiKey      string
	model       string
	client      *http.Client
	enabled     bool
}

type APIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type APIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func New() *Summarizer {
	apiEndpoint := os.Getenv("API_ENDPOINT")
	apiKey := os.Getenv("API_KEY")
	model := os.Getenv("MODEL_NAME")

	enabled := apiEndpoint != "" && apiKey != "" && model != ""

	return &Summarizer{
		apiEndpoint: apiEndpoint,
		apiKey:      apiKey,
		model:       model,
		enabled:     enabled,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (s *Summarizer) IsEnabled() bool {
	return s.enabled
}

func (s *Summarizer) Summarize(ctx context.Context, title, description string) (string, error) {
	if !s.enabled {
		return "", fmt.Errorf("summarizer is not enabled")
	}

	// 构建提示词
	content := buildPrompt(title, description)

	// 构建请求
	reqBody := APIRequest{
		Model: s.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
		Temperature: 0.7,
		MaxTokens:   500, // 限制总结长度
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// 先检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		// 尝试解析错误响应
		var errorResp struct {
			Error *struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code,omitempty"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &errorResp) == nil && errorResp.Error != nil {
			return "", fmt.Errorf("API error (status %d): %s (type: %s)", resp.StatusCode, errorResp.Error.Message, errorResp.Error.Type)
		}
		// 如果无法解析错误响应，只显示状态码和响应长度，不显示内容以避免泄露敏感信息
		return "", fmt.Errorf("API returned status %d (response length: %d bytes)", resp.StatusCode, len(body))
	}

	// 检查响应体是否为空
	if len(body) == 0 {
		return "", fmt.Errorf("empty response body")
	}

	// 解析响应
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		// 不打印响应体内容，避免泄露敏感信息，只显示解析错误和响应长度
		return "", fmt.Errorf("failed to parse response (invalid JSON): %w (response length: %d bytes)", err, len(body))
	}

	// 检查错误字段
	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s (type: %s)", apiResp.Error.Message, apiResp.Error.Type)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in API response")
	}

	summary := strings.TrimSpace(apiResp.Choices[0].Message.Content)
	if summary == "" {
		return "", fmt.Errorf("empty summary from API")
	}

	return summary, nil
}

func buildPrompt(title, description string) string {
	// 限制内容长度
	content := title
	if description != "" {
		content += "\n\n" + description
	}

	// 如果内容太长，截断
	runes := []rune(content)
	if len(runes) > maxContentLen {
		content = string(runes[:maxContentLen]) + "..."
	}

	prompt := fmt.Sprintf(`请为以下文章生成一个简洁的中文总结，要求：
1. 总结长度控制在100字以内
2. 突出文章的核心观点和关键信息
3. 使用简洁明了的语言
4. 如果原文是英文，请翻译成中文

文章标题：%s

文章内容：
%s

请直接输出总结内容，不要添加任何前缀或说明。`, title, content)

	return prompt
}
