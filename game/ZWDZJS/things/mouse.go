package things

import (
	"ZWDZJS/consts"
	"ZWDZJS/inter"
	"github.com/hajimehoshi/ebiten/v2"
)

type Mouse struct {
	layers         int
	gameIdentity   consts.GameIdentityKey
	gameThingsKey  int64
	nowMousePuppet inter.MousePuppet
}

func (m *Mouse) GetMouseType() consts.MouseTypeKey {
	if m.nowMousePuppet != nil {
		return m.nowMousePuppet.GetMouseType()
	}
	return consts.MouseTypeNil
}

func (m *Mouse) BindMousePuppet(mousePuppet inter.MousePuppet) {
	if m.nowMousePuppet != nil {
		m.nowMousePuppet.Cancel()
	}
	m.nowMousePuppet = mousePuppet
}

func (m *Mouse) CancelMousePuppet() {
	if m.nowMousePuppet != nil {
		m.nowMousePuppet.Cancel()
	}
	m.nowMousePuppet = nil
}

func (m *Mouse) ConfirmMousePuppet() bool {
	if m.nowMousePuppet != nil {
		status := m.nowMousePuppet.Confirm()
		m.nowMousePuppet = nil
		return status
	}
	return false
}

func (m *Mouse) GetKey() int64 {
	return m.gameThingsKey
}

func (m *Mouse) SetKey(key int64) {
	m.gameThingsKey = key
}

func (m *Mouse) Update() error {
	return nil
}

func (m *Mouse) Draw(screen *ebiten.Image) {
	if m.nowMousePuppet != nil {
		m.nowMousePuppet.MouseEffects(screen)
	}
}

func (m *Mouse) GetLayers() int {
	return m.layers
}

func NewMouse(layer int) *Mouse {
	m := &Mouse{
		layers: layer,
	}
	inter.GameThings.Put(m)
	return m
}

var MOUSE = NewMouse(consts.LayersMouseKey)
