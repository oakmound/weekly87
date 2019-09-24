// Package vfx provides visual effects for use across package lines
package vfx

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
)

var (
	SmallShaker = oak.ScreenShaker{
			Random: true,
			Magnitude: floatgeom.Point2{
				8,
				12,
			},
	}
	VerySmallShaker = oak.ScreenShaker{
			Random: true,
			Magnitude: floatgeom.Point2{
				3,
				4,
			},
	}
)