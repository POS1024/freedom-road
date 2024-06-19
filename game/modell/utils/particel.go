package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math/rand"
)

const (
	maxParticles = 1000
)

type Particle struct {
	X, Y   float64
	VX, VY float64
	Life   float64
	Color  color.Color
}

type ParticleSystem struct {
	particles []Particle
	x         float64
	y         float64
}

func NewParticleSystem(x, y float64) *ParticleSystem {
	return &ParticleSystem{
		x:         x,
		y:         y,
		particles: make([]Particle, 0, maxParticles),
	}
}

func (ps *ParticleSystem) Update() error {
	for i := 0; i < 10; i++ {
		ps.AddParticle()
	}

	aliveParticles := ps.particles[:0]
	for _, p := range ps.particles {
		p.X += p.VX
		p.Y += p.VY
		p.Life -= 0.05

		if p.Life > 0 {
			aliveParticles = append(aliveParticles, p)
		}
	}
	ps.particles = aliveParticles
	return nil
}

func (ps *ParticleSystem) AddParticle() {
	vx := rand.Float64()*4 - 2
	vy := -rand.Float64()*3 - 1
	life := rand.Float64()*2 + 1
	color := color.RGBA{R: 255, G: uint8(rand.Intn(156)), B: 0, A: 255}

	ps.particles = append(ps.particles, Particle{
		X:     ps.x,
		Y:     ps.y,
		VX:    vx,
		VY:    vy,
		Life:  life,
		Color: color,
	})
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for _, p := range ps.particles {
		ebitenutil.DrawRect(screen, p.X, p.Y, 3, 3, p.Color)
	}
}
