// Package vfx provides visual effects for use across package lines
package vfx

import (
	"sync"
	"time"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
)

var (
	SmallShaker Shaker
)

func init() {
	SmallShaker = Shaker{
		mu: &sync.Mutex{},
		shaker: oak.ScreenShaker{
			Random: true,
			Magnitude: floatgeom.Point2{
				8,
				12,
			},
		},
		shakingEnd: time.Now(),
	}
}

// Shaker exposes a simple oak shaker that refuses shakes if already shaking
// This should be updated later
type Shaker struct {
	mu         *sync.Mutex
	shaker     oak.ScreenShaker
	shakingEnd time.Time
}

func (s *Shaker) Shake(shaking time.Duration) {
	// s.mu.Lock()
	// if time.Now().Before(s.shakingEnd) {
	// 	s.mu.Unlock()
	// 	return
	// }
	// s.shakingEnd = time.Now().Add(shaking).Add(time.Second)
	// s.mu.Unlock()

	s.shaker.Shake(shaking)
}
