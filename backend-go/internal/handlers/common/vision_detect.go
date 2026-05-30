package common

import (
	"encoding/json"
	"strings"

	"github.com/BenedictKing/ccx/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const visionDetectedContextKey = "ccx_has_image_content"

// HasImageContent 检测请求体是否包含图片内容（覆盖 Claude/OpenAI/Responses/Gemini 四种协议格式）。
// 结果缓存在 gin.Context 中，failover 重试时不重复解析。
func HasImageContent(c *gin.Context, bodyBytes []byte) bool {
	if cached, exists := c.Get(visionDetectedContextKey); exists {
		return cached.(bool)
	}
	detected := detectImageInBody(bodyBytes)
	c.Set(visionDetectedContextKey, detected)
	return detected
}

func detectImageInBody(body []byte) bool {
	if len(body) == 0 {
		return false
	}

	// Claude Messages / OpenAI Chat: messages[*].content[*].type == "image" | "image_url"
	messages := gjson.GetBytes(body, "messages")
	if messages.Exists() && messages.IsArray() {
		for _, msg := range messages.Array() {
			content := msg.Get("content")
			if content.IsArray() {
				for _, block := range content.Array() {
					t := block.Get("type").String()
					if t == "image" || t == "image_url" {
						return true
					}
				}
			}
		}
	}

	// Responses API: input[*].type == "input_image" 或嵌套 content 中的 input_image
	input := gjson.GetBytes(body, "input")
	if input.Exists() && input.IsArray() {
		for _, item := range input.Array() {
			t := item.Get("type").String()
			if t == "input_image" {
				return true
			}
			// 嵌套 content 数组（如 input_message.content）
			itemContent := item.Get("content")
			if itemContent.IsArray() {
				for _, block := range itemContent.Array() {
					if block.Get("type").String() == "input_image" {
						return true
					}
				}
			}
		}
	}

	// Gemini: contents[*].parts[*].inlineData 或 fileData（含 image MIME）
	contents := gjson.GetBytes(body, "contents")
	if contents.Exists() && contents.IsArray() {
		for _, c := range contents.Array() {
			parts := c.Get("parts")
			if parts.IsArray() {
				for _, part := range parts.Array() {
					if part.Get("inlineData").Exists() || part.Get("fileData").Exists() {
						return true
					}
				}
			}
		}
	}

	return false
}

// isNoVisionModel 检查模型是否在渠道的 NoVisionModels 列表中（精确匹配）。

// isDeepSeekUpstream 检测是否为 DeepSeek 上游（名称或 baseURL 含 deepseek）
func isDeepSeekUpstream(upstream *config.UpstreamConfig) bool {
	return strings.Contains(strings.ToLower(upstream.Name), "deepseek") ||
		strings.Contains(strings.ToLower(upstream.BaseURL), "deepseek")
}

func isNoVisionModel(upstream *config.UpstreamConfig, model string) bool {
	for _, m := range upstream.NoVisionModels {
		if m == model {
			return true
		}
	}
	return false
}

// StripImageContentFromChatBody 从 OpenAI Chat Completions 请求体中移除 image_url 内容块。
// 用于上游模型不支持视觉输入时，剥离图片保留纯文本继续对话。
// 如果消息内容仅有图片（无文本），则内容设为空字符串。
func StripImageContentFromChatBody(body []byte) []byte {
	if len(body) == 0 {
		return body
	}

	root := gjson.ParseBytes(body)
	messages := root.Get("messages")
	if !messages.Exists() || !messages.IsArray() {
		return body
	}

	modified := false
	out := string(body)

	messages.ForEach(func(i, msg gjson.Result) bool {
		content := msg.Get("content")
		if !content.IsArray() {
			return true
		}

		var textBlocks []json.RawMessage
		hasImage := false

		content.ForEach(func(_, block gjson.Result) bool {
			blockType := block.Get("type").String()
			if blockType == "image_url" || blockType == "image" {
				hasImage = true
				return true
			}
			textBlocks = append(textBlocks, json.RawMessage(block.Raw))
			return true
		})

		if !hasImage {
			return true
		}

		modified = true
		path := "messages." + i.String() + ".content"

		if len(textBlocks) == 0 {
			out, _ = sjson.Set(out, path, "")
		} else if len(textBlocks) == 1 {
			var single map[string]interface{}
			if json.Unmarshal(textBlocks[0], &single) == nil {
				if t, ok := single["text"].(string); ok {
					out, _ = sjson.Set(out, path, t)
					return true
				}
			}
			out, _ = sjson.SetRaw(out, path, string(textBlocks[0]))
		} else {
			var arr []json.RawMessage
			for _, b := range textBlocks {
				arr = append(arr, b)
			}
			out, _ = sjson.Set(out, path, arr)
		}
		return true
	})

	if !modified {
		return body
	}
	return []byte(out)
}

// ConvertChatImageToDeepSeekFormat 将 OpenAI Chat 格式的 image_url 内容块转换为
// DeepSeek V4 的顶层字段格式。
//
// OpenAI 格式（不被 DeepSeek 接受）：
//
//	{"role":"user","content":[{"type":"text","text":"..."},{"type":"image_url","image_url":{"url":"..."}}]}
//
// DeepSeek 格式：
//
//	{"role":"user","content":"...","image_url":"https://..."}       (URL)
//	{"role":"user","content":"...","image_data":"iVBORw..."}        (Base64)
//
// 注意：
//   - DeepSeek 每条消息只支持一张图片，多余图片会被丢弃
//   - data: URI 会被提取为纯 base64 放入 image_data
//   - 非 user 角色的消息不动
//   - 没有图片的消息不动
func ConvertChatImageToDeepSeekFormat(body []byte) []byte {
	if len(body) == 0 {
		return body
	}

	root := gjson.ParseBytes(body)
	messages := root.Get("messages")
	if !messages.Exists() || !messages.IsArray() {
		return body
	}

	modified := false
	out := string(body)

	messages.ForEach(func(i, msg gjson.Result) bool {
		role := msg.Get("role").String()
		if role != "user" {
			return true
		}

		content := msg.Get("content")
		if !content.IsArray() {
			return true
		}

		var textParts []string
		var imageURL string
		var imageData string
		hasImage := false

		content.ForEach(func(_, block gjson.Result) bool {
			blockType := block.Get("type").String()
			switch blockType {
			case "text":
				if t := block.Get("text").String(); t != "" {
					textParts = append(textParts, t)
				}
			case "image_url":
				hasImage = true
				img := block.Get("image_url")
				if img.IsObject() {
					if url := img.Get("url").String(); url != "" {
						if data, ok := extractBase64FromDataURI(url); ok {
							imageData = data
						} else {
							imageURL = url
						}
					}
				} else if img.Type == gjson.String {
					if data, ok := extractBase64FromDataURI(img.String()); ok {
						imageData = data
					} else {
						imageURL = img.String()
					}
				}
			case "image":
				hasImage = true
				// Claude 格式的 image block：优先 source.url，其次 source.data
				source := block.Get("source")
				if url := source.Get("url").String(); url != "" {
					imageURL = url
				} else if data := source.Get("data").String(); data != "" {
					mediaType := source.Get("media_type").String()
					if mediaType == "" {
						mediaType = "image/png"
					}
					imageData = data
				}
			}
			return true
		})

		if !hasImage {
			return true
		}

		modified = true
		basePath := "messages." + i.String()

		// 设置文本内容
		textContent := strings.Join(textParts, "\n")
		out, _ = sjson.Set(out, basePath+".content", textContent)

		// 设置图片字段（DeepSeek 每条消息只支持一张图）
		if imageData != "" {
			out, _ = sjson.Set(out, basePath+".image_data", imageData)
			out, _ = sjson.Delete(out, basePath+".image_url")
		} else if imageURL != "" {
			out, _ = sjson.Set(out, basePath+".image_url", imageURL)
			out, _ = sjson.Delete(out, basePath+".image_data")
		}
		return true
	})

	if !modified {
		return body
	}
	return []byte(out)
}

// extractBase64FromDataURI 从 data URI 中提取纯 base64 数据。
// 例如 "data:image/png;base64,iVBORw..." → ("iVBORw...", true)
func extractBase64FromDataURI(uri string) (string, bool) {
	const prefix = "data:"
	if !strings.HasPrefix(uri, prefix) {
		return "", false
	}
	// 找到 base64, 标记
	idx := strings.Index(uri, ";base64,")
	if idx == -1 {
		return "", false
	}
	return uri[idx+len(";base64,"):], true
}
