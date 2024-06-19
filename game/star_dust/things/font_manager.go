package things

import (
	"golang.org/x/image/font"
)

type FontCacheInfo struct {
	Face font.Face
}

var FontCache = make(map[string]FontCacheInfo, 0)
