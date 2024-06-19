package things

import (
	"ZWDZJS/caches"
	"ZWDZJS/consts"
	"ZWDZJS/inter"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"math/rand"
	"time"
)

type Placement struct {
	layers                              int
	gameIdentity                        consts.GameIdentityKey
	gameThingsKey                       int64
	uuid                                string
	X, Y, Width, Height                 int
	isShow                              bool
	totalFrame                          int64
	canThroughput, canShot              bool
	pictureFrame                        int
	pictureFrameCount                   int
	isIdle                              bool
	isShot                              bool
	image                               *ebiten.Image
	imageFrame, imageFrameSpeed         int
	idleImage                           *ebiten.Image
	idleImageFrame, idleImageFrameSpeed int
	shotImage                           *ebiten.Image
	shotImageFrame, shotImageFrameSpeed int
	Health                              int
	isShowBorder                        bool
	ThroughputFunc                      func(*Placement)
	ShotFunc                            func(*Placement)
	landBlock                           *LandBlock
	shotScope                           []*LandBlock
}

func (m *Placement) DamageCalculation(damage int, face int) (die bool, ex int) {
	if m.Health-damage <= 0 {
		m.Health = 0
	} else {
		m.Health = m.Health - damage
	}
	return false, 0
}

func (m *Placement) GetImage() *ebiten.Image {
	subRect := image.Rect(m.pictureFrame*m.Width, 0, m.pictureFrame*m.Width+m.Width, m.Height)
	subImg := m.image.SubImage(subRect).(*ebiten.Image)
	return subImg
}

func (m *Placement) Move(nextX int, nextY int) {
	if m.X != nextX || m.Y != nextY {
		m.X, m.Y = nextX, nextY
		QM.MoveGameObject(m)
	}
}

func (m *Placement) GetPictureFrame() int {
	return m.pictureFrame
}

func (m *Placement) GetXY() (int, int) {
	return m.X, m.Y
}

func (m *Placement) SetXY(x int, y int) {
	m.X = x
	m.Y = y
}

func (m *Placement) IsShot() (isShot bool) {
	for _, info := range m.shotScope {
		_, num := info.GetEnemies()
		if num > 0 {
			isShot = true
			return
		}
	}
	return
}

func (m *Placement) Delete() {
	m.isShow = false
	QM.RemoveGameObject(m)
	inter.GameThings.Delete(m.gameThingsKey)
}

func (m *Placement) GetWidth() int {
	return m.Width
}

func (m *Placement) GetHeight() int {
	return m.Height
}

func (m *Placement) GetUuid() string {
	return m.uuid
}

func (m *Placement) GetKey() int64 {
	return m.gameThingsKey
}

func (m *Placement) SetKey(key int64) {
	m.gameThingsKey = key
}

func (m *Placement) GetGameIdentity() consts.GameIdentityKey {
	return m.gameIdentity
}

