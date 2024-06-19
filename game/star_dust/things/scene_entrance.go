package things

import (
	"github.com/hajimehoshi/ebiten/v2"
	"star_dust/consts"
	"star_dust/inter"
)

type Entrance struct {
	X, Y          int
	gameIdentity  consts.GameIdentityKey
	width, height int
	uuid          string
	isInWindow    bool
	image         *ebiten.Image
}

func (e *Entrance) GetImage() *ebiten.Image {
	return e.image
}

func (e *Entrance) Move(nextX int, nextY int) {
	if e.X != nextX || e.Y != nextY {
		e.X, e.Y = nextX, nextY
		QM.MoveGameObject(e)
	}
}

func (e *Entrance) GetPictureFrame() int {
	return 0
}

func (e *Entrance) GetXY() (int, int) {
	return e.X, e.Y
}

func (e *Entrance) SetXY(x int, y int) {
	e.X, e.Y = x, y
}

func (e *Entrance) GetWidth() int {
	return e.width
}

func (e *Entrance) GetHeight() int {
	return e.height
}

func (e *Entrance) GetUuid() string {
	return e.uuid
}

func (e *Entrance) GetGameIdentity() consts.GameIdentityKey {
	return e.gameIdentity
}

func (e *Entrance) Update() error {
	if e.isInWindow {
		if !(e.X+2*e.width > 0 && e.Y+2*e.height > 0 && e.X-e.width < inter.ScreenWidth && e.Y-e.height < inter.ScreenHeight) {
			e.isInWindow = false
		}
	} else {
		if e.X+2*e.width > 0 && e.Y+2*e.height > 0 && e.X-e.width < inter.ScreenWidth && e.Y-e.height < inter.ScreenHeight {
			e.isInWindow = true
		}
	}

	if e.isInWindow {
		// 特效
	}
	return nil
}

func (e *Entrance) Draw(screen *ebiten.Image) {
	if e.isInWindow {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(e.X), float64(e.Y))
		screen.DrawImage(e.image, opts)
	}
}
