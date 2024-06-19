package things

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/color"
	"star_dust/consts"
	"star_dust/inter"
	"star_dust/utils"
)

type MainCharacter struct {
	layers              int
	gameIdentity        consts.GameIdentityKey
	gameThingsKey       int64
	uuid                string
	X, Y, Width, Height int
	isShow              bool
	pictureFrame        int
	pictureFrameCount   int
	isRun               bool
	isIdle              bool
	isShot              bool
	isShoting           bool
	faceStatus          bool
	moveButtons         map[consts.MoveDirectionKey]ebiten.Key
	image               *ebiten.Image
	imageFrame          int
	imageFrameSpeed     int
	leftRunImage        *ebiten.Image
	rightRunImage       *ebiten.Image
	runImageFrame       int
	leftIdleImage       *ebiten.Image
	rightIdleImage      *ebiten.Image
	idleImageFrame      int
	leftShotImage       *ebiten.Image
	rightShotImage      *ebiten.Image
	shotImageFrame      int
	Health              int
	MaxHealth           int
	ex, lv              int
}

func (m *MainCharacter) AbsorbingExperience(ex int) {
	m.ex = m.ex + ex
	nextLv := m.ex/100 + 1
	if m.lv != nextLv {
		m.lv = nextLv
	}
}

func (m *MainCharacter) DamageCalculation(damage int, face int) (die bool, ex int) {
	if m.Health-damage <= 0 {
		m.Health = 0
	} else {
		m.Health = m.Health - damage
	}
	return false, 0
}

func (m *MainCharacter) GetImage() *ebiten.Image {
	subRect := image.Rect(m.pictureFrame*m.Width, 0, m.pictureFrame*m.Width+m.Width, m.Height)
	subImg := m.image.SubImage(subRect).(*ebiten.Image)
	return subImg
}

func (m *MainCharacter) Move(nextX int, nextY int) {
	if m.X != nextX || m.Y != nextY {
		m.X, m.Y = nextX, nextY
		QM.MoveGameObject(m)
	}
}

func (m *MainCharacter) GetPictureFrame() int {
	return m.pictureFrame
}

func (m *MainCharacter) GetXY() (int, int) {
	return m.X, m.Y
}

func (m *MainCharacter) SetXY(x int, y int) {
	m.X = x
	m.Y = y
}

func (m *MainCharacter) GetWidth() int {
	return m.Width
}

func (m *MainCharacter) GetHeight() int {
	return m.Height
}

func (m *MainCharacter) GetUuid() string {
	return m.uuid
}

func (m *MainCharacter) GetKey() int64 {
	return m.gameThingsKey
}

func (m *MainCharacter) SetKey(key int64) {
	m.gameThingsKey = key
}

func (m *MainCharacter) GetGameIdentity() consts.GameIdentityKey {
	return m.gameIdentity
}

