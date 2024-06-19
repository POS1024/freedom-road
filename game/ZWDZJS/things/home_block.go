package things

import (
	"ZWDZJS/consts"
	"ZWDZJS/inter"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type HomeBlock struct {
	layers              int
	gameIdentity        consts.GameIdentityKey
	gameThingsKey       int64
	X, Y, Width, Height int
	isShow              bool
	cor                 color.RGBA
	isShowBorder        bool
}

func (l *HomeBlock) GetKey() int64 {
	return l.gameThingsKey
}

func (l *HomeBlock) SetKey(key int64) {
	l.gameThingsKey = key
}

func (l *HomeBlock) Update() error {
	return nil
}

func (l *HomeBlock) Draw(screen *ebiten.Image) {
	if l.isShow {
		if l.isShowBorder {
			l.drawRectOutline(screen, l.X, l.Y, l.Width, l.Height, l.cor)
		}
	}
}

func (l *HomeBlock) drawRectOutline(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x+width), float64(y), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y+height), float64(x+width), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x+width), float64(y), float64(x+width), float64(y+height), clr)
}

func (l *HomeBlock) GetLayers() int {
	return l.layers
}

func NewHomeBlock(X, Y, Width, Height int, layer int) *HomeBlock {
	l := &HomeBlock{
		X:      X,
		Y:      Y,
		Width:  Width,
		Height: Height,
		isShow: true,
		layers: layer,
		cor:    color.RGBA{255, 255, 0, 255},
	}
	inter.GameThings.Put(l)
	return l
}

var HomeLand = NewHomeBlock(consts.HomeLandX, consts.HomeLandY, consts.HomeLandWidth, consts.HomeLandHeight, consts.LayersHomeKey)
