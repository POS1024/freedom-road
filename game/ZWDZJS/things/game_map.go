package things

import (
	"ZWDZJS/caches"
	"ZWDZJS/consts"
	"ZWDZJS/inter"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type GameMap struct {
	layers        int
	gameIdentity  consts.GameIdentityKey
	gameThingsKey int64
	offsetX       int
	image         *ebiten.Image
	isShow        bool
	enemy         []*Enemy
}

func (g *GameMap) GetKey() int64 {
	return g.gameThingsKey
}

func (g *GameMap) SetKey(key int64) {
	g.gameThingsKey = key
}

func (g *GameMap) Update() error {
	return nil
}

func (g *GameMap) Draw(screen *ebiten.Image) {
	if g.isShow {
		subRect := image.Rect(g.offsetX, 0, g.offsetX+consts.ScreenWidth, consts.ScreenHeight)
		subImg := g.image.SubImage(subRect).(*ebiten.Image)
		screen.DrawImage(subImg, nil)
	}
}

func (g *GameMap) GetLayers() int {
	return g.layers
}

func (g *GameMap) bindImage(path string) error {
	if imageCache, ok := caches.ImageCache[path]; ok {
		g.image = imageCache.LeftImage
	} else {
		image, _, _ := ebitenutil.NewImageFromFile(path)
		caches.ImageCache[path] = caches.ImageCacheInfo{
			LeftImage: image,
		}
		g.image = image
	}
	return nil
}

func NewGameMap() *GameMap {
	m := &GameMap{
		offsetX: 0,
		isShow:  false,
		layers:  consts.LayersMapKey,
	}
	m.bindImage("/Users/admin/Desktop/ZWDZJS/image/map/background1.png")
	inter.GameThings.Put(m)
	return m
}

var GameMapper = NewGameMap()
