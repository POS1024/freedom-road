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
	"sync"
	"time"
)

type EnemyCharacter struct {
	layers                int
	gameIdentity          consts.GameIdentityKey
	gameThingsKey         int64
	uuid                  string
	X, Y, Width, Height   int
	faceX, faceY          int
	nowFaceX              int
	isShow                bool
	showHealthBarTimeLeft int
	isRun                 bool
	isBeAttacked          bool
	isBeAttacking         bool
	isAttacking           bool
	pictureFrame          int
	pictureFrameCount     int
	faceStatus            bool
	image                 *ebiten.Image
	imageFrame            int
	imageFrameSpeed       int
	leftRunImage          *ebiten.Image
	rightRunImage         *ebiten.Image
	runImageFrame         int
	leftHurtImage         *ebiten.Image
	rightHurtImage        *ebiten.Image
	hurtImageFrame        int
	leftAttackImage       *ebiten.Image
	rightAttackImage      *ebiten.Image
	attackImageFrame      int
	Health                int
	MaxHealth             int
	moveSpeed             int
}

func (m *EnemyCharacter) DamageCalculation(damage int, face int) (die bool, ex int) {
	m.isBeAttacked = true
	m.Move(m.X+face*3, m.Y)
	if m.Health-damage <= 0 {
		NewDamageFont(m.X+m.Width/2-3, m.Y, m.Health)
		m.Health = 0
		return true, 30
	} else {
		m.Health = m.Health - damage
		NewDamageFont(m.X+m.Width/2-3, m.Y, damage)
		m.showHealthBarTimeLeft = 600
		return false, 0
	}
}

func (m *EnemyCharacter) Move(nextX int, nextY int) {
	if m.X != nextX || m.Y != nextY {
		m.X, m.Y = nextX, nextY
		QM.MoveGameObject(m)
	}
}

func (m *EnemyCharacter) GetImage() *ebiten.Image {
	subRect := image.Rect(m.pictureFrame*m.Width, 0, m.pictureFrame*m.Width+m.Width, m.Height)
	subImg := m.image.SubImage(subRect).(*ebiten.Image)
	return subImg
}

func (m *EnemyCharacter) GetPictureFrame() int {
	return m.pictureFrame
}

func (m *EnemyCharacter) GetXY() (int, int) {
	return m.X, m.Y
}

func (m *EnemyCharacter) SetXY(x int, y int) {
	m.X = x
	m.Y = y
}

func (m *EnemyCharacter) GetWidth() int {
	return m.Width
}

func (m *EnemyCharacter) GetHeight() int {
	return m.Height
}

func (m *EnemyCharacter) GetUuid() string {
	return m.uuid
}

func (m *EnemyCharacter) GetKey() int64 {
	return m.gameThingsKey
}

func (m *EnemyCharacter) SetKey(key int64) {
	m.gameThingsKey = key
}

func (m *EnemyCharacter) GetGameIdentity() consts.GameIdentityKey {
	return m.gameIdentity
}

