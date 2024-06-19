package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

func cropImage(img image.Image, rect image.Rectangle) image.Image {
	// 创建一个新的 RGBA 图像，其大小为裁剪区域的大小
	cropped := image.NewRGBA(rect)

	// 在新的图像上绘制裁剪区域
	draw.Draw(cropped, rect, img, rect.Min, draw.Src)

	return cropped
}

func main() {
	// 打开 PNG 文件
	file, err := os.Open("/Users/admin/Desktop/Crab_shadow1.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 解码 PNG 文件
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	// 指定裁剪区域 (x0, y0, x1, y1)
	rect := image.Rect(16, 0, 16+104, 104)

	// 裁剪图片
	croppedImg := cropImage(img, rect)

	// 创建输出文件
	outFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// 将裁剪后的图片编码为 PNG 并保存到文件
	err = png.Encode(outFile, croppedImg)
	if err != nil {
		panic(err)
	}

	println("Image cropped and saved as output.png")
}
