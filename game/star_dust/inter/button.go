package inter

import (
	"github.com/hajimehoshi/ebiten/v2"
	"star_dust/consts"
)

type Button interface {
	BindTTF(path string) error
	SetCoordinate(x int, y int)
	Show()
	Hide()
	HoverJudgment()
	ClickJudgment()
	Status() bool
	Draw(screen *ebiten.Image)
}

func HideButtons(consoleButtons map[consts.ConsoleButtonKey]Button, keys []consts.ConsoleButtonKey) {
	for _, buttonKey := range keys {
		if buttonInfo, ok := consoleButtons[buttonKey]; ok {
			buttonInfo.Hide()
		}
	}
}

func ShowButtons(consoleButtons map[consts.ConsoleButtonKey]Button, keys []consts.ConsoleButtonKey) {
	for _, buttonKey := range keys {
		if buttonInfo, ok := consoleButtons[buttonKey]; ok {
			buttonInfo.Show()
		}
	}
}