func (m *Placement) Update() error {
	if m.isShow {
		isShot := m.IsShot()
		if m.totalFrame >= 6000000 {
			m.totalFrame = 0
		}
		m.totalFrame++

		if m.Health <= 0 {
			m.isShow = false
			m.landBlock.P = nil
			QM.RemoveGameObject(m)
			inter.GameThings.Delete(m.gameThingsKey)
			return nil
		}

		if m.canThroughput {
			m.ThroughputFunc(m)
		}

		if m.isShot != isShot {
			m.pictureFrame = 0
			m.pictureFrameCount = 0
		}

		m.isShot = isShot
		m.isIdle = !m.isShot

		if m.isIdle {
			m.imageFrame = m.idleImageFrame
			m.imageFrameSpeed = m.idleImageFrameSpeed
			m.image = m.idleImage
		} else {
			m.imageFrame = m.shotImageFrame
			m.imageFrameSpeed = m.shotImageFrameSpeed
			m.image = m.shotImage
			if m.pictureFrame == 0 && m.pictureFrameCount%m.imageFrameSpeed == m.imageFrameSpeed-1 {
				//NewBullet(m.X+m.Width, m.Y+(m.Height/2)-1+15, 1, m.lv*10, m.uuid, m)
				if m.canShot {
					m.ShotFunc(m)
				}
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

func (m *Placement) Draw(screen *ebiten.Image) {
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

func (m *Placement) drawRectOutline(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x+width), float64(y), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y+height), float64(x+width), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x+width), float64(y), float64(x+width), float64(y+height), clr)
}

func (m *Placement) GetLayers() int {
	return m.layers
}

func (m *Placement) bindIdleImage(info MainImageInfo) error {
	if info.IdleImagePath == "" {
		return nil
	}
	if imageCache, ok := caches.ImageCache[info.IdleImagePath]; ok {
		m.idleImage = imageCache.RightImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(info.IdleImagePath)
		caches.ImageCache[info.IdleImagePath] = caches.ImageCacheInfo{
			RightImage: image,
		}
		m.idleImage = image
	}
	m.idleImageFrame = info.IdleImageFrame
	m.idleImageFrameSpeed = info.IdleImageFrameSpeed

	m.image = m.idleImage
	m.imageFrame = info.IdleImageFrame
	m.imageFrameSpeed = info.IdleImageFrameSpeed

	return nil
}

func (m *Placement) bindShotImage(info MainImageInfo) error {
	if info.ShotImagePath == "" {
		return nil
	}
	if imageCache, ok := caches.ImageCache[info.ShotImagePath]; ok {
		m.shotImage = imageCache.RightImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(info.ShotImagePath)
		caches.ImageCache[info.ShotImagePath] = caches.ImageCacheInfo{
			RightImage: image,
		}
		m.shotImage = image
	}
	m.shotImageFrame = info.ShotImageFrame
	m.shotImageFrameSpeed = info.ShotImageFrameSpeed
	return nil
}

type MainImageInfo struct {
	IdleImagePath       string
	IdleImageFrame      int
	IdleImageFrameSpeed int
	ShotImagePath       string
	ShotImageFrame      int
	ShotImageFrameSpeed int
}

func NewPlacement(X, Y, Width, Height, Health int, landBlock *LandBlock, image MainImageInfo, canThroughput bool, throughputFunc func(placement *Placement), canShot bool, shotFunc func(placement *Placement), shotScope []*LandBlock) *Placement {
	m := &Placement{
		X:              X,
		Y:              Y,
		layers:         consts.LayersPlacementKey,
		Width:          Width,
		Height:         Height,
		uuid:           uuid.New().String(),
		isIdle:         true,
		Health:         Health,
		gameIdentity:   consts.GameIdentityPlacement,
		isShow:         true,
		canThroughput:  canThroughput,
		ThroughputFunc: throughputFunc,
		canShot:        canShot,
		ShotFunc:       shotFunc,
		landBlock:      landBlock,
		isShowBorder:   true,
		shotScope:      shotScope,
	}
	m.bindIdleImage(image)
	m.bindShotImage(image)
	inter.GameThings.Put(m)
	QM.AddGameObject(m, true)
	return m
}

func NewSoldier(X, Y int, landBlock *LandBlock) *Placement {
	X = X + consts.PlantLandWidth/2 - 69
	Y = Y - 128 + consts.PlantLandHeight - 10
	shotScope := make([]*LandBlock, 0)
	line, row := landBlock.GetLineRow()
	for row <= consts.LandRowNineKey {
		shotScope = append(shotScope, LandBlocker[line][row])
		row++
	}
	return NewPlacement(X, Y, 128, 128, 1000, landBlock, MainImageInfo{
		IdleImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/soldier_idle.png",
		IdleImageFrame:      7,
		IdleImageFrameSpeed: 7,
		ShotImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/soldier_shot.png",
		ShotImageFrame:      4,
		ShotImageFrameSpeed: 2,
	}, false, func(placement *Placement) {

	}, true, func(placement *Placement) {
		for _, info := range placement.shotScope {
			enemies, num := info.GetEnemies()
			if num > 0 {
				for _, enemyInfo := range enemies {
					enemyInfo.DamageCalculation(100, 1)
				}
			}
		}
	}, shotScope)
}

func NewWarrior(X, Y int, landBlock *LandBlock) *Placement {
	X = X + consts.PlantLandWidth/2 - 69
	Y = Y - 128 + consts.PlantLandHeight - 10
	shotScope := make([]*LandBlock, 0)
	line, row := landBlock.GetLineRow()
	shotScope = append(shotScope, LandBlocker[line][row])

	return NewPlacement(X, Y, 128, 128, 4000, landBlock, MainImageInfo{
		IdleImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/warrior_idle.png",
		IdleImageFrame:      9,
		IdleImageFrameSpeed: 9,
		ShotImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/warrior_shot.png",
		ShotImageFrame:      5,
		ShotImageFrameSpeed: 7,
	}, false, func(placement *Placement) {

	}, true, func(placement *Placement) {
		for _, info := range placement.shotScope {
			enemies, num := info.GetEnemies()
			if num > 0 {
				for _, enemyInfo := range enemies {
					enemyInfo.DamageCalculation(100, 1)
				}
			}
		}
	}, shotScope)
}

func NewRedHair(X, Y int, landBlock *LandBlock) *Placement {
	X = X + consts.PlantLandWidth/2 - 69
	Y = Y - 128 + consts.PlantLandHeight - 10
	shotScope := make([]*LandBlock, 0)
	line, row := landBlock.GetLineRow()
	shotScope = append(shotScope, LandBlocker[line][row])
	return NewPlacement(X, Y, 128, 128, 4000, landBlock, MainImageInfo{
		IdleImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/red_hair_idle.png",
		IdleImageFrame:      5,
		IdleImageFrameSpeed: 5,
		ShotImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/red_hair_shot.png",
		ShotImageFrame:      5,
		ShotImageFrameSpeed: 7,
	}, false, func(placement *Placement) {

	}, true, func(placement *Placement) {
		for _, info := range placement.shotScope {
			enemies, num := info.GetEnemies()
			if num > 0 {
				for _, enemyInfo := range enemies {
					enemyInfo.DamageCalculation(100, 1)
				}
			}
		}
	}, shotScope)
}

func NewRedHat(X, Y int, landBlock *LandBlock) *Placement {
	X = X + consts.PlantLandWidth/2 - 69
	Y = Y - 128 + consts.PlantLandHeight - 10
	shotScope := make([]*LandBlock, 0)
	line, row := landBlock.GetLineRow()
	shotScope = append(shotScope, LandBlocker[line][row])
	return NewPlacement(X, Y, 128, 128, 4000, landBlock, MainImageInfo{
		IdleImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/red_hat_idle.png",
		IdleImageFrame:      5,
		IdleImageFrameSpeed: 5,
		ShotImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/red_hat_shot.png",
		ShotImageFrame:      6,
		ShotImageFrameSpeed: 8,
	}, false, func(placement *Placement) {

	}, true, func(placement *Placement) {
		for _, info := range placement.shotScope {
			enemies, num := info.GetEnemies()
			if num > 0 {
				for _, enemyInfo := range enemies {
					enemyInfo.DamageCalculation(100, 1)
				}
			}
		}
	}, shotScope)
}

func NewEnergyFlower(X, Y int, landBlock *LandBlock) *Placement {
	X = X + consts.PlantLandWidth/2 - 69
	Y = Y - 128 + consts.PlantLandHeight - 10
	shotScope := make([]*LandBlock, 0)
	return NewPlacement(X, Y, 128, 128, 4000, landBlock, MainImageInfo{
		IdleImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/energy_flower_idle.png",
		IdleImageFrame:      11,
		IdleImageFrameSpeed: 11,
		ShotImagePath:       "/Users/admin/Desktop/ZWDZJS/image/placement/energy_flower_shot.png",
		ShotImageFrame:      7,
		ShotImageFrameSpeed: 7,
	}, true, func(placement *Placement) {
		if placement.totalFrame%600 == 599 {
			NewEnergy(float64(placement.X), float64(placement.Y), float64(placement.X)-40, float64(placement.Y)+80, 100)
		}
	}, false, func(placement *Placement) {

	}, shotScope)
}

func NewRandomPlacement(X, Y int, landBlock *LandBlock) *Placement {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(5)
	switch {
	case randomInt == 0:
		return NewWarrior(X, Y, landBlock)
	case randomInt == 1:
		return NewRedHair(X, Y, landBlock)
	case randomInt == 2:
		return NewRedHat(X, Y, landBlock)
	case randomInt == 3:
		return NewEnergyFlower(X, Y, landBlock)
	default:
		return NewSoldier(X, Y, landBlock)
	}
}
