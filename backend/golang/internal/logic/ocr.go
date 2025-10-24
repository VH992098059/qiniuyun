package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func OcrLogic(ctx context.Context, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建 multipart 表单
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", "screenshot.png")
	if err != nil {
		panic(err)
	}
	if _, err = io.Copy(part, file); err != nil {
		panic(err)
	}
	writer.Close()

	// 发送 POST 请求到 Flask 接口
	resp, err := http.Post("http://localhost:5000/ocr", writer.FormDataContentType(), &buf)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析响应，优先处理顶层数组的首元素（与你之前的要求一致）
	var value interface{}
	var arr []json.RawMessage
	if err := json.Unmarshal(body, &arr); err == nil && len(arr) > 0 {
		// 取数组首元素为最终写入内容
		if err := json.Unmarshal(arr[0], &value); err != nil {
			// 若无法解析为结构，则作为原始 JSON 片段
			value = arr[0]
		}
	} else {
		// 非数组或解析失败时，尝试直接解析为对象/其他类型
		if err := json.Unmarshal(body, &value); err != nil {
			// 作为最后回退：去掉最外层 [] 并直接写原文
			trimmed := bytes.TrimSpace(body)
			if len(trimmed) >= 2 && trimmed[0] == '[' && trimmed[len(trimmed)-1] == ']' {
				trimmed = bytes.TrimSpace(trimmed[1 : len(trimmed)-1])
			}
			if err := os.WriteFile("ocr_result.json", trimmed, 0644); err != nil {
				panic(err)
			}
			fmt.Println("OCR Response 已保存到 ocr_result.json（原始文本回退）")
			fmt.Println(string(trimmed))
			return
		}
	}

	// 以 UTF-8 且不转义中文写入文件（避免 \uXXXX）
	out, err := os.Create("ocr_result.json")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		panic(err)
	}

	// 同样打印到控制台为可读中文
	/*var pretty bytes.Buffer
	enc2 := json.NewEncoder(&pretty)
	enc2.SetEscapeHTML(false)
	enc2.SetIndent("", "  ")
	if err := enc2.Encode(value); err == nil {
		fmt.Println(pretty.String())
	}*/

	fmt.Println("OCR Response 已保存到 ocr_result.json（已解除 \\uXXXX 中文转义）")
}
