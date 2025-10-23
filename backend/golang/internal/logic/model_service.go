package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

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
		Level   string `json:"level"`
		Message string `json:"message"`
	} `json:"safetyCheck"`
	Data map[string][]string `json:"data"`
}

// ActionModelReq 意图操作请求
type ActionModelReq struct {
	FileImage       multipart.File `json:"file"`
	Text            string         `json:"prompt"`
	PreviousActions string         `json:"previous_actions"`
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
type ActionParamMap map[string]interface{}

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

	return
}

func ActionModel(ctx context.Context, actionReq *ActionModelReq) (result *ActionModelInTaskData, err error) {
	modelURL := "http://localhost:8081/api/ai/multi/analyze"
	post, err := g.Client().Post(ctx, modelURL, actionReq)
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
	//解析外层结构
	var outer actionModelOuterRes
	if err = json.Unmarshal(raw, &outer); err != nil {
		fmt.Println("解析外层失败:", err)
		return
	}
	if err = json.Unmarshal([]byte(outer.Data), &result); err != nil {
		fmt.Println("解析内层失败:", err)
		return
	}
	return
}
