package providers

import (
	"testing"

	"github.com/BenedictKing/ccx/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestGeminiProvider_ConvertMessage_ToolResultArray(t *testing.T) {
	provider := &GeminiProvider{}

	// 测试场景：tool_result 的 content 是一个 Content Blocks 数组
	msg := types.ClaudeMessage{
		Role: "user",
		Content: []interface{}{
			map[string]interface{}{
				"type":        "tool_result",
				"tool_use_id": "toolu_0",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": "Tokyo is sunny.",
					},
					map[string]interface{}{
						"type": "text",
						"text": "Temperature is 22C.",
					},
				},
			},
		},
	}

	geminiMsg := provider.convertMessage(msg)
	assert.NotNil(t, geminiMsg)
	assert.Equal(t, "user", geminiMsg["role"])

	parts, ok := geminiMsg["parts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, parts, 1)

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)

	funcResp, ok := part["functionResponse"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "toolu_0", funcResp["name"])

	response, ok := funcResp["response"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Tokyo is sunny.\nTemperature is 22C.", response["result"])
}

func TestGeminiProvider_ConvertMessage_ToolResultString(t *testing.T) {
	provider := &GeminiProvider{}

	// 测试场景：tool_result 的 content 是一个简单字符串
	msg := types.ClaudeMessage{
		Role: "user",
		Content: []interface{}{
			map[string]interface{}{
				"type":        "tool_result",
				"tool_use_id": "toolu_1",
				"content":     "Tokyo is sunny.",
			},
		},
	}

	geminiMsg := provider.convertMessage(msg)
	assert.NotNil(t, geminiMsg)

	parts, ok := geminiMsg["parts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, parts, 1)

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)

	funcResp, ok := part["functionResponse"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "toolu_1", funcResp["name"])

	response, ok := funcResp["response"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Tokyo is sunny.", response["result"])
}

func TestGeminiProvider_ConvertMessage_ToolResultObject(t *testing.T) {
	provider := &GeminiProvider{}

	// 测试场景：tool_result 的 content 是一个 JSON 对象
	msg := types.ClaudeMessage{
		Role: "user",
		Content: []interface{}{
			map[string]interface{}{
				"type":        "tool_result",
				"tool_use_id": "toolu_2",
				"content": map[string]interface{}{
					"temperature": 22,
					"condition":   "sunny",
				},
			},
		},
	}

	geminiMsg := provider.convertMessage(msg)
	assert.NotNil(t, geminiMsg)

	parts, ok := geminiMsg["parts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, parts, 1)

	part, ok := parts[0].(map[string]interface{})
	assert.True(t, ok)

	funcResp, ok := part["functionResponse"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "toolu_2", funcResp["name"])

	response, ok := funcResp["response"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 22, response["temperature"])
	assert.Equal(t, "sunny", response["condition"])
}
