package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

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
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      string `json:"data"` // 注意：这是嵌套的 JSON 字符串
	Timestamp int64  `json:"timestamp"`
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
	ActionName string         `json:"action_name"`
	Parameters ActionParamMap `json:"parameters"`
}
type ActionParamMap struct {
	X           string `json:"x"`
	Y           string `json:"y"`
	Comment     string `json:"comment"`
	DoubleClick string `json:"double_click"`
	Button      string `json:"button"`
	Keyword     string `json:"keyword"`
}

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

	// 发送 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", modelURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 校验响应类型
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		return nil, fmt.Errorf("响应 Content-Type 非 application/json: %s", ct)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析外层
	var outer actionModelOuterRes
	if err = json.Unmarshal(raw, &outer); err != nil {
		fmt.Println("解析外层失败:", err)
		return nil, fmt.Errorf("解析外层响应失败: %v", err)
	}
	// 解析内层
	if err = json.Unmarshal([]byte(outer.Data), &result); err != nil {
		fmt.Println("解析内层失败:", err)
		return nil, fmt.Errorf("解析内层数据失败: %v", err)
	}
	return result, nil
}
