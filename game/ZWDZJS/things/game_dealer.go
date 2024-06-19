package things

import (
	"ZWDZJS/consts"
	"ZWDZJS/inter"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math/rand"
	"time"
)

type GameDealer struct {
	layers              int
	totalFrame          int64
	gameIdentity        consts.GameIdentityKey
	gameThingsKey       int64
	image               *ebiten.Image
	isShow              bool
	X, Y, Width, Height int
	energy              int
	shovel              *ShovelCard
	cards               []*GameCard
	cor                 color.RGBA
	isShowBorder        bool
}

func (g *GameDealer) GetKey() int64 {
	return g.gameThingsKey
}

func (g *GameDealer) SetKey(key int64) {
	g.gameThingsKey = key
}

func (g *GameDealer) Update() error {
	if g.isShow {
		if g.totalFrame == 120000 {
			g.totalFrame = 0
		}
		if g.totalFrame%600 == 599 {
			rand.Seed(time.Now().UnixNano())
			randomInt := rand.Intn(g.Width)
			NewEnergy(float64(g.X+randomInt), float64(g.Y), float64(g.X+randomInt), float64(g.Y+consts.ScreenHeight-50), 300)
		}
		if g.totalFrame%300 == 299 {
			rand.Seed(time.Now().UnixNano())
			randomInt := rand.Intn(5)
			NewRandomEnemy(consts.LandKey(randomInt+1), consts.EnemyBeginX)
		}
		g.totalFrame++
		for _, card := range g.cards {
			card.Update()
		}
		g.shovel.Update()
	}
	return nil
}

func (g *GameDealer) Draw(screen *ebiten.Image) {
	if g.isShow {
		if g.isShowBorder {
			g.drawRectOutline(screen, g.X, g.Y, g.Width, g.Height, g.cor)
		}
		for _, card := range g.cards {
			card.Draw(screen)
		}
		g.shovel.Draw(screen)
		g.drawHealthBar(screen)
	}
}

func (g *GameDealer) drawHealthBar(screen *ebiten.Image) {
	energyText := fmt.Sprintf("%d", g.energy)
	font := basicfont.Face7x13
	text.Draw(screen, energyText, font, consts.GameDealerEnergyX, consts.GameDealerEnergyY, color.White)
}

func (g *GameDealer) GetLayers() int {
	return g.layers
}

func (g *GameDealer) GetEnergy() int {
	return g.energy
}

func (g *GameDealer) AddEnergy(energy int) bool {
	g.energy = g.energy + energy
	return true
}

func (g *GameDealer) ConsumeEnergy(energy int) bool {
	if g.energy < energy {
		return false
	}
	g.energy = g.energy - energy
	return true
}

func (g *GameDealer) SearchEnergy() int {
	return g.energy
}

func (g *GameDealer) drawRectOutline(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x+width), float64(y), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y+height), float64(x+width), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x+width), float64(y), float64(x+width), float64(y+height), clr)
}

func NewGameDealer() *GameDealer {
	m := &GameDealer{
		X:            consts.GameDealerX,
		Y:            consts.GameDealerY,
		Width:        consts.GameDealerWidth,
		Height:       consts.GameDealerHeight,
		isShow:       true,
		layers:       consts.LayersConsoleKey,
		cards:        NewGameCards(consts.GameCardX, consts.GameCardY, consts.GameCardWidth, consts.GameCardHeight),
		cor:          color.RGBA{255, 255, 0, 255},
		isShowBorder: true,
		energy:       2000,
		shovel:       NewShovelCard(consts.ShovelCardX, consts.ShovelCardY),
	}
	inter.GameThings.Put(m)
	return m
}

var GameDealerMR = NewGameDealer()