func (m *MainCharacter) Update() error {
	if m.isShow {
		if m.Health <= 0 {
			m.isShow = false
			QM.RemoveGameObject(m)
			inter.GameThings.Delete(m.gameThingsKey)
			return nil
		}
		isSwitch := false
		isIdle := true
		isRun := false
		isShot := false
		preX, preY := m.X, m.Y

		if ebiten.IsKeyPressed(m.moveButtons[consts.MoveDirectionLeft]) {
			isRun = true
			isIdle = false
			m.isShoting = false
			m.faceStatus = false
			m.X = m.X + 5*-1
		}
		if ebiten.IsKeyPressed(m.moveButtons[consts.MoveDirectionRight]) {
			isRun = true
			isIdle = false
			m.faceStatus = true
			m.isShoting = false
			m.X = m.X + 5*1
		}
		if ebiten.IsKeyPressed(m.moveButtons[consts.MoveDirectionUp]) {
			isRun = true
			isIdle = false
			m.isShoting = false
			m.Y = m.Y + 3*-1
		}
		if ebiten.IsKeyPressed(m.moveButtons[consts.MoveDirectionDown]) {
			isRun = true
			isIdle = false
			m.isShoting = false
			m.Y = m.Y + 3*1
		}
		if (ebiten.IsKeyPressed(m.moveButtons[consts.MoveDirectionShot]) || m.isShoting) && !isRun {
			isShot = true
			isIdle = false
		}

		if isRun && m.isRun != isRun {
			isSwitch = true
		} else if isShot && m.isShot != isShot {
			isSwitch = true
		} else if isIdle && m.isIdle != isIdle {
			isSwitch = true
		}
		m.isIdle = isIdle
		m.isShot = isShot
		m.isRun = isRun

		if isSwitch {
			m.pictureFrame = 0
			m.pictureFrameCount = 0
		}

		if m.isRun {
			m.imageFrame = m.runImageFrame
			m.imageFrameSpeed = m.runImageFrame
			if m.faceStatus {
				m.image = m.rightRunImage
			} else {
				m.image = m.leftRunImage
			}
		} else if m.isIdle {
			m.imageFrame = m.idleImageFrame
			m.imageFrameSpeed = m.idleImageFrame
			if m.faceStatus {
				m.image = m.rightIdleImage
			} else {
				m.image = m.leftIdleImage
			}
		} else if m.isShot {
			m.isShoting = true
			m.imageFrame = m.shotImageFrame
			m.imageFrameSpeed = m.shotImageFrame / 4
			if m.faceStatus {
				m.image = m.rightShotImage
			} else {
				m.image = m.leftShotImage
			}
			if m.pictureFrame == 0 && m.pictureFrameCount%m.imageFrameSpeed == m.imageFrameSpeed-1 {
				if m.faceStatus {
					NewBullet(m.X+m.Width, m.Y+(m.Height/2)-1+15, 1, m.lv*10, m.uuid, m)
				} else {
					NewBullet(m.X-7, m.Y+(m.Height/2)-1+15, -1, m.lv*10, m.uuid, m)
				}
			}
			if m.pictureFrame == m.imageFrame-1 {
				m.isShoting = false
			}
		}

		if m.X != preX || m.Y != preY {
			objs := QM.SameAreaGameObjects(m)
			isCollision := false
			for _, otherObj := range objs {
				if otherObj.GetGameIdentity() == consts.GameIdentityEntrance && QM.PixelCollision(m, otherObj) {
					isCollision = true
					break
				}
			}
			if isCollision {
				canOut := false
				checkCanOut := 0
				for !canOut && checkCanOut <= 10 {
					m.X = m.X - (m.X-preX)*3
					m.Y = m.Y - (m.Y-preY)*3
					canOut = true
					for _, otherObj := range objs {
						if QM.PixelCollision(m, otherObj) {
							canOut = false
							break
						}
					}
					checkCanOut++
					if checkCanOut > 10 && !canOut {
						m.X, m.Y = preX, preY
					}
				}

			} else {
				if (m.faceStatus && m.X+m.Width >= inter.ScreenWidth-100) || (!m.faceStatus && m.X <= 100) {
					if SM.Offset(m.X-preX, 0, m) {
						QM.MoveAreaGameObjects(-1*(m.X-preX), 0)
					} else {
						if m.X < 0 || m.X+m.Width > inter.ScreenWidth {
							m.X = preX
						}
					}
				}
			}
			if m.X != preX || m.Y != preY {
				QM.MoveGameObject(m)
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

func (m *MainCharacter) Show() {
	m.isShow = true
}

func (m *MainCharacter) Hide() {
	m.isShow = false
}

func (m *MainCharacter) Draw(screen *ebiten.Image) {
	if m.isShow {
		subRect := image.Rect(m.pictureFrame*m.Width, 0, m.pictureFrame*m.Width+m.Width, m.Height)
		subImg := m.image.SubImage(subRect).(*ebiten.Image)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(m.X), float64(m.Y))
		screen.DrawImage(subImg, opts)
		m.drawHealthBar(screen)
	}
}

func (m *MainCharacter) drawHealthBar(screen *ebiten.Image) {
	barWidth := 100
	barHeight := 10
	offsetX := m.Width/2 - barWidth/2
	barX := m.X + offsetX
	barY := m.Y + 20

	// 绘制血条背景
	ebitenutil.DrawRect(screen, float64(barX), float64(barY), float64(barWidth), float64(barHeight), color.RGBA{100, 100, 100, 255})

	// 绘制当前血量条
	healthWidth := barWidth * m.Health / m.MaxHealth
	ebitenutil.DrawRect(screen, float64(barX), float64(barY), float64(healthWidth), float64(barHeight), color.RGBA{255, 0, 0, 255})

	// 绘制分隔线
	segments := 10
	segmentWidth := barWidth / segments
	for i := 1; i < segments; i++ {
		lineX := barX + segmentWidth*i
		ebitenutil.DrawLine(screen, float64(lineX), float64(barY), float64(lineX), float64(barY+barHeight), color.RGBA{0, 0, 0, 255})
	}

	// 绘制血量数值
	healthText := fmt.Sprintf("%d/%d", m.Health, m.MaxHealth)
	font := basicfont.Face7x13
	text.Draw(screen, healthText, font, barX+barWidth+5, barY+barHeight, color.White)

	lvText := fmt.Sprintf("LV %d", m.lv)
	text.Draw(screen, lvText, font, barX-50, barY+barHeight, color.White)
}

func (m *MainCharacter) GetLayers() int {
	return m.layers
}

func (m *MainCharacter) bindRunImage(path string, imageFrame int) error {
	image, _, _ := ebitenutil.NewImageFromFile(path)
	m.rightRunImage = image
	m.leftRunImage = utils.FlipImageSegments(image, 128)
	m.runImageFrame = imageFrame
	return nil
}

func (m *MainCharacter) bindIdleImage(path string, imageFrame int) error {
	image, _, _ := ebitenutil.NewImageFromFile(path)
	m.rightIdleImage = image
	m.leftIdleImage = utils.FlipImageSegments(image, 128)
	m.idleImageFrame = imageFrame
	m.imageFrame = m.idleImageFrame
	m.imageFrameSpeed = m.idleImageFrame
	m.image = image
	return nil
}

func (m *MainCharacter) bindShotImage(path string, imageFrame int) error {
	image, _, _ := ebitenutil.NewImageFromFile(path)
	m.rightShotImage = image
	m.leftShotImage = utils.FlipImageSegments(image, 128)
	m.shotImageFrame = imageFrame
	return nil
}

type MainImageInfo struct {
	RunImagePath   string
	RunImageFrame  int
	IdleImagePath  string
	IdleImageFrame int
	ShotImagePath  string
	ShotImageFrame int
}

func NewMainCharacter(x, y int, moveButtons map[consts.MoveDirectionKey]ebiten.Key, image MainImageInfo) *MainCharacter {
	m := &MainCharacter{
		X:            x,
		Y:            y,
		layers:       consts.LayersMainKey,
		Width:        128,
		Height:       128,
		uuid:         uuid.New().String(),
		moveButtons:  moveButtons,
		isIdle:       true,
		faceStatus:   true,
		Health:       100,
		MaxHealth:    100,
		gameIdentity: consts.GameIdentityUser,
		lv:           1,
	}
	m.bindRunImage(image.RunImagePath, image.RunImageFrame)
	m.bindIdleImage(image.IdleImagePath, image.IdleImageFrame)
	m.bindShotImage(image.ShotImagePath, image.ShotImageFrame)
	inter.GameThings.Put(m)
	QM.AddGameObject(m, true)
	return m
}

func NewMainCharacters() *MainCharacter {
	u1 := NewMainCharacter(0, 128, map[consts.MoveDirectionKey]ebiten.Key{
		consts.MoveDirectionLeft:  ebiten.KeyLeft,
		consts.MoveDirectionRight: ebiten.KeyRight,
		consts.MoveDirectionUp:    ebiten.KeyUp,
		consts.MoveDirectionDown:  ebiten.KeyDown,
		consts.MoveDirectionShot:  ebiten.KeySpace,
	}, MainImageInfo{
		RunImagePath:   "/Users/admin/Desktop/star_dust/image/human_run_4.png",
		RunImageFrame:  8,
		IdleImagePath:  "/Users/admin/Desktop/star_dust/image/human_idle_4.png",
		IdleImageFrame: 7,
		ShotImagePath:  "/Users/admin/Desktop/star_dust/image/human_shot_4.png",
		ShotImageFrame: 4,
	})
	return u1
}

var MCS = NewMainCharacters()

func RemakeMainCharacters() {
	QM.RemoveGameObject(MCS)
	inter.GameThings.Delete(MCS.gameThingsKey)
	MCS = NewMainCharacters()
}

func ShowMCS() {
	MCS.Show()
}

func HideMCS() {
	MCS.Hide()
}
