package inter

import "github.com/hajimehoshi/ebiten/v2"

type InteractiveObjects interface {
	ClickJudgment()
	Update() error
	Draw(screen *ebiten.Image)
}
