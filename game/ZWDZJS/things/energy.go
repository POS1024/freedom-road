package things

import (
	"ZWDZJS/caches"
	"ZWDZJS/consts"
	"ZWDZJS/inter"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
)

type Energy struct {
	layers                             int
	gameIdentity                       consts.GameIdentityKey
	gameThingsKey                      int64
	X, Y, Width, Height                float64
	isMove, isShow, isHover, isClicked bool
	xLock, yLock                       bool
	xSpeed, ySpeed                     float64
	cor                                color.RGBA
	totalFrame                         int64
	pictureFrame                       int
	pictureFrameCount                  int
	image                              *ebiten.Image
	imageFrame, imageFrameSpeed        int
	energy                             int
	endX, endY                         float64
	isShowBorder                       bool
}

func (l *Energy) IsHover() bool {
	xInt, yInt := ebiten.CursorPosition()
	x, y := float64(xInt), float64(yInt)
	if x > l.X && x < l.X+l.Width && y > l.Y && y < l.Y+l.Height {
		l.cor = color.RGBA{255, 0, 0, 255}
		l.isHover = true
		return true
	} else {
		l.cor = color.RGBA{255, 255, 0, 255}
		l.isHover = false
		return false
	}
}

func (l *Energy) ChangeClick() {
	l.isMove = false
	centerX, centerY := l.X+l.Width/2, l.Y+l.Height/2
	xDistance, yDistance := centerX-float64(consts.GameDealerEnergyX), centerY-float64(consts.GameDealerEnergyY)
	if xDistance > yDistance {
		l.xSpeed = consts.GameDealerEnergySpeed
		l.ySpeed = consts.GameDealerEnergySpeed * yDistance / xDistance
	} else {
		l.ySpeed = consts.GameDealerEnergySpeed
		l.xSpeed = consts.GameDealerEnergySpeed * xDistance / yDistance
	}
	l.isClicked = true

}

func (l *Energy) IsClick() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && l.isHover && !l.isClicked {
		l.ChangeClick()
		return true
	} else {
		return false
	}
}

func (l *Energy) IsDown() bool {
	return true
}

func (l *Energy) GetKey() int64 {
	return l.gameThingsKey
}

func (l *Energy) SetKey(key int64) {
	l.gameThingsKey = key
}

func (l *Energy) Update() error {
	if l.isShow {

		if l.isClicked {
			centerX, centerY := l.X+l.Width/2, l.Y+l.Height/2
			if (centerX - consts.GameDealerEnergyX) < 10 {
				l.xLock = true
			}
			if (centerY - consts.GameDealerEnergyY) < 10 {
				l.yLock = true
			}
			if (centerX-consts.GameDealerEnergyX) < 10 && (centerY-consts.GameDealerEnergyY) < 10 {
				GameDealerMR.AddEnergy(l.energy)
				inter.GameThings.Delete(l.gameThingsKey)
				return nil
			}

			if !l.xLock {
				if centerX > consts.GameDealerEnergyX {
					l.X = l.X - l.xSpeed
				} else if centerX < consts.GameDealerEnergyX {
					l.X = l.X + l.xSpeed
				}
			}

			if !l.yLock {
				if centerY > consts.GameDealerEnergyY {
					l.Y = l.Y - l.ySpeed
				} else if centerY < consts.GameDealerEnergyY {
					l.Y = l.Y + l.ySpeed
				}
			}

		}

		if l.isMove {
			endMoveX, endMoveY := false, false
			if (l.X-l.endX) > 10 && l.X > l.endX {
				l.X -= consts.GameDealerEnergySpeed / 6
			} else {
				endMoveX = true
			}
			if (l.endY-l.Y) > 10 && l.Y < l.endY {
				l.Y += consts.GameDealerEnergySpeed / 6
			} else {
				endMoveY = true
			}
			if endMoveX && endMoveY {
				l.totalFrame++
			}
			if l.totalFrame > 240 {
				l.ChangeClick()
			}
		}

		l.IsHover()
		l.IsClick()
		if l.pictureFrameCount%l.imageFrameSpeed == l.imageFrameSpeed-1 {

			if l.pictureFrame == l.imageFrame-1 {
				l.pictureFrame = 0
			} else {
				l.pictureFrame++
			}
		}
		if l.pictureFrameCount == l.imageFrameSpeed*100 {
			l.pictureFrameCount = 0
		}
		l.pictureFrameCount++
	}
	return nil
}

func (l *Energy) Draw(screen *ebiten.Image) {
	if l.isShow {
		subRect := image.Rect(l.pictureFrame*int(l.Width), 0, l.pictureFrame*int(l.Width)+int(l.Width), int(l.Height))
		subImg := l.image.SubImage(subRect).(*ebiten.Image)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(l.X, l.Y)
		screen.DrawImage(subImg, opts)
		if l.isShowBorder {
			l.drawRectOutline(screen, int(l.X), int(l.Y), int(l.Width), int(l.Height), l.cor)
		}
	}
}

func (l *Energy) drawRectOutline(screen *ebiten.Image, x, y, width, height int, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x+width), float64(y), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y+height), float64(x+width), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x), float64(y+height), clr)
	ebitenutil.DrawLine(screen, float64(x+width), float64(y), float64(x+width), float64(y+height), clr)
}

func (l *Energy) GetLayers() int {
	return l.layers
}

func (l *Energy) bindImage(path string, imageFrame, imageFrameSpeed int) error {
	if imageCache, ok := caches.ImageCache[path]; ok {
		l.image = imageCache.LeftImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(path)
		caches.ImageCache[path] = caches.ImageCacheInfo{
			LeftImage: image,
		}
		l.image = image
	}
	l.imageFrame = imageFrame
	l.imageFrameSpeed = imageFrameSpeed
	return nil
}

func NewEnergy(X, Y, endX, endY float64, energy int) *Energy {
	l := &Energy{
		X:      X,
		Y:      Y,
		Width:  78,
		Height: 78,
		isShow: true,
		layers: consts.LayersEnergyKey,
		cor:    color.RGBA{255, 255, 0, 255},
		energy: energy,
		isMove: true,
		endX:   endX,
		endY:   endY,
	}
	l.bindImage("/Users/admin/Desktop/ZWDZJS/image/items/sun.png", 21, 2)
	inter.GameThings.Put(l)
	return l
}
