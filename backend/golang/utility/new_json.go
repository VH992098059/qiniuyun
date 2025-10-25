package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type CoordinateResult struct {
	Text string `json:"text"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

func NewJsonFile(jsonFile string) string {
	result, err := FindAllCenterCoordinatesByText(jsonFile, "")
	if err != nil {
		log.Fatalf("FATAL: %v", err)
	}

	// 转换为输出格式
	var coordinates []CoordinateResult
	for _, v := range result {
		//fmt.Printf("%s: 坐标 %d,%d\n", v.Text, v.CenterX, v.CenterY)
		coordinates = append(coordinates, CoordinateResult{
			Text: fmt.Sprintf("X:%d\nY:%d", v.CenterX, v.CenterY), // 添加换行符分隔 X 和 Y 坐标
			X:    v.CenterX,
			Y:    v.CenterY,
		})
	}

	// 使用 filepath.Abs 获取绝对路径
	abs, err := filepath.Abs("files/json/coordinates_result.json")
	if err != nil {
		log.Fatalf("FATAL: Failed to get absolute path: %v", err)
	}
	outputFile := abs

	// 确保目录存在
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("FATAL: Failed to create directory: %v", err)
	}

	jsonData, err := json.MarshalIndent(coordinates, "", "  ")
	if err != nil {
		log.Fatalf("FATAL: Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("FATAL: Failed to write JSON file: %v", err)
	}

	fmt.Printf("\n已将 %d 个坐标结果保存到: %s\n", len(coordinates), outputFile)
	return outputFile
}
