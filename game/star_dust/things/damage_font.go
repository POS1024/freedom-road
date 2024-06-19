package things

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"io/ioutil"
	"star_dust/consts"
	"star_dust/inter"
)

type DamageFont struct {
	layers        int
	gameIdentity  consts.GameIdentityKey
	gameThingsKey int64
	mainId        string
	uuid          string
	X, Y, Number  int
	Frame         int
	FrameSpeed    int
	isShow        bool
	fontFace      font.Face
}

func (d *DamageFont) GetKey() int64 {
	return d.gameThingsKey
}

func (d *DamageFont) SetKey(key int64) {
	d.gameThingsKey = key
}

func (d *DamageFont) Update() error {
	if d.isShow {
		if d.Frame > 20 {
			d.isShow = false
			inter.GameThings.Delete(d.gameThingsKey)
			return nil
		}
		if d.FrameSpeed%2 == 1 {
			d.Y = d.Y - 2
			d.Frame++
		}
		d.FrameSpeed++

	}
	return nil
}

func (d *DamageFont) Draw(screen *ebiten.Image) {
	if d.isShow {
		healthText := fmt.Sprintf("-%d", d.Number)
		text.Draw(screen, healthText, d.fontFace, d.X, d.Y, color.RGBA{255, 0, 0, 255})
	}
}

func (d *DamageFont) GetLayers() int {
	return d.layers
}

func (d *DamageFont) BindTTF(path string) error {
	if fontCache, ok := FontCache[path]; ok {
		d.fontFace = fontCache.Face
	} else {
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

		FontCache[path] = FontCacheInfo{
			Face: regularFont,
		}
		d.fontFace = regularFont
	}
	return nil
}

func NewDamageFont(X, Y, Number int) *DamageFont {
	d := &DamageFont{
		X:      X,
		Y:      Y,
		Number: Number,
		isShow: true,
	}
	d.BindTTF("/Users/admin/Desktop/star_dust/ttf/shitou.ttf")
	inter.GameThings.GoPut(d)
	return d
}
