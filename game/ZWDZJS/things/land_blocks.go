package things

import (
	"ZWDZJS/consts"
	"ZWDZJS/inter"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type LandBlock struct {
	layers                  int
	gameIdentity            consts.GameIdentityKey
	gameThingsKey           int64
	X, Y, Width, Height     int
	isShow, isHover, isDown bool
	P                       *Placement
	cor                     color.RGBA
	isShowBorder            bool
	MouseType               consts.MouseTypeKey
	line                    consts.LandKey
	row                     consts.LandKey
	enemies                 []*Enemy
}

func (l *LandBlock) GetLineRow() (line consts.LandKey, row consts.LandKey) {
	return l.line, l.row
}

func (l *LandBlock) GetMouseType() consts.MouseTypeKey {
	return l.MouseType
}

func (l *LandBlock) Cancel() {

}

func (l *LandBlock) Confirm() bool {
	return true
}

func (l *LandBlock) MouseEffects(screen *ebiten.Image) {

}

func (l *LandBlock) CollectEnemies() {
	enemies := Enemies[l.line]
	l.enemies = make([]*Enemy, 0)
	for _, enemyInfo := range enemies {
		X, _ := enemyInfo.GetXY()
		W := enemyInfo.GetWidth()
		if l.X <= (X+W) && (l.X+l.Width) >= X {
			l.enemies = append(l.enemies, enemyInfo)
		}
	}
}

func (l *LandBlock) GetEnemies() ([]*Enemy, int) {
	return l.enemies, len(l.enemies)
}

func (l *LandBlock) IsHover() bool {
	x, y := ebiten.CursorPosition()
	if x > l.X && x < l.X+l.Width && y > l.Y && y < l.Y+l.Height {
		l.cor = color.RGBA{255, 0, 0, 255}
		l.isHover = true
		return true
	} else {
		l.cor = color.RGBA{255, 255, 0, 255}
		l.isHover = false
		l.isDown = false
		return false
	}
}

func (l *LandBlock) IsClick() bool {
	if l.isHover {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			l.isDown = true
		} else {
			if l.isDown {
				l.isDown = false
				if MOUSE.GetMouseType() == consts.MouseTypeNil {
					return false
				}
				if MOUSE.GetMouseType() == consts.MouseTypeAddCard && l.P == nil {
					if MOUSE.ConfirmMousePuppet() {
						l.P = NewRandomPlacement(l.X, l.Y, l)
						return true
					}
				}
				if MOUSE.GetMouseType() == consts.MouseTypeRemoveCard && l.P != nil {
					if MOUSE.ConfirmMousePuppet() {
						l.P.Delete()
						l.P = nil
						return true
					}
				}
			}
		}
	}
	return false
}

func (l *LandBlock) IsDown() bool {
	return true
}

func (l *LandBlock) GetKey() int64 {
	return l.gameThingsKey
}

func (l *LandBlock) SetKey(key int64) {
	l.gameThingsKey = key
}

func (l *LandBlock) Update() error {
	if l.isShow {
		l.IsHover()
		l.IsClick()
		l.CollectEnemies()
	}
	return nil
}

func (l *LandBlock) Draw(screen *ebiten.Image) {
	if l.isShow {
		if l.isShowBorder {
			l.drawRectOutline(screen, l.X, l.Y, l.Width, l.Height, l.cor)
		}
	}
}

func (l *LandBlock) drawRectOutline(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x+width), float64(y), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y+height), float64(x+width), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x+width), float64(y), float64(x+width), float64(y+height), clr)
}

func (l *LandBlock) GetLayers() int {
	return l.layers
}

func NewLandBlock(line, row consts.LandKey, X, Y, Width, Height int, layer int) *LandBlock {
	l := &LandBlock{
		X:            X,
		Y:            Y,
		Width:        Width,
		Height:       Height,
		isShow:       true,
		layers:       layer,
		cor:          color.RGBA{255, 255, 0, 255},
		isShowBorder: true,
		line:         line,
		row:          row,
		enemies:      make([]*Enemy, 0),
	}
	inter.GameThings.Put(l)
	return l
}

func NewLandBlocks() map[consts.LandKey]map[consts.LandKey]*LandBlock {
	lb := make(map[consts.LandKey]map[consts.LandKey]*LandBlock, 0)
	for i := 1; i < 6; i++ {
		lbInfo := make(map[consts.LandKey]*LandBlock, 0)
		for j := 1; j < 10; j++ {
			lbInfo[consts.LandKey(j)] = NewLandBlock(consts.LandKey(i), consts.LandKey(j), consts.PlantLandX+(j-1)*consts.PlantLandWidth+5, consts.PlantLandY+(i-1)*consts.PlantLandHeight+5, consts.PlantLandWidth-10, consts.PlantLandHeight-10, consts.LayersLandKey)
		}
		lb[consts.LandKey(i)] = lbInfo
	}

	return lb
}

func MatchLandBlocks(line consts.LandKey, X, Width int) (matchBlocks []*LandBlock) {
	blocks := LandBlocker[line]
	for _, blockInfo := range blocks {
		blockX := blockInfo.X + 30
		blockW := blockInfo.Width - 60
		if blockX <= (X+Width) && (blockX+blockW) >= X {
			matchBlocks = append(matchBlocks, blockInfo)
		}
	}
	return
}

var LandBlocker = NewLandBlocks()
