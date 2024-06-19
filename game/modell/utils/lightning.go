package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math/rand"
	"time"
)

type lightning struct {
	x      float64
	y      float64
	width  float64
	points []struct {
		x float64
		y float64
	}
	pointsNum  int
	isBeingDel bool
	endX       float64
	endY       float64
	newNums    *int
	nowFx      float64
}

func (l *lightning) Update() error {
	for {
		<-time.After(time.Second / 18000)
		if l.isBeingDel {
			if l.pointsNum > 0 {
				l.points = l.points[1:]
				l.pointsNum--
			} else {
				return nil
			}
		} else {
			if l.endY <= 600 {
				rand.Seed(time.Now().UnixNano())
				randomNumber := rand.Intn(20)
				randomNew := false
				if randomNumber == 0 {
					l.nowFx = -1 * l.nowFx
				}
				if *l.newNums < 20 {
					randomNumber2 := rand.Intn(80)
					if randomNumber2 == 0 {
						randomNew = true
					}
				}
				if randomNew {
					*l.newNums++
					ls := NewLightning(l.endX, l.endY, l.width, l.nowFx*-1, l.newNums)
					AllThings = append(AllThings, ls)
					go ls.Update()
				}
				l.points = append(l.points, struct {
					x float64
					y float64
				}{
					x: l.endX + l.nowFx*l.width,
					y: l.endY + l.width,
				})
				l.endX = l.endX + l.nowFx*l.width
				l.endY = l.endY + l.width
				l.pointsNum++
			} else {
				l.isBeingDel = true
			}
		}

	}

	return nil
}

func (l *lightning) Draw(screen *ebiten.Image) {
	for _, s := range l.points {
		ebitenutil.DrawRect(screen, s.x, s.y, l.width, l.width, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	}
}

func NewLightning(x, y, width, nowFx float64, newNums *int) *lightning {
	return &lightning{
		x:       x,
		y:       y,
		width:   width,
		endX:    x,
		endY:    y,
		newNums: newNums,
		nowFx:   nowFx,
	}
}
