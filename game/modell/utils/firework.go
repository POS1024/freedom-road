package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
	"math/rand"
)

const (
	fireMaxParticles = 500
)

type Fireworks struct {
	screenWidth  float64
	screenHeight float64
	fireworks    []*Firework
}

func (f *Fireworks) Update() error {
	if rand.Float64() < 0.05 {
		x := rand.Float64() * f.screenWidth
		y := rand.Float64() * f.screenHeight / 2
		f.fireworks = append(f.fireworks, NewFirework(x, y))
	}

	for _, fw := range f.fireworks {
		fw.Update()
	}

	aliveFireworks := f.fireworks[:0]
	for _, fw := range f.fireworks {
		if len(fw.Particles) > 0 {
			aliveFireworks = append(aliveFireworks, fw)
		}
	}
	f.fireworks = aliveFireworks
	return nil
}

func (f *Fireworks) Draw(screen *ebiten.Image) {
	for _, fw := range f.fireworks {
		fw.Draw(screen)
	}
}

func NewFireworks(screenWidth, screenHeight float64) *Fireworks {
	return &Fireworks{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		fireworks:    []*Firework{},
	}
}

type FireParticle struct {
	X, Y   float64
	VX, VY float64
	Life   float64
	Color  color.Color
}

type Firework struct {
	Particles []FireParticle
}

func NewFirework(x, y float64) *Firework {
	f := &Firework{
		Particles: make([]FireParticle, 0, fireMaxParticles),
	}

	// 创建爆炸点的粒子
	numParticles := rand.Intn(100) + 100 // 每个爆炸点的粒子数量
	for i := 0; i < numParticles; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*2 + 2
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed
		life := rand.Float64()*2 + 2
		color := color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}

		f.Particles = append(f.Particles, FireParticle{
			X:     x,
			Y:     y,
			VX:    vx,
			VY:    vy,
			Life:  life,
			Color: color,
		})
	}

	return f
}

func (f *Firework) Update() {
	aliveParticles := f.Particles[:0]
	for _, p := range f.Particles {
		p.X += p.VX
		p.Y += p.VY
		p.VY += 0.05 // 模拟重力
		p.Life -= 0.05

		if p.Life > 0 {
			aliveParticles = append(aliveParticles, p)
		}
	}
	f.Particles = aliveParticles
}

func (f *Firework) Draw(screen *ebiten.Image) {
	for _, p := range f.Particles {
		ebitenutil.DrawRect(screen, p.X, p.Y, 2, 2, p.Color)
	}
}
