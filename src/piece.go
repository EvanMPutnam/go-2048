package main

import "image/color"

type piece struct {
	value     int
	hasMerged bool
}

func (p *piece) determineColor() color.Color {
	switch p.value {
	case 0:
		return color.RGBA{R: 187, G: 173, B: 160}
	case 2:
		return color.RGBA{R: 238, G: 228, B: 218}
	case 4:
		return color.RGBA{R: 237, G: 224, B: 200}
	case 8:
		return color.RGBA{R: 242, G: 217, B: 121}
	case 16:
		return color.RGBA{R: 245, G: 149, B: 99}
	case 32:
		return color.RGBA{R: 246, G: 124, B: 95}
	case 64:
		return color.RGBA{R: 246, G: 94, B: 59}
	case 128:
		return color.RGBA{R: 237, G: 207, B: 114}
	case 256:
		return color.RGBA{R: 237, G: 204, B: 97}
	case 512:
		return color.RGBA{R: 237, G: 200, B: 80}
	case 1024:
		return color.RGBA{R: 237, G: 197, B: 63}
	case 2048:
		return color.RGBA{R: 237, G: 194, B: 46}
	default:
		return color.RGBA{R: 237, G: 194, B: 46}
	}
}