func (m *EnemyCharacter) Update() error {
	if m.isShow {
		if m.showHealthBarTimeLeft > 0 {
			m.showHealthBarTimeLeft--
		}
		if m.Health <= 0 {
			QM.RemoveGameObject(m)
			RemoveEnemyCharacter(m.uuid)
			inter.GameThings.Delete(m.gameThingsKey)
			if ECSMaxNumber > 1 {
				ECSMaxNumber--
			}
			return nil
		}
		if m.isBeAttacked {
			m.isBeAttacked = false
			m.isRun = false
			m.isAttacking = false
			m.pictureFrame = 0
			m.pictureFrameCount = 0
			m.isBeAttacking = true
			if m.nowFaceX == -1 {
				m.image = m.leftHurtImage
			} else {
				m.image = m.rightHurtImage
			}
			m.imageFrame = m.hurtImageFrame
			m.imageFrameSpeed = m.hurtImageFrame + 3
		}

		if m.isRun {
			nextX, nextY := m.X, m.Y
			if m.X > inter.ScreenWidth-100 || m.X < 100 || m.Y < 150 || m.Y > inter.ScreenHeight-140 {
				if m.X > inter.ScreenWidth-100 {
					m.faceX = -1
				}
				if m.X < 100 {
					m.faceX = 1
				}
				if m.Y < 150 {
					m.faceY = 1
				}
				if m.Y > inter.ScreenHeight-140 {
					m.faceY = -1
				}
			} else {
				userX, userY := MCS.GetXY()
				userW, userH := MCS.GetWidth(), MCS.GetHeight()
				if (m.X+m.Width/2)-(userX+userW/2) > 0 {
					m.faceX = -1
				} else {
					m.faceX = 1
				}
				if (m.Y+m.Height/2)-(userY+userH/2) > 0 {
					m.faceY = -1
				} else {
					m.faceY = 1
				}
			}

			nextX, nextY = m.X+m.moveSpeed*m.faceX, m.Y+m.moveSpeed*m.faceY
			if m.faceX != 0 {
				m.nowFaceX = m.faceX
			}
			m.Move(nextX, nextY)
			objs := QM.SameAreaGameObjects(m)

			if m.nowFaceX == -1 {
				m.image = m.leftRunImage
			} else {
				m.image = m.rightRunImage
			}

			for _, otherObj := range objs {
				if otherObj.GetGameIdentity() == consts.GameIdentityUser && QM.PixelCollision(m, otherObj) {
					m.isRun = false
					m.isAttacking = true
					if m.nowFaceX == -1 {
						m.image = m.leftAttackImage
					} else {
						m.image = m.rightAttackImage
					}
					m.imageFrame = m.attackImageFrame
					m.imageFrameSpeed = m.attackImageFrame + 8
					m.pictureFrameCount = 0
					m.pictureFrame = 0
					break
				}
			}
		}

		if m.pictureFrameCount%m.imageFrameSpeed == m.imageFrameSpeed-1 {

			if m.pictureFrame == m.imageFrame-1 {
				if m.isBeAttacking {
					m.isBeAttacking = false
					m.isRun = true
					m.imageFrame = m.runImageFrame
					m.imageFrameSpeed = m.runImageFrame + 3
					m.pictureFrameCount = 0
				}
				if m.isAttacking {
					objs := QM.SameAreaGameObjects(m)
					for _, otherObj := range objs {
						if otherObj.GetGameIdentity() == consts.GameIdentityUser && QM.PixelCollision(m, otherObj) {
							middObj := otherObj
							go func() {
								middObj.(inter.DamageableItems).DamageCalculation(10, m.nowFaceX)
							}()
						}
					}
					m.isAttacking = false
					m.isRun = true
					m.imageFrame = m.runImageFrame
					m.imageFrameSpeed = m.runImageFrame + 3
					m.pictureFrameCount = 0
				}
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

func (m *EnemyCharacter) Show() {
	m.isShow = true
}

func (m *EnemyCharacter) Hide() {
	m.isShow = false
}

func (m *EnemyCharacter) Draw(screen *ebiten.Image) {
	if m.isShow {
		subRect := image.Rect(m.pictureFrame*m.Width, 0, m.pictureFrame*m.Width+m.Width, m.Height)
		subImg := m.image.SubImage(subRect).(*ebiten.Image)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(m.X), float64(m.Y))
		screen.DrawImage(subImg, opts)
		if m.showHealthBarTimeLeft > 0 {
			m.drawHealthBar(screen)
		}
	}
}

func (m *EnemyCharacter) drawHealthBar(screen *ebiten.Image) {
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
}

func (m *EnemyCharacter) GetLayers() int {
	return m.layers
}

func (m *EnemyCharacter) bindRunImage(path string, imageFrame int) error {
	if imageCache, ok := ImageCache[path]; ok {
		m.image = imageCache.leftImage
		m.leftRunImage = imageCache.leftImage
		m.rightRunImage = imageCache.rightImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(path)
		m.image = image
		m.leftRunImage = m.image
		m.rightRunImage = utils.FlipImageSegments(image, 72)
		ImageCache[path] = ImageCacheInfo{
			leftImage:  m.leftRunImage,
			rightImage: m.rightRunImage,
		}
	}
	m.runImageFrame = imageFrame
	m.imageFrame = imageFrame
	m.imageFrameSpeed = imageFrame + 3
	return nil
}

func (m *EnemyCharacter) bindHurtImage(path string, imageFrame int) error {
	if imageCache, ok := ImageCache[path]; ok {
		m.leftHurtImage = imageCache.leftImage
		m.rightHurtImage = imageCache.rightImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(path)
		m.leftHurtImage = image
		m.rightHurtImage = utils.FlipImageSegments(image, 72)
		ImageCache[path] = ImageCacheInfo{
			leftImage:  m.leftHurtImage,
			rightImage: m.rightHurtImage,
		}
	}
	m.hurtImageFrame = imageFrame
	return nil
}

func (m *EnemyCharacter) bindAttackImage(path string, imageFrame int) error {
	if imageCache, ok := ImageCache[path]; ok {
		m.leftAttackImage = imageCache.leftImage
		m.rightAttackImage = imageCache.rightImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(path)
		m.leftAttackImage = image
		m.rightAttackImage = utils.FlipImageSegments(image, 72)
		ImageCache[path] = ImageCacheInfo{
			leftImage:  m.leftAttackImage,
			rightImage: m.rightAttackImage,
		}
	}
	m.attackImageFrame = imageFrame
	return nil
}

type EnemyImageInfo struct {
	RunImagePath     string
	RunImageFrame    int
	HurtImagePath    string
	HurtImageFrame   int
	AttackImagePath  string
	AttackImageFrame int
}

func NewEnemyCharacter(x, y int, image EnemyImageInfo) *EnemyCharacter {
	m := &EnemyCharacter{
		X:            x,
		Y:            y,
		layers:       consts.LayersMainKey,
		Width:        72,
		Height:       72,
		uuid:         uuid.New().String(),
		faceStatus:   false,
		Health:       1000,
		MaxHealth:    1000,
		gameIdentity: consts.GameIdentityEnemy,
		isShow:       true,
		isRun:        true,
		faceX:        -1,
		faceY:        0,
		nowFaceX:     -1,
		moveSpeed:    1,
	}
	m.bindRunImage(image.RunImagePath, image.RunImageFrame)
	m.bindHurtImage(image.HurtImagePath, image.HurtImageFrame)
	m.bindAttackImage(image.AttackImagePath, image.AttackImageFrame)
	inter.GameThings.GoPut(m)
	QM.AddGameObject(m, true)
	return m
}

func NewEnemyCharacters() (characters map[string]*EnemyCharacter) {
	characters = make(map[string]*EnemyCharacter, 0)
	go func() {
		for {
			<-time.After(time.Second * 20)
			if ECSStatus && ECSMaxNumber < 30 {
				ECSMaxNumber++
			}
		}
	}()
	go func() {
		number := 0
		for {
			if number >= 4 {
				number = 0
			}
			if ECSStatus && len(characters) < ECSMaxNumber {
				ECSLock.Lock()
				u1 := NewEnemyCharacter(inter.ScreenWidth-100, 150+number*60, EnemyImageInfo{
					RunImagePath:     "/Users/admin/Desktop/star_dust/image/enemy/enemy_run_1.png",
					RunImageFrame:    4,
					HurtImagePath:    "/Users/admin/Desktop/star_dust/image/enemy/enemy_hurt_1.png",
					HurtImageFrame:   2,
					AttackImagePath:  "/Users/admin/Desktop/star_dust/image/enemy/enemy_attack_1.png",
					AttackImageFrame: 4,
				})
				characters[u1.GetUuid()] = u1
				number++
				ECSLock.Unlock()
			}
			<-time.After(time.Second)

		}
	}()
	return
}

var ECS = NewEnemyCharacters()

var ECSStatus = false

var ECSLock sync.Mutex

var ECSMaxNumber = 4

func RemakeEnemyCharacters() {
	ECSLock.Lock()
	defer ECSLock.Unlock()
	ECSStatus = false
	for _, info := range ECS {
		QM.RemoveGameObject(info)
		inter.GameThings.Delete(info.gameThingsKey)
		delete(ECS, info.GetUuid())
	}
	ECSMaxNumber = 4
	ECSStatus = true
}

func RemoveEnemyCharacter(mainUuid string) {
	ECSLock.Lock()
	defer ECSLock.Unlock()
	delete(ECS, mainUuid)
}

func ShowECS() {
	ECSLock.Lock()
	defer ECSLock.Unlock()
	ECSStatus = true
	for _, info := range ECS {
		info.Show()
	}
	return
}

func HideECS() {
	ECSLock.Lock()
	defer ECSLock.Unlock()
	ECSStatus = false
	for _, info := range ECS {
		info.Hide()
	}
	return
}
