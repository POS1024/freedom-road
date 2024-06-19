package things

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"io/ioutil"
	"star_dust/config"
	"star_dust/consts"
	"star_dust/inter"
)

type VolumeButton struct {
	x               int
	y               int
	volumeBarX      int
	volumeBarY      int
	volumeBarWidth  int
	volumeBarHeight int
	sliderX         int
	sliderY         int
	sliderWidth     int
	sliderHeight    int
	preMouseX       int
	preSliderX      int
	numberX         int
	numberY         int
	volumeNumber    int64
	isPressed       bool
	isShow          bool
	isHoverSlider   bool

	message  string
	fontFace font.Face
}

func (h *VolumeButton) BindTTF(path string) error {
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

func (h *VolumeButton) SetCoordinate(x int, y int) {
	h.x = x
	h.y = y
}

func (h *VolumeButton) Show() {
	var volumeNumber int64 = 0
	volumeNumberInter, ok := config.GameConfigurationConfigurator.Get("Volume")
	if ok {
		volumeNumber = volumeNumberInter.(int64)
	}
	h.volumeNumber = volumeNumber
	h.sliderX = h.volumeBarX + int(volumeNumber*3)
	h.isShow = true
}

func (h *VolumeButton) Hide() {
	var volumeNumber int64 = 0
	volumeNumberInter, ok := config.GameConfigurationConfigurator.Get("Volume")
	if ok {
		volumeNumber = volumeNumberInter.(int64)
	}
	h.volumeNumber = volumeNumber
	h.sliderX = h.volumeBarX + int(volumeNumber*3)
	h.isShow = false
}

func (h *VolumeButton) HoverJudgment() {
	x, y := ebiten.CursorPosition()
	if x > (h.sliderX) && x < (h.sliderX+h.sliderWidth) && y > (h.sliderY) && y < (h.sliderY+h.sliderHeight) {
		h.isHoverSlider = true
	} else {
		h.isHoverSlider = false
	}
}

func (h *VolumeButton) ClickJudgment() {
	if h.isHoverSlider || h.isPressed {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if !h.isPressed {
				h.preMouseX = x
				h.preSliderX = h.sliderX
			} else {
				movedSliderX := h.preSliderX + x - h.preMouseX
				if movedSliderX < h.volumeBarX {
					movedSliderX = h.volumeBarX
				} else if movedSliderX > (h.volumeBarX + h.volumeBarWidth) {
					movedSliderX = h.volumeBarX + h.volumeBarWidth
				}
				h.sliderX = movedSliderX
				h.volumeNumber = int64((h.sliderX - h.volumeBarX) / 3)
			}
			if !h.isPressed && x > (h.sliderX) && x < (h.sliderX+20) && y > (h.sliderY) && y < (h.sliderY+40) {
				h.isPressed = true
			}
		} else if h.isPressed {
			h.isPressed = false
			ok := config.GameConfigurationConfigurator.Set("Volume", h.volumeNumber)
			if !ok {
				consts.ClosingCommand()
			}
		}
	} else {
		h.isPressed = false
	}

}

func (h *VolumeButton) Status() bool {
	return h.isShow
}

func (h *VolumeButton) Draw(screen *ebiten.Image) {
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

		for dx := -outlineOffset; dx <= outlineOffset; dx++ {
			for dy := -outlineOffset; dy <= outlineOffset; dy++ {
				if dx != 0 || dy != 0 {
					text.Draw(screen, fmt.Sprintf("%d", h.volumeNumber), h.fontFace, h.numberX+dx, h.numberY+dy, color.RGBA{0, 0, 0, 255})
				}
			}
		}
		text.Draw(screen, fmt.Sprintf("%d", h.volumeNumber), h.fontFace, h.numberX, h.numberY, color.RGBA{255, 255, 255, 255})

		ebitenutil.DrawRect(screen, float64(h.volumeBarX), float64(h.volumeBarY), float64(h.volumeBarWidth), float64(h.volumeBarHeight), color.RGBA{0xcc, 0xcc, 0xcc, 0xff})
		ebitenutil.DrawRect(screen, float64(h.sliderX), float64(h.sliderY), float64(h.sliderWidth), float64(h.sliderHeight), color.RGBA{0x00, 0x00, 0x00, 0xff})
	}
}

func NewVolumeButton(x, y int, message string, isShow bool) inter.Button {
	var volumeNumber int64 = 0
	volumeNumberInter, ok := config.GameConfigurationConfigurator.Get("Volume")
	if ok {
		volumeNumber = volumeNumberInter.(int64)
	}
	v := &VolumeButton{
		message:         message,
		x:               x,
		y:               y,
		isShow:          isShow,
		volumeBarX:      x + 100,
		volumeBarY:      y - 8,
		volumeBarWidth:  300,
		volumeBarHeight: 5,
		sliderX:         x + 100 + int(volumeNumber*3),
		sliderY:         y - 16,
		sliderWidth:     10,
		sliderHeight:    20,
		numberX:         x + 430,
		numberY:         y,
		volumeNumber:    volumeNumber,
	}
	v.BindTTF("/Users/admin/Desktop/star_dust/ttf/shitou.ttf")
	return v
}
