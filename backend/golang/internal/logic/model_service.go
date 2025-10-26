package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type modelChatReq struct {
	Text      string `json:"text"`
	SessionId string `json:"session_id"`
}

// 外层响应结构
type modelChatRes struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      string `json:"data"` // 注意：这是 JSON 字符串，不是结构体
	Timestamp int64  `json:"timestamp"`
}

// DataAnalysis 内部响应结构
type DataAnalysis struct {
	SafetyCheck struct {
		Level     string `json:"level"`
		Message   string `json:"message"`
		Intention string `json:"intention"`
	} `json:"safetyCheck"`
	Data map[string][]string `json:"data"`
}

// ActionModelReq 意图操作请求
type ActionModelReq struct {
	FileImage multipart.File `json:"file"`
	Text      string         `json:"prompt"`
}

// 意图操作外层响应
type actionModelOuterRes struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

// ActionModelInTaskData 意图操作内部响应
type ActionModelInTaskData struct {
	OverallGoal  string            `json:"overall_goal"`
	GoalAnalysis string            `json:"goal_analysis"`
	Status       string            `json:"status"`
	Thought      string            `json:"thought"`
	Actions      []taskDataActions `json:"actions"`
}
type taskDataActions struct {
	ActionName string        `json:"action_name"`
	Parameters ParametersMap `json:"parameters"`
}
type ParametersMap map[string]any

func ModelService(ctx context.Context, text string) (result *DataAnalysis, err error) {
	modelURL := "http://localhost:8081/api/ai/chat/do"
	jsonReq := modelChatReq{
		Text:      text,
		SessionId: gtime.Now().String(),
	}

	jsonModelReq, err := json.Marshal(jsonReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	post, err := g.Client().Post(ctx, modelURL, jsonModelReq)
	if err != nil {
		return nil, err
	}
	defer post.Body.Close()
	if post.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request")
	}
	raw, err := io.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}

	// 第一步：解析外层结构
	var outer modelChatRes
	if err = json.Unmarshal(raw, &outer); err != nil {
		fmt.Println("解析外层失败:", err)
		return
	}

	// 第二步：解析内层 data 字符串
	if err = json.Unmarshal([]byte(outer.Data), &result); err != nil {
		fmt.Println("解析内层失败:", err)
		return
	}
	log.Println(result)
	return
}

func ActionModel(ctx context.Context, actionReq *ActionModelReq) (result *ActionModelInTaskData, err error) {
	modelURL := "http://localhost:8081/api/ai/multi/analyze"

	// 使用 multipart.Writer 构建 form-data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 文件字段（可选）
	if actionReq.FileImage != nil {
		fileBytes, err := io.ReadAll(actionReq.FileImage)
		if err != nil {
			return nil, fmt.Errorf("读取文件失败: %v", err)
		}

		// 嗅探 MIME 类型
		contentType := http.DetectContentType(fileBytes)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// 基于类型构造文件名（仅用于服务端识别）
		filename := "image.bin"
		switch contentType {
		case "image/jpeg":
			filename = "image.jpg"
		case "image/png":
			filename = "image.png"
		case "image/gif":
			filename = "image.gif"
		}

		// 显式设置文件 part 的 Content-Type
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
		h.Set("Content-Type", contentType)
		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, fmt.Errorf("创建文件字段失败: %v", err)
		}
		if _, err := part.Write(fileBytes); err != nil {
			return nil, fmt.Errorf("写入文件内容失败: %v", err)
		}
	}

	// 文本字段
	if err := writer.WriteField("prompt", actionReq.Text); err != nil {
		return nil, fmt.Errorf("添加文本字段失败: %v", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭 writer 失败: %v", err)
	}

	// 独立超时上下文，避免上游请求结束导致被取消
	reqCtx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()

	// 发送 HTTP 请求
	req, err := http.NewRequestWithContext(reqCtx, "POST", modelURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 200 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("发送请求失败: 上游 context 已取消或生命周期已结束: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("发送请求失败: 请求超时: %w", err)
		}
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 校验响应类型（仅告警，不拦截）
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		log.Printf("警告: 响应 Content-Type 非 application/json: %s", ct)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析外层（兼容 data 为字符串或对象）
	var outer actionModelOuterRes
	if err = json.Unmarshal(raw, &outer); err == nil && len(outer.Data) > 0 {
		// 先尝试将 data 解析为字符串（服务端可能返回字符串包裹的 JSON）
		var dataStr string
		if err2 := json.Unmarshal(outer.Data, &dataStr); err2 == nil {
			extractJSON, err := CleanAndExtractJSON(dataStr)
			if err != nil {
				return nil, err
			}
			if err = json.Unmarshal([]byte(extractJSON), &result); err != nil {
				return nil, fmt.Errorf("解析内层数据失败: %v", err)
			}
			return result, nil
		}
		// 若不是字符串，则直接按对象解析
		if err2 := json.Unmarshal(outer.Data, &result); err2 == nil {
			return result, nil
		}
		// 解析失败则返回错误
		return nil, fmt.Errorf("解析内层数据失败: %v", err)
	}

	// 降级：直接解析为内部结构（用于服务端直接返回内部对象的情况）
	if err = json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	return result, nil
}

// CleanAndExtractJSON 从可能被Markdown包裹的字符串中提取出纯JSON部分
func CleanAndExtractJSON(rawResponse string) (string, error) {
	// 找到第一个 '{' 的位置
	startIndex := strings.Index(rawResponse, "{")
	if startIndex == -1 {
		return "", errors.New("响应中未找到JSON起始括号 '{'")
	}

	// 找到最后一个 '}' 的位置
	endIndex := strings.LastIndex(rawResponse, "}")
	if endIndex == -1 || endIndex < startIndex {
		return "", errors.New("响应中未找到JSON结束括号 '}' 或括号顺序错误")
	}

	// 提取从第一个 '{' 到最后一个 '}' 的子字符串
	jsonStr := rawResponse[startIndex : endIndex+1]

	// 验证一下它是否是有效的JSON（可选但推荐）
	var js map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &js); err != nil {
		return "", errors.New("提取的字符串不是有效的JSON")
	}

	return jsonStr, nil
}
