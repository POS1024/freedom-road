package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	var images []image.Image
	var totalWidth, maxHeight int

	// 读取并解码图片
	for i := 1; ; i++ {
		fileName := fmt.Sprintf("/Users/admin/Desktop/sun.png/%d.png", i)
		file, err := os.Open(fileName)
		if err != nil {
			// 假设当无法打开文件时，循环结束
			break
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			fmt.Println("Error decoding image:", err)
			return
		}

		images = append(images, img)

		// 累加总宽度和最大高度
		totalWidth += img.Bounds().Dx()
		if img.Bounds().Dy() > maxHeight {
			maxHeight = img.Bounds().Dy()
		}
	}

	// 创建一个新的图像来拼接所有图像
	newImg := image.NewRGBA(image.Rect(0, 0, totalWidth, maxHeight))

	// 当前x坐标位置
	currentX := 0
	for _, img := range images {
		draw.Draw(newImg, image.Rect(currentX, 0, currentX+img.Bounds().Dx(), img.Bounds().Dy()), img, image.Point{0, 0}, draw.Src)
		currentX += img.Bounds().Dx()
	}

	// 创建输出文件
	outputFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// 将拼接后的图像编码为PNG格式并写入文件
	err = png.Encode(outputFile, newImg)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Images have been successfully concatenated and saved as output.png")
}
