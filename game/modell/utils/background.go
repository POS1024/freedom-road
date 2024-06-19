package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Background struct {
	image *ebiten.Image
}

func (b *Background) Update() error {
	return nil
}

func (b *Background) Draw(screen *ebiten.Image) {
	screen.DrawImage(b.image, nil)
}

func NewBackground(imagePath string) *Background {
	image, _, _ := ebitenutil.NewImageFromFile(imagePath)
	return &Background{
		image: image,
	}
}
