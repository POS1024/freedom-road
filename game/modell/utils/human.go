package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type Human struct {
	x               float64
	y               float64
	image           *ebiten.Image
	beginX          float64
	imageNum        int
	waitTimes       int
	maxImageNum     int
	idleImage       *ebiten.Image
	idleMaxImageNum int
	jumpImage       *ebiten.Image
	jumpMaxImageNum int
	isJump          bool
}

func (h *Human) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !h.isJump {
		h.isJump = true
		h.image = h.jumpImage
		h.maxImageNum = h.jumpMaxImageNum
		h.beginX = 0
		h.waitTimes = 0
		h.imageNum = 0
	}
	if h.isJump {
		if h.imageNum >= 4 && h.imageNum <= 6 {
			h.x = h.x + 3
		}
		if h.imageNum >= 7 && h.imageNum <= 9 {
			h.x = h.x + 2
		}
		if h.imageNum == 4 {
			h.y = h.y - 5
		}
		if h.imageNum == 6 {
			h.y = h.y + 5
		}
		if h.imageNum == h.maxImageNum {
			h.isJump = false
			h.image = h.idleImage
			h.maxImageNum = h.idleMaxImageNum
			h.beginX = 0
			h.waitTimes = 0
			h.imageNum = 0
		}
	}
	if h.waitTimes < 8 {
		h.waitTimes++
	} else {
		h.waitTimes = 0
		if h.imageNum == h.maxImageNum {
			h.beginX = 0
			h.imageNum = 0
		} else {
			h.beginX = h.beginX + 128
			h.imageNum++
		}
	}

	return nil
}

func (h *Human) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(h.x, h.y)
	subImage := h.image.SubImage(image.Rect(int(h.beginX), 0, int(h.beginX)+128, 128)).(*ebiten.Image)
	screen.DrawImage(subImage, op)
}

func NewHuman(x, y float64, idleImagePath string, idleMaxImageNum int, jumpImagePath string, jumpXaxImageNum int) *Human {
	idleImage, _, _ := ebitenutil.NewImageFromFile(idleImagePath)
	jumpImage, _, _ := ebitenutil.NewImageFromFile(jumpImagePath)

	return &Human{
		x:               x,
		y:               y,
		image:           idleImage,
		maxImageNum:     idleMaxImageNum,
		idleImage:       idleImage,
		idleMaxImageNum: idleMaxImageNum,
		jumpImage:       jumpImage,
		jumpMaxImageNum: jumpXaxImageNum,
	}
}
