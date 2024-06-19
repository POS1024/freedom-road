package utils

import "github.com/hajimehoshi/ebiten/v2"

type Thing interface {
	Update() error
	Draw(screen *ebiten.Image)
}

var AllThings []Thing
