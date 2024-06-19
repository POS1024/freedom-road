package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

func FlipImageHorizontally(img *ebiten.Image) *ebiten.Image {
	w, h := img.Size()
	flippedImg := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(w), 0)
	flippedImg.DrawImage(img, op)
	return flippedImg
}

func FlipImageSegments(img *ebiten.Image, segmentWidth int) *ebiten.Image {
	w, h := img.Size()
	flippedImg := ebiten.NewImage(w, h)

	for x := 0; x < w; x += segmentWidth {
		width := segmentWidth
		if x+segmentWidth > w {
			width = w - x
		}

		// Extract the segment
		segment := img.SubImage(image.Rect(x, 0, x+width, h)).(*ebiten.Image)

		// Flip the segment horizontally
		flippedSegment := FlipImageHorizontally(segment)

		// Draw the flipped segment to the result image
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), 0)
		flippedImg.DrawImage(flippedSegment, op)
	}

	return flippedImg
}
