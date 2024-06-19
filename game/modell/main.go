// main.go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/jpeg"
	"log"
	"modee/utils"
	"os"
)

func test() {
	// 打开输入图像文件
	inputFile, err := os.Open("./utils/game_background_4.png")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// 解码图像文件
	inputImg, _, err := image.Decode(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// 创建输出图像，尺寸为原始图像的 1/4
	outputImg := image.NewRGBA(image.Rect(0, 0, inputImg.Bounds().Dx()/4, inputImg.Bounds().Dy()/4))

	// 将原始图像缩小到输出图像中
	for y := 0; y < outputImg.Bounds().Dy(); y++ {
		for x := 0; x < outputImg.Bounds().Dx(); x++ {
			// 计算输入图像中的对应像素位置
			inputX := x * 4
			inputY := y * 4

			// 取输入图像中对应像素的颜色
			outputImg.Set(x, y, inputImg.At(inputX, inputY))
		}
	}

	// 创建输出图像文件
	outputFile, err := os.Create("./utils/game_background_4.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 将输出图像编码为 JPEG 格式，并写入文件
	err = jpeg.Encode(outputFile, outputImg, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Image resized and saved as output.jpg")
}

const (
	screenWidth  = 960
	screenHeight = 540
)

type Game struct {
}

func (g *Game) Update() error {
	for _, th := range utils.AllThings {
		td := th
		td.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for _, th := range utils.AllThings {
		td := th
		td.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	//game := &Game{}
	//
	//ebiten.SetWindowSize(screenWidth, screenHeight)
	//ebiten.SetWindowTitle("Snake Game")
	//utils.AllThings = append(utils.AllThings, utils.NewBackground("./utils/game_background_4.png"))
	//utils.AllThings = append(utils.AllThings, utils.NewFireworks(screenWidth, screenHeight))
	//utils.AllThings = append(utils.AllThings, utils.NewHuman(320, 230, "./utils/1_idle_2.png", 10, "./utils/1_jump.png", 15))
	//if err := ebiten.RunGame(game); err != nil {
	//	log.Fatal(err)
	//}

}
