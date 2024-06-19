package inter

import "github.com/hajimehoshi/ebiten/v2"

type Entrance interface {
	TriggerJudgment()
	Update() error
	Draw(screen *ebiten.Image)
	GetOffsetXY() (int int)
}
