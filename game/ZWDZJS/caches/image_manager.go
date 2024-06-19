package caches

import "github.com/hajimehoshi/ebiten/v2"

type ImageCacheInfo struct {
	LeftImage  *ebiten.Image
	RightImage *ebiten.Image
}

var ImageCache = make(map[string]ImageCacheInfo, 0)
