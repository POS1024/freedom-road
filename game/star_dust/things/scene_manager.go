package things

import (
	"github.com/hajimehoshi/ebiten/v2"
	"star_dust/consts"
	"star_dust/inter"
)

type SceneManager struct {
	layers           int
	isShow           bool
	gameThingsKey    int64
	showingSceneMap  *SceneMap
	showingScenesKey consts.ScenesKey
	sceneMaps        map[consts.ScenesKey]*SceneMap
}

func (s *SceneManager) GetKey() int64 {
	return s.gameThingsKey
}

func (s *SceneManager) SetKey(key int64) {
	s.gameThingsKey = key
}

func (s *SceneManager) Update() error {
	if s.isShow {
		return s.showingSceneMap.Update()
	}
	return nil
}

func (s *SceneManager) Draw(screen *ebiten.Image) {
	if s.isShow {
		s.showingSceneMap.Draw(screen)
	}
}

func (s *SceneManager) GetLayers() int {
	return s.layers
}

func (s *SceneManager) Offset(offsetX, offsetY int, obj inter.GameObject) bool {
	if s.isShow {
		return s.showingSceneMap.Offset(offsetX, offsetY, obj)
	}
	return false
}

func (s *SceneManager) Show() {
	s.isShow = true
}

func (s *SceneManager) Hide() {
	s.isShow = false
}

func NewSceneManager(scenesKey consts.ScenesKey, isShow bool) *SceneManager {
	sm := NewSceneMap(0, 0, 9530, 540, "/Users/admin/Desktop/star_dust/image/2.png")
	s := &SceneManager{
		layers:           consts.LayersSceneKey,
		showingScenesKey: scenesKey,
		isShow:           isShow,
		showingSceneMap:  sm,
		sceneMaps:        make(map[consts.ScenesKey]*SceneMap, 0),
	}
	s.sceneMaps[scenesKey] = sm
	inter.GameThings.Put(s)
	return s
}

var SM = NewSceneManager(consts.ScenesHome, false)
