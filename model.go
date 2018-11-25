package main

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Model is an Ising model
type Model struct {
	sync.Mutex
	Alpha float64
	Beta  float64

	de [][]int // Delta Energy for a flip at each site
	l  [][]int // Lattice
	n  int     // Lattice size (n x n)
}

// NewModel returns a NxN Ising model
func NewModel(N int) (*Model, error) {
	if !powerOfTwo(N) {
		return nil, errors.New("Model must have dimension power of 2")
	}

	m := Model{Alpha: 0.3, Beta: 6.0}
	m.n = N
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

// Step advances the Ising Model by one Metropolis step.
func (m *Model) Step() {
	m.Lock()
	defer m.Unlock()

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

// Run steps the model once per 100 ms for ever.
func (m *Model) Run() {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			m.Step()
		}
	}()
}

// powerOfTwo checks if a number is a power of two or not.
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
