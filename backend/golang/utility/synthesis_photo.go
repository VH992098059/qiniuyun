package common

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

type OCRItem struct {
	Text string `json:"text"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

func SynthesisPhoto(jsonPath, photo string) string {

	// 读取JSON文件内容
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("无法读取JSON文件: %v", err)
	}
	var ocrData []OCRItem
	if err := json.Unmarshal([]byte(jsonData), &ocrData); err != nil {
		log.Fatalf("无法解析JSON数据: %v", err)
	}
	// 3. 加载原始图片
	backgroundImage, err := gg.LoadImage(photo) // 请替换为你的图片路径
	if err != nil {
		log.Fatalf("无法加载图片: %v", err)
	}
	width := backgroundImage.Bounds().Dx()
	height := backgroundImage.Bounds().Dy()

	// 5. 基于原始图片创建绘图上下文
	dc := gg.NewContext(width, height)
	dc.DrawImage(backgroundImage, 0, 0)

	// 4. 加载字体文件
	// 请确保你有一个中文字体文件（如 .ttf），并替换路径
	absFont, _ := filepath.Abs("files/fonts/font.ttf")
	if err := dc.LoadFontFace(absFont, 12); err != nil { // 减小字体大小从20改为12
		log.Fatalf("无法加载字体: %v", err)
	}

	// 6. 遍历数据并进行绘制
	for _, item := range ocrData {
		// 处理多行文本
		lines := strings.Split(item.Text, "\n")

		// 绘制文本
		dc.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255}) // 设置文本颜色为红色

		// 绘制每一行文本
		for i, line := range lines {
			yOffset := float64(i * 15) // 每行间距15像素
			dc.DrawStringAnchored(line, float64(item.X), float64(item.Y)+yOffset, 0.3, 1.3)
		}

		// 绘制一个小圆点来标记坐标点
		dc.SetColor(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // 设置圆点颜色为绿色
		dc.DrawCircle(float64(item.X), float64(item.Y), 3)  // 绘制半径为3的小圆点
		dc.Fill()
	}

	abs, _ := filepath.Abs("files/photos/output_image.png")
	// 7. 保存绘制完成的新图片
	fmt.Println(abs)
	err = dc.SavePNG(abs)
	if err != nil {
		log.Fatalf("无法保存图片: %v", err)
	}

	fmt.Println("图片绘制完成，已保存为 output_image.png")
	return abs
}
