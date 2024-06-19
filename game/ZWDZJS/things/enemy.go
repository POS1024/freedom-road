package things

import (
	"ZWDZJS/caches"
	"ZWDZJS/consts"
	"ZWDZJS/inter"
	"ZWDZJS/utils"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"math/rand"
	"time"
)

type Enemy struct {
	layers                              int
	gameIdentity                        consts.GameIdentityKey
	line                                consts.LandKey
	gameThingsKey                       int64
	uuid                                string
	X, Y, Width, Height, Speed, harm    int
	isShow                              bool
	pictureFrame                        int
	pictureFrameCount                   int
	isRun                               bool
	isShot                              bool
	image                               *ebiten.Image
	imageFrame, imageFrameSpeed         int
	runImage                            *ebiten.Image
	runImageFrame, runImageFrameSpeed   int
	shotImage                           *ebiten.Image
	shotImageFrame, shotImageFrameSpeed int
	Health                              int
	isShowBorder                        bool
	ShotFunc                            func(*Enemy)
}

func (m *Enemy) DamageCalculation(damage int, face int) (die bool, ex int) {
	if m.Health-damage <= 0 {
		m.Health = 0
	} else {
		m.Health = m.Health - damage
	}
	return false, 0
}

func (m *Enemy) GetImage() *ebiten.Image {
	subRect := image.Rect(m.pictureFrame*m.Width, 0, m.pictureFrame*m.Width+m.Width, m.Height)
	subImg := m.image.SubImage(subRect).(*ebiten.Image)
	return subImg
}

func (m *Enemy) Move(nextX int, nextY int) {
	if m.X != nextX || m.Y != nextY {
		m.X, m.Y = nextX, nextY
		QM.MoveGameObject(m)
	}
}

func (m *Enemy) GetPictureFrame() int {
	return m.pictureFrame
}

func (m *Enemy) GetXY() (int, int) {
	return m.X, m.Y
}

func (m *Enemy) SetXY(x int, y int) {
	m.X = x
	m.Y = y
}

func (m *Enemy) IsShot() (isShot bool) {
	blocks := MatchLandBlocks(m.line, m.X, m.Width)
	for _, blockInfo := range blocks {
		if blockInfo.P != nil {
			isShot = true
			return
		}
	}
	return
}

func (m *Enemy) Delete() {
	m.isShow = false
	delete(Enemies[m.line], m.uuid)
	QM.RemoveGameObject(m)
	inter.GameThings.Delete(m.gameThingsKey)
}

func (m *Enemy) GetWidth() int {
	return m.Width
}

func (m *Enemy) GetHeight() int {
	return m.Height
}

func (m *Enemy) GetUuid() string {
	return m.uuid
}

func (m *Enemy) GetKey() int64 {
	return m.gameThingsKey
}

func (m *Enemy) SetKey(key int64) {
	m.gameThingsKey = key
}

func (m *Enemy) GetGameIdentity() consts.GameIdentityKey {
	return m.gameIdentity
}

func (m *Enemy) Update() error {
	if m.isShow {

		if m.Health <= 0 || m.X < 0 {
			m.Delete()
			return nil
		}

		isShot := m.IsShot()

		if m.isShot != isShot {
			m.pictureFrame = 0
			m.pictureFrameCount = 0
		}
		m.isShot = isShot
		m.isRun = !m.isShot

		if m.isRun {
			m.imageFrame = m.runImageFrame
			m.imageFrameSpeed = m.runImageFrameSpeed
			m.image = m.runImage
			m.Move(m.X-m.Speed, m.Y)
		} else {
			m.imageFrame = m.shotImageFrame
			m.imageFrameSpeed = m.shotImageFrameSpeed
			m.image = m.shotImage
			if m.pictureFrame == 0 && m.pictureFrameCount%m.imageFrameSpeed == m.imageFrameSpeed-1 {
				blocks := MatchLandBlocks(m.line, m.X, m.Width)
				for _, blockInfo := range blocks {
					if blockInfo.P != nil {
						blockInfo.P.DamageCalculation(m.harm, -1)
						break
					}
				}
				//NewBullet(m.X+m.Width, m.Y+(m.Height/2)-1+15, 1, m.lv*10, m.uuid, m)
			}
		}
		if m.pictureFrameCount%m.imageFrameSpeed == m.imageFrameSpeed-1 {

			if m.pictureFrame == m.imageFrame-1 {
				m.pictureFrame = 0
			} else {
				m.pictureFrame++
			}
		}
		if m.pictureFrameCount == m.imageFrameSpeed*100 {
			m.pictureFrameCount = 0
		}
		m.pictureFrameCount++
	}

	return nil
}

