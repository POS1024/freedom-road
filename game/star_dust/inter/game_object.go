package inter

import (
	"github.com/hajimehoshi/ebiten/v2"
	"star_dust/consts"
)

type GameObject interface {
	GetImage() *ebiten.Image
	GetXY() (int, int)
	SetXY(int, int)
	GetWidth() int
	GetHeight() int
	GetUuid() string
	GetPictureFrame() int
	GetGameIdentity() consts.GameIdentityKey
	Move(int, int)
}
