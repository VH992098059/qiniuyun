package common

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ImgLogic(ctx context.Context, filePath string) (ob image.Rectangle, lb image.Rectangle, rb image.Rectangle) {
	// 打开原图
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("FATAL: Failed to open file '%s': %v", filePath, err)
	}
	defer f.Close()

	// 解码，支持 JPEG/PNG/GIF（通过导入注册解码器）
	img, format, err := image.Decode(f)
	if err != nil {
		log.Fatalf("FATAL: Failed to decode image: %v", err)
	}
	formatLower := strings.ToLower(format)

	// 分别裁剪左右半幅
	left := CropLeftHalf(img)
	right := CropRightHalf(img)

	// 输出文件名（尽量保持与输入格式一致；未知格式回退为 JPG）
	ext := strings.ToLower(filepath.Ext(filePath))
	base := strings.TrimSuffix(filePath, ext)
	var outLeft, outRight string
	switch formatLower {
	case "png":
		outLeft = base + "_left.png"
		outRight = base + "_right.png"
	case "jpeg", "jpg":
		outLeft = base + "_left.jpg"
		outRight = base + "_right.jpg"
	default:
		outLeft = base + "_left.jpg"
		outRight = base + "_right.jpg"
		log.Printf("WARN: Unsupported format '%s' for encoding, defaulting to JPEG.", format)
	}

	// 保存为对应格式
	lf, err := os.Create(outLeft)
	if err != nil {
		log.Fatalf("FATAL: Failed to create output file '%s': %v", outLeft, err)
	}
	defer lf.Close()
	switch formatLower {
	case "png":
		if err := png.Encode(lf, left); err != nil {
			log.Fatalf("FATAL: Failed to encode PNG (left): %v", err)
		}
	default:
		if err := jpeg.Encode(lf, left, &jpeg.Options{Quality: 95}); err != nil {
			log.Fatalf("FATAL: Failed to encode JPEG (left): %v", err)
		}
	}

	rf, err := os.Create(outRight)
	if err != nil {
		log.Fatalf("FATAL: Failed to create output file '%s': %v", outRight, err)
	}
	defer rf.Close()
	switch formatLower {
	case "png":
		if err := png.Encode(rf, right); err != nil {
			log.Fatalf("FATAL: Failed to encode PNG (right): %v", err)
		}
	default:
		if err := jpeg.Encode(rf, right, &jpeg.Options{Quality: 95}); err != nil {
			log.Fatalf("FATAL: Failed to encode JPEG (right): %v", err)
		}
	}

	// 打印尺寸信息
	ob = img.Bounds()
	lb = left.Bounds()
	rb = right.Bounds()
	fmt.Println("---------------------------------")
	fmt.Printf("Original: %dx%d (format=%s)\n", ob.Dx(), ob.Dy(), format)
	fmt.Printf("Left:     %dx%d -> saved: %s\n", lb.Dx(), lb.Dy(), outLeft)
	fmt.Printf("Right:    %dx%d -> saved: %s\n", rb.Dx(), rb.Dy(), outRight)
	fmt.Println("---------------------------------")
	return
}
