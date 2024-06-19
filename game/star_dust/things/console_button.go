package things

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"io/ioutil"
	"star_dust/inter"
)

type ConsoleButton struct {
	x         int
	y         int
	width     int
	height    int
	message   string
	isHover   bool
	isPressed bool
	isShow    bool
	clickRun  func() error
	fontFace  font.Face
}

func (h *ConsoleButton) BindTTF(path string) error {
	ttfData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	ttf, err := opentype.Parse(ttfData)
	if err != nil {
		return err
	}
	const dpi = 72
	regularFont, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}
	h.fontFace = regularFont
	return nil
}

func (h *ConsoleButton) SetCoordinate(x int, y int) {
	h.x = x
	h.y = y
}

func (h *ConsoleButton) Show() {
	h.isShow = true
}

func (h *ConsoleButton) Hide() {
	h.isShow = false
}

func (h *ConsoleButton) HoverJudgment() {
	x, y := ebiten.CursorPosition()
	if x > (h.x-25) && x < (h.x-25+h.width) && y > (h.y-24) && y < (h.y-24+h.height) {
		h.isHover = true
	} else {
		h.isHover = false
	}
}

func (h *ConsoleButton) ClickJudgment() {
	if h.isHover {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if !h.isPressed && x > (h.x-25) && x < (h.x-25+h.width) && y > (h.y-24) && y < (h.y-24+h.height) {
				h.isPressed = true
			}
		} else if h.isPressed {
			h.isPressed = false
			h.clickRun()
			h.isHover = false
		}
	} else {
		h.isPressed = false
	}

}

func (h *ConsoleButton) Status() bool {
	return h.isShow
}

func (h *ConsoleButton) Draw(screen *ebiten.Image) {
	if h.isShow {
		outlineOffset := 2
		for dx := -outlineOffset; dx <= outlineOffset; dx++ {
			for dy := -outlineOffset; dy <= outlineOffset; dy++ {
				if dx != 0 || dy != 0 {
					text.Draw(screen, h.message, h.fontFace, h.x+dx, h.y+dy, color.RGBA{0, 0, 0, 255})
				}
			}
		}
		text.Draw(screen, h.message, h.fontFace, h.x, h.y, color.RGBA{255, 255, 255, 255})
		if h.isHover {
			borderColor := color.RGBA{255, 255, 255, 255}
			if h.isPressed {
				borderColor = color.RGBA{255, 0, 0, 255}
			}
			x := h.x - 25
			y := h.y - 25
			h.drawRectBorder(screen, x+1, y+1, h.width, h.height, 2, color.RGBA{0, 0, 0, 90})
			h.drawRectBorder(screen, x, y, h.width, h.height, 3, borderColor)
		}
	}
}

func (h *ConsoleButton) drawRectBorder(screen *ebiten.Image, x, y, width, height, thickness int, clr color.Color) {
	for i := 0; i < thickness; i++ {
		for j := 0; j < width; j++ {
			screen.Set(x+j, y+i, clr)
			screen.Set(x+j, y+height-1-i, clr)
		}
	}
	for i := 0; i < thickness; i++ {
		for j := 0; j < height; j++ {
			screen.Set(x+i, y+j, clr)
			screen.Set(x+width-1-i, y+j, clr)
		}
	}
}

func newConsoleButton(x, y, width, height int, message string, isShow bool, clickRun func() error) inter.Button {
	h := &ConsoleButton{
		message:  message,
		x:        x,
		y:        y,
		width:    width,
		height:   height,
		isShow:   isShow,
		clickRun: clickRun,
	}
	h.BindTTF("/Users/admin/Desktop/star_dust/ttf/shitou.ttf")
	return h
}
