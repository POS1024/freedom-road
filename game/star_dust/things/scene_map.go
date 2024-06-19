package things

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"star_dust/inter"
)

type SceneMap struct {
	offsetX, offsetY int
	width, height    int
	backgroundImage  *ebiten.Image
	innerObjs        inter.GameObject
	entrances        []Entrance
}

func (s *SceneMap) Update() error {
	for _, entrance := range s.entrances {
		entrance.Update()
	}
	return nil
}

func (s *SceneMap) Draw(screen *ebiten.Image) {
	subRect := image.Rect(s.offsetX, 0, s.offsetX+inter.ScreenWidth, s.height)
	subImg := s.backgroundImage.SubImage(subRect).(*ebiten.Image)
	screen.DrawImage(subImg, nil)
	for _, entrance := range s.entrances {
		entrance.Draw(screen)
	}
}

func (s *SceneMap) Offset(offsetX, offsetY int, obj inter.GameObject) bool {
	if (offsetX > 0 && s.offsetX+offsetX+inter.ScreenWidth > s.width-10) || (offsetX < 0 && s.offsetX+offsetX < 10) {
		return false
	}
	s.offsetX, s.offsetY = s.offsetX+offsetX, s.offsetY+offsetY
	s.innerObjs = obj
	return true
}

func NewSceneMap(offsetX, offsetY int, width, height int, backgroundImagePath string) *SceneMap {
	image, _, _ := ebitenutil.NewImageFromFile(backgroundImagePath)
	s := &SceneMap{
		offsetX:         offsetX,
		offsetY:         offsetY,
		width:           width,
		height:          height,
		backgroundImage: image,
	}
	return s
}
