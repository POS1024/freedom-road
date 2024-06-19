package inter

import (
	"ZWDZJS/consts"
	"github.com/hajimehoshi/ebiten/v2"
)

type MousePuppet interface {
	GetMouseType() consts.MouseTypeKey
	Cancel()
	Confirm() bool
	MouseEffects(screen *ebiten.Image)
	IsHover() bool
	IsDown() bool
}
