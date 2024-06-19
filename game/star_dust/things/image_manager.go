package things

import "github.com/hajimehoshi/ebiten/v2"

type ImageCacheInfo struct {
	leftImage  *ebiten.Image
	rightImage *ebiten.Image
}

var ImageCache = make(map[string]ImageCacheInfo, 0)
