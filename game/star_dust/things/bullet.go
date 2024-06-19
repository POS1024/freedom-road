package things

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"star_dust/consts"
	"star_dust/inter"
)

type Bullet struct {
	layers                         int
	gameIdentity                   consts.GameIdentityKey
	gameThingsKey                  int64
	mainId                         string
	MainRole                       inter.Role
	uuid                           string
	X, Y, Width, Height, DrawWidth int
	isShow                         bool
	pictureFrame                   int
	pictureFrameCount              int
	image                          *ebiten.Image
	imageFrame                     int
	imageFrameSpeed                int
	faceStatus                     int
	damage                         int
}

func (b *Bullet) GetDamage() int {
	return b.damage
}

func (b *Bullet) Move(nextX int, nextY int) {
	if b.X != nextX || b.Y != nextY {
		b.X, b.Y = nextX, nextY
		QM.MoveGameObject(b)
	}
}

func (b *Bullet) GetImage() *ebiten.Image {
	subRect := image.Rect(b.pictureFrame*b.Width, 0, b.pictureFrame*b.Width+b.Width, b.Height)
	subImg := b.image.SubImage(subRect).(*ebiten.Image)
	return subImg
}

func (b *Bullet) GetPictureFrame() int {
	return b.pictureFrame
}

func (b *Bullet) GetXY() (int, int) {
	return b.X, b.Y
}

func (b *Bullet) SetXY(x int, y int) {
	b.X = x
	b.Y = y
}

func (b *Bullet) GetWidth() int {
	return b.Width
}

func (b *Bullet) GetHeight() int {
	return b.Height
}

func (b *Bullet) GetUuid() string {
	return b.uuid
}

func (b *Bullet) GetKey() int64 {
	return b.gameThingsKey
}

func (b *Bullet) SetKey(key int64) {
	b.gameThingsKey = key
}

func (b *Bullet) Update() error {
	if b.isShow {
		objs := QM.SameAreaGameObjects(b)
		for _, otherObj := range objs {
			if otherObj.GetUuid() != b.mainId && (otherObj.GetGameIdentity() == consts.GameIdentityUser || otherObj.GetGameIdentity() == consts.GameIdentityEnemy) && QM.PixelCollision(b, otherObj) {
				QM.RemoveGameObject(b)
				inter.GameThings.Delete(b.gameThingsKey)
				middObj := otherObj
				go func() {
					if status, ex := middObj.(inter.DamageableItems).DamageCalculation(b.damage, b.faceStatus); status {
						b.MainRole.AbsorbingExperience(ex)
					}
				}()
				otherX, otherY := otherObj.GetXY()
				otherW, otherH := otherObj.GetWidth(), otherObj.GetHeight()
				NewLightningBeginning(otherX+otherW/2-32, otherY+otherH-193, b.mainId, func() {
					NewCircleExplosion(otherX+otherW/2-128, otherY+otherH-140, b.mainId, func() {

					})
				})
				return nil
			}
		}
		if b.X >= inter.ScreenWidth-5 || b.X <= -5 {
			QM.RemoveGameObject(b)
			inter.GameThings.Delete(b.gameThingsKey)
		}
		b.Move(b.X+b.Width*b.faceStatus, b.Y)
		if b.pictureFrameCount%b.imageFrameSpeed == b.imageFrameSpeed-1 {

			if b.pictureFrame == b.imageFrame-1 {
				b.pictureFrame = 0
			} else {
				b.pictureFrame++
			}
		}
		if b.pictureFrameCount == b.imageFrameSpeed*100 {
			b.pictureFrameCount = 0
		}
		b.pictureFrameCount++
	}
	return nil
}

func (b *Bullet) Show() {
	b.isShow = true
}

func (b *Bullet) Hide() {
	b.isShow = false
}

func (b *Bullet) GetGameIdentity() consts.GameIdentityKey {
	return b.gameIdentity
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.isShow {
		subRect := image.Rect(b.pictureFrame*b.DrawWidth, 0, b.pictureFrame*b.DrawWidth+b.DrawWidth, b.Height)
		subImg := b.image.SubImage(subRect).(*ebiten.Image)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(b.X), float64(b.Y))
		screen.DrawImage(subImg, opts)
	}
}

func (b *Bullet) GetLayers() int {
	return b.layers
}

func (b *Bullet) bindImage(path string, imageFrame int) error {
	if imageCache, ok := ImageCache[path]; ok {
		b.image = imageCache.leftImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(path)
		ImageCache[path] = ImageCacheInfo{
			leftImage: image,
		}
		b.image = image
	}
	b.imageFrame = imageFrame
	b.imageFrameSpeed = 6
	return nil
}

func NewBullet(x, y int, faceStatus int, damage int, mainId string, role inter.Role) *Bullet {
	s := &Bullet{
		mainId:       mainId,
		MainRole:     role,
		X:            x,
		Y:            y,
		layers:       consts.LayersSkillKey,
		DrawWidth:    7,
		Width:        16,
		Height:       2,
		uuid:         uuid.New().String(),
		isShow:       true,
		faceStatus:   faceStatus,
		gameIdentity: consts.GameIdentityArms,
		damage:       damage,
	}
	s.bindImage("/Users/admin/Desktop/star_dust/image/bullet/bullet_1.png", 1)
	inter.GameThings.Put(s)
	QM.AddGameObject(s, true)
	return s
}
