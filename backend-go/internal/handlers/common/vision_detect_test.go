package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func newTestContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", nil)
	return c
}

func TestHasImageContent_ClaudeMessages(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{
			name:     "claude image block base64",
			body:     `{"messages":[{"role":"user","content":[{"type":"image","source":{"type":"base64","data":"abc"}}]}]}`,
			expected: true,
		},
		{
			name:     "claude image block url",
			body:     `{"messages":[{"role":"user","content":[{"type":"image","source":{"type":"url","url":"https://example.com/img.png"}}]}]}`,
			expected: true,
		},
		{
			name:     "claude text only",
			body:     `{"messages":[{"role":"user","content":[{"type":"text","text":"hello"}]}]}`,
			expected: false,
		},
		{
			name:     "claude string content",
			body:     `{"messages":[{"role":"user","content":"hello"}]}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestContext()
			got := HasImageContent(c, []byte(tt.body))
			if got != tt.expected {
				t.Errorf("HasImageContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasImageContent_OpenAIChat(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{
			name:     "openai image_url block",
			body:     `{"messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"https://example.com/img.png"}}]}]}`,
			expected: true,
		},
		{
			name:     "openai text only",
			body:     `{"messages":[{"role":"user","content":[{"type":"text","text":"hello"}]}]}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestContext()
			got := HasImageContent(c, []byte(tt.body))
			if got != tt.expected {
				t.Errorf("HasImageContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasImageContent_Responses(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{
			name:     "responses input_image top level",
			body:     `{"input":[{"type":"input_image","image_url":"https://example.com/img.png"}]}`,
			expected: true,
		},
		{
			name:     "responses input_image nested in content",
			body:     `{"input":[{"type":"message","role":"user","content":[{"type":"input_image","image_url":"https://example.com/img.png"}]}]}`,
			expected: true,
		},
		{
			name:     "responses text only",
			body:     `{"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"}]}]}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestContext()
			got := HasImageContent(c, []byte(tt.body))
			if got != tt.expected {
				t.Errorf("HasImageContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasImageContent_Gemini(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{
			name:     "gemini inlineData",
			body:     `{"contents":[{"parts":[{"inlineData":{"mimeType":"image/png","data":"abc"}}]}]}`,
			expected: true,
		},
		{
			name:     "gemini fileData",
			body:     `{"contents":[{"parts":[{"fileData":{"mimeType":"image/jpeg","fileUri":"gs://bucket/img.jpg"}}]}]}`,
			expected: true,
		},
		{
			name:     "gemini text only",
			body:     `{"contents":[{"parts":[{"text":"hello"}]}]}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestContext()
			got := HasImageContent(c, []byte(tt.body))
			if got != tt.expected {
				t.Errorf("HasImageContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasImageContent_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{
			name:     "empty body",
			body:     "",
			expected: false,
		},
		{
			name:     "empty json",
			body:     "{}",
			expected: false,
		},
		{
			name:     "malformed json",
			body:     "{invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestContext()
			got := HasImageContent(c, []byte(tt.body))
			if got != tt.expected {
				t.Errorf("HasImageContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasImageContent_ContextCaching(t *testing.T) {
	c := newTestContext()
	body := []byte(`{"messages":[{"role":"user","content":[{"type":"image","source":{"type":"base64","data":"abc"}}]}]}`)

	result1 := HasImageContent(c, body)
	if !result1 {
		t.Fatal("first call should detect image")
	}

	// 第二次调用即使传空 body 也应返回缓存结果
	result2 := HasImageContent(c, nil)
	if !result2 {
		t.Fatal("second call should return cached result")
	}
}

func TestStripImageContentFromChatBody(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, result []byte)
	}{
		{
			name:  "strips image_url and keeps text in multimodal message",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"text","text":"describe this image"},{"type":"image_url","image_url":{"url":"https://example.com/img.png"}}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.Type != gjson.String {
					t.Fatalf("expected string content after single text+image strip, got %s", content.Type)
				}
				if content.String() != "describe this image" {
					t.Fatalf("expected text content 'describe this image', got %q", content.String())
				}
			},
		},
		{
			name:  "strips multiple image_url blocks",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"https://a.com/1.png"}},{"type":"text","text":"compare"},{"type":"image_url","image_url":{"url":"https://a.com/2.png"}}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.Type != gjson.String {
					t.Fatalf("expected string content after multi-image strip, got %s", content.Type)
				}
				if content.String() != "compare" {
					t.Fatalf("expected text content 'compare', got %q", content.String())
				}
			},
		},
		{
			name:  "image-only message gets empty string content",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"https://example.com/img.png"}}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.String() != "" {
					t.Fatalf("expected empty string for image-only message, got %q", content.String())
				}
			},
		},
		{
			name:  "text-only message is unchanged",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"text","text":"hello"}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if !content.IsArray() {
					t.Fatal("expected array content for text-only message")
				}
			},
		},
		{
			name:  "string content message is unchanged",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":"hello"}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.String() != "hello" {
					t.Fatalf("expected 'hello', got %q", content.String())
				}
			},
		},
		{
			name:  "no messages field is unchanged",
			input: `{"model":"deepseek-v4"}`,
			check: func(t *testing.T, result []byte) {
				var m map[string]interface{}
				if err := json.Unmarshal(result, &m); err != nil {
					t.Fatal(err)
				}
				if m["model"] != "deepseek-v4" {
					t.Fatal("model should be unchanged")
				}
			},
		},
		{
			name:  "empty body",
			input: ``,
			check: func(t *testing.T, result []byte) {
				if len(result) != 0 {
					t.Fatal("expected empty result")
				}
			},
		},
		{
			name:  "preserves system and assistant messages",
			input: `{"model":"deepseek-v4","messages":[{"role":"system","content":"You are helpful"},{"role":"user","content":[{"type":"text","text":"hi"},{"type":"image_url","image_url":{"url":"https://example.com/img.png"}}]},{"role":"assistant","content":"Hello!"}]}`,
			check: func(t *testing.T, result []byte) {
				if msg := gjson.GetBytes(result, "messages.0.content"); msg.String() != "You are helpful" {
					t.Fatalf("system message changed: %q", msg.String())
				}
				if msg := gjson.GetBytes(result, "messages.2.content"); msg.String() != "Hello!" {
					t.Fatalf("assistant message changed: %q", msg.String())
				}
			},
		},
		{
			name:  "claude image block is also stripped",
			input: `{"messages":[{"role":"user","content":[{"type":"text","text":"look"},{"type":"image","source":{"type":"base64","data":"abc"}}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.Type != gjson.String || content.String() != "look" {
					t.Fatalf("expected 'look' string after claude image strip, got %s=%q", content.Type, content.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripImageContentFromChatBody([]byte(tt.input))
			tt.check(t, result)
		})
	}
}

func TestConvertChatImageToDeepSeekFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, result []byte)
	}{
		{
			name:  "url image_url converted to top-level image_url",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"text","text":"describe this"},{"type":"image_url","image_url":{"url":"https://example.com/img.png"}}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.String() != "describe this" {
					t.Fatalf("content = %q, want 'describe this'", content.String())
				}
				imageURL := gjson.GetBytes(result, "messages.0.image_url")
				if imageURL.String() != "https://example.com/img.png" {
					t.Fatalf("image_url = %q, want 'https://example.com/img.png'", imageURL.String())
				}
				imageData := gjson.GetBytes(result, "messages.0.image_data")
				if imageData.Exists() {
					t.Fatalf("unexpected image_data: %s", imageData.String())
				}
			},
		},
		{
			name:  "data URI converted to image_data with base64 extracted",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"text","text":"look"},{"type":"image_url","image_url":{"url":"data:image/png;base64,iVBORw0KGgo="}}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.String() != "look" {
					t.Fatalf("content = %q, want 'look'", content.String())
				}
				imageData := gjson.GetBytes(result, "messages.0.image_data")
				if imageData.String() != "iVBORw0KGgo=" {
					t.Fatalf("image_data = %q, want 'iVBORw0KGgo='", imageData.String())
				}
				imageURL := gjson.GetBytes(result, "messages.0.image_url")
				if imageURL.Exists() {
					t.Fatalf("unexpected image_url: %s", imageURL.String())
				}
			},
		},
		{
			name:  "image_url as plain string URL",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"text","text":"hi"},{"type":"image_url","image_url":"https://example.com/img.png"}]}]}`,
			check: func(t *testing.T, result []byte) {
				if u := gjson.GetBytes(result, "messages.0.image_url"); u.String() != "https://example.com/img.png" {
					t.Fatalf("image_url = %q", u.String())
				}
			},
		},
		{
			name:  "only image no text",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"image_url","image_url":{"url":"https://example.com/img.png"}}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.String() != "" {
					t.Fatalf("content = %q, want empty", content.String())
				}
				if u := gjson.GetBytes(result, "messages.0.image_url"); u.String() != "https://example.com/img.png" {
					t.Fatalf("image_url = %q", u.String())
				}
			},
		},
		{
			name:  "text only unchanged",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"text","text":"hello"}]}]}`,
			check: func(t *testing.T, result []byte) {
				if gjson.GetBytes(result, "messages.0.image_url").Exists() {
					t.Fatal("should not have image_url")
				}
				if gjson.GetBytes(result, "messages.0.image_data").Exists() {
					t.Fatal("should not have image_data")
				}
			},
		},
		{
			name:  "assistant message preserved",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"text","text":"hi"},{"type":"image_url","image_url":{"url":"https://a.com/1.png"}}]},{"role":"assistant","content":"Hello!"}]}`,
			check: func(t *testing.T, result []byte) {
				if msg := gjson.GetBytes(result, "messages.1.content"); msg.String() != "Hello!" {
					t.Fatalf("assistant message changed: %q", msg.String())
				}
				if gjson.GetBytes(result, "messages.1.image_url").Exists() {
					t.Fatal("assistant message should not have image_url")
				}
			},
		},
		{
			name:  "system message preserved",
			input: `{"model":"deepseek-v4","messages":[{"role":"system","content":"You are helpful"},{"role":"user","content":[{"type":"text","text":"hi"},{"type":"image_url","image_url":{"url":"https://a.com/1.png"}}]}]}`,
			check: func(t *testing.T, result []byte) {
				if msg := gjson.GetBytes(result, "messages.0.content"); msg.String() != "You are helpful" {
					t.Fatalf("system message changed: %q", msg.String())
				}
			},
		},
		{
			name:  "empty body",
			input: ``,
			check: func(t *testing.T, result []byte) {
				if len(result) != 0 {
					t.Fatal("expected empty result")
				}
			},
		},
		{
			name:  "claude image block converted",
			input: `{"model":"deepseek-v4","messages":[{"role":"user","content":[{"type":"text","text":"look"},{"type":"image","source":{"type":"base64","media_type":"image/jpeg","data":"abcdef"}}]}]}`,
			check: func(t *testing.T, result []byte) {
				content := gjson.GetBytes(result, "messages.0.content")
				if content.String() != "look" {
					t.Fatalf("content = %q, want 'look'", content.String())
				}
				imageData := gjson.GetBytes(result, "messages.0.image_data")
				if imageData.String() != "abcdef" {
					t.Fatalf("image_data = %q, want 'abcdef'", imageData.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertChatImageToDeepSeekFormat([]byte(tt.input))
			tt.check(t, result)
		})
	}
}
