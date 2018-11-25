package main

import (
	"bytes"
	"image"
	"image/color"
)

func (m *Model) String() string {
	m.Lock()
	defer m.Unlock()

	var b bytes.Buffer
	for i := 0; i < m.n; i++ {
		for j := 0; j < m.n; j++ {
			if m.l[i][j] > 0 {
				b.WriteRune('■')
			} else {
				b.WriteRune('□')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

// Image renders a black and white image of the Ising model state.
func (m *Model) Image() image.Image {
	m.Lock()
	defer m.Unlock()

	img := image.NewRGBA(image.Rect(0, 0, m.n, m.n))
	for i := 0; i < m.n; i++ {
		for j := 0; j < m.n; j++ {
			if m.l[i][j] > 0 {
				img.Set(i, j, color.RGBA{68, 68, 68, 255})
			} else {
				img.Set(i, j, color.RGBA{255, 255, 249, 255})
			}
		}
	}
	return img
}
