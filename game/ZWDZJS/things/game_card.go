package things

import "C"
import (
	"ZWDZJS/caches"
	"ZWDZJS/consts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type GameCard struct {
	X, Y, Width, Height, energy       int
	image                             *ebiten.Image
	cor                               color.RGBA
	isShow, isHover, isSelect, isDown bool
	isShowBorder                      bool
	MouseType                         consts.MouseTypeKey
}

func (c *GameCard) GetMouseType() consts.MouseTypeKey {
	return c.MouseType
}

func (c *GameCard) Cancel() {
	c.isSelect = false
}

func (c *GameCard) Confirm() bool {
	c.isSelect = false
	GameDealerMR.ConsumeEnergy(c.energy)
	return true
}

func (c *GameCard) MouseEffects(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(c.image, opts)
}

func (c *GameCard) IsHover() bool {
	x, y := ebiten.CursorPosition()
	if x > c.X && x < c.X+c.Width && y > c.Y && y < c.Y+c.Height {
		c.isHover = true
		c.cor = color.RGBA{255, 0, 0, 255}
		return true
	} else {
		c.isHover = false
		c.isDown = false
		c.cor = color.RGBA{255, 255, 0, 255}
		return false
	}
}

func (c *GameCard) IsDown() bool {
	if c.isHover {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			c.isDown = true
		} else {
			if c.isDown {
				c.isDown = false
				if GameDealerMR.GetEnergy() >= c.energy {
					if !c.isSelect {
						c.isSelect = true
						MOUSE.BindMousePuppet(c)
						return true
					}
					c.isSelect = false
					MOUSE.CancelMousePuppet()
				}
			}
		}
	}
	return false
}

func (c *GameCard) Update() error {
	if c.isShow {
		c.IsHover()
		c.IsDown()
	}
	return nil
}

func (c *GameCard) Draw(screen *ebiten.Image) {
	if c.isShow {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(c.X), float64(c.Y))
		screen.DrawImage(c.image, opts)
		if c.isSelect {
			c.cor = color.RGBA{255, 255, 255, 255}
		}
		if c.isShowBorder {
			c.drawRectOutline(screen, c.X, c.Y, c.Width, c.Height, c.cor)
		}
	}
}

func (c *GameCard) drawRectOutline(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x+width), float64(y), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y+height), float64(x+width), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x+width), float64(y), float64(x+width), float64(y+height), clr)
}

func (c *GameCard) bindImage(path string) error {
	if imageCache, ok := caches.ImageCache[path]; ok {
		c.image = imageCache.LeftImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(path)
		caches.ImageCache[path] = caches.ImageCacheInfo{
			LeftImage: image,
		}
		c.image = image
	}
	return nil
}

func NewGameCard(X, Y, Width, Height int) *GameCard {
	m := &GameCard{
		X:            X,
		Y:            Y,
		Width:        Width,
		Height:       Height,
		cor:          color.RGBA{255, 255, 0, 255},
		isShow:       true,
		isShowBorder: true,
		energy:       400,
		MouseType:    consts.MouseTypeAddCard,
	}
	m.bindImage("/Users/admin/Desktop/ZWDZJS/image/card/card1.png")
	return m
}

func NewGameCards(X, Y, Width, Height int) []*GameCard {
	cards := make([]*GameCard, 0)
	for i := 0; i < 4; i++ {
		cards = append(cards, NewGameCard(X+i*100, Y, Width, Height))
	}
	return cards
}
