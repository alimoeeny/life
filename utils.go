package main

import "image/color"

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// neighbourCount calculates the Moore neighborhood of (x, y).
func neighbourCount(a []color.RGBA, width, height, x, y int) (int, color.RGBA) {
	c := 0
	r := int(lifeColor.R)
	g := int(lifeColor.G)
	b := int(lifeColor.B)
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}
			x2 := x + i
			y2 := y + j
			if x2 < 0 || y2 < 0 || width <= x2 || height <= y2 {
				continue
			}
			index := y2*width + x2
			if a[index] != bgColor {
				c++
				r += int(a[index].R)
				g += int(a[index].G)
				b += int(a[index].B)
			}
		}
	}
	if c > 0 {
		r /= c
		g /= c
		b /= c
	}

	return c, color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}
