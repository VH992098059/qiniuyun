package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// PaddleOCRResponse 定义与JSON结构匹配的Go结构体
type PaddleOCRResponse struct {
	DtPolys  [][][]int `json:"dt_polys"`
	RecTexts []string  `json:"rec_texts"`
}

// MatchResult 新增：多匹配结果结构
type MatchResult struct {
	Index   int
	Text    string
	CenterX int
	CenterY int
}

// FindAllCenterCoordinatesByText 解析PaddleOCR的JSON输出，并查找目标文本的中心坐标
// jsonFilePath: PaddleOCR输出的JSON文件路径
// targetText: 你想要查找的文本，例如 "搜索音乐"
// returns: centerX, centerY, error
// 新增：返回所有匹配项的中心坐标
func FindAllCenterCoordinatesByText(jsonFilePath, targetText string) ([]MatchResult, error) {
	// 1. 读取JSON文件
	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %v", err)
	}

	// 2. 解析JSON到我们的结构体
	var response PaddleOCRResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// 检查文本和坐标列表的长度是否匹配
	if len(response.RecTexts) != len(response.DtPolys) {
		return nil, fmt.Errorf("mismatch between number of texts (%d) and polygons (%d)", len(response.RecTexts), len(response.DtPolys))
	}

	// 3. 遍历所有识别出的文本，收集全部匹配
	var results []MatchResult
	//log.Printf("INFO: Searching for text containing '%s' among %d results...", targetText, len(response.RecTexts))
	for i, text := range response.RecTexts {
		if strings.Contains(text, targetText) {
			poly := response.DtPolys[i]
			if len(poly) != 4 {
				log.Printf("WARN: Found text but polygon has %d points, skipping.", len(poly))
				continue
			}

			// 计算中心点坐标（取左上角与右下角）
			topLeft := poly[0]
			bottomRight := poly[2]
			centerX := (topLeft[0] + bottomRight[0]) / 2
			centerY := (topLeft[1] + bottomRight[1]) / 2

			results = append(results, MatchResult{
				Index:   i,
				Text:    text,
				CenterX: centerX,
				CenterY: centerY,
			})
			//log.Printf("INFO: Matched '%s' at index %d with center (%d, %d)", text, i, centerX, centerY)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("target text '%s' not found in the OCR results", targetText)
	}
	return results, nil
}

// 保留原函数：仍只返回第一个匹配，兼容已有调用
func FindCenterCoordinateByText(jsonFilePath, targetText string) (int, int, error) {
	results, err := FindAllCenterCoordinatesByText(jsonFilePath, targetText)
	if err != nil {
		return 0, 0, err
	}
	first := results[0]
	return first.CenterX, first.CenterY, nil
}
