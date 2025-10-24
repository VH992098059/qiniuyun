package common

import (
	"image"
	"image/draw"
)

// 裁剪函数对具体格式无关：只要能解码为 image.Image（如 JPEG/PNG/GIF），即可裁剪。
// 解码器在上层逻辑中注册，传入的 img 将按其内存像素进行裁剪。

// CropLeftHalf 返回图像左半部分（宽度的一半）
func CropLeftHalf(img image.Image) image.Image {
	b := img.Bounds()
	halfW := b.Dx() / 2
	cropRect := image.Rect(b.Min.X, b.Min.Y, b.Min.X+halfW, b.Max.Y)

	if si, ok := img.(interface {
		SubImage(image.Rectangle) image.Image
	}); ok {
		return si.SubImage(cropRect)
	}
	dst := image.NewRGBA(image.Rect(0, 0, halfW, b.Dy()))
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Src)
	return dst
}

// CropRightHalf 返回图像右半部分（宽度的一半）
func CropRightHalf(img image.Image) image.Image {
	b := img.Bounds()
	halfW := b.Dx() / 2
	cropRect := image.Rect(b.Min.X+halfW, b.Min.Y, b.Max.X, b.Max.Y)

	if si, ok := img.(interface {
		SubImage(image.Rectangle) image.Image
	}); ok {
		return si.SubImage(cropRect)
	}
	dst := image.NewRGBA(image.Rect(0, 0, halfW, b.Dy()))
	draw.Draw(dst, dst.Bounds(), img, image.Point{X: b.Min.X + halfW, Y: b.Min.Y}, draw.Src)
	return dst
}
