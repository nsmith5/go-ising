package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"code.nfsmith.ca/nsmith/mjpeg"
)

type Model struct {
	Alpha float64
	Beta  float64

	de [][]int // Delta Energy for a flip at each site
	l  [][]int // Lattice
	n  int     // Lattice size (n x n)
}

func NewModel(N int) (*Model, error) {
	if !powerOfTwo(N) {
		return nil, errors.New("Model must have dimension power of 2")
	}

	m := Model{Alpha: 0.3, Beta: 3.0, n: N}
	m.de = make([][]int, N)
	m.l = make([][]int, N)

	for i := 0; i < N; i++ {
		m.de[i] = make([]int, N)
		m.l[i] = make([]int, N)
		for j := 0; j < N; j++ {
			m.l[i][j] = (rand.Intn(2) * 2) - 1 // +1 or -1
		}
	}

	return &m, nil
}

func (m *Model) String() string {
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

func (m *Model) Image() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, m.n, m.n))
	for i := 0; i < m.n; i++ {
		for j := 0; j < m.n; j++ {
			if m.l[i][j] > 0 {
				img.Set(i, j, color.Black)
			} else {
				img.Set(i, j, color.RGBA{255, 255, 255, 0})
			}
		}
	}
	return img
}

func (m *Model) Step() {
	// Compute the change in energy from a flip at each site
	for i := 0; i < m.n; i++ {
		for j := 0; j < m.n; j++ {
			m.de[i][j] = 2 * m.l[i][j] * (m.l[(i+1)&(m.n-1)][j] +
				m.l[(i-1)&(m.n-1)][j] +
				m.l[i][(j+1)&(m.n-1)] +
				m.l[i][(j-1)&(m.n-1)])
		}
	}

	// Flipping loop
	for i := 0; i < m.n; i++ {
		go func(i int) {
			for j := 0; j < m.n; j++ {
				if m.de[i][j] < 0 {
					if rand.Float64() < m.Alpha {
						m.l[i][j] = -m.l[i][j]
					}
				} else {
					if rand.Float64() < m.Alpha*math.Exp(-m.Beta*float64(m.de[i][j])) {
						m.l[i][j] = -m.l[i][j]
					}
				}
			}
		}(i)
	}

	return
}

func powerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	for n > 1 {
		if (n & 1) == 1 {
			return false
		}
		n = n >> 1
	}
	return true
}

const index = `
<head>
<style>
body {
	background-color: red;
}
</style>
</head>
<body>
<img width="512px", src="/ising"></img>
</body
`

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, index)
}

func main() {
	m, err := NewModel(1 << 8)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			m.Step()
		}
	}()

	handler := mjpeg.Handler{
		Stream: func() image.Image {
			time.Sleep(100 * time.Millisecond)
			return m.Image()
		},
		Options: nil,
	}

	m.Beta = 1.0

	mux := http.NewServeMux()
	mux.Handle("/ising", handler)
	mux.HandleFunc("/", IndexHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))

}
