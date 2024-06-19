package things

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"star_dust/consts"
	"star_dust/inter"
)

type ConsoleManager struct {
	layers           int
	gameThingsKey    int64
	image            *ebiten.Image
	isShowBackground bool
	consoleButtons   map[consts.ConsoleButtonKey]inter.Button
}

func (c *ConsoleManager) bindBackground(path string) error {
	image, _, _ := ebitenutil.NewImageFromFile(path)
	c.image = image
	return nil
}

func (c *ConsoleManager) initButtons() {
	c.consoleButtons = make(map[consts.ConsoleButtonKey]inter.Button)
	InitHomeButtons(c.consoleButtons, c)
	InitGameButtons(c.consoleButtons, c)
	InitSettingButtons(c.consoleButtons, c)
}

func (c *ConsoleManager) updateButtonHoverCapture() {
	for _, info := range c.consoleButtons {
		if info.Status() {
			info.HoverJudgment()
		}
	}
}

func (c *ConsoleManager) updateButtonPressed() {
	for _, info := range c.consoleButtons {
		if info.Status() {
			info.ClickJudgment()
		}
	}
}

func (c *ConsoleManager) drawBackground(screen *ebiten.Image) {
	if c.isShowBackground {
		screen.DrawImage(c.image, nil)
	}
}

func (c *ConsoleManager) drawOverlay(screen *ebiten.Image) {
	overlay := ebiten.NewImage(screen.Size())
	overlay.Fill(color.RGBA{0, 0, 0, 30})
	screen.DrawImage(overlay, nil)
}

func (c *ConsoleManager) drawTextButton(screen *ebiten.Image) {
	for _, info := range c.consoleButtons {
		if info.Status() {
			info.Draw(screen)
		}
	}
}

func (c *ConsoleManager) GetKey() int64 {
	return c.gameThingsKey
}

func (c *ConsoleManager) SetKey(key int64) {
	c.gameThingsKey = key
	return
}

func (c *ConsoleManager) Update() error {
	c.updateButtonHoverCapture()
	c.updateButtonPressed()
	return nil
}

func (c *ConsoleManager) GetLayers() int {
	return c.layers
}

func (c *ConsoleManager) Draw(screen *ebiten.Image) {
	c.drawBackground(screen)
	c.drawTextButton(screen)
}

func (c *ConsoleManager) ShowBackground() {
	c.isShowBackground = true
}

func (c *ConsoleManager) HideBackground() {
	c.isShowBackground = false
}

func NewConsoleManager() *ConsoleManager {
	c := &ConsoleManager{
		layers: consts.LayersConsoleKey,
	}
	c.bindBackground("/Users/admin/Desktop/star_dust/image/back_tree.png")
	c.initButtons()
	c.ShowBackground()
	inter.GameThings.Put(c)
	return c
}

var CM = NewConsoleManager()