func (m *Enemy) Draw(screen *ebiten.Image) {
	if m.isShow {
		subRect := image.Rect(m.pictureFrame*m.Width, 0, m.pictureFrame*m.Width+m.Width, m.Height)
		subImg := m.image.SubImage(subRect).(*ebiten.Image)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(m.X), float64(m.Y))
		screen.DrawImage(subImg, opts)
		if m.isShowBorder {
			m.drawRectOutline(screen, m.X, m.Y, m.Width, m.Height, color.RGBA{255, 255, 255, 255})
		}
	}
}

func (m *Enemy) drawRectOutline(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x+width), float64(y), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y+height), float64(x+width), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x+width), float64(y), float64(x+width), float64(y+height), clr)
}

func (m *Enemy) GetLayers() int {
	return m.layers
}

func (m *Enemy) bindRunImage(info EnemyImageInfo) error {
	if info.RunImagePath == "" {
		return nil
	}
	if imageCache, ok := caches.ImageCache[info.RunImagePath]; ok {
		m.runImage = imageCache.LeftImage
	} else {
		imageRight, _, _ := ebitenutil.NewImageFromFile(info.RunImagePath)
		image := utils.FlipImageSegments(imageRight, 128)
		caches.ImageCache[info.RunImagePath] = caches.ImageCacheInfo{
			LeftImage: image,
		}
		m.runImage = image
	}
	m.runImageFrame = info.RunImageFrame
	m.runImageFrameSpeed = info.RunImageFrameSpeed

	m.image = m.runImage
	m.imageFrame = info.RunImageFrame
	m.imageFrameSpeed = info.RunImageFrameSpeed

	return nil
}

func (m *Enemy) bindShotImage(info EnemyImageInfo) error {
	if info.ShotImagePath == "" {
		return nil
	}
	if imageCache, ok := caches.ImageCache[info.ShotImagePath]; ok {
		m.shotImage = imageCache.LeftImage
	} else {
		imageRight, _, _ := ebitenutil.NewImageFromFile(info.ShotImagePath)
		image := utils.FlipImageSegments(imageRight, 128)
		caches.ImageCache[info.ShotImagePath] = caches.ImageCacheInfo{
			LeftImage: image,
		}
		m.shotImage = image
	}
	m.shotImageFrame = info.ShotImageFrame
	m.shotImageFrameSpeed = info.ShotImageFrameSpeed
	return nil
}

type EnemyImageInfo struct {
	RunImagePath        string
	RunImageFrame       int
	RunImageFrameSpeed  int
	ShotImagePath       string
	ShotImageFrame      int
	ShotImageFrameSpeed int
}

func NewEnemy(line consts.LandKey, X, Width, Height, Health, Speed, harm int, image EnemyImageInfo, shotFunc func(placement *Enemy)) *Enemy {
	Y := consts.PlantLandY + (int(line)-1)*consts.PlantLandHeight + consts.PlantLandHeight - 133
	m := &Enemy{
		X:            X,
		Y:            Y,
		layers:       consts.LayersEnemyKey,
		Width:        Width,
		Height:       Height,
		uuid:         uuid.New().String(),
		isRun:        true,
		Health:       Health,
		gameIdentity: consts.GameIdentityEnemy,
		isShow:       true,
		ShotFunc:     shotFunc,
		line:         line,
		isShowBorder: true,
		Speed:        Speed,
		harm:         harm,
	}
	m.bindRunImage(image)
	m.bindShotImage(image)
	inter.GameThings.Put(m)
	QM.AddGameObject(m, true)
	Enemies[line][m.uuid] = m
	return m
}

func NewSkeleton(line consts.LandKey, X int) *Enemy {
	X = X - 64
	return NewEnemy(line, X, 128, 128, 1000, 1, 100, EnemyImageInfo{
		RunImagePath:        "/Users/admin/Desktop/ZWDZJS/image/enemy/skeleton_run.png",
		RunImageFrame:       6,
		RunImageFrameSpeed:  10,
		ShotImagePath:       "/Users/admin/Desktop/ZWDZJS/image/enemy/skeleton_shot.png",
		ShotImageFrame:      4,
		ShotImageFrameSpeed: 8,
	}, func(placement *Enemy) {

	})
}

func NewRandomEnemy(line consts.LandKey, X int) *Enemy {
	rand.Seed(time.Now().UnixNano())
	//randomInt := rand.Intn(1)
	return NewSkeleton(line, X)
}

func NewEnemies() map[consts.LandKey]map[string]*Enemy {
	enemies := make(map[consts.LandKey]map[string]*Enemy, 0)
	for i := 1; i <= 5; i++ {
		enemies[consts.LandKey(i)] = make(map[string]*Enemy, 0)
	}
	return enemies
}

var Enemies = NewEnemies()
