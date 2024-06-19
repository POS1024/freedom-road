package things

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"star_dust/consts"
	"star_dust/inter"
)

type SkillEffects struct {
	layers               int
	gameIdentity         consts.GameIdentityKey
	gameThingsKey        int64
	mainId               string
	uuid                 string
	X, Y, Width, Height  int
	isShow               bool
	pictureFrame         int
	pictureFrameCount    int
	skillImage           *ebiten.Image
	skillImageFrame      int
	skillImageFrameSpeed int
	After                func()
	afterFrame           int
}

func (s *SkillEffects) GetImage() *ebiten.Image {
	subRect := image.Rect(s.pictureFrame*s.Width, 0, s.pictureFrame*s.Width+s.Width, s.Height)
	subImg := s.skillImage.SubImage(subRect).(*ebiten.Image)
	return subImg
}

func (s *SkillEffects) Move(nextX int, nextY int) {
	if s.X != nextX || s.Y != nextY {
		s.X, s.Y = nextX, nextY
		QM.MoveGameObject(s)
	}
}

func (s *SkillEffects) GetPictureFrame() int {
	return s.pictureFrame
}

func (s *SkillEffects) GetXY() (int, int) {
	return s.X, s.Y
}

func (s *SkillEffects) SetXY(x int, y int) {
	s.X = x
	s.Y = y
}

func (s *SkillEffects) GetWidth() int {
	return s.Width
}

func (s *SkillEffects) GetHeight() int {
	return s.Height
}

func (s *SkillEffects) GetUuid() string {
	return s.uuid
}

func (s *SkillEffects) GetKey() int64 {
	return s.gameThingsKey
}

func (s *SkillEffects) SetKey(key int64) {
	s.gameThingsKey = key
}

func (s *SkillEffects) GetGameIdentity() consts.GameIdentityKey {
	return s.gameIdentity
}

func (s *SkillEffects) Update() error {
	if s.isShow {
		if s.pictureFrameCount%s.skillImageFrameSpeed == s.skillImageFrameSpeed-1 {
			if s.pictureFrame == s.afterFrame {
				s.After()
			}
			if s.pictureFrame == s.skillImageFrame-1 {
				s.pictureFrame = 0
				s.isShow = false
				QM.RemoveGameObject(s)
				inter.GameThings.Delete(s.gameThingsKey)
			} else {
				s.pictureFrame++
			}
		}
		if s.pictureFrameCount == s.skillImageFrameSpeed*100 {
			s.pictureFrameCount = 0
		}
		s.pictureFrameCount++

	}
	return nil
}

func (s *SkillEffects) Show() {
	s.isShow = true
}

func (s *SkillEffects) Hide() {
	s.isShow = false
}

func (s *SkillEffects) Draw(screen *ebiten.Image) {
	if s.isShow {
		subRect := image.Rect(s.pictureFrame*s.Width, 0, s.pictureFrame*s.Width+s.Width, s.Height)
		subImg := s.skillImage.SubImage(subRect).(*ebiten.Image)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(s.X), float64(s.Y))
		screen.DrawImage(subImg, opts)
	}
}

func (s *SkillEffects) GetLayers() int {
	return s.layers
}

func (s *SkillEffects) bindSkillImage(path string, imageFrame int) error {
	if imageCache, ok := ImageCache[path]; ok {
		s.skillImage = imageCache.leftImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(path)
		ImageCache[path] = ImageCacheInfo{
			leftImage: image,
		}
		s.skillImage = image
	}
	s.skillImageFrame = imageFrame
	s.skillImageFrameSpeed = 6
	return nil
}

func NewCircleExplosion(x, y int, mainId string, after func()) *SkillEffects {
	s := &SkillEffects{
		mainId:       mainId,
		X:            x,
		Y:            y,
		layers:       consts.LayersSkillKey,
		Width:        256,
		Height:       256,
		uuid:         uuid.New().String(),
		isShow:       true,
		gameIdentity: consts.GameIdentityEffect,
		After:        after,
		afterFrame:   0,
	}
	s.bindSkillImage("/Users/admin/Desktop/star_dust/image/skill_effects/circle_explosion.png", 10)
	QM.AddGameObject(s, true)
	inter.GameThings.Put(s)
	return s
}

func NewLightningBeginning(x, y int, mainId string, after func()) *SkillEffects {
	s := &SkillEffects{
		mainId:       mainId,
		X:            x,
		Y:            y,
		layers:       consts.LayersSkillKey,
		Width:        64,
		Height:       193,
		uuid:         uuid.New().String(),
		isShow:       true,
		gameIdentity: consts.GameIdentityEffect,
		After:        after,
		afterFrame:   2,
	}
	s.bindSkillImage("/Users/admin/Desktop/star_dust/image/skill_effects/lightning_beginning.png", 7)
	QM.AddGameObject(s, true)
	inter.GameThings.Put(s)
	return s
}
