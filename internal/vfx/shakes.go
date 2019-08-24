// Package vfx provides visual effects for use across package lines
package vfx

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
)

var (
	SmallShaker oak.ScreenShaker
)

func Init() {
	SmallShaker = oak.ScreenShaker{
		Random: false,
		Magnitude: floatgeom.Point2{
			3,
			3,
		},
	}
}
