package inter

import (
	"ZWDZJS/consts"
	"github.com/hajimehoshi/ebiten/v2"
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
